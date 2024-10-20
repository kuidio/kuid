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

package indexserver

/*

import (
	"context"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	builderrest "github.com/henderiw/apiserver-builder/pkg/builder/rest"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/henderiw/apiserver-store/pkg/storebackend"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	"github.com/kuidio/kuid/pkg/kuidserver/store"
	"go.opentelemetry.io/otel"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServerObjContext struct {
	TracerString   string
	Obj            resource.Object
	ConversionFunc runtime.FieldLabelConversionFunc
	TableConverter func(gr schema.GroupResource) registry.TableConvertor
}

func NewProvider(ctx context.Context, client client.Client, serverObjContext *ServerObjContext, storeConfig *store.Config, be bebackend.Backend) builderrest.ResourceHandlerProvider {
	return func(ctx context.Context, scheme *runtime.Scheme, getter generic.RESTOptionsGetter) (rest.Storage, error) {
		return NewREST(ctx, scheme, getter, client, serverObjContext, storeConfig, be)
	}
}

// NewPackageRevisionREST returns a RESTStorage object that will work against API services.
func NewREST(ctx context.Context, scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter, client client.Client, serverObjContext *ServerObjContext, storeConfig *store.Config, be bebackend.Backend) (rest.Storage, error) {
	scheme.AddFieldLabelConversionFunc(
		serverObjContext.Obj.GetObjectKind().GroupVersionKind(),
		serverObjContext.ConversionFunc,
	)

	var configStore storebackend.Storer[runtime.Object]
	var err error
	switch storeConfig.Type {
	case store.StorageType_File:
		configStore, err = store.CreateFileStore(ctx, scheme, serverObjContext.Obj, storeConfig.Prefix)
		if err != nil {
			return nil, err
		}
	case store.StorageType_KV:
		configStore, err = store.CreateKVStore(ctx, storeConfig.DB, scheme, serverObjContext.Obj)
		if err != nil {
			return nil, err
		}
	default:
		configStore = store.CreateMemStore(ctx)
	}

	singlularResource := serverObjContext.Obj.GetGroupVersionResource().GroupResource()
	singlularResource.Resource = serverObjContext.Obj.GetSingularName()
	strategy := NewStrategy(ctx, scheme, client, serverObjContext, configStore, be)

	// overwrite the default table convertor if specified by the user
	tableConvertor := DefaultTableConvertor(serverObjContext.Obj.GetGroupVersionResource().GroupResource())
	if serverObjContext.TableConverter != nil {
		tableConvertor = serverObjContext.TableConverter(serverObjContext.Obj.GetGroupVersionResource().GroupResource())
	}

	// This is the etcd store
	store := &registry.Store{
		Tracer:                    otel.Tracer(serverObjContext.TracerString),
		NewFunc:                   serverObjContext.Obj.New,
		NewListFunc:               serverObjContext.Obj.NewList,
		PredicateFunc:             Match,
		DefaultQualifiedResource:  serverObjContext.Obj.GetGroupVersionResource().GroupResource(),
		SingularQualifiedResource: singlularResource,
		GetStrategy:               strategy,
		ListStrategy:              strategy,
		CreateStrategy:            strategy,
		UpdateStrategy:            strategy,
		DeleteStrategy:            strategy,
		WatchStrategy:             strategy,

		TableConvertor: tableConvertor,
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
*/