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

package vlanclaim

import (
	"context"
	"fmt"
	"reflect"
	"strings"

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
	reconcilers.Register(vlan.GroupName, vlanbev1alpha1.VLANClaimKind, &reconciler{})
}

const (
	reconcilerName = "VLANClaimController"
	finalizer      = "vlanclaim.vlan.be.kuid.dev/finalizer"
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
		For(&vlanbev1alpha1.VLANClaim{}).
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

	claim := &vlanbev1alpha1.VLANClaim{}
	if err := r.Get(ctx, req.NamespacedName, claim); err != nil {
		// if the resource no longer exists the reconcile loop is done
		if resource.IgnoreNotFound(err) != nil {
			log.Error(errGetCr, "error", err)
			return ctrl.Result{}, errors.Wrap(resource.IgnoreNotFound(err), errGetCr)
		}
		return ctrl.Result{}, nil
	}
	claimOrig := claim.DeepCopy()
	log.Debug("reconcile", "status orig", claimOrig.Status)

	if !claim.GetDeletionTimestamp().IsZero() {
		intClaim := &vlan.VLANClaim{}
		if err := vlanbev1alpha1.Convert_v1alpha1_VLANClaim_To_vlan_VLANClaim(claim, intClaim, nil); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, claimOrig, "cannot convert claim before delete claim", err), errUpdateStatus)
		}

		if err := r.be.Release(ctx, intClaim, false); err != nil {
			if !strings.Contains(err.Error(), "not initialized") {
				return ctrl.Result{Requeue: true},
					errors.Wrap(r.handleError(ctx, claimOrig, "cannot delete claim", err), errUpdateStatus)
			}
		}
		if err := vlanbev1alpha1.Convert_vlan_VLANClaim_To_v1alpha1_VLANClaim(intClaim, claim, nil); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, claimOrig, "cannot convert claim after delete claim", err), errUpdateStatus)
		}

		if err := r.finalizer.RemoveFinalizer(ctx, claim); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, claimOrig, "cannot delete finalizer", err), errUpdateStatus)
		}
		return ctrl.Result{}, nil
	}

	if err := r.finalizer.AddFinalizer(ctx, claim); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, claimOrig, "cannot add finalizer", err), errUpdateStatus)
	}

	intClaim := &vlan.VLANClaim{}
	if err := vlanbev1alpha1.Convert_v1alpha1_VLANClaim_To_vlan_VLANClaim(claim, intClaim, nil); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, claimOrig, "cannot convert claim before claim", err), errUpdateStatus)
	}
	if err := r.be.Claim(ctx, intClaim, false); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, claimOrig, "cannot claim", err), errUpdateStatus)
	}
	if err := vlanbev1alpha1.Convert_vlan_VLANClaim_To_v1alpha1_VLANClaim(intClaim, claim, nil); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, claimOrig, "cannot convert claim after claim", err), errUpdateStatus)
	}

	return ctrl.Result{}, errors.Wrap(r.handleSuccess(ctx, claimOrig), errUpdateStatus)
}

func (r *reconciler) handleSuccess(ctx context.Context, claim *vlanbev1alpha1.VLANClaim) error {
	log := log.FromContext(ctx)
	log.Debug("handleSuccess", "key", claim.GetNamespacedName(), "status old", claim.DeepCopy().Status)
	// take a snapshot of the current object
	patch := client.MergeFrom(claim.DeepCopy())
	// update status
	claim.SetConditions(condv1alpha1.Ready())
	r.recorder.Eventf(claim, corev1.EventTypeNormal, vlanbev1alpha1.VLANClaimKind, "ready")

	
	log.Debug("handleSuccess", "key", claim.GetNamespacedName(), "status new", claim.Status)

	return r.Client.Status().Patch(ctx, claim, patch, &client.SubResourcePatchOptions{
		PatchOptions: client.PatchOptions{
			FieldManager: "backend",
		},
	})
}

func (r *reconciler) handleError(ctx context.Context, claim *vlanbev1alpha1.VLANClaim, msg string, err error) error {
	log := log.FromContext(ctx)
	// take a snapshot of the current object
	patch := client.MergeFrom(claim.DeepCopy())

	if err != nil {
		msg = fmt.Sprintf("%s err %s", msg, err.Error())
	}
	claim.SetConditions(condv1alpha1.Failed(msg))
	log.Error(msg)
	r.recorder.Eventf(claim, corev1.EventTypeWarning, vlanbev1alpha1.VLANClaimKind, msg)

	return r.Client.Status().Patch(ctx, claim, patch, &client.SubResourcePatchOptions{
		PatchOptions: client.PatchOptions{
			FieldManager: "backend",
		},
	})
}
