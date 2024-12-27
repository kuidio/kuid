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

package vlan

import (
	"reflect"

	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VLANIndexSpec defines the desired state of VLANIndex
type VLANIndexSpec struct {
	// MinID defines the min VLAN ID the index supports
	// +optional
	MinID *uint32 `json:"minID,omitempty" protobuf:"bytes,1,opt,name=minID"`
	// MaxID defines the max VLAN ID the index supports
	// +optional
	MaxID *uint32 `json:"maxID,omitempty" protobuf:"bytes,2,opt,name=maxID"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	common.UserDefinedLabels `json:",inline" protobuf:"bytes,3,opt,name=userDefinedLabels"`
	// Claims define the embedded claims in the Index
	Claims []VLANIndexClaim `json:"claims,omitempty" protobuf:"bytes,4,rep,name=claims"`
}

type VLANIndexClaim struct {
	// Name of the Claim
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// ID defines the id of the resource
	ID *uint32 `json:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	// Range defines the range of the resource
	// The following notation is used: start-end <start-ID>-<end-ID>
	// the IDs in the range must be consecutive
	Range *string `json:"range,omitempty" protobuf:"bytes,3,opt,name=range"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	common.UserDefinedLabels `json:",inline" protobuf:"bytes,4,opt,name=userDefinedLabels"`
}

// VLANIndexStatus defines the observed state of VLANIndex
type VLANIndexStatus struct {
	// MinID defines the min VLAN ID the index supports
	// +optional
	MinID *uint32 `json:"minID,omitempty" protobuf:"bytes,1,opt,name=minID"`
	// MaxID defines the max VLAN ID the index supports
	// +optional
	MaxID *uint32 `json:"maxID,omitempty" protobuf:"bytes,2,opt,name=maxID"`
	// ConditionedStatus provides the status of the VLANIndex using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" protobuf:"bytes,3,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// VLANIndex is the Schema for the VLANIndex API
type VLANIndex struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   VLANIndexSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status VLANIndexStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// VLANIndexList contains a list of VLANIndexs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VLANIndexList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []VLANIndex `json:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	VLANIndexKind     = reflect.TypeOf(VLANIndex{}).Name()
	VLANIndexListKind = reflect.TypeOf(VLANIndexList{}).Name()
)
