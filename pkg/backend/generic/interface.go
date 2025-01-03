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

package generic

import (
	"context"
	"errors"
	"reflect"

	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

type BackendStorage interface {
	ListEntries(ctx context.Context, k store.Key) ([]backend.EntryObject, error)
	CreateEntry(ctx context.Context, obj backend.EntryObject) error
	UpdateEntry(ctx context.Context, obj, old backend.EntryObject) error
	DeleteEntry(ctx context.Context, obj backend.EntryObject) error
	ListClaims(ctx context.Context, k store.Key, opts ...ListOption) (map[string]backend.ClaimObject, error)
	CreateClaim(ctx context.Context, obj backend.ClaimObject) error
	UpdateClaim(ctx context.Context, obj, old backend.ClaimObject) error
	DeleteClaim(ctx context.Context, obj backend.ClaimObject) error
}

type ListOption interface {
	// ApplyToGet applies this configuration to the given get options.
	ApplyToList(*ListOptions)
}

var _ ListOption = &ListOptions{}

type ListOptions struct {
	OwnerKind string
}

func (o *ListOptions) ApplyToList(lo *ListOptions) {
	lo.OwnerKind = o.OwnerKind
}

// ApplyOptions applies the given get options on these options,
// and then returns itself (for convenient chaining).
func (o *ListOptions) ApplyOptions(opts []ListOption) *ListOptions {
	for _, opt := range opts {
		opt.ApplyToList(o)
	}
	return o
}

func NewKuidBackendstorage(entryStorage, claimStorage *registry.Store) BackendStorage {
	return &kuidbe{
		entryStorage: entryStorage,
		claimStorage: claimStorage,
	}
}

type kuidbe struct {
	entryStorage *registry.Store
	claimStorage *registry.Store
}

func (r *kuidbe) ListEntries(ctx context.Context, k store.Key) ([]backend.EntryObject, error) {
	log := log.FromContext(ctx).With("key", k.String())
	list, err := r.entryStorage.List(ctx, &internalversion.ListOptions{})
	if err != nil {
		return nil, err
	}

	items, err := meta.ExtractList(list)
	if err != nil {
		return nil, err
	}

	entryList := make([]backend.EntryObject, 0)
	var errm error
	for _, obj := range items {
		entryObj, ok := obj.(backend.EntryObject)
		if !ok {
			log.Error("obj is not an EntryObject", "obj", reflect.TypeOf(obj).Name())
			errm = errors.Join(errm, err)
			continue
		}
		if entryObj.GetIndex() == k.Name {
			entryList = append(entryList, entryObj)
		}
	}
	return entryList, errm
}

func (r *kuidbe) CreateEntry(ctx context.Context, obj backend.EntryObject) error {
	ctx = genericapirequest.WithNamespace(ctx, obj.GetNamespace())
	if _, err := r.entryStorage.Create(ctx, obj, nil, &metav1.CreateOptions{
		FieldManager: "backend",
	}); err != nil {
		return err
	}
	return nil
}
func (r *kuidbe) UpdateEntry(ctx context.Context, obj, old backend.EntryObject) error {
	ctx = genericapirequest.WithNamespace(ctx, obj.GetNamespace())
	defaultObjInfo := rest.DefaultUpdatedObjectInfo(old, EntryTransformer)
	if _, _, err := r.entryStorage.Update(ctx, old.GetName(), defaultObjInfo, nil, nil, false, &metav1.UpdateOptions{
		FieldManager: "backend",
	}); err != nil {
		return err
	}
	return nil
}

func (r *kuidbe) DeleteEntry(ctx context.Context, obj backend.EntryObject) error {
	ctx = genericapirequest.WithNamespace(ctx, obj.GetNamespace())
	if _, _, err := r.entryStorage.Delete(ctx, obj.GetName(), nil, &metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func (r *kuidbe) ListClaims(ctx context.Context, k store.Key, opts ...ListOption) (map[string]backend.ClaimObject, error) {
	o := &ListOptions{}
	o.ApplyOptions(opts)

	log := log.FromContext(ctx).With("key", k.String())
	log.Debug("list claims from storage")
	list, err := r.claimStorage.List(ctx, &internalversion.ListOptions{})
	if err != nil {
		return nil, err
	}

	items, err := meta.ExtractList(list)
	if err != nil {
		return nil, err
	}

	claimMap := make(map[string]backend.ClaimObject)
	var errm error
	for _, obj := range items {
		claimObj, ok := obj.(backend.ClaimObject)
		if !ok {
			log.Error("obj is not an ClaimObject", "obj", reflect.TypeOf(obj).Name())
			errm = errors.Join(errm, err)
			continue
		}
		if claimObj.GetIndex() != k.Name {
			continue
		}
		if o.OwnerKind != "" {
			for _, ownerref := range claimObj.GetOwnerReferences() {
				if ownerref.Kind == o.OwnerKind {
					claimMap[claimObj.GetNamespacedName().String()] = claimObj
				}
			}
		} else {
			claimMap[claimObj.GetNamespacedName().String()] = claimObj
		}
	}
	return claimMap, errm
}

func (r *kuidbe) CreateClaim(ctx context.Context, obj backend.ClaimObject) error {
	log := log.FromContext(ctx)
	ctx = genericapirequest.WithNamespace(ctx, obj.GetNamespace())
	if _, err := r.claimStorage.Create(ctx, obj, nil, &metav1.CreateOptions{
		FieldManager: "backend",
		DryRun:       []string{"recursion"},
	}); err != nil {
		log.Error("create claim failed", "name", obj.GetName(), "error", err.Error())
		return err
	}
	return nil
}

func (r *kuidbe) UpdateClaim(ctx context.Context, obj, old backend.ClaimObject) error {
	log := log.FromContext(ctx)
	defaultObjInfo := rest.DefaultUpdatedObjectInfo(old, ClaimTransformer)
	if _, _, err := r.claimStorage.Update(ctx, old.GetName(), defaultObjInfo, nil, nil, false, &metav1.UpdateOptions{
		FieldManager: "backend",
		DryRun:       []string{"recursion"},
	}); err != nil {
		log.Error("update claim failed", "name", obj.GetName(), "error", err.Error())
		return err
	}
	return nil
}

func (r *kuidbe) DeleteClaim(ctx context.Context, obj backend.ClaimObject) error {
	log := log.FromContext(ctx)
	ctx = genericapirequest.WithNamespace(ctx, obj.GetNamespace())
	if _, _, err := r.claimStorage.Delete(ctx, obj.GetName(), nil, &metav1.DeleteOptions{
		DryRun: []string{"recursion"},
	}); err != nil {
		log.Error("cannot delete claim", "error", err)
		return err
	}
	return nil
}
