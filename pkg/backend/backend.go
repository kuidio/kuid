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

	"github.com/henderiw/store"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Backend[T1 any] interface {
	// CreateIndex creates a backend index
	CreateIndex(ctx context.Context, obj runtime.Object) error
	// DeleteIndex deletes a backend index
	DeleteIndex(ctx context.Context, obj runtime.Object) error
	// ValidateClaimSyntax validates the claim
	ValidateClaimSyntax(ctx context.Context, obj runtime.Object) field.ErrorList
	// ValidateClaim validates the claim
	ValidateClaim(ctx context.Context, obj runtime.Object) error
	// Claim claims an entry in the backend index
	Claim(ctx context.Context, obj runtime.Object) error
	// DeleteClaim delete a claim in the backend index
	DeleteClaim(ctx context.Context, obj runtime.Object) error
	// GetCache returns the cache
	GetCache(ctx context.Context, k store.Key) (T1, error)
}
