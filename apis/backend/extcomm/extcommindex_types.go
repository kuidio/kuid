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

// EXTCOMMIndexSpec defines the dEXTCOMMred state of EXTCOMMIndex
type EXTCOMMIndexSpec struct {
	// MinID defines the min EXTCOMM ID the index supports
	// +optional
	MinID *uint64 `json:"minID,omitempty" yaml:"minID,omitempty" protobuf:"bytes,1,opt,name=minID"`
	// MaxID defines the max EXTCOMM ID the index supports
	// +optional
	MaxID *uint64 `json:"maxID,omitempty" yaml:"maxID,omitempty" protobuf:"bytes,2,opt,name=maxID"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	common.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,3,opt,name=userDefinedLabels"`
	// Transitive defines the transative nature of the extended community
	Transitive bool `json:"transitive,omitempty" yaml:"transitive,omitempty" protobuf:"bytes,4,opt,name=transitive"`
	// Type defines the type of the extended community
	// 2byteAS, 4byteAS, ipv4Address, opaque
	Type string `json:"type" yaml:"type" protobuf:"bytes,5,opt,name=type"`
	// SubType defines the subTyoe of the extended community
	// routeTarget, routeOrigin;
	SubType string `json:"subType" yaml:"subType" protobuf:"bytes,6,opt,name=subType"`
	// GlobalID is interpreted dependeing on the type
	// AS in case of 2byteAS, 4byteAS
	// IPV4 addrress
	// irrelevant for the opaque type
	GlobalID string `json:"globalID,omitempty" yaml:"globalID,omitempty" protobuf:"bytes,7,opt,name=globalID"`
}

// EXTCOMMIndexStatus defines the observed state of EXTCOMMIndex
type EXTCOMMIndexStatus struct {
	// MinID defines the min EXTCOMM ID the index supports
	// +optional
	MinID *int64 `json:"minID,omitempty" yaml:"minID,omitempty" protobuf:"bytes,1,opt,name=minID"`
	// MaxID defines the max EXTCOMM ID the index supports
	// +optional
	MaxID *int64 `json:"maxID,omitempty" yaml:"maxID,omitempty" protobuf:"bytes,2,opt,name=maxID"`
	// ConditionedStatus provides the status of the EXTCOMMIndex using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,3,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// EXTCOMMIndex is the Schema for the EXTCOMMIndex API
type EXTCOMMIndex struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   EXTCOMMIndexSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status EXTCOMMIndexStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// EXTCOMMIndexList contains a list of EXTCOMMIndexs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EXTCOMMIndexList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []EXTCOMMIndex `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	EXTCOMMIndexKind     = reflect.TypeOf(EXTCOMMIndex{}).Name()
	EXTCOMMIndexListKind = reflect.TypeOf(EXTCOMMIndexList{}).Name()
)
