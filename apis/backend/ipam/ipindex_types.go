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

package ipam

import (
	"reflect"

	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPIndexSpec defines the desired state of IPIndex
type IPIndexSpec struct {
	// Prefixes define the aggregate prefixes for the network instance
	// A Network instance needs at least 1 prefix to be defined to become operational
	Prefixes []Prefix `json:"prefixes" yaml:"prefixes" protobuf:"bytes,1,opt,name=prefixes"`
}

type Prefix struct {
	// Prefix defines the ip cidr in prefix notation.
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])/(([0-9])|([1-2][0-9])|(3[0-2]))|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))`
	Prefix string `json:"prefix" yaml:"prefix" protobuf:"bytes,1,opt,name=prefix"`
	// PrefixType defines the prefixtype of IPIndex;
	// - network kind is used for physical, virtual nics on a device
	// - pool kind is used for allocating dedicated IP addresses
	// - aggregate kind is used for claiming an aggregate prefix; only used for networkInstance prefixes
	// +kubebuilder:validation:Enum=network;aggregate;pool;
	// +kubebuilder:default:=aggregate
	// +optional
	PrefixType *IPPrefixType `json:"prefixType,omitempty" yaml:"prefixType,omitempty" protobuf:"bytes,2,opt,name=prefixType"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	common.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,2,opt,name=userDefinedLabels"`
}

// IPIndexStatus defines the observed state of IPIndex
type IPIndexStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// Prefixes defines the prefixes, claimed through the IPAM backend
	Prefixes []Prefix `json:"prefixes,omitempty" yaml:"prefixes,omitempty" protobuf:"bytes,2,rep,name=prefixes"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// IPIndex is the Schema for the IPIndex API
type IPIndex struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   IPIndexSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status IPIndexStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// IPIndexList contains a list of IPIndexs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IPIndexList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []IPIndex `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	IPIndexKind     = reflect.TypeOf(IPIndex{}).Name()
	IPIndexListKind = reflect.TypeOf(IPIndexList{}).Name()
)
