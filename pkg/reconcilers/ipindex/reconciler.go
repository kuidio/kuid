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

package ipindex

import (
	"context"
	"fmt"
	"reflect"

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
	reconcilers.Register(ipambev1alpha1.Group, ipambev1alpha1.IPIndexKind, &reconciler{})
}

const (
	reconcilerName = "IPIndexController"
	finalizer      = "ipindex.ipam.be.kuid.dev/finalizer"
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
		For(&ipambev1alpha1.IPIndex{}).
		//Watches(&ipambev1alpha1.IPEntry{},
		//	&eventhandler.IPEntryEventHandler{
		//		Client:  mgr.GetClient(),
		//		ObjList: &ipambev1alpha1.IPIndexList{},
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

	ipIndex := &ipambev1alpha1.IPIndex{}
	if err := r.Get(ctx, req.NamespacedName, ipIndex); err != nil {
		// if the resource no longer exists the reconcile loop is done
		if resource.IgnoreNotFound(err) != nil {
			log.Error(errGetCr, "error", err)
			return ctrl.Result{}, errors.Wrap(resource.IgnoreNotFound(err), errGetCr)
		}
		return ctrl.Result{}, nil
	}
	ipIndexOrig := ipIndex.DeepCopy()
	log.Debug("reconcile", "status orig", ipIndexOrig.Status)

	if !ipIndex.GetDeletionTimestamp().IsZero() {
		// if prefixes are provided from the network instance we treat them as
		// ipindex prefixes.
		/*
			for _, prefix := range index.Spec.Prefixes {
				if err := r.deleteIPClaim(ctx, cr, prefix); err != nil {
					if resource.IgnoreNotFound(err) != nil {
						return ctrl.Result{Requeue: true}, errors.Wrap(r.Update(ctx, cr), errUpdateStatus)
					}
				}
			}
		*/
		intIPIndex := &ipam.IPIndex{}
		if err := ipambev1alpha1.Convert_v1alpha1_IPIndex_To_ipam_IPIndex(ipIndex, intIPIndex, nil); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, ipIndexOrig, "cannot convert ipIndex before delete", err), errUpdateStatus)
		}
		if err := r.be.DeleteIndex(ctx, intIPIndex); err != nil {
			if resource.IgnoreNotFound(err) != nil {
				return ctrl.Result{Requeue: true},
					errors.Wrap(r.handleError(ctx, ipIndexOrig, "cannot delete index", err), errUpdateStatus)
			}
		}
		if err := ipambev1alpha1.Convert_ipam_IPIndex_To_v1alpha1_IPIndex(intIPIndex, ipIndex, nil); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, ipIndexOrig, "cannot convert ipIndex after delete", err), errUpdateStatus)
		}

		// We use owner reference so the k8s garbage collector takes care of the cleanup
		if err := r.finalizer.RemoveFinalizer(ctx, ipIndex); err != nil {
			return ctrl.Result{Requeue: true},
				errors.Wrap(r.handleError(ctx, ipIndexOrig, "cannot remove finalizer", err), errUpdateStatus)
		}
		return ctrl.Result{}, nil
	}

	if err := r.finalizer.AddFinalizer(ctx, ipIndex); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, ipIndexOrig, "cannot add finalizer", err), errUpdateStatus)
	}

	// create ip index
	intIPIndex := &ipam.IPIndex{}
	if err := ipambev1alpha1.Convert_v1alpha1_IPIndex_To_ipam_IPIndex(ipIndex, intIPIndex, nil); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, ipIndexOrig, "cannot convert ipIndex before create", err), errUpdateStatus)
	}
	if err := r.be.CreateIndex(ctx, intIPIndex); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, ipIndexOrig, "cannot apply index", err), errUpdateStatus)
	}
	if err := ipambev1alpha1.Convert_ipam_IPIndex_To_v1alpha1_IPIndex(intIPIndex, ipIndex, nil); err != nil {
		return ctrl.Result{Requeue: true},
			errors.Wrap(r.handleError(ctx, ipIndexOrig, "cannot convert ipIndex after create", err), errUpdateStatus)
	}

	// first validate if the previously claimed prefixes are still needed
	// if not delete them
	/*
		for _, claimedPrefix := range cr.Status.Prefixes {
			found := false
			for _, prefix := range cr.Spec.Prefixes {
				if claimedPrefix.Prefix == prefix.Prefix {
					found = true
					break
				}
			}
			if !found {
				if err := r.deleteIPClaim(ctx, cr, claimedPrefix); err != nil {
					return ctrl.Result{RequeueAfter: 1 * time.Second}, errors.Wrap(r.Update(ctx, cr), errUpdateStatus)
				}
			}
		}
		// if prefixes are provided from the network instance we treat them as
		// aggregate prefixes.
		for _, prefix := range cr.Spec.Prefixes {
			prefix := prefix
			if err := r.applyIPClaim(ctx, cr, prefix); err != nil {
				return ctrl.Result{RequeueAfter: 1 * time.Second}, errors.Wrap(r.Update(ctx, cr), errUpdateStatus)
			}
		}
	*/

	return ctrl.Result{}, errors.Wrap(r.handleSuccess(ctx, ipIndexOrig), errUpdateStatus)
}

