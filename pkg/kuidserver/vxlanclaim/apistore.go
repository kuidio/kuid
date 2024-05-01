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

package vxlanclaim

import (
	"context"

	builderrest "github.com/henderiw/apiserver-builder/pkg/builder/rest"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/henderiw/apiserver-store/pkg/storebackend"
	vxlanbe1v1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
	"github.com/kuidio/kuid/pkg/backend"
	"github.com/kuidio/kuid/pkg/backend/vxlan"
	"github.com/kuidio/kuid/pkg/kuidserver/store"
	"go.opentelemetry.io/otel"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewProvider(ctx context.Context, client client.Client, storeConfig *store.Config, be backend.Backend[*vxlan.CacheContext]) builderrest.ResourceHandlerProvider {
	return func(ctx context.Context, scheme *runtime.Scheme, getter generic.RESTOptionsGetter) (rest.Storage, error) {
		return NewREST(ctx, scheme, getter, client, storeConfig, be)
	}
}

// NewPackageRevisionREST returns a RESTStorage object that will work against API services.
func NewREST(ctx context.Context, scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter, client client.Client, storeConfig *store.Config, be backend.Backend[*vxlan.CacheContext]) (rest.Storage, error) {
	scheme.AddFieldLabelConversionFunc(
		schema.GroupVersionKind{
			Group:   vxlanbe1v1alpha1.Group,
			Version: vxlanbe1v1alpha1.Version,
			Kind:    vxlanbe1v1alpha1.VXLANClaimKind,
		},
		vxlanbe1v1alpha1.ConvertVXLANClaimFieldSelector,
	)

	var configStore storebackend.Storer[runtime.Object]
	var err error
	switch storeConfig.Type {
	case store.StorageType_File:
		configStore, err = store.CreateFileStore(ctx, scheme, &vxlanbe1v1alpha1.VXLANClaim{}, storeConfig.Prefix)
		if err != nil {
			return nil, err
		}
	case store.StorageType_KV:
		configStore, err = store.CreateKVStore(ctx, storeConfig.DB, scheme, &vxlanbe1v1alpha1.VXLANClaim{})
		if err != nil {
			return nil, err
		}
	default:
		configStore = store.CreateMemStore(ctx)
	}

	gr := vxlanbe1v1alpha1.Resource(vxlanbe1v1alpha1.VXLANClaimPlural)
	strategy := NewStrategy(ctx, scheme, client, configStore, be)

	// This is the etcd store
	store := &registry.Store{
		Tracer:                    otel.Tracer("vxlanclaim-server"),
		NewFunc:                   func() runtime.Object { return &vxlanbe1v1alpha1.VXLANClaim{} },
		NewListFunc:               func() runtime.Object { return &vxlanbe1v1alpha1.VXLANClaimList{} },
		PredicateFunc:             Match,
		DefaultQualifiedResource:  gr,
		SingularQualifiedResource: vxlanbe1v1alpha1.Resource(vxlanbe1v1alpha1.VXLANClaimSingular),
		GetStrategy:               strategy,
		ListStrategy:              strategy,
		CreateStrategy:            strategy,
		UpdateStrategy:            strategy,
		DeleteStrategy:            strategy,
		WatchStrategy:             strategy,

		TableConvertor: NewTableConvertor(vxlanbe1v1alpha1.Resource(vxlanbe1v1alpha1.VXLANClaimPlural)),
	}
	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    GetAttrs,
	}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return store, nil
}
