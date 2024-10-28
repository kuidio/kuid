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
	"fmt"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-builder/pkg/builder/rest"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/kuidio/kuid/apis/backend/ipam"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	ipambe "github.com/kuidio/kuid/pkg/backend/ipam"
	"github.com/kuidio/kuid/pkg/config"
	genericregistry "github.com/kuidio/kuid/pkg/registry/generic"
	"github.com/kuidio/kuid/pkg/registry/options"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
)

func init() {
	config.Register(
		ipam.SchemeGroupVersion.Group,
		ipambev1alpha1.AddToScheme,
		NewBackend,
		ApplyStorageToBackend,
		[]*config.ResourceConfig{
			{StorageProviderFn: NewIndexStorageProvider, Internal: &ipam.IPIndex{}, ResourceVersions: []resource.Object{&ipam.IPIndex{}, &ipambev1alpha1.IPIndex{}}},
			{StorageProviderFn: NewClaimStorageProvider, Internal: &ipam.IPClaim{}, ResourceVersions: []resource.Object{&ipam.IPClaim{}, &ipambev1alpha1.IPClaim{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &ipam.IPEntry{}, ResourceVersions: []resource.Object{&ipam.IPEntry{}, &ipambev1alpha1.IPEntry{}}},
		},
	)
}

func NewBackend() bebackend.Backend {
	return ipambe.New()
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
		Group:    ipam.SchemeGroupVersion.Group,
		Resource: ipam.IPClaimPlural,
	}]

	claimStorage, err := claimStorageProvider.Get(ctx, apiServer.Schemes[0], &Getter{})
	if err != nil {
		return err
	}
	claimStore, ok := claimStorage.(*registry.Store)
	if !ok {
		return fmt.Errorf("claimStorage is not a registry store")
	}

	entryStorageProvider := apiServer.StorageProvider[schema.GroupResource{
		Group:    ipam.SchemeGroupVersion.Group,
		Resource: ipam.IPEntryPlural,
	}]

	entryStorage, err := entryStorageProvider.Get(ctx, apiServer.Schemes[0], &Getter{})
	if err != nil {
		return err
	}
	entryStore, ok := entryStorage.(*registry.Store)
	if !ok {
		return fmt.Errorf("entryStorage is not a registry store")
	}

	return be.AddStorageInterfaces(ipambe.NewKuidBackendstorage(entryStore, claimStore))
}

var _ generic.RESTOptionsGetter = &Getter{}

type Getter struct{}

func (r *Getter) GetRESTOptions(resource schema.GroupResource, example runtime.Object) (generic.RESTOptions, error) {
	return generic.RESTOptions{}, nil
}
