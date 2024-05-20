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

// NodeSetSetSpec defines the desired state of NodeSet
type NodeSetSpec struct {
	// NodeGroupName identifies the nodeGroup this resource belongs to
	// E.g. a NodeSet in a NodeSet belongs to a nodeGroup where the name of the nodeGroup is the NodeSet
	// E.g. a Virtual Node, belongs to a nodeGroup where the name of the nodeGroup represents the topology this node is deployed in
	NodeGroup string `json:"nodeGroup" yaml:"nodeGroup" protobuf:"bytes,1,opt,name=nodeGroup"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	commonv1alpha1.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,2,opt,name=userDefinedLabels"`
}

// NodeSetStatus defines the observed state of NodeSet
type NodeSetStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A NodeSet represents a set of nodes.
// E.g. it can be used to model a set of nodes in a NodeSet that share the same
// charecteristics wrt, Numa, interfaces, etc.
// Another usage of NodeSet is the representation of a virtual Node that consists of multiple nodes.
// +k8s:openapi-gen=true
type NodeSet struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   NodeSetSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status NodeSetStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// NodeSetList contains a list of NodeSets
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NodeSetList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []NodeSet `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	NodeSetKind     = reflect.TypeOf(NodeSet{}).Name()
	NodeSetKindList = reflect.TypeOf(NodeSetList{}).Name()
)