func (r *reconciler) handleSuccess(ctx context.Context, ipIndex *ipambev1alpha1.IPIndex) error {
	// take a snapshot of the current object
	patch := client.MergeFrom(ipIndex.DeepCopy())
	// update status
	ipIndex.SetConditions(condv1alpha1.Ready())
	r.recorder.Eventf(ipIndex, corev1.EventTypeNormal, ipambev1alpha1.IPIndexKind, "ready")

	return r.Client.Status().Patch(ctx, ipIndex, patch, &client.SubResourcePatchOptions{
		PatchOptions: client.PatchOptions{
			FieldManager: "backend",
		},
	})
}

func (r *reconciler) handleError(ctx context.Context, ipIndex *ipambev1alpha1.IPIndex, msg string, err error) error {
	log := log.FromContext(ctx)
	// take a snapshot of the current object
	patch := client.MergeFrom(ipIndex.DeepCopy())

	if err != nil {
		msg = fmt.Sprintf("%s err %s", msg, err.Error())
	}
	ipIndex.SetConditions(condv1alpha1.Failed(msg))
	log.Error(msg)
	r.recorder.Eventf(ipIndex, corev1.EventTypeWarning, ipambev1alpha1.IPIndexKind, msg)

	return r.Client.Status().Patch(ctx, ipIndex, patch, &client.SubResourcePatchOptions{
		PatchOptions: client.PatchOptions{
			FieldManager: "backend",
		},
	})
}

/*
func (r *reconciler) deleteIndex(ctx context.Context, cr *ipambev1alpha1.IPIndex) error {
	if err := r.be.DeleteIndex(ctx, cr); err != nil {
		r.handleError(ctx, cr, "cannot delete index", err)
		return err
	}
	return nil
}

func (r *reconciler) applyIndex(ctx context.Context, cr *ipambev1alpha1.IPIndex) error {
	if err := r.be.CreateIndex(ctx, cr); err != nil {
		r.handleError(ctx, cr, "cannot create index", err)
		return err
	}
	return nil
}
*/

/*
func (r *reconciler) applyIPClaim(ctx context.Context, cr *ipambev1alpha1.IPIndex, prefix ipambev1alpha1.Prefix) error {
	ipclaim, err := cr.GetClaim(prefix)
	if err != nil { // strange if this happens since the prefix was already processed
		r.handleError(ctx, cr, "build ipclaim", err)
		return err
	}
	if err := r.be.Claim(ctx, ipclaim); err != nil {
		r.handleError(ctx, cr, "cannot claim ip", err)
		return err
	}
	if ipclaim.Status.Prefix == nil || *ipclaim.Status.Prefix != prefix.Prefix {
		//we got a different prefix than requested one
		msg := fmt.Sprintf("ip prefix not ready: req/rsp %s/%v", prefix.Prefix, ipclaim.Status.Prefix)
		r.handleError(ctx, cr, msg, nil)
		return err
	}
	return nil
}

func (r *reconciler) deleteIPClaim(ctx context.Context, cr *ipambev1alpha1.IPIndex, prefix ipambev1alpha1.Prefix) error {
	//ipclaim, err := buildIPClaim(ctx, cr, prefix)
	ipclaim, err := cr.GetClaim(prefix)
	if err != nil { // strange if this happens since the prefix was already processed
		r.handleError(ctx, cr, "cannot build ipclaim", err)
		return err
	}
	if err := r.be.Release(ctx, ipclaim); err != nil {
		if !strings.Contains(err.Error(), "not initialized") {
			r.handleError(ctx, cr, "cannot delete ipclaim", err)
		}
	}
	return nil
}
*/
