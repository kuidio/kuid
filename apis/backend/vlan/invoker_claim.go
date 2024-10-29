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

package vlan

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kuidio/kuid/pkg/backend"
	"github.com/kuidio/kuid/pkg/registry/options"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewChoreoClaimInvoker(be backend.Backend) options.BackendInvoker {
	return &claiminvoker{
		be: be,
	}
}

type claiminvoker struct {
	be backend.Backend
}

func claimConvertToInternal(obj runtime.Object) (*VLANClaim, error) {
	ru, ok := obj.(runtime.Unstructured)
	if !ok {
		return nil, fmt.Errorf("not an unstructured obj, got: %s", reflect.TypeOf(obj).Name())
	}
	claim := &VLANClaim{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ru.UnstructuredContent(), claim); err != nil {
		return nil, fmt.Errorf("unable to convert unstructured object to VLANClaim: %v", err)
	}
	return claim, nil
}

func claimConvertFromInternal(obj runtime.Object) (runtime.Unstructured, error) {
	claim, ok := obj.(*VLANClaim)
	if !ok {
		return nil, fmt.Errorf("not an unstructured obj, got: %s", reflect.TypeOf(obj).Name())
	}

	uobj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(claim)
	if err != nil {
		return nil, fmt.Errorf("unable to convert to unstructured: %v", err)
	}
	return &unstructured.Unstructured{Object: uobj}, nil
}

func (r *claiminvoker) InvokeCreate(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error) {
	claim, err := claimConvertToInternal(obj)
	if err != nil {
		return obj, err
	}
	if err := r.be.Claim(ctx, claim, recursion); err != nil {
		return obj, err
	}
	newClaim, err := claimConvertFromInternal(claim)
	if err != nil {
		return obj, err
	}
	return newClaim, nil
}

func (r *claiminvoker) InvokeUpdate(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error) {
	claim, err := claimConvertToInternal(obj)
	if err != nil {
		return obj, err
	}
	if err := r.be.Claim(ctx, claim, recursion); err != nil {
		return obj, err
	}
	newClaim, err := claimConvertFromInternal(claim)
	if err != nil {
		return obj, err
	}
	return newClaim, nil
}

func (r *claiminvoker) InvokeDelete(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error) {
	claim, err := claimConvertToInternal(obj)
	if err != nil {
		return obj, err
	}
	if err := r.be.Release(ctx, claim, recursion); err != nil {
		return obj,err
	}
	newClaim, err := claimConvertFromInternal(claim)
	if err != nil {
		return obj, err
	}
	return newClaim, nil
}
