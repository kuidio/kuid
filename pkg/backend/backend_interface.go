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

package backend

import (
	"context"

	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

type Backend interface {
	AddStorage(entryStorage, claimStorage rest.Storage) error
	// CreateIndex creates a backend index
	CreateIndex(ctx context.Context, obj runtime.Object) error
	// DeleteIndex deletes a backend index
	DeleteIndex(ctx context.Context, obj runtime.Object) error
	// Claim claims an entry in the backend index
	Claim(ctx context.Context, obj runtime.Object) error
	// Release a claim in the backend
	Release(ctx context.Context, obj runtime.Object) error
	// PrintEntries prints the entries of the cache
	PrintEntries(ctx context.Context, index string)

	GetClaimStorage() *registry.Store
}
