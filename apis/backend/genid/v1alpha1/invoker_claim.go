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

package v1alpha1

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kuidio/kuid/apis/backend/genid"
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

func claimConvertToInternal(obj runtime.Object) (*genid.GENIDClaim, error) {
	ru, ok := obj.(runtime.Unstructured)
	if !ok {
		return nil, fmt.Errorf("not an unstructured obj, got: %s", reflect.TypeOf(obj).Name())
	}
	claim := &GENIDClaim{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ru.UnstructuredContent(), claim); err != nil {
		return nil, fmt.Errorf("unable to convert unstructured object to ipclaim: %v", err)
	}
	claimInternal := &genid.GENIDClaim{}
	if err := Convert_v1alpha1_GENIDClaim_To_genid_GENIDClaim(claim, claimInternal, nil); err != nil {
		return nil, fmt.Errorf("unable to convert unstructured object to ipclaim: %v", err)
	}
	return claimInternal, nil
}

func claimConvertFromInternal(obj runtime.Object) (runtime.Unstructured, error) {
	claimInternal, ok := obj.(*genid.GENIDClaim)
	if !ok {
		return nil, fmt.Errorf("not an unstructured obj, got: %s", reflect.TypeOf(obj).Name())
	}
	claim := &GENIDClaim{}
	if err := Convert_genid_GENIDClaim_To_v1alpha1_GENIDClaim(claimInternal, claim, nil); err != nil {
		return nil, fmt.Errorf("unable to convert unstructured object to ipclaim: %v", err)
	}

	uobj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(claim)
	if err != nil {
		return nil, fmt.Errorf("unable to convert to unstructured: %v", err)
	}

	return &unstructured.Unstructured{Object: uobj}, nil
}

func (r *claiminvoker) InvokeCreate(ctx context.Context, obj runtime.Object, recursion bool) error {
	claim, err := claimConvertToInternal(obj)
	if err != nil {
		return err
	}
	if err := r.be.Claim(ctx, claim, recursion); err != nil {
		return err
	}
	newClaim, err := claimConvertFromInternal(claim)
	if err != nil {
		return err
	}
	obj = newClaim
	return nil
}

func (r *claiminvoker) InvokeUpdate(ctx context.Context, obj, old runtime.Object, recursion bool) error {
	claim, err := claimConvertToInternal(obj)
	if err != nil {
		return err
	}
	if err := r.be.Claim(ctx, claim, recursion); err != nil {
		return err
	}
	newClaim, err := claimConvertFromInternal(claim)
	if err != nil {
		return err
	}
	obj = newClaim
	return nil
}

func (r *claiminvoker) InvokeDelete(ctx context.Context, obj runtime.Object, recursion bool) error {
	claim, err := claimConvertToInternal(obj)
	if err != nil {
		return err
	}
	if err := r.be.Release(ctx, claim, recursion); err != nil {
		return err
	}
	newClaim, err := claimConvertFromInternal(claim)
	if err != nil {
		return err
	}
	obj = newClaim
	return nil
}
