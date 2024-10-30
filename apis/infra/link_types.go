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

package infra

import (
	"reflect"

	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/id"
	"github.com/kuidio/kuid/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LinkSpec defines the desired state of Link
type LinkSpec struct {
	// Endpoints define the 2 endpoint identifiers of the link
	// Can only have 2 endpoints
	Endpoints []*id.PartitionEndpointID `json:"endpoints" yaml:"endpoints" protobuf:"bytes,1,opt,name=endpoints"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	common.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,2,opt,name=userDefinedLabels"`
	// BFD defines the BFD specific parameters on the link
	// +optional
	//BFD *BFDLinkParameters `json:"bfd,omitempty" yaml:"bfd,omitempty" protobuf:"bytes,3,opt,name=bfd"`
	// OSPF defines the OSPF specific parameters on the link
	// +optional
	//OSPF *OSPFLinkParameters `json:"ospf,omitempty" yaml:"ospf,omitempty" protobuf:"bytes,4,opt,name=ospf"`
	// ISIS defines the ISIS specific parameters on the link
	// +optional
	//ISIS *ISISLinkParameters `json:"isis,omitempty" yaml:"isis,omitempty" protobuf:"bytes,5,opt,name=isis"`
	// BGP defines the BGP specific parameters on the link
	// +optional
	//BGP *BGPLinkParameters `json:"bgp,omitempty" yaml:"bgp,omitempty" protobuf:"bytes,6,opt,name=bgp"`
}

// LinkStatus defines the observed state of Link
type LinkStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A link represents a physical/logical connection that enables communication and data transfer
// between 2 endpoints of a node.
type Link struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   LinkSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status LinkStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// LinkList contains a list of Links
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LinkList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Link `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	LinkKind     = reflect.TypeOf(Link{}).Name()
	LinkKindList = reflect.TypeOf(LinkList{}).Name()
)
