/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package genidindex

import (
	"context"
	"fmt"
	"reflect"

	"github.com/henderiw/logger/log"
	genidbev1alpha1 "github.com/kuidio/kuid/apis/backend/genid/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	"github.com/kuidio/kuid/pkg/backend/backend"
	"github.com/kuidio/kuid/pkg/reconcilers"
	"github.com/kuidio/kuid/pkg/reconcilers/ctrlconfig"
	"github.com/kuidio/kuid/pkg/reconcilers/eventhandler"
	"github.com/kuidio/kuid/pkg/reconcilers/resource"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func init() {
	reconcilers.Register("genidindex", &reconciler{})
}

const (
	crName         = "genidindex"
	controllerName = "GENIDIndexController"
	finalizer      = "genidindex.genid.res.kuid.dev/finalizer"
	// errors
	errGetCr        = "cannot get cr"
	errUpdateStatus = "cannot update status"
)

// SetupWithManager sets up the controller with the Manager.
func (r *reconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, c interface{}) (map[schema.GroupVersionKind]chan event.GenericEvent, error) {

	cfg, ok := c.(*ctrlconfig.ControllerConfig)
	if !ok {
		return nil, fmt.Errorf("cannot initialize, expecting controllerConfig, got: %s", reflect.TypeOf(c).Name())
	}

	r.Client = mgr.GetClient()
	r.finalizer = resource.NewAPIFinalizer(mgr.GetClient(), finalizer)
	r.recorder = mgr.GetEventRecorderFor(controllerName)
	r.be = cfg.GENIDBackend

	return nil, ctrl.NewControllerManagedBy(mgr).
		Named(controllerName).
		For(&genidbev1alpha1.GENIDIndex{}).
		Watches(&genidbev1alpha1.GENIDIndex{},
			&eventhandler.GENIDEntryEventHandler{
				Client:  mgr.GetClient(),
				ObjList: &genidbev1alpha1.GENIDIndexList{},
			}).
		Complete(r)
}

type reconciler struct {
	client.Client
	finalizer *resource.APIFinalizer
	recorder  record.EventRecorder
	be        backend.Backend
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = ctrlconfig.InitContext(ctx, controllerName, req.NamespacedName)
	log := log.FromContext(ctx)
	log.Info("reconcile")

	cr := &genidbev1alpha1.GENIDIndex{}
	if err := r.Get(ctx, req.NamespacedName, cr); err != nil {
		// if the resource no longer exists the reconcile loop is done
		if resource.IgnoreNotFound(err) != nil {
			log.Error(errGetCr, "error", err)
			return ctrl.Result{}, errors.Wrap(resource.IgnoreNotFound(err), errGetCr)
		}
		return ctrl.Result{}, nil
	}
	cr = cr.DeepCopy()

	if !cr.GetDeletionTimestamp().IsZero() {
		if err := r.deleteIndex(ctx, cr); err != nil {
			return ctrl.Result{Requeue: true}, errors.Wrap(r.Update(ctx, cr), errUpdateStatus)
		}

		if err := r.finalizer.RemoveFinalizer(ctx, cr); err != nil {
			r.handleError(ctx, cr, "cannot remove finalizer", err)
			return ctrl.Result{Requeue: true}, errors.Wrap(r.Update(ctx, cr), errUpdateStatus)
		}
		log.Debug("Successfully deleted resource")
		return ctrl.Result{}, nil
	}

	if err := r.finalizer.AddFinalizer(ctx, cr); err != nil {
		r.handleError(ctx, cr, "cannot add finalizer", err)
		return ctrl.Result{Requeue: true}, errors.Wrap(r.Update(ctx, cr), errUpdateStatus)
	}

	if r.hasMinMaxRangeChanged(cr) {
		// delete index
		if err := r.deleteIndex(ctx, cr); err != nil {
			return ctrl.Result{Requeue: true}, errors.Wrap(r.Update(ctx, cr), errUpdateStatus)
		}
	}
	// create index
	if err := r.applyIndex(ctx, cr); err != nil {
		return ctrl.Result{Requeue: true}, errors.Wrap(r.Update(ctx, cr), errUpdateStatus)
	}

	if err := r.applyMinMaxRange(ctx, cr); err != nil {
		return ctrl.Result{Requeue: true}, errors.Wrap(r.Update(ctx, cr), errUpdateStatus)
	}

	cr.SetConditions(conditionv1alpha1.Ready())
	cr.Status.MinID = cr.Spec.MinID
	cr.Status.MaxID = cr.Spec.MaxID
	r.recorder.Eventf(cr, corev1.EventTypeNormal, crName, "ready")
	return ctrl.Result{}, errors.Wrap(r.Update(ctx, cr), errUpdateStatus)
}

func (r *reconciler) handleError(ctx context.Context, cr *genidbev1alpha1.GENIDIndex, msg string, err error) {
	log := log.FromContext(ctx)
	if err == nil {
		cr.SetConditions(conditionv1alpha1.Failed(msg))
		log.Error(msg)
		r.recorder.Eventf(cr, corev1.EventTypeWarning, crName, msg)
	} else {
		cr.SetConditions(conditionv1alpha1.Failed(err.Error()))
		log.Error(msg, "error", err)
		r.recorder.Eventf(cr, corev1.EventTypeWarning, crName, fmt.Sprintf("%s, err: %s", msg, err.Error()))
	}
}

func (r *reconciler) deleteIndex(ctx context.Context, cr *genidbev1alpha1.GENIDIndex) error {
	if err := r.be.DeleteIndex(ctx, cr); err != nil {
		r.handleError(ctx, cr, "cannot delete index", err)
		return err
	}
	return nil
}

func (r *reconciler) applyIndex(ctx context.Context, cr *genidbev1alpha1.GENIDIndex) error {
	if err := r.be.CreateIndex(ctx, cr); err != nil {
		r.handleError(ctx, cr, "cannot create index", err)
		return err
	}
	return nil
}

func (r *reconciler) hasMinMaxRangeChanged(cr *genidbev1alpha1.GENIDIndex) bool {
	return changed(cr.Status.MinID, cr.Spec.MinID) || changed(cr.Status.MaxID, cr.Spec.MaxID)
}

func changed(status, spec *int64) bool {
	if status != nil {
		if spec == nil {
			return true
		} else {
			if *status != *spec {
				return true
			}
		}
	}
	return false
}

func (r *reconciler) applyMinMaxRange(ctx context.Context, cr *genidbev1alpha1.GENIDIndex) error {
	if cr.Spec.MinID != nil && *cr.Spec.MinID != genidbev1alpha1.GENIDID_Min {
		claim := cr.GetMinClaim()
		if err := r.be.Claim(ctx, claim); err != nil {
			r.handleError(ctx, cr, "cannot claim min reserved range", err)
			return err
		}
	}
	if cr.Spec.MaxID != nil && *cr.Spec.MaxID != genidbev1alpha1.GENIDID_MaxValue[genidbev1alpha1.GetGenIDType(cr.Spec.Type)] {
		claim := cr.GetMaxClaim()
		if err := r.be.Claim(ctx, claim); err != nil {
			r.handleError(ctx, cr, "cannot claim max reserved range", err)
			return err
		}
	}
	return nil
}
