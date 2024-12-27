/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "VLAN IS" BVLANIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vlan

import (
	"reflect"

	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/common"
	"github.com/kuidio/kuid/apis/backend"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VLANEntrySpec defines the desired state of VLANEntry
type VLANEntrySpec struct {
	// Index defines the index for the resource
	Index string `json:"index" protobuf:"bytes,1,opt,name=index"`
	// IndexEntry identifies if the entry is originated from an IP Index
	IndexEntry bool `json:"indexEntry" protobuf:"bytes,2,opt,name=indexEntry"`
	// ClaimType defines the claimType of the resource
	ClaimType backend.ClaimType `json:"claimType,omitempty" protobuf:"bytes,3,opt,name=claimType"`
	// ID defines the id of the resource in the tree
	ID string `json:"id,omitempty" protobuf:"bytes,4,opt,name=id"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	common.ClaimLabels `json:",inline" protobuf:"bytes,5,opt,name=claimLabels"`
}

// VLANEntryStatus defines the observed state of VLANEntry
type VLANEntryStatus struct {
	// ConditionedStatus provides the status of the VLANEntry using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// VLANEntry is the Schema for the VLANentry API
type VLANEntry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   VLANEntrySpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status VLANEntryStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// VLANEntryList contains a list of VLANEntries
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VLANEntryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []VLANEntry `json:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	VLANEntryKind     = reflect.TypeOf(VLANEntry{}).Name()
	VLANEntryListKind = reflect.TypeOf(VLANEntryList{}).Name()
)
