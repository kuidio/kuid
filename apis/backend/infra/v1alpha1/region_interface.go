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
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const RegionPlural = "regions"
const RegionSingular = "region"

// +k8s:deepcopy-gen=false
var _ resource.Object = &Region{}
var _ resource.ObjectList = &RegionList{}

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

func (r *Region) ValidateSyntax() field.ErrorList {
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
