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
	"reflect"

	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPEntrySpec defines the desired state of IPEntry
type VLANEntrySpec struct {
	// VLANIndex defines the vlan index for the VLAN Claim
	Index string `json:"index" yaml:"index" protobuf:"bytes,1,opt,name=index"`
	// ClaimType defines the claimType of the VLAN Entry
	ClaimType VLANClaimType `json:"claimType,omitempty" yaml:"claimType,omitempty" protobuf:"bytes,2,opt,name=claimType"`
	// VLANID defines the vlan id for the VLAN entry
	ID uint32 `json:"id,omitempty" yaml:"id,omitempty" protobuf:"bytes,3,opt,name=id"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	commonv1alpha1.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,4,opt,name=claimLabels"`
	// Claim defines the name of the claim that is the origin of this  entry
	Claim string `json:"claim" yaml:"claim" protobuf:"bytes,5,opt,name=claim"`
	// Owner defines the ownerReference of the VLANClaim
	// Allow for different namesapces, hence it is part of the spec
	Owner *commonv1alpha1.OwnerReference `json:"owner,omitempty" yaml:"owner,omitempty" protobuf:"bytes,6,opt,name=owner"`
}

// VLANEntryStatus defines the observed state of VLANEntry
type VLANEntryStatus struct {
	// ConditionedStatus provides the status of the VLANEntry using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VLANEntry is the Schema for the vlanentry API
//
// +k8s:openapi-gen=true
type VLANEntry struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   VLANEntrySpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status VLANEntryStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// VLANEntryList contains a list of VLANEntries
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VLANEntryList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []VLANEntry `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	VLANEntryKind = reflect.TypeOf(VLANEntry{}).Name()
)