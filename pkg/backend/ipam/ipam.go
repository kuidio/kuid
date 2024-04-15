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

package ipam

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/hansthienpondt/nipam/pkg/table"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"

	"github.com/kuidio/kuid/pkg/backend"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func New(c client.Client) backend.Backend[*table.RIB] {
	cache := backend.NewCache[*table.RIB]()
	return &be{
		cache:  cache,
		client: c,
		store:  NewStore(c, cache),
	}
}

type be struct {
	cache  backend.Cache[*table.RIB]
	client client.Client
	store  backend.Store
}

func initIndexContext(ctx context.Context, op string, idx *ipambev1alpha1.IPIndex) context.Context {
	l := log.FromContext(ctx).
		With(
			"op", fmt.Sprintf("%s index", op),
			"nsn", idx.GetNamespacedName().String(),
		)
	return log.IntoContext(ctx, l)
}

func initClaimContext(ctx context.Context, op string, claim *ipambev1alpha1.IPClaim) context.Context {
	var l *slog.Logger
	if claim.Spec.Prefix != nil {
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s prefix claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"ni", claim.Spec.NetworkInstance,
				"prefix", claim.Spec.Prefix,
			)
	} else {
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s dynamic claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"ni", claim.Spec.NetworkInstance,
			)
	}
	return log.IntoContext(ctx, l)
}

/*
log := log.FromContext(ctx).With("name", claim.GetName(), "networkInstance", claim.Spec.NetworkInstance, "prefix", claim.Spec.Prefix)
	log.Info("ipclaim create")
*/

// CreateIndex creates a backend index
func (r *be) CreateIndex(ctx context.Context, obj runtime.Object) error {
	cr, ok := obj.(*ipambev1alpha1.IPIndex)
	if !ok {
		return fmt.Errorf("cannot create index expecting %s, got %s", ipambev1alpha1.IPIndexKind, reflect.TypeOf(obj).Name())
	}
	ctx = initIndexContext(ctx, "create", cr)
	log := log.FromContext(ctx)
	log.Info("start")
	key := cr.GetKey()
	//log := log.FromContext(ctx).With("key", key)

	log.Info("start", "isInitialized", r.cache.IsInitialized(ctx, key))
	// if the Cache is not initialized -> restore the cache
	// this happens upon initialization or backend restart
	r.cache.Create(ctx, key, table.NewRIB())
	if r.cache.IsInitialized(ctx, key) {
		log.Info("already initialized")
		return nil
	}
	if err := r.store.Restore(ctx, key); err != nil {
		log.Error("cannot restore index", "error", err.Error())
		return err
	}
	log.Info("finished")
	return r.cache.SetInitialized(ctx, key)
}

// DeleteIndex deletes a backend index
func (r *be) DeleteIndex(ctx context.Context, obj runtime.Object) error {
	cr, ok := obj.(*ipambev1alpha1.IPIndex)
	if !ok {
		return fmt.Errorf("cannot delete index expecting %s, got %s", ipambev1alpha1.IPIndexKind, reflect.TypeOf(obj).Name())
	}
	ctx = initIndexContext(ctx, "delete", cr)
	log := log.FromContext(ctx)
	log.Debug("start")
	key := cr.GetKey()

	log.Debug("start", "isInitialized", r.cache.IsInitialized(ctx, key))
	// delete the data from the backend
	if err := r.store.Destroy(ctx, key); err != nil {
		log.Error("cannot delete Index", "error", err.Error())
		return err
	}
	r.cache.Delete(ctx, key)

	log.Debug("finished")
	return nil

}

// Claim claims an entry in the backend index
func (r *be) Claim(ctx context.Context, obj runtime.Object) error {
	claim, ok := obj.(*ipambev1alpha1.IPClaim)
	if !ok {
		return fmt.Errorf("cannot claim ip expecting %s, got %s", ipambev1alpha1.IPIndexKind, reflect.TypeOf(obj).Name())
	}
	ctx = initClaimContext(ctx, "create", claim)
	log := log.FromContext(ctx)
	log.Debug("start")

	a, err := r.getApplicator(ctx, claim)
	if err != nil {
		return err
	}
	if err := a.Apply(ctx, claim); err != nil {
		return err
	}

	// store the resources in the backend
	return r.store.SaveAll(ctx, claim.GetKey())
}

