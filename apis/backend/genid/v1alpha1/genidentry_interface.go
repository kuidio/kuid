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
	strings "strings"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
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

const GENIDEntryPlural = "genidentries"
const GENIDEntrySingular = "genidentry"

// +k8s:deepcopy-gen=false
var _ resource.Object = &GENIDEntry{}
var _ resource.ObjectList = &GENIDEntryList{}

// GetListMeta returns the ListMeta
func (r *GENIDEntryList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *GENIDEntryList) GetItems() []backend.Object {
	entries := make([]backend.Object, 0, len(r.Items))
	for _, entry := range r.Items {
		entries = append(entries, &entry)
	}
	return entries
}

func (r *GENIDEntry) GetSingularName() string {
	return GENIDEntrySingular
}

func (GENIDEntry) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: GENIDEntryPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (GENIDEntry) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *GENIDEntry) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (GENIDEntry) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (GENIDEntry) New() runtime.Object {
	return &GENIDEntry{}
}

// NewList implements resource.Object
func (GENIDEntry) NewList() runtime.Object {
	return &GENIDEntryList{}
}

// GetCondition returns the condition based on the condition kind
func (r *GENIDEntry) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *GENIDEntry) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// GENIDEntryConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for GENIDEntry
func GENIDEntryConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	case "spec.index":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *GENIDEntry) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *GENIDEntry) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *GENIDEntry) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return r.Spec.Owner
}

func (r *GENIDEntry) GetClaimName() string {
	return r.Spec.Claim
}

func (r *GENIDEntry) GetClaimType() backend.ClaimType {
	return r.Spec.ClaimType
}

func (r *GENIDEntry) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

func (r *GENIDEntry) GetIndex() string {
	return r.Spec.Index
}

func (r *GENIDEntry) GetOwnerGVK() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   r.Spec.Owner.Group,
		Version: r.Spec.Owner.Version,
		Kind:    r.Spec.Owner.Kind,
	}
}

func (r *GENIDEntry) GetOwnerNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Spec.Owner.Namespace,
		Name:      r.Spec.Owner.Name,
	}
}

func (r *GENIDEntry) GetSpec() any {
	return r.Spec
}

func (r *GENIDEntry) GetSpecID() string {
	return r.Spec.ID
}

func (r *GENIDEntry) SetSpec(s any) {
	if spec, ok := s.(GENIDEntrySpec); ok {
		r.Spec = spec
	}
}

func GetGENIDEntry(k store.Key, vrange, id string, labels map[string]string) backend.EntryObject {
	//log := log.FromContext(ctx)

	index := k.Name
	ns := k.Namespace

	spec := &GENIDEntrySpec{
		Index:     index,
		ClaimType: backend.GetClaimTypeFromString(labels[backend.KuidClaimTypeKey]),
		Claim:     labels[backend.KuidClaimNameKey],
		ID:        id,
		Owner: &commonv1alpha1.OwnerReference{
			Group:     labels[backend.KuidOwnerGroupKey],
			Version:   labels[backend.KuidOwnerVersionKey],
			Kind:      labels[backend.KuidOwnerKindKey],
			Namespace: labels[backend.KuidOwnerNamespaceKey],
			Name:      labels[backend.KuidOwnerNameKey],
		},
	}
	// filter the system defined labels from the labels to prepare for the user defined labels
	udLabels := map[string]string{}
	for k, v := range labels {
		if !backend.BackendSystemKeys.Has(k) {
			udLabels[k] = v
		}
	}
	spec.UserDefinedLabels.Labels = udLabels

	status := &GENIDEntryStatus{}
	status.SetConditions(conditionv1alpha1.Ready())

	id = strings.ReplaceAll(id, "/", "-")
	name := fmt.Sprintf("%s.%s", index, id)
	if vrange != "" {
		name = fmt.Sprintf("%s.%s", vrange, id)
	}

	return BuildGENIDEntry(
		metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		spec,
		status,
	)
}

// BuildGENIDEntry returns a reource from a client Object a Spec/Status
func BuildGENIDEntry(meta metav1.ObjectMeta, spec *GENIDEntrySpec, status *GENIDEntryStatus) *GENIDEntry {
	aspec := GENIDEntrySpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := GENIDEntryStatus{}
	if status != nil {
		astatus = *status
	}
	return &GENIDEntry{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       GENIDEntryKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func (r *GENIDEntry) ValidateSyntax(_ string) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

func (r *GENIDEntry) NewObjList() backend.GenericObjectList {
	return &GENIDEntryList{
		TypeMeta: metav1.TypeMeta{APIVersion: SchemeGroupVersion.Identifier(), Kind: GENIDEntryListKind},
	}
}

func (r *GENIDEntryList) GetObjects() []backend.GenericObject {
	objs := []backend.GenericObject{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}