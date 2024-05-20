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

const RegionPlural = "regions"
const RegionSingular = "region"

// +k8s:deepcopy-gen=false
var _ resource.Object = &Region{}
var _ resource.ObjectList = &RegionList{}
var _ backend.ObjectList = &RegionList{}
var _ backend.GenericObject = &Region{}
var _ backend.GenericObjectList = &RegionList{}

// GetListMeta returns the ListMeta
func (r *RegionList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *Region) GetSingularName() string {
	return RegionSingular
}

func (Region) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: RegionPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (Region) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *Region) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (Region) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (Region) New() runtime.Object {
	return &Region{}
}

// NewList implements resource.Object
func (Region) NewList() runtime.Object {
	return &RegionList{}
}

func (r *Region) NewObjList() backend.GenericObjectList {
	return &RegionList{
		TypeMeta: metav1.TypeMeta{APIVersion: SchemeGroupVersion.Identifier(), Kind: RegionKindList},
	}
}

func (r *Region) SchemaGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind(RegionKind)
}

// GetCondition returns the condition based on the condition kind
func (r *Region) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *Region) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// RegionConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for Region
func RegionConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *RegionList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *RegionList) GetObjects() []backend.GenericObject {
	objs := []backend.GenericObject{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *Region) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *Region) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *Region) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *Region) GetRegion() string {
	return r.Name
}

func (r *Region) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      RegionKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *Region) ValidateSyntax(_ string) field.ErrorList {
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

func (r *Region) GetSpec() any {
	return r.Spec
}

func (r *Region) SetSpec(s any) {
	if spec, ok := s.(RegionSpec); ok {
		r.Spec = spec
	}
}

// BuildRegion returns a reource from a client Object a Spec/Status
func BuildRegion(meta metav1.ObjectMeta, spec *RegionSpec, status *RegionStatus) *Region {
	aspec := RegionSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := RegionStatus{}
	if status != nil {
		astatus = *status
	}
	return &Region{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       RegionKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func RegionTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			r, ok := obj.(*Region)
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

func RegionParseFieldSelector(ctx context.Context, fieldSelector fields.Selector) (backend.Filter, error) {
	var filter *RegionFilter

	// add the namespace to the list
	namespace, ok := genericapirequest.NamespaceFrom(ctx)
	if fieldSelector == nil {
		if ok {
			return &RegionFilter{Namespace: namespace}, nil
		}
		return filter, nil
	}
	requirements := fieldSelector.Requirements()
	for _, requirement := range requirements {
		filter = &RegionFilter{}
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
			filter = &RegionFilter{Namespace: namespace}
		}
	}

	return &RegionFilter{}, nil
}

type RegionFilter struct {
	// Name filters by the name of the objects
	Name string `protobuf:"bytes,1,opt,name=name"`

	// Namespace filters by the namespace of the objects
	Namespace string `protobuf:"bytes,2,opt,name=namespace"`
}

func (r *RegionFilter) Filter(ctx context.Context, obj runtime.Object) bool {
	f := false // result of the previous filter
	o, ok := obj.(*Region)
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

func (r *Region) GetUserDefinedLabels() map[string]string {
	return r.Spec.GetUserDefinedLabels()
}

func (r *Region) GetProvider() string {
	return ""
}
