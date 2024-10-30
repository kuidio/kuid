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

package ipclaim

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/henderiw/logger/log"
	condv1alpha1 "github.com/kform-dev/choreo/apis/condition/v1alpha1"
	"github.com/kuidio/kuid/apis/backend/ipam"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
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
	reconcilers.Register(ipambev1alpha1.Group, ipambev1alpha1.IPClaimKind, &reconciler{})
}

const (
	reconcilerName = "IPClaimController"
	finalizer      = "ipclaim.ipam.be.kuid.dev/finalizer"
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
	r.be = cfg.Backends[ipambev1alpha1.SchemeGroupVersion.Group]

	return nil, ctrl.NewControllerManagedBy(mgr).
		Named(reconcilerName).
		For(&ipambev1alpha1.IPClaim{}).
		//Watches(&ipambev1alpha1.IPEntry{},
		//	&eventhandler.IPEntryEventHandler{
		//		Client:  mgr.GetClient(),
		//		ObjList: &ipambev1alpha1.IPClaimList{},
		//	}).
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

	ipclaim := &ipambev1alpha1.IPClaim{}
	if err := r.Get(ctx, req.NamespacedName, ipclaim); err != nil {
		// if the resource no longer exists the reconcile loop is done
		if resource.IgnoreNotFound(err) != nil {
			log.Error(errGetCr, "error", err)
			return ctrl.Result{}, errors.Wrap(resource.IgnoreNotFound(err), errGetCr)
		}
		return ctrl.Result{}, nil
	}
	ipclaimOrig := ipclaim.DeepCopy()
	log.Debug("reconcile", "status orig", ipclaimOrig.Status)

	if !ipclaim.GetDeletionTimestamp().IsZero() {
		intIPClaim := &ipam.IPClaim{}
		if err := ipambev1alpha1.Convert_v1alpha1_IPClaim_To_ipam_IPClaim(ipclaim, intIPClaim, nil); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, ipclaimOrig, "cannot convert ipclaim before delete claim", err), errUpdateStatus)
		}

		if err := r.be.Release(ctx, intIPClaim, false); err != nil {
			if !strings.Contains(err.Error(), "not initialized") {
				return ctrl.Result{Requeue: true},
					errors.Wrap(r.handleError(ctx, ipclaimOrig, "cannot delete ipclaim", err), errUpdateStatus)
			}
		}
		if err := ipambev1alpha1.Convert_ipam_IPClaim_To_v1alpha1_IPClaim(intIPClaim, ipclaim, nil); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, ipclaimOrig, "cannot convert ipclaim after delete claim", err), errUpdateStatus)
		}

		if err := r.finalizer.RemoveFinalizer(ctx, ipclaim); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, ipclaimOrig, "cannot delete finalizer", err), errUpdateStatus)
		}
		return ctrl.Result{}, nil
	}

	if err := r.finalizer.AddFinalizer(ctx, ipclaim); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, ipclaimOrig, "cannot add finalizer", err), errUpdateStatus)
	}

	intIPClaim := &ipam.IPClaim{}
	if err := ipambev1alpha1.Convert_v1alpha1_IPClaim_To_ipam_IPClaim(ipclaim, intIPClaim, nil); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, ipclaimOrig, "cannot convert ipclaim before claim", err), errUpdateStatus)
	}
	if err := r.be.Claim(ctx, intIPClaim, false); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, ipclaimOrig, "cannot claim ip", err), errUpdateStatus)
	}
	if err := ipambev1alpha1.Convert_ipam_IPClaim_To_v1alpha1_IPClaim(intIPClaim, ipclaim, nil); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, ipclaimOrig, "cannot convert ipclaim after claim", err), errUpdateStatus)
	}

	return ctrl.Result{}, errors.Wrap(r.handleSuccess(ctx, ipclaimOrig), errUpdateStatus)
}

func (r *reconciler) handleSuccess(ctx context.Context, ipClaim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx)
	log.Debug("handleSuccess", "key", ipClaim.GetNamespacedName(), "status old", ipClaim.DeepCopy().Status)
	// take a snapshot of the current object
	patch := client.MergeFrom(ipClaim.DeepCopy())
	// update status
	ipClaim.SetConditions(condv1alpha1.Ready())
	r.recorder.Eventf(ipClaim, corev1.EventTypeNormal, ipambev1alpha1.IPClaimKind, "ready")

	
	log.Debug("handleSuccess", "key", ipClaim.GetNamespacedName(), "status new", ipClaim.Status)

	return r.Client.Status().Patch(ctx, ipClaim, patch, &client.SubResourcePatchOptions{
		PatchOptions: client.PatchOptions{
			FieldManager: "backend",
		},
	})
}

func (r *reconciler) handleError(ctx context.Context, ipClaim *ipambev1alpha1.IPClaim, msg string, err error) error {
	log := log.FromContext(ctx)
	// take a snapshot of the current object
	patch := client.MergeFrom(ipClaim.DeepCopy())

	if err != nil {
		msg = fmt.Sprintf("%s err %s", msg, err.Error())
	}
	ipClaim.SetConditions(condv1alpha1.Failed(msg))
	log.Error(msg)
	r.recorder.Eventf(ipClaim, corev1.EventTypeWarning, ipambev1alpha1.IPClaimKind, msg)

	return r.Client.Status().Patch(ctx, ipClaim, patch, &client.SubResourcePatchOptions{
		PatchOptions: client.PatchOptions{
			FieldManager: "backend",
		},
	})
}
