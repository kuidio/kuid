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

package genericserver

/*

import (
	"context"
	"fmt"

	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/apis/backend"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
)

// parseFieldSelector parses client-provided fields.Selector into a packageFilter
func parseFieldSelector(ctx context.Context, fieldSelector fields.Selector) (backend.Filter, error) {
	var filter *storerFilter

	//log := log.FromContext(ctx)

	// add the namespace to the list
	namespace, ok := genericapirequest.NamespaceFrom(ctx)
	if fieldSelector == nil {
		if ok {
			return &storerFilter{namespace: namespace}, nil
		}
		return filter, nil
	}

	requirements := fieldSelector.Requirements()
	for _, requirement := range requirements {
		filter = &storerFilter{}
		switch requirement.Operator {
		case selection.Equals, selection.DoesNotExist:
			if requirement.Value == "" {
				return filter, apierrors.NewBadRequest(fmt.Sprintf("unsupported fieldSelector value %q for field %q with operator %q", requirement.Value, requirement.Field, requirement.Operator))
			}
		default:
			return filter, apierrors.NewBadRequest(fmt.Sprintf("unsupported fieldSelector operator %q for field %q", requirement.Operator, requirement.Field))
		}

		switch requirement.Field {
		case "metadata.name":
			filter.name = requirement.Value
		case "metadata.namespace":
			filter.namespace = requirement.Value
		default:
			return filter, apierrors.NewBadRequest(fmt.Sprintf("unknown fieldSelector field %q", requirement.Field))
		}
	}
	// add namespace to the filter selector if specified
	if ok {
		if filter != nil {
			filter.namespace = namespace
		} else {
			filter = &storerFilter{namespace: namespace}
		}
	}

	return filter, nil
}

// Filter
type storerFilter struct {
	// Name filters by the name of the objects
	name string

	// Namespace filters by the namespace of the objects
	namespace string
}

func (r *storerFilter) Filter(ctx context.Context, obj runtime.Object) bool {
	log := log.FromContext(ctx)
	f := false // this is the result of the previous filtering
	accessor, err := meta.Accessor(obj)
	if err != nil {
		log.Error("cannot get meta from object", "error", err.Error())
		return f
	}
	if r.name != "" {
		if accessor.GetName() == r.name {
			f = false
		} else {
			f = true
		}
	}
	if r.namespace != "" {
		if accessor.GetNamespace() == r.namespace {
			f = false
		} else {
			f = true
		}
	}
	return f
}
*/