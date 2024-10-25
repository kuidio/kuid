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

	condv1alpha1 "github.com/kform-dev/choreo/apis/condition/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	idv1alpha1 "github.com/kuidio/kuid/apis/id/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeSpec defines the desired state of Node
type NodeSpec struct {
	// NodeGroupNodeID identifies the nodeGroup identity this resource belongs to
	idv1alpha1.PartitionNodeID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=nodeID"`
	// Rack defines the rack in which the node is deployed
	// +optional
	Rack *string `json:"rack,omitempty" yaml:"rack,omitempty" protobuf:"bytes,2,opt,name=rack"`
	// relative position in the rack
	// +optional
	Position *string `json:"position,omitempty" yaml:"position,omitempty" protobuf:"bytes,3,opt,name=position"`
	// Location defines the location information where this resource is located
	// in lon/lat coordinates
	// +optional
	Location *Location `json:"location,omitempty" yaml:"location,omitempty" protobuf:"bytes,4,opt,name=location"`
	// Provider defines the provider implementing this resource.
	Provider string `json:"provider" yaml:"provider" protobuf:"bytes,5,opt,name=provider"`
	// PlatformType define the type of platform implementing the nodespec
	PlatformType string `json:"platformType" yaml:"platformType" protobuf:"bytes,6,opt,name=platformType"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	commonv1alpha1.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,7,opt,name=userDefinedLabels"`

	// TBD
	// Serial number
	// Node config
	// Initial config
	// IPAddress: IPv4 or IPv6
	// OOB IPAddress
}

// NodeStatus defines the observed state of Node
type NodeStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// System ID define the unique system id of the node
	// +optional
	SystemID *string `json:"systemID,omitempty" yaml:"systemID,omitempty" protobuf:"bytes,2,opt,name=systemID"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A Node represents a fundamental unit that implements compute, storage, and/or networking within your environment.
// Nodes can embody physical, virtual, or containerized entities, offering versatility in deployment options to suit
// diverse infrastructure requirements.
// Nodes are logically organized within racks and sites/regions, establishing a hierarchical structure for efficient
// resource management and organization. Additionally, Nodes are associated with nodeGroups, facilitating centralized
// management and control within defined administrative boundaries.
// Each Node is assigned a provider, representing the entity responsible for implementing the specifics of the Node.
// +k8s:openapi-gen=true
type Node struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   NodeSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status NodeStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// NodeList contains a list of Nodes
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NodeList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Node `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	NodeKind     = reflect.TypeOf(Node{}).Name()
	NodeKindList = reflect.TypeOf(NodeList{}).Name()
)
