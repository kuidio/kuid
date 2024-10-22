// Copyright 2022 The kpt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package as

import (
	"context"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/henderiw/apiserver-builder/pkg/builder/rest"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	genericbackend "github.com/kuidio/kuid/pkg/backend/generic"
	genericregistry "github.com/kuidio/kuid/pkg/registry/generic"
	"github.com/kuidio/kuid/pkg/registry/options"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
)

func NewStorageProviders(ctx context.Context, sync bool, options *options.Options) bebackend.StorageProviders {
	r := &storageProviders{
		be: genericbackend.New(
			ASClaimKind,
			ASIndexFromRuntime,
			ASClaimFromRuntime,
			ASEntryFromRuntime,
			GetASEntry,
		),
	}

	if sync {
		opts := *options
		opts.BackendInvoker = bebackend.NewIndexInvoker(r.be)
		r.indexStorageProvider = genericregistry.NewStorageProvider(ctx, &ASIndex{}, &opts)
	}
	if sync {
		opts := *options
		opts.BackendInvoker = bebackend.NewClaimInvoker(r.be)
		r.claimStorageProvider = genericregistry.NewStorageProvider(ctx, &ASClaim{}, &opts)
	}
	r.entryStorageProvider = genericregistry.NewStorageProvider(ctx, &ASEntry{}, options)

	return r
}

type storageProviders struct {
	be                   bebackend.Backend
	indexStorageProvider *rest.StorageProvider
	claimStorageProvider *rest.StorageProvider
	entryStorageProvider *rest.StorageProvider
}

func (r *storageProviders) GetIndexStorageProvider() *rest.StorageProvider {
	return r.indexStorageProvider

}
func (r *storageProviders) GetClaimStorageProvider() *rest.StorageProvider {
	return r.claimStorageProvider

}
func (r *storageProviders) GetEntryStorageProvider() *rest.StorageProvider {
	return r.entryStorageProvider
}

func (r *storageProviders) ApplyStorageToBackend(ctx context.Context, apiServer *builder.Server) error {
	claimStorageProvider := apiServer.StorageProvider[schema.GroupResource{
		Group:    SchemeGroupVersion.Group,
		Resource: ASClaimPlural,
	}]

	claimStorage, err := claimStorageProvider.Get(ctx, apiServer.Schemes[0], &ClaimGetter{})
	if err != nil {
		return err
	}

	entryStorageProvider := apiServer.StorageProvider[schema.GroupResource{
		Group:    SchemeGroupVersion.Group,
		Resource: ASEntryPlural,
	}]

	entryStorage, err := entryStorageProvider.Get(ctx, apiServer.Schemes[0], &EntryGetter{})
	if err != nil {
		return err
	}

	r.be.AddStorage(entryStorage, claimStorage)
	return nil
}

func (r *storageProviders) GetBackend() bebackend.Backend {
	return r.be
}

var _ generic.RESTOptionsGetter = &ClaimGetter{}

type ClaimGetter struct{}

func (r *ClaimGetter) GetRESTOptions(resource schema.GroupResource, example runtime.Object) (generic.RESTOptions, error) {
	return generic.RESTOptions{}, nil
}

var _ generic.RESTOptionsGetter = &EntryGetter{}

type EntryGetter struct{}

func (r *EntryGetter) GetRESTOptions(resource schema.GroupResource, example runtime.Object) (generic.RESTOptions, error) {
	return generic.RESTOptions{}, nil
}
