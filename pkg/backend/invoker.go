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

	"github.com/kuidio/kuid/pkg/registry/options"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewClaimInvoker(be Backend) options.BackendInvoker {
	return &claimInvoker{
		be: be,
	}
}

type claimInvoker struct {
	be Backend
}

func (r *claimInvoker) InvokeCreate(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error) {
	if err := r.be.Claim(ctx, obj, recursion); err != nil {
		return obj, err
	}
	return obj, nil
}

func (r *claimInvoker) InvokeUpdate(ctx context.Context, obj, old runtime.Object, recursion bool) (runtime.Object, runtime.Object, error) {
	if err := r.be.Claim(ctx, obj, recursion); err != nil {
		return obj, old, err
	}
	return obj, old, nil
}

func (r *claimInvoker) InvokeDelete(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error) {
	if err := r.be.Release(ctx, obj, recursion); err != nil {
		return obj, err
	}
	return obj, nil
}

func NewIndexInvoker(be Backend) options.BackendInvoker {
	return &indexPreparator{
		be: be,
	}
}

type indexPreparator struct {
	be Backend
}

func (r *indexPreparator) InvokeCreate(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error) {
	if err := r.be.CreateIndex(ctx, obj); err != nil {
		return obj, err
	}
	return obj, nil
}

func (r *indexPreparator) InvokeUpdate(ctx context.Context, obj, old runtime.Object, recursion bool) (runtime.Object, runtime.Object, error) {
	if err := r.be.CreateIndex(ctx, obj); err != nil {
		return obj, old, err
	}
	return obj, old, nil
}

func (r *indexPreparator) InvokeDelete(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error) {
	if err := r.be.DeleteIndex(ctx, obj); err != nil {
		return obj, err
	}
	return obj, nil
}
