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

package register

import (
	"context"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-builder/pkg/builder/rest"
	"github.com/kuidio/kuid/apis/backend/as"
	asbev1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	genericbackend "github.com/kuidio/kuid/pkg/backend/generic"
	"github.com/kuidio/kuid/pkg/config"
	genericregistry "github.com/kuidio/kuid/pkg/registry/generic"
	"github.com/kuidio/kuid/pkg/registry/options"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
)

func init() {
	config.Register(
		as.SchemeGroupVersion.Group,
		asbev1alpha1.AddToScheme,
		NewBackend,
		ApplyStorageToBackend,
		[]*config.ResourceConfig{
			{StorageProviderFn: NewIndexStorageProvider, Internal: &as.ASIndex{}, ResourceVersions: []resource.Object{&as.ASIndex{}, &asbev1alpha1.ASIndex{}}},
			{StorageProviderFn: NewClaimStorageProvider, Internal: &as.ASClaim{}, ResourceVersions: []resource.Object{&as.ASClaim{}, &asbev1alpha1.ASClaim{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &as.ASEntry{}, ResourceVersions: []resource.Object{&as.ASEntry{}, &asbev1alpha1.ASEntry{}}},
		},
	)
}

func NewBackend() bebackend.Backend {
	return genericbackend.New(
		as.ASClaimKind,
		as.ASIndexFromRuntime,
		as.ASClaimFromRuntime,
		as.ASEntryFromRuntime,
		as.GetASEntry,
	)
}

func NewIndexStorageProvider(ctx context.Context, obj resource.InternalObject, be bebackend.Backend, sync bool, options *options.Options) *rest.StorageProvider {
	opts := *options
	if sync {
		opts.BackendInvoker = bebackend.NewIndexInvoker(be)
		return genericregistry.NewStorageProvider(ctx, obj, &opts)
	}
	return genericregistry.NewStorageProvider(ctx, obj, &opts)
}

func NewClaimStorageProvider(ctx context.Context, obj resource.InternalObject, be bebackend.Backend, sync bool, options *options.Options) *rest.StorageProvider {
	opts := *options
	if sync {
		opts.BackendInvoker = bebackend.NewClaimInvoker(be)
		return genericregistry.NewStorageProvider(ctx, obj, &opts)
	}
	return genericregistry.NewStorageProvider(ctx, obj, &opts)
}

func NewStorageProvider(ctx context.Context, obj resource.InternalObject, be bebackend.Backend, sync bool, options *options.Options) *rest.StorageProvider {
	return genericregistry.NewStorageProvider(ctx, obj, options)
}

func ApplyStorageToBackend(ctx context.Context, be bebackend.Backend, apiServer *builder.Server) error {
	claimStorageProvider := apiServer.StorageProvider[schema.GroupResource{
		Group:    as.SchemeGroupVersion.Group,
		Resource: as.ASClaimPlural,
	}]

	claimStorage, err := claimStorageProvider.Get(ctx, apiServer.Schemes[0], &Getter{})
	if err != nil {
		return err
	}

	entryStorageProvider := apiServer.StorageProvider[schema.GroupResource{
		Group:    as.SchemeGroupVersion.Group,
		Resource: as.ASEntryPlural,
	}]

	entryStorage, err := entryStorageProvider.Get(ctx, apiServer.Schemes[0], &Getter{})
	if err != nil {
		return err
	}

	return be.AddStorage(entryStorage, claimStorage)
}

var _ generic.RESTOptionsGetter = &Getter{}

type Getter struct{}

func (r *Getter) GetRESTOptions(resource schema.GroupResource, example runtime.Object) (generic.RESTOptions, error) {
	return generic.RESTOptions{}, nil
}
