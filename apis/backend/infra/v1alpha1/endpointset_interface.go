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
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
)

const EndpointSetPlural = "endpointsets"
const EndpointSetSingular = "endpointset"

// +k8s:deepcopy-gen=false
var _ resource.Object = &EndpointSet{}
var _ resource.ObjectList = &EndpointSetList{}
var _ backend.ObjectList = &EndpointSetList{}
var _ backend.GenericObject = &EndpointSet{}
var _ backend.GenericObjectList = &EndpointSetList{}

// GetListMeta returns the ListMeta
func (r *EndpointSetList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *EndpointSet) GetSingularName() string {
	return EndpointSetSingular
}

func (EndpointSet) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: EndpointSetPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (EndpointSet) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *EndpointSet) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (EndpointSet) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (EndpointSet) New() runtime.Object {
	return &EndpointSet{}
}

// NewList implements resource.Object
func (EndpointSet) NewList() runtime.Object {
	return &EndpointSetList{}
}

func (r *EndpointSet) NewObjList() backend.GenericObjectList {
	return &EndpointSetList{
		TypeMeta: metav1.TypeMeta{APIVersion: SchemeGroupVersion.Identifier(), Kind: EndpointSetKindList},
	}
}

func (r *EndpointSet) SchemaGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind(EndpointSetKind)
}

// GetCondition returns the condition based on the condition kind
func (r *EndpointSet) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *EndpointSet) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// EndpointSetConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for EndpointSet
func EndpointSetConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *EndpointSetList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}


func (r *EndpointSetList) GetObjects() []backend.GenericObject {
	objs := []backend.GenericObject{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *EndpointSet) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *EndpointSet) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *EndpointSet) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *EndpointSet) GetEndpointSet() string {
	return r.Name
}

func (r *EndpointSet) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      EndpointSetKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *EndpointSet) ValidateSyntax() field.ErrorList {
	var allErrs field.ErrorList

	/*
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.type"),
				r,
				fmt.Errorf("invalid GENID Type %s", r.Spec.Type).Error(),
			))
		}
	*/
	return allErrs
}

func (r *EndpointSet) GetSpec() any {
	return r.Spec
}

func (r *EndpointSet) SetSpec(s any) {
	if spec, ok := s.(EndpointSetSpec); ok {
		r.Spec = spec
	}
}

// BuildEndpointSet returns a reource from a client Object a Spec/Status
func BuildEndpointSet(meta metav1.ObjectMeta, spec *EndpointSetSpec, status *EndpointSetStatus) *EndpointSet {
	aspec := EndpointSetSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := EndpointSetStatus{}
	if status != nil {
		astatus = *status
	}
	return &EndpointSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       EndpointSetKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func EndpointSetTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			r, ok := obj.(*EndpointSet)
			if !ok {
				return nil
			}
			return []interface{}{
				r.GetName(),
				r.GetCondition(conditionv1alpha1.ConditionTypeReady).Status,
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
		},
	}
}

func EndpointSetParseFieldSelector(ctx context.Context, fieldSelector fields.Selector) (backend.Filter, error) {
	var filter *EndpointSetFilter

	// add the namespace to the list
	namespace, ok := genericapirequest.NamespaceFrom(ctx)
	if fieldSelector == nil {
		if ok {
			return &EndpointSetFilter{Namespace: namespace}, nil
		}
		return filter, nil
	}
	requirements := fieldSelector.Requirements()
	for _, requirement := range requirements {
		filter = &EndpointSetFilter{}
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
			filter.Name = requirement.Value
		case "metadata.namespace":
			filter.Namespace = requirement.Value
		default:
			return filter, apierrors.NewBadRequest(fmt.Sprintf("unknown fieldSelector field %q", requirement.Field))
		}
	}
	// add namespace to the filter selector if specified
	if ok {
		if filter != nil {
			filter.Namespace = namespace
		} else {
			filter = &EndpointSetFilter{Namespace: namespace}
		}
	}

	return &EndpointSetFilter{}, nil
}

type EndpointSetFilter struct {
	// Name filters by the name of the objects
	Name string `protobuf:"bytes,1,opt,name=name"`

	// Namespace filters by the namespace of the objects
	Namespace string `protobuf:"bytes,2,opt,name=namespace"`
}

func (r *EndpointSetFilter) Filter(ctx context.Context, obj runtime.Object) bool {
	f := true
	o, ok := obj.(*EndpointSet)
	if !ok {
		return f
	}
	if r.Name != "" {
		if o.GetName() == r.Name {
			f = false
		} else {
			f = true
		}
	}
	if r.Namespace != "" {
		if o.GetNamespace() == r.Namespace {
			f = false
		} else {
			f = true
		}
	}
	return f
}
