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
	"github.com/henderiw/store"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	rresource "github.com/kuidio/kuid/pkg/reconcilers/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

const ASIndexPlural = "asindices"
const ASIndexSingular = "asindex"

// +k8s:deepcopy-gen=false
var _ resource.Object = &ASIndex{}
var _ resource.ObjectList = &ASIndexList{}

// GetListMeta returns the ListMeta
func (r *ASIndexList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *ASIndex) GetSingularName() string {
	return ASIndexSingular
}

func (ASIndex) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: ASIndexPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (ASIndex) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *ASIndex) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (ASIndex) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (ASIndex) New() runtime.Object {
	return &ASIndex{}
}

// NewList implements resource.Object
func (ASIndex) NewList() runtime.Object {
	return &ASIndexList{}
}

// GetCondition returns the condition based on the condition kind
func (r *ASIndex) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *ASIndex) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// ConvertASIndexFieldSelector is the schema conversion function for normalizing the FieldSelector for ASIndex
func ConvertASIndexFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *ASIndexList) GetItems() []rresource.Object {
	objs := []rresource.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *ASIndex) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *ASIndex) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *ASIndex) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *ASIndex) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      ASIndexKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

// BuildASIndex returns a reource from a client Object a Spec/Status
func BuildASIndex(meta metav1.ObjectMeta, spec *ASIndexSpec, status *ASIndexStatus) *ASIndex {
	aspec := ASIndexSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := ASIndexStatus{}
	if status != nil {
		astatus = *status
	}
	return &ASIndex{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       ASIndexKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}
