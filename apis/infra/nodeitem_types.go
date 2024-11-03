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

// NodeItemSpec defines the desired state of NodeItem
type NodeItemSpec struct {
	// NodeID identifies the node identity this resource belongs to
	id.PartitionNodeID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=nodeID"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	common.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,2,opt,name=userDefinedLabels"`
}

// NodeItemStatus defines the observed state of NodeItem
type NodeItemStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:skipversion
// A NodeItem represents a specific hardware component or accessory associated with a Node.
// NodeItems represent a wide range of hardware elements, e.g Fan(s), PowerUnit(s), CPU(s),
// and other peripheral devices essential for the operation of the Node.
// NodeItem is used to represent the modular components of a node.
type NodeItem struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   NodeItemSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status NodeItemStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// NodeItemList contains a list of NodeItems
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NodeItemList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []NodeItem `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	NodeItemKind     = reflect.TypeOf(NodeItem{}).Name()
	NodeItemKindList = reflect.TypeOf(NodeItemList{}).Name()
)
