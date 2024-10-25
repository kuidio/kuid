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

// Claims can be expressed in 3 ways
// OnwerReference, Finalizer or Status with reference to the claim onwer -> finalizer seem the best option for this

// EndpointSpec defines the desired state of Endpoint
type EndpointSpec struct {
	// NodeGroupEndpointID identifies the endpoint identity this resource belongs to
	id.PartitionEndpointID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=nodeGroupEndpointID"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	common.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,3,opt,name=userDefinedLabels"`
	// (Gbps)
	Speed *string `json:"speed,omitempty" yaml:"speed,omitempty" protobuf:"bytes,4,opt,name=speed"`
	// VLANTagging defines if VLAN tagging is enabled or disabled on the interface
	VLANTagging bool `json:"vlanTagging,omitempty" yaml:"vlanTagging,omitempty" protobuf:"bytes,5,opt,name=vlanTagging"`
}

// EndpointStatus defines the observed state of Endpoint
type EndpointStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// An Endpoint represents a communication interface or connection point within a Node,
// facilitating network communication and data transfer between different components
// or systems within the environment. `Endpoints` serve as gateways for transmitting and
// receiving data, enabling seamless communication between Nodes.
type Endpoint struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   EndpointSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status EndpointStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// EndpointList contains a list of Endpoints
type EndpointList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Endpoint `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	EndpointKind     = reflect.TypeOf(Endpoint{}).Name()
	EndpointKindList = reflect.TypeOf(EndpointList{}).Name()
)
