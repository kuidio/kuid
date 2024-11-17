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

package genid

import (
	"reflect"

	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GENIDEntrySpec defines the desired state of GENIDEntry
type GENIDEntrySpec struct {
	// Index defines the index for the resource
	Index string `json:"index" yaml:"index" protobuf:"bytes,1,opt,name=index"`
	// IndexEntry identifies if the entry is originated from an IP Index
	IndexEntry bool `json:"indexEntry" yaml:"indexEntry" protobuf:"bytes,2,opt,name=indexEntry"`
	// ClaimType defines the claimType of the resource
	ClaimType backend.ClaimType `json:"claimType,omitempty" yaml:"claimType,omitempty" protobuf:"bytes,3,opt,name=claimType"`
	// ID defines the id of the resource in the tree
	ID string `json:"id,omitempty" yaml:"id,omitempty" protobuf:"bytes,4,opt,name=id"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	common.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,5,opt,name=claimLabels"`
}

// GENIDEntryStatus defines the observed state of GENIDEntry
type GENIDEntryStatus struct {
	// ConditionedStatus provides the status of the GENIDEntry using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// GENIDEntry is the Schema for the GENIDEntry API
type GENIDEntry struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   GENIDEntrySpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status GENIDEntryStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// GENIDEntryList contains a list of VLANEntries
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type GENIDEntryList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []GENIDEntry `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	GENIDEntryKind     = reflect.TypeOf(GENIDEntry{}).Name()
	GENIDEntryListKind = reflect.TypeOf(GENIDEntryList{}).Name()
)
