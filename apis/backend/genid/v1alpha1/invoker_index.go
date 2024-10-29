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

func NewChoreoIndexInvoker(be backend.Backend) options.BackendInvoker {
	return &idxinvoker{
		be: be,
	}
}

type idxinvoker struct {
	be backend.Backend
}

func indexConvertToInternal(obj runtime.Object) (*genid.GENIDIndex, error) {
	ru, ok := obj.(runtime.Unstructured)
	if !ok {
		return nil, fmt.Errorf("not an unstructured obj, got: %s", reflect.TypeOf(obj).Name())
	}
	index := &GENIDIndex{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ru.UnstructuredContent(), index); err != nil {
		return nil, fmt.Errorf("unable to convert unstructured object to index: %v", err)
	}
	indexInternal := &genid.GENIDIndex{}
	if err := Convert_v1alpha1_GENIDIndex_To_genid_GENIDIndex(index, indexInternal, nil); err != nil {
		return nil, fmt.Errorf("unable to convert unstructured object to index: %v", err)
	}
	return indexInternal, nil
}

func indexConvertFromInternal(obj runtime.Object) (runtime.Unstructured, error) {
	indexInternal, ok := obj.(*genid.GENIDIndex)
	if !ok {
		return nil, fmt.Errorf("not an unstructured obj, got: %s", reflect.TypeOf(obj).Name())
	}
	index := &GENIDIndex{}
	if err := Convert_genid_GENIDIndex_To_v1alpha1_GENIDIndex(indexInternal, index, nil); err != nil {
		return nil, fmt.Errorf("unable to convert unstructured object to index: %v", err)
	}

	uobj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(index)
	if err != nil {
		return nil, fmt.Errorf("unable to convert to unstructured: %v", err)
	}

	return &unstructured.Unstructured{Object: uobj}, nil
}

func (r *idxinvoker) InvokeCreate(ctx context.Context, obj runtime.Object, recursion bool) error {
	index, err := indexConvertToInternal(obj)
	if err != nil {
		return err
	}
	if err := r.be.CreateIndex(ctx, index); err != nil {
		return err
	}
	newIndex, err := indexConvertFromInternal(index)
	if err != nil {
		return err
	}
	obj = newIndex
	return nil
}

func (r *idxinvoker) InvokeUpdate(ctx context.Context, obj, old runtime.Object, recursion bool) error {
	index, err := indexConvertToInternal(obj)
	if err != nil {
		return err
	}
	if err := r.be.CreateIndex(ctx, index); err != nil {
		return err
	}
	newIndex, err := indexConvertFromInternal(index)
	if err != nil {
		return err
	}
	obj = newIndex
	return nil
}

func (r *idxinvoker) InvokeDelete(ctx context.Context, obj runtime.Object, recursion bool) error {
	index, err := indexConvertToInternal(obj)
	if err != nil {
		return err
	}
	if err := r.be.DeleteIndex(ctx, index); err != nil {
		return err
	}
	newIndex, err := indexConvertFromInternal(index)
	if err != nil {
		return err
	}
	obj = newIndex
	return nil
}
