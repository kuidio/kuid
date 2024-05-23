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

// ESIClaimSpec defines the desired state of ESIClaim
type ESIClaimSpec struct {
	// ESIIndex defines the ESI index for the ESI Claim
	Index string `json:"index" yaml:"index" protobuf:"bytes,1,opt,name=index"`
	// ESIID defines the ESI for the ESI claim
	ID *uint64 `json:"id,omitempty" yaml:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	// Range defines the ESI range for the ESI claim
	// The following notation is used: start-end <start-ESIID>-<end-ESIID>
	// the ESIs in the range must be consecutive
	Range *string `json:"range,omitempty" yaml:"range,omitempty" protobuf:"bytes,3,opt,name=range"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	commonv1alpha1.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,4,opt,name=claimLabels"`
	// Owner defines the ownerReference of the ESIClaim
	// Allow for different namesapces, hence it is part of the spec
	Owner *commonv1alpha1.OwnerReference `json:"owner,omitempty" yaml:"owner,omitempty" protobuf:"bytes,5,opt,name=owner"`
}

// ESIClaimStatus defines the observed state of ESIClaim
type ESIClaimStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// ESIID defines the ESI for the ESI claim
	// +optional
	ID *uint64 `json:"id,omitempty" yaml:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	// ESIRange defines the ESI range for the ESI claim
	// +optional
	Range *string `json:"range,omitempty" yaml:"range,omitempty" protobuf:"bytes,3,opt,name=range"`
	// ExpiryTime defines when the claim expires
	// +kubebuilder:validation:Optional
	// +optional
	ExpiryTime *string `json:"expiryTime,omitempty" yaml:"expiryTime,omitempty" protobuf:"bytes,4,opt,name=expiryTime"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ESIClaim is the Schema for the ESIClaim API
//
// +k8s:openapi-gen=true
type ESIClaim struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ESIClaimSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status ESIClaimStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ESIClaimList contains a list of ESIClaims
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ESIClaimList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ESIClaim `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	ESIClaimKind     = reflect.TypeOf(ESIClaim{}).Name()
	ESIClaimListKind = reflect.TypeOf(ESIClaimList{}).Name()
)
