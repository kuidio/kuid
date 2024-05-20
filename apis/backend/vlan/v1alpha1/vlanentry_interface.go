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

const VLANEntryPlural = "vlanentries"
const VLANEntrySingular = "vlanentry"

// +k8s:deepcopy-gen=false
var _ resource.Object = &VLANEntry{}
var _ resource.ObjectList = &VLANEntryList{}

// GetListMeta returns the ListMeta
func (r *VLANEntryList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *VLANEntryList) GetItems() []backend.Object {
	entries := make([]backend.Object, 0, len(r.Items))
	for _, entry := range r.Items {
		entries = append(entries, &entry)
	}
	return entries
}

func (r *VLANEntry) GetSingularName() string {
	return VLANEntrySingular
}

func (VLANEntry) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: VLANEntryPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (VLANEntry) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *VLANEntry) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (VLANEntry) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (VLANEntry) New() runtime.Object {
	return &VLANEntry{}
}

// NewList implements resource.Object
func (VLANEntry) NewList() runtime.Object {
	return &VLANEntryList{}
}

// GetCondition returns the condition based on the condition kind
func (r *VLANEntry) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *VLANEntry) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// VLANEntryConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for VLANEntry
func VLANEntryConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
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

func (r *VLANEntry) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *VLANEntry) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *VLANEntry) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return r.Spec.Owner
}

func (r *VLANEntry) GetClaimName() string {
	return r.Spec.Claim
}

func (r *VLANEntry) GetClaimType() backend.ClaimType {
	return r.Spec.ClaimType
}

func (r *VLANEntry) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

func (r *VLANEntry) GetIndex() string {
	return r.Spec.Index
}

func (r *VLANEntry) GetOwnerGVK() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   r.Spec.Owner.Group,
		Version: r.Spec.Owner.Version,
		Kind:    r.Spec.Owner.Kind,
	}
}

func (r *VLANEntry) GetOwnerNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Spec.Owner.Namespace,
		Name:      r.Spec.Owner.Name,
	}
}

func (r *VLANEntry) GetSpec() any {
	return r.Spec
}

func (r *VLANEntry) GetSpecID() string {
	return r.Spec.ID
}

func (r *VLANEntry) SetSpec(spec any) {
	s, ok := spec.(VLANEntrySpec)
	if ok {
		r.Spec = s
	}
}

func GetVLANEntry(k store.Key, vrange, id string, labels map[string]string) backend.EntryObject {
	//log := log.FromContext(ctx)

	index := k.Name
	ns := k.Namespace

	spec := &VLANEntrySpec{
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

	status := &VLANEntryStatus{}
	status.SetConditions(conditionv1alpha1.Ready())

	id = strings.ReplaceAll(id, "/", "-")
	name := fmt.Sprintf("%s.%s", index, id)
	if vrange != "" {
		name = fmt.Sprintf("%s.%s", vrange, id)
	}

	return BuildVLANEntry(
		metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		spec,
		status,
	)
}

// BuildVLANEntry returns a reource from a client Object a Spec/Status
func BuildVLANEntry(meta metav1.ObjectMeta, spec *VLANEntrySpec, status *VLANEntryStatus) backend.EntryObject {
	aspec := VLANEntrySpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := VLANEntryStatus{}
	if status != nil {
		astatus = *status
	}
	return &VLANEntry{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       VLANEntryKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func (r *VLANEntry) ValidateSyntax(_ string) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

func (r *VLANEntry) NewObjList() backend.GenericObjectList {
	return &VLANEntryList{
		TypeMeta: metav1.TypeMeta{APIVersion: SchemeGroupVersion.Identifier(), Kind: VLANEntryListKind},
	}
}

func (r *VLANEntryList) GetObjects() []backend.GenericObject {
	objs := []backend.GenericObject{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}