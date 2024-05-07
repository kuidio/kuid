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

	builderrest "github.com/henderiw/apiserver-builder/pkg/builder/rest"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/henderiw/apiserver-store/pkg/storebackend"
	vlanbe1v1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	"github.com/kuidio/kuid/pkg/backend/backend"
	"github.com/kuidio/kuid/pkg/kuidserver/store"
	"go.opentelemetry.io/otel"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewProvider(ctx context.Context, client client.Client, storeConfig *store.Config, be backend.Backend) builderrest.ResourceHandlerProvider {
	return func(ctx context.Context, scheme *runtime.Scheme, getter generic.RESTOptionsGetter) (rest.Storage, error) {
		return NewREST(ctx, scheme, getter, client, storeConfig, be)
	}
}

// NewPackageRevisionREST returns a RESTStorage object that will work against API services.
func NewREST(ctx context.Context, scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter, client client.Client, storeConfig *store.Config, be backend.Backend) (rest.Storage, error) {
	scheme.AddFieldLabelConversionFunc(
		schema.GroupVersionKind{
			Group:   vlanbe1v1alpha1.Group,
			Version: vlanbe1v1alpha1.Version,
			Kind:    vlanbe1v1alpha1.VLANClaimKind,
		},
		vlanbe1v1alpha1.ConvertVLANClaimFieldSelector,
	)

	var configStore storebackend.Storer[runtime.Object]
	var err error
	switch storeConfig.Type {
	case store.StorageType_File:
		configStore, err = store.CreateFileStore(ctx, scheme, &vlanbe1v1alpha1.VLANClaim{}, storeConfig.Prefix)
		if err != nil {
			return nil, err
		}
	case store.StorageType_KV:
		configStore, err = store.CreateKVStore(ctx, storeConfig.DB, scheme, &vlanbe1v1alpha1.VLANClaim{})
		if err != nil {
			return nil, err
		}
	default:
		configStore = store.CreateMemStore(ctx)
	}

	gr := vlanbe1v1alpha1.Resource(vlanbe1v1alpha1.VLANClaimPlural)
	strategy := NewStrategy(ctx, scheme, client, configStore, be)

	// This is the etcd store
	store := &registry.Store{
		Tracer:                    otel.Tracer("vlanclaim-server"),
		NewFunc:                   func() runtime.Object { return &vlanbe1v1alpha1.VLANClaim{} },
		NewListFunc:               func() runtime.Object { return &vlanbe1v1alpha1.VLANClaimList{} },
		PredicateFunc:             Match,
		DefaultQualifiedResource:  gr,
		SingularQualifiedResource: vlanbe1v1alpha1.Resource(vlanbe1v1alpha1.VLANClaimSingular),
		GetStrategy:               strategy,
		ListStrategy:              strategy,
		CreateStrategy:            strategy,
		UpdateStrategy:            strategy,
		DeleteStrategy:            strategy,
		WatchStrategy:             strategy,

		TableConvertor: NewTableConvertor(vlanbe1v1alpha1.Resource(vlanbe1v1alpha1.VLANClaimPlural)),
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
