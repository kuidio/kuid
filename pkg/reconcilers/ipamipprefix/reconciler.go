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

package ipamipprefix

import (
	"context"
	"fmt"

	"github.com/henderiw/logger/log"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	ipamresv1alpha1 "github.com/kuidio/kuid/apis/resource/ipam/v1alpha1"
	"github.com/kuidio/kuid/pkg/reconcilers"
	"github.com/kuidio/kuid/pkg/reconcilers/ctrlconfig"
	"github.com/kuidio/kuid/pkg/reconcilers/resource"
	"github.com/kuidio/kuid/pkg/reconcilers/ipentryeventhandler"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func init() {
	reconcilers.Register("ipamipprefix", &reconciler{})
}

const (
	crName         = "ipPrefix"
	controllerName = "IPAMIPPrefixController"
	finalizer      = "ipprefix.ipam.res.kuid.dev/finalizer"
	// errors
	errGetCr        = "cannot get cr"
	errUpdateStatus = "cannot update status"
)

/*
type adder interface {
	Add(item interface{})
}
*/

//+kubebuilder:rbac:groups=ipprefix.ipam.res.kuid.dev,resources=ipprefixes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ipprefix.ipam.res.kuid.dev,resources=ipprefixes/status,verbs=get;update;patch

// SetupWithManager sets up the controller with the Manager.
func (r *reconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, c interface{}) (map[schema.GroupVersionKind]chan event.GenericEvent, error) {
	r.Client = mgr.GetClient()
	r.finalizer = resource.NewAPIFinalizer(mgr.GetClient(), finalizer)
	r.recorder = mgr.GetEventRecorderFor(controllerName)

	return nil, ctrl.NewControllerManagedBy(mgr).
		Named(controllerName).
		For(&ipamresv1alpha1.IPPrefix{}).
		Watches(&ipambev1alpha1.IPEntry{},
			&ipentryeventhandler.IPEntryEventHandler{
				Client:  mgr.GetClient(),
				ObjList: &ipamresv1alpha1.IPPrefixList{},
			}).
		Complete(r)
}

type reconciler struct {
	client.Client
	finalizer *resource.APIFinalizer
	recorder  record.EventRecorder
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = ctrlconfig.InitContext(ctx, controllerName, req.NamespacedName)
	log := log.FromContext(ctx)
	log.Info("reconcile")

	cr := &ipamresv1alpha1.IPPrefix{}
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
		if err := r.deleteIPClaim(ctx, cr); err != nil {
			// error already handled
			return ctrl.Result{Requeue: true}, errors.Wrap(r.Status().Update(ctx, cr), errUpdateStatus)
		}

		if err := r.finalizer.RemoveFinalizer(ctx, cr); err != nil {
			r.handleError(ctx, cr, "cannot remove finalizer", err)
			return ctrl.Result{Requeue: true}, errors.Wrap(r.Status().Update(ctx, cr), errUpdateStatus)
		}
		log.Debug("Successfully deleted resource")
		return ctrl.Result{}, nil
	}

	if err := r.finalizer.AddFinalizer(ctx, cr); err != nil {
		r.handleError(ctx, cr, "cannot add finalizer", err)
		return ctrl.Result{Requeue: true}, errors.Wrap(r.Status().Update(ctx, cr), errUpdateStatus)
	}

	if err := r.applyIPClaim(ctx, cr); err != nil {
		// error already handled
		return ctrl.Result{Requeue: true}, errors.Wrap(r.Status().Update(ctx, cr), errUpdateStatus)
	}

	cr.SetConditions(conditionv1alpha1.Ready())
	r.recorder.Eventf(cr, corev1.EventTypeNormal, crName, "ready")
	return ctrl.Result{}, errors.Wrap(r.Status().Update(ctx, cr), errUpdateStatus)
}

func (r *reconciler) handleError(ctx context.Context, cr *ipamresv1alpha1.IPPrefix, msg string, err error) {
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

func (r *reconciler) deleteIPClaim(ctx context.Context, cr *ipamresv1alpha1.IPPrefix) error {
	ipclaim, err := ipambev1alpha1.GetIPClaimFromPrefix(cr.Spec.Kind, cr.Spec.NetworkInstance, cr.Spec.Prefix, cr.Spec.UserDefinedLabels, cr)
	if err != nil { // strange if this happens since the prefix was already processed
		r.handleError(ctx, cr, "cannot build ipclaim", err)
		return err
	}
	if err := r.Client.Get(ctx, ipclaim.GetNamespacedName(), ipclaim); err != nil {
		if resource.IgnoreNotFound(err) != nil {
			msg := fmt.Sprintf("cannot get ipclaim %s", ipclaim.GetNamespacedName())
			r.handleError(ctx, cr, msg, err)
			return err
		}
		return nil
	}
	if err := r.Client.Delete(ctx, ipclaim); err != nil {
		if resource.IgnoreNotFound(err) != nil {
			msg := fmt.Sprintf("cannot get ipclaim %s", ipclaim.GetNamespacedName())
			r.handleError(ctx, cr, msg, err)
			return err
		}
	}
	return nil
}

func (r *reconciler) applyIPClaim(ctx context.Context, cr *ipamresv1alpha1.IPPrefix) error {
	ipclaim, err := ipambev1alpha1.GetIPClaimFromPrefix(cr.Spec.Kind, cr.Spec.NetworkInstance, cr.Spec.Prefix, cr.Spec.UserDefinedLabels, cr)
	if err != nil { // strange if this happens since the prefix was already processed
		r.handleError(ctx, cr, "build ipclaim", err)
		return err
	}
	if err := r.Client.Get(ctx, ipclaim.GetNamespacedName(), ipclaim); err != nil {
		if resource.IgnoreNotFound(err) != nil {
			msg := fmt.Sprintf("cannot get ipclaim %s", ipclaim.GetNamespacedName())
			r.handleError(ctx, cr, msg, err)
			return err
		}
		if err := r.Client.Create(ctx, ipclaim); err != nil {
			msg := fmt.Sprintf("cannot create ipclaim %s", ipclaim.GetNamespacedName())
			r.handleError(ctx, cr, msg, err)
			return err
		}
		return nil
	}
	// update the ipClaim with the latest spec information
	if err = ipclaim.UpdateSpecFromPrefix(cr.Spec.Kind, cr.Spec.NetworkInstance, cr.Spec.Prefix, cr.Spec.UserDefinedLabels, cr); err != nil {
		msg := fmt.Sprintf("cannot create ipclaim %s", ipclaim.GetNamespacedName())
		r.handleError(ctx, cr, msg, err)
		return err
	}
	if err := r.Client.Update(ctx, ipclaim); err != nil {
		msg := fmt.Sprintf("cannot update ipclaim %s", ipclaim.GetNamespacedName())
		r.handleError(ctx, cr, msg, err)
		return err
	}

	if ipclaim.Status.Prefix == nil || *ipclaim.Status.Prefix != cr.Spec.Prefix {
		//we got a different prefix than requested one
		msg := fmt.Sprintf("ip prefix not ready: req/rsp %s/%v", cr.Spec.Prefix, ipclaim.Status.Prefix)
		r.handleError(ctx, cr, msg, nil)
		return err
	}
	return nil
}
