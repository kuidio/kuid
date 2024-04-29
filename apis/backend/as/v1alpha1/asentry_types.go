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

// ASEntrySpec defines the desired state of ASEntry
type ASEntrySpec struct {
	// ASIndex defines the AS index for the AS Claim
	Index string `json:"index" yaml:"index" protobuf:"bytes,1,opt,name=index"`
	// ClaimType defines the claimType of the AS Entry
	ClaimType ASClaimType `json:"claimType,omitempty" yaml:"claimType,omitempty" protobuf:"bytes,2,opt,name=claimType"`
	// ID defines the id of the AS entry in the tree
	ID string `json:"id,omitempty" yaml:"id,omitempty" protobuf:"bytes,3,opt,name=id"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	commonv1alpha1.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,4,opt,name=claimLabels"`
	// Claim defines the name of the claim that is the origin of this  entry
	Claim string `json:"claim" yaml:"claim" protobuf:"bytes,5,opt,name=claim"`
	// Owner defines the ownerReference of the ASClaim
	// Allow for different namesapces, hence it is part of the spec
	Owner *commonv1alpha1.OwnerReference `json:"owner,omitempty" yaml:"owner,omitempty" protobuf:"bytes,6,opt,name=owner"`
}

// ASEntryStatus defines the observed state of ASEntry
type ASEntryStatus struct {
	// ConditionedStatus provides the status of the ASEntry using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ASEntry is the Schema for the ASentry API
//
// +k8s:openapi-gen=true
type ASEntry struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ASEntrySpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status ASEntryStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ASEntryList contains a list of ASEntries
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ASEntryList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ASEntry `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	ASEntryKind = reflect.TypeOf(ASEntry{}).Name()
)
