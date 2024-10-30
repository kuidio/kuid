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

package as

import (
	"context"
	"fmt"
	"reflect"

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

func indexConvertToInternal(obj runtime.Object) (*ASIndex, error) {
	ru, ok := obj.(runtime.Unstructured)
	if !ok {
		return nil, fmt.Errorf("not an unstructured obj, got: %s", reflect.TypeOf(obj).Name())
	}
	index := &ASIndex{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(ru.UnstructuredContent(), index); err != nil {
		return nil, fmt.Errorf("unable to convert unstructured object to index: %v", err)
	}
	return index, nil
}

func indexConvertFromInternal(obj runtime.Object) (runtime.Unstructured, error) {
	index, ok := obj.(*ASIndex)
	if !ok {
		return nil, fmt.Errorf("not an unstructured obj, got: %s", reflect.TypeOf(obj).Name())
	}

	uobj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(index)
	if err != nil {
		return nil, fmt.Errorf("unable to convert to unstructured: %v", err)
	}

	return &unstructured.Unstructured{Object: uobj}, nil
}

func (r *idxinvoker) InvokeCreate(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error) {
	index, err := indexConvertToInternal(obj)
	if err != nil {
		return obj, err
	}
	if err := r.be.CreateIndex(ctx, index); err != nil {
		return obj, err
	}
	newIndex, err := indexConvertFromInternal(index)
	if err != nil {
		return obj, err
	}
	return newIndex, nil
}

func (r *idxinvoker) InvokeUpdate(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error) {
	index, err := indexConvertToInternal(obj)
	if err != nil {
		return obj, err
	}
	if err := r.be.CreateIndex(ctx, index); err != nil {
		return obj, err
	}
	newIndex, err := indexConvertFromInternal(index)
	if err != nil {
		return obj, err
	}
	return newIndex, nil
}

func (r *idxinvoker) InvokeDelete(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error) {
	index, err := indexConvertToInternal(obj)
	if err != nil {
		return obj, err
	}
	if err := r.be.DeleteIndex(ctx, index); err != nil {
		return obj, err
	}
	newIndex, err := indexConvertFromInternal(index)
	if err != nil {
		return obj, err
	}
	return newIndex, nil
}
