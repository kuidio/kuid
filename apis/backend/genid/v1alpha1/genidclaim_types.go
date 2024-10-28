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
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GENIDClaimSpec defines the desired state of GENIDClaim
type GENIDClaimSpec struct {
	// Index defines the index for the resource
	Index string `json:"index" yaml:"index" protobuf:"bytes,1,opt,name=index"`
	// ID defines the id of the resource
	ID *uint64 `json:"id,omitempty" yaml:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	// Range defines the VLAN range of the resource
	// The following notation is used: start-end <start-VLANID>-<end-VLANID>
	// the VLANs in the range must be consecutive
	Range *string `json:"range,omitempty" yaml:"range,omitempty" protobuf:"bytes,3,opt,name=range"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	commonv1alpha1.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,4,opt,name=claimLabels"`
}

// GENIDClaimStatus defines the observed state of GENIDClaim
type GENIDClaimStatus struct {
	// ConditionedStatus provides the status of the VLANClain using conditions
	// - a ready condition indicates the overall status of the resource
	condv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// ID defines the ID assigned to the resource
	// +optional
	ID *uint64 `json:"id,omitempty" yaml:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	// Range defines the range assigned to the resource
	// +optional
	Range *string `json:"range,omitempty" yaml:"range,omitempty" protobuf:"bytes,3,opt,name=range"`
	// ExpiryTime defines when the claim expires
	// +kubebuilder:validation:Optional
	// +optional
	ExpiryTime *string `json:"expiryTime,omitempty" yaml:"expiryTime,omitempty" protobuf:"bytes,4,opt,name=expiryTime"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GENIDClaim is the Schema for the GENIDClaim API
type GENIDClaim struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   GENIDClaimSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status GENIDClaimStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// GENIDClaimList contains a list of GENIDClaims
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type GENIDClaimList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []GENIDClaim `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	GENIDClaimKind     = reflect.TypeOf(GENIDClaim{}).Name()
	GENIDClaimListKind = reflect.TypeOf(GENIDClaimList{}).Name()
)
