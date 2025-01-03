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

package v1alpha1

import (
	"reflect"

	condv1alpha1 "github.com/kform-dev/choreo/apis/condition/v1alpha1"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EXTCOMMEntrySpec defines the dEXTCOMMred state of EXTCOMMEntry
type EXTCOMMEntrySpec struct {
	// EXTCOMMIndex defines the EXTCOMM index for the EXTCOMM Claim
	Index string `json:"index" protobuf:"bytes,1,opt,name=index"`
	// IndexEntry identifies if the entry is originated from an IP Index
	IndexEntry bool `json:"indexEntry" protobuf:"bytes,2,opt,name=indexEntry"`
	// ClaimType defines the claimType of the EXTCOMM Entry
	ClaimType backend.ClaimType `json:"claimType,omitempty" protobuf:"bytes,3,opt,name=claimType"`
	// ID defines the id of the EXTCOMM entry in the tree
	ID string `json:"id,omitempty" protobuf:"bytes,4,opt,name=id"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	commonv1alpha1.ClaimLabels `json:",inline" protobuf:"bytes,5,opt,name=claimLabels"`
	// Claim defines the name of the claim that is the origin of this  entry
	Claim string `json:"claim" protobuf:"bytes,6,opt,name=claim"`
}

// EXTCOMMEntryStatus defines the observed state of EXTCOMMEntry
type EXTCOMMEntryStatus struct {
	// ConditionedStatus provides the status of the EXTCOMMEntry using conditions
	// - a ready condition indicates the overall status of the resource
	condv1alpha1.ConditionedStatus `json:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// EXTCOMMEntry is the Schema for the EXTCOMMentry API
type EXTCOMMEntry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   EXTCOMMEntrySpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status EXTCOMMEntryStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// EXTCOMMEntryList contains a list of EXTCOMMEntries
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EXTCOMMEntryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []EXTCOMMEntry `json:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	EXTCOMMEntryKind     = reflect.TypeOf(EXTCOMMEntry{}).Name()
	EXTCOMMEntryListKind = reflect.TypeOf(EXTCOMMEntryList{}).Name()
)