// DeleteClaim delete a claim in the backend index
func (r *be) DeleteClaim(ctx context.Context, obj runtime.Object) error {
	claim, ok := obj.(*ipambev1alpha1.IPClaim)
	if !ok {
		return fmt.Errorf("cannot delete ip cliam expecting %s, got %s", ipambev1alpha1.IPIndexKind, reflect.TypeOf(obj).Name())
	}
	ctx = initClaimContext(ctx, "delete", claim)
	log := log.FromContext(ctx)
	log.Debug("start")

	// ip claim delete and store
	a, err := r.getApplicator(ctx, claim)
	if err != nil {
		// error gets returned when rib is not initialized -> this means we can safely return
		// and pretend nothing is wrong (hence return nil) since the cleanup already happened
		return nil
	}
	if err := a.Delete(ctx, claim); err != nil {
		return err
	}

	return r.store.SaveAll(ctx, claim.GetKey())
}

func (r *be) GetCache(ctx context.Context, key store.Key) (*table.RIB, error) {
	return r.cache.Get(ctx, key, false)
}

func (r *be) ValidateClaimSyntax(ctx context.Context, obj runtime.Object) field.ErrorList {
	var allErrs field.ErrorList
	claim, ok := obj.(*ipambev1alpha1.IPClaim)
	if !ok {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath(""),
			obj,
			fmt.Errorf("unexpected new object, expecting: %s, got: %s", ipambev1alpha1.IPClaimKind, reflect.TypeOf(obj)).Error(),
		))
		return allErrs
	}
	gv, err := schema.ParseGroupVersion(claim.APIVersion)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("apiVersion"),
			obj,
			fmt.Errorf("invalid apiVersion: err: %s", err.Error()).Error(),
		))
		return allErrs
	}

	if claim.Spec.Owner == nil {
		claim.Spec.Owner = &commonv1alpha1.OwnerReference{
			Group:     gv.Group,
			Version:   gv.Version,
			Kind:      claim.Kind,
			Namespace: claim.Namespace,
			Name:      claim.Name,
		}
	}

	var v SyntaxValidator
	if claim.Spec.Prefix == nil {
		v = &dynamicClaimSyntaxValidator{}
	} else {
		v = &prefixClaimSyntaxValidator{}

	}
	return v.ValidateSyntax(ctx, claim)
}

func (r *be) ValidateClaim(ctx context.Context, obj runtime.Object) error {
	claim, ok := obj.(*ipambev1alpha1.IPClaim)
	if !ok {
		return fmt.Errorf("unexpected new object, expecting: %s, got: %s", ipambev1alpha1.IPClaimKind, reflect.TypeOf(obj))
	}

	rib, err := r.GetCache(ctx, claim.GetKey())
	if err != nil {
		return fmt.Errorf("rib not ready, initializing: err: %s", err.Error())
	}

	if claim.Spec.Prefix != nil {
		v := &prefixClaimValidator{rib: rib}
		return v.Validate(ctx, claim)
	}
	return nil
}

func (r *be) getApplicator(ctx context.Context, claim *ipambev1alpha1.IPClaim) (Applicator, error) {
	rib, err := r.cache.Get(ctx, claim.GetKey(), false)
	if err != nil {
		return nil, err
	}

	// validate - happened before
	if claim.Spec.Prefix != nil {
		/*
			pi, err := iputil.New(*claim.Spec.Prefix)
			if err != nil {
				return nil, err
			}
		*/
		return &prefixApplicator{
			applicator: applicator{
				rib: rib,
				//pi:  pi,
			},
		}, nil
	}
	return &dynamicApplicator{
		applicator: applicator{
			rib: rib,
		},
	}, nil
}
