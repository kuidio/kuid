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
)

const VXLANEntryPlural = "vxlanentries"
const VXLANEntrySingular = "vxlanentry"

// +k8s:deepcopy-gen=false
var _ resource.Object = &VXLANEntry{}
var _ resource.ObjectList = &VXLANEntryList{}

// GetListMeta returns the ListMeta
func (r *VXLANEntryList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *VXLANEntry) GetSingularName() string {
	return VXLANEntrySingular
}

func (VXLANEntry) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: VXLANEntryPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (VXLANEntry) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *VXLANEntry) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (VXLANEntry) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (VXLANEntry) New() runtime.Object {
	return &VXLANEntry{}
}

// NewList implements resource.Object
func (VXLANEntry) NewList() runtime.Object {
	return &VXLANEntryList{}
}

// GetCondition returns the condition based on the condition kind
func (r *VXLANEntry) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *VXLANEntry) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// ConvertVXLANEntryFieldSelector is the schema conversion function for normalizing the FieldSelector for VXLANEntry
func ConvertVXLANEntryFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
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

func (r *VXLANEntry) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *VXLANEntry) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *VXLANEntry) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return r.Spec.Owner
}

func (r *VXLANEntry) GetClaimName() string {
	return r.Spec.Claim
}

func GetVXLANEntry(ctx context.Context, k store.Key, vrange, id string, labels map[string]string) *VXLANEntry {
	//log := log.FromContext(ctx)

	index := k.Name
	ns := k.Namespace

	spec := &VXLANEntrySpec{
		Index:     index,
		ClaimType: GetClaimTypeFromString(labels[backend.KuidVXLANClaimTypeKey]),
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
		if !backend.BackendSystemKeys.Has(k) && !backend.BackendVXLANSystemKeys.Has(k) {
			udLabels[k] = v
		}
	}
	spec.UserDefinedLabels.Labels = udLabels

	status := &VXLANEntryStatus{}
	status.SetConditions(conditionv1alpha1.Ready())

	id = strings.ReplaceAll(id, "/", "-")
	name := fmt.Sprintf("%s.%s", index, id)
	if vrange != "" {
		name = fmt.Sprintf("%s.%s", vrange, id)
	}

	return BuildVXLANEntry(
		metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		spec,
		status,
	)
}

// BuildVXLANEntry returns a reource from a client Object a Spec/Status
func BuildVXLANEntry(meta metav1.ObjectMeta, spec *VXLANEntrySpec, status *VXLANEntryStatus) *VXLANEntry {
	aspec := VXLANEntrySpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := VXLANEntryStatus{}
	if status != nil {
		astatus = *status
	}
	return &VXLANEntry{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       VXLANEntryKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}
