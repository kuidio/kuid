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

package vlanindex

import (
	"context"
	"fmt"
	"reflect"

	"github.com/henderiw/logger/log"
	condv1alpha1 "github.com/kform-dev/choreo/apis/condition/v1alpha1"
	"github.com/kuidio/kuid/apis/backend/vlan"
	vlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	"github.com/kuidio/kuid/pkg/backend"
	"github.com/kuidio/kuid/pkg/reconcilers"
	"github.com/kuidio/kuid/pkg/reconcilers/ctrlconfig"
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
	reconcilers.Register(vlan.GroupName, vlanbev1alpha1.VLANIndexKind, &reconciler{})
}

const (
	reconcilerName = "VLANIndexController"
	finalizer      = "vlanindex.vlan.be.kuid.dev/finalizer"
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
	r.finalizer = resource.NewAPIFinalizer(mgr.GetClient(), finalizer, reconcilerName)
	r.recorder = mgr.GetEventRecorderFor(reconcilerName)
	r.be = cfg.Backends[vlanbev1alpha1.SchemeGroupVersion.Group]

	return nil, ctrl.NewControllerManagedBy(mgr).
		Named(reconcilerName).
		For(&vlanbev1alpha1.VLANIndex{}).
		Complete(r)
}

type reconciler struct {
	client.Client
	finalizer *resource.APIFinalizer
	recorder  record.EventRecorder
	be        backend.Backend
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = ctrlconfig.InitContext(ctx, reconcilerName, req.NamespacedName)
	log := log.FromContext(ctx)
	log.Info("reconcile")

	index := &vlanbev1alpha1.VLANIndex{}
	if err := r.Get(ctx, req.NamespacedName, index); err != nil {
		// if the resource no longer exists the reconcile loop is done
		if resource.IgnoreNotFound(err) != nil {
			log.Error(errGetCr, "error", err)
			return ctrl.Result{}, errors.Wrap(resource.IgnoreNotFound(err), errGetCr)
		}
		return ctrl.Result{}, nil
	}
	indexOrig := index.DeepCopy()
	log.Debug("reconcile", "status orig", index.Status)

	if !index.GetDeletionTimestamp().IsZero() {
		// Prefixes are not to be deleted as the sync delete index takes care and garbage collector 
		// takes care of this
		intIndex := &vlan.VLANIndex{}
		if err := vlanbev1alpha1.Convert_v1alpha1_VLANIndex_To_vlan_VLANIndex(index, intIndex, nil); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, indexOrig, "cannot convert index before delete", err), errUpdateStatus)
		}
		if err := r.be.DeleteIndex(ctx, intIndex); err != nil {
			if resource.IgnoreNotFound(err) != nil {
				return ctrl.Result{Requeue: true},
					errors.Wrap(r.handleError(ctx, indexOrig, "cannot delete index", err), errUpdateStatus)
			}
		}
		if err := vlanbev1alpha1.Convert_vlan_VLANIndex_To_v1alpha1_VLANIndex(intIndex, index, nil); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, indexOrig, "cannot convert index after delete", err), errUpdateStatus)
		}

		// We use owner reference so the k8s garbage collector takes care of the cleanup
		if err := r.finalizer.RemoveFinalizer(ctx, index); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, indexOrig, "cannot remove finalizer", err), errUpdateStatus)
		}
		return ctrl.Result{}, nil
	}

	if err := r.finalizer.AddFinalizer(ctx, index); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, indexOrig, "cannot add finalizer", err), errUpdateStatus)
	}

	// create index
	intIndex := &vlan.VLANIndex{}
	if err := vlanbev1alpha1.Convert_v1alpha1_VLANIndex_To_vlan_VLANIndex(index, intIndex, nil); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, indexOrig, "cannot convert index before create", err), errUpdateStatus)
	}
	if err := r.be.CreateIndex(ctx, intIndex); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, indexOrig, "cannot apply index", err), errUpdateStatus)
	}
	if err := vlanbev1alpha1.Convert_vlan_VLANIndex_To_v1alpha1_VLANIndex(intIndex, index, nil); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, indexOrig, "cannot convert index after create", err), errUpdateStatus)
	}

	// updating the index is taken care of by the createIndex code

	return ctrl.Result{}, errors.Wrap(r.handleSuccess(ctx, indexOrig), errUpdateStatus)
}

func (r *reconciler) handleSuccess(ctx context.Context, index *vlanbev1alpha1.VLANIndex) error {
	// take a snapshot of the current object
	patch := client.MergeFrom(index.DeepCopy())
	// update status
	index.SetConditions(condv1alpha1.Ready())
	r.recorder.Eventf(index, corev1.EventTypeNormal, vlanbev1alpha1.VLANIndexKind, "ready")

	return r.Client.Status().Patch(ctx, index, patch, &client.SubResourcePatchOptions{
		PatchOptions: client.PatchOptions{
			FieldManager: "backend",
		},
	})
}

func (r *reconciler) handleError(ctx context.Context, index *vlanbev1alpha1.VLANIndex, msg string, err error) error {
	log := log.FromContext(ctx)
	// take a snapshot of the current object
	patch := client.MergeFrom(index.DeepCopy())

	if err != nil {
		msg = fmt.Sprintf("%s err %s", msg, err.Error())
	}
	index.SetConditions(condv1alpha1.Failed(msg))
	log.Error(msg)
	r.recorder.Eventf(index, corev1.EventTypeWarning, vlanbev1alpha1.VLANIndexKind, msg)

	return r.Client.Status().Patch(ctx, index, patch, &client.SubResourcePatchOptions{
		PatchOptions: client.PatchOptions{
			FieldManager: "backend",
		},
	})
}