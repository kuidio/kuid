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

package as

import (
	"reflect"

	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ASClaimSpec defines the desired state of ASClaim
type ASClaimSpec struct {
	// Index defines the index for the AS Claim
	Index string `json:"index" yaml:"index" protobuf:"bytes,1,opt,name=index"`
	// ASID defines the AS for the AS claim
	ID *uint32 `json:"id,omitempty" yaml:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	// Range defines the AS range for the AS claim
	// The following notation is used: start-end <start-ASID>-<end-ASID>
	// the ASs in the range must be consecutive
	Range *string `json:"range,omitempty" yaml:"range,omitempty" protobuf:"bytes,3,opt,name=range"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	common.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,4,opt,name=claimLabels"`
}

// ASClaimStatus defines the observed state of ASClaim
type ASClaimStatus struct {
	// ConditionedStatus provides the status of the ASClain using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// ASID defines the AS for the AS claim
	// +optional
	ID *uint32 `json:"id,omitempty" yaml:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	// ASRange defines the AS range for the AS claim
	// +optional
	Range *string `json:"range,omitempty" yaml:"range,omitempty" protobuf:"bytes,3,opt,name=range"`
	// ExpiryTime defines when the claim expires
	// +optional
	ExpiryTime *string `json:"expiryTime,omitempty" yaml:"expiryTime,omitempty" protobuf:"bytes,4,opt,name=expiryTime"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ASClaim is the Schema for the ASClaim API
type ASClaim struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ASClaimSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status ASClaimStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ASClaimList contains a list of ASClaims
type ASClaimList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ASClaim `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	ASClaimKind     = reflect.TypeOf(ASClaim{}).Name()
	ASClaimListKind = reflect.TypeOf(ASClaimList{}).Name()
)
