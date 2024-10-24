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

/*

import (
	"context"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/henderiw/apiserver-builder/pkg/builder/rest"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	genericbackend "github.com/kuidio/kuid/pkg/backend/generic"
	genericregistry "github.com/kuidio/kuid/pkg/registry/generic"
	"github.com/kuidio/kuid/pkg/registry/options"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// backend -> can be initialized at init
//

func NewBackend() bebackend.Backend {
	return genericbackend.New(
		ASClaimKind,
		ASIndexFromRuntime,
		ASClaimFromRuntime,
		ASEntryFromRuntime,
		GetASEntry,
	)
}

func NewIndexStorageProvider(ctx context.Context, be bebackend.Backend, sync bool, options *options.Options) *rest.StorageProvider {
	opts := *options
	if sync {
		opts.BackendInvoker = bebackend.NewIndexInvoker(be)
		return genericregistry.NewStorageProvider(ctx, &ASIndex{}, &opts)
	}
	return genericregistry.NewStorageProvider(ctx, &ASIndex{}, &opts)
}

func NewClaimStorageProvider(ctx context.Context, be bebackend.Backend, sync bool, options *options.Options) *rest.StorageProvider {
	opts := *options
	if sync {
		opts.BackendInvoker = bebackend.NewClaimInvoker(be)
		return genericregistry.NewStorageProvider(ctx, &ASClaim{}, &opts)
	}
	return genericregistry.NewStorageProvider(ctx, &ASClaim{}, &opts)
}

func NewEntryStorageProvider(ctx context.Context, be bebackend.Backend, sync bool, options *options.Options) *rest.StorageProvider {
	return genericregistry.NewStorageProvider(ctx, &ASEntry{}, options)
}

func ApplyStorageToBackend(ctx context.Context, be bebackend.Backend, apiServer *builder.Server) error {
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

	return be.AddStorage(entryStorage, claimStorage)
}
*/

/*
func NewStorageProviders(ctx context.Context, sync bool, options *options.Options) bebackend.StorageProviders {
	r := &StorageProviders{
		be: CreateBackend(),
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

type StorageProviders struct {
	be                   bebackend.Backend
	indexStorageProvider *rest.StorageProvider
	claimStorageProvider *rest.StorageProvider
	entryStorageProvider *rest.StorageProvider
}

func (r *StorageProviders) GetIndexStorageProvider() *rest.StorageProvider {
	return r.indexStorageProvider

}
func (r *StorageProviders) GetClaimStorageProvider() *rest.StorageProvider {
	return r.claimStorageProvider

}
func (r *StorageProviders) GetEntryStorageProvider() *rest.StorageProvider {
	return r.entryStorageProvider
}


func (r *StorageProviders) ApplyStorageToBackend(ctx context.Context, apiServer *builder.Server) error {
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


func (r *StorageProviders) GetBackend() bebackend.Backend {
	return r.be
}

*/
