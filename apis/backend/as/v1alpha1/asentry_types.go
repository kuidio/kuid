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

	condv1alpha1 "github.com/kform-dev/choreo/apis/condition/v1alpha1"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ASEntrySpec defines the desired state of ASEntry
type ASEntrySpec struct {
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
	commonv1alpha1.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,5,opt,name=claimLabels"`
}

// ASEntryStatus defines the observed state of ASEntry
type ASEntryStatus struct {
	// ConditionedStatus provides the status of the ASEntry using conditions
	// - a ready condition indicates the overall status of the resource
	condv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// ASEntry is the Schema for the ASentry API
type ASEntry struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ASEntrySpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status ASEntryStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// ASEntryList contains a list of ASEntries
type ASEntryList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ASEntry `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	ASEntryKind     = reflect.TypeOf(ASEntry{}).Name()
	ASEntryListKind = reflect.TypeOf(ASEntryList{}).Name()
)
