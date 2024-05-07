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

	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VXLANEntrySpec defines the desired state of VXLANEntry
type VXLANEntrySpec struct {
	// VXLANIndex defines the VXLAN index for the VXLAN Claim
	Index string `json:"index" yaml:"index" protobuf:"bytes,1,opt,name=index"`
	// ClaimType defines the claimType of the VXLAN Entry
	ClaimType backend.ClaimType `json:"claimType,omitempty" yaml:"claimType,omitempty" protobuf:"bytes,2,opt,name=claimType"`
	// ID defines the id of the VXLAN entry in the tree
	ID string `json:"id,omitempty" yaml:"id,omitempty" protobuf:"bytes,3,opt,name=id"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	commonv1alpha1.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,4,opt,name=claimLabels"`
	// Claim defines the name of the claim that is the origin of this  entry
	Claim string `json:"claim" yaml:"claim" protobuf:"bytes,5,opt,name=claim"`
	// Owner defines the ownerReference of the VXLANClaim
	// Allow for different namesapces, hence it is part of the spec
	Owner *commonv1alpha1.OwnerReference `json:"owner,omitempty" yaml:"owner,omitempty" protobuf:"bytes,6,opt,name=owner"`
}

// VXLANEntryStatus defines the observed state of VXLANEntry
type VXLANEntryStatus struct {
	// ConditionedStatus provides the status of the VXLANEntry using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VXLANEntry is the Schema for the VXLANentry API
//
// +k8s:openapi-gen=true
type VXLANEntry struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   VXLANEntrySpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status VXLANEntryStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// VXLANEntryList contains a list of VXLANEntries
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VXLANEntryList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []VXLANEntry `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	VXLANEntryKind = reflect.TypeOf(VXLANEntry{}).Name()
)
