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

const ESIEntryPlural = "esientries"
const ESIEntrySingular = "esientry"

// +k8s:deepcopy-gen=false
var _ resource.Object = &ESIEntry{}
var _ resource.ObjectList = &ESIEntryList{}

// GetListMeta returns the ListMeta
func (r *ESIEntryList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *ESIEntryList) GetItems() []*ESIEntry {
	entries := make([]*ESIEntry, 0, len(r.Items))
	for _, entry := range r.Items {
		entries = append(entries, &entry)
	}
	return entries
}

func (r *ESIEntry) GetSingularName() string {
	return ESIEntrySingular
}

func (ESIEntry) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: ESIEntryPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (ESIEntry) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *ESIEntry) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (ESIEntry) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (ESIEntry) New() runtime.Object {
	return &ESIEntry{}
}

// NewList implements resource.Object
func (ESIEntry) NewList() runtime.Object {
	return &ESIEntryList{}
}

// GetCondition returns the condition based on the condition kind
func (r *ESIEntry) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *ESIEntry) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// ConvertESIEntryFieldSelector is the schema conversion function for normalizing the FieldSelector for ESIEntry
func ConvertESIEntryFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
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

func (r *ESIEntry) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *ESIEntry) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *ESIEntry) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

func (r *ESIEntry) GetIndex() string {
	return r.Spec.Index
}

func (r *ESIEntry) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return r.Spec.Owner
}

func (r *ESIEntry) GetClaimName() string {
	return r.Spec.Claim
}

func (r *ESIEntry) GetSpecID() string {
	return r.Spec.ID
}

func GetESIEntry(ctx context.Context, k store.Key, vrange, id string, labels map[string]string) *ESIEntry {
	//log := log.FromContext(ctx)

	index := k.Name
	ns := k.Namespace

	spec := &ESIEntrySpec{
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

	status := &ESIEntryStatus{}
	status.SetConditions(conditionv1alpha1.Ready())

	id = strings.ReplaceAll(id, "/", "-")
	name := fmt.Sprintf("%s.%s", index, id)
	if vrange != "" {
		name = fmt.Sprintf("%s.%s", vrange, id)
	}

	return BuildESIEntry(
		metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		spec,
		status,
	)
}

// BuildESIEntry returns a reource from a client Object a Spec/Status
func BuildESIEntry(meta metav1.ObjectMeta, spec *ESIEntrySpec, status *ESIEntryStatus) *ESIEntry {
	aspec := ESIEntrySpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := ESIEntryStatus{}
	if status != nil {
		astatus = *status
	}
	return &ESIEntry{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       ESIEntryKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func (r *ESIEntry) ValidateSyntax(_ string) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

func (r *ESIEntry) GetSpec() any {
	return r.Spec
}

func (r *ESIEntry) SetSpec(s any) {
	if spec, ok := s.(ESIEntrySpec); ok {
		r.Spec = spec
	}
}

func (r *ESIEntry) NewObjList() backend.GenericObjectList {
	return &ESIEntryList{
		TypeMeta: metav1.TypeMeta{APIVersion: SchemeGroupVersion.Identifier(), Kind: ESIEntryListKind},
	}
}

func (r *ESIEntryList) GetObjects() []backend.GenericObject {
	objs := []backend.GenericObject{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}
