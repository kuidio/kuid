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

package extcomm

import (
	"reflect"

	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EXTCOMMClaimSpec defines the dEXTCOMMred state of EXTCOMMClaim
type EXTCOMMClaimSpec struct {
	// EXTCOMMIndex defines the EXTCOMM index for the EXTCOMM Claim
	Index string `json:"index" protobuf:"bytes,1,opt,name=index"`
	// EXTCOMMID defines the EXTCOMM for the EXTCOMM claim
	ID *uint64 `json:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	// Range defines the range of the resource
	// The following notation is used: start-end <start-ID>-<end-ID>
	// the IDs in the range must be consecutive
	Range *string `json:"range,omitempty" protobuf:"bytes,3,opt,name=range"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	common.ClaimLabels `json:",inline" protobuf:"bytes,4,opt,name=claimLabels"`
}

// EXTCOMMClaimStatus defines the observed state of EXTCOMMClaim
type EXTCOMMClaimStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// EXTCOMMID defines the EXTCOMM for the EXTCOMM claim
	// +optional
	ID *uint64 `json:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	// EXTCOMMRange defines the EXTCOMM range for the EXTCOMM claim
	// +optional
	Range *string `json:"range,omitempty" protobuf:"bytes,3,opt,name=range"`
	// ExpiryTime defines when the claim expires
	// +kubebuilder:validation:Optional
	// +optional
	ExpiryTime *string `json:"expiryTime,omitempty" protobuf:"bytes,4,opt,name=expiryTime"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// EXTCOMMClaim is the Schema for the EXTCOMMClaim API
type EXTCOMMClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   EXTCOMMClaimSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status EXTCOMMClaimStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// EXTCOMMClaimList contains a list of EXTCOMMClaims
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EXTCOMMClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []EXTCOMMClaim `json:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	EXTCOMMClaimKind     = reflect.TypeOf(EXTCOMMClaim{}).Name()
	EXTCOMMClaimListKind = reflect.TypeOf(EXTCOMMClaimList{}).Name()
)
