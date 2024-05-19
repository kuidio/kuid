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

// NodeGroupSpec defines the desired state of NodeGroup
// E.g. A nodeGroup can be a NodeGroup
// E.g. A nodeGroup can be a topology like a DC fabric (frontend and backend could be a different nodeGroup)
// A Node Group is a global unique identifier within the system e.g. representing a topology, a NodeGroup or
// another set of elements that are managed together by a single entity
type NodeGroupSpec struct {
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	commonv1alpha1.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=userDefinedLabels"`
}

// NodeGroupStatus defines the observed state of NodeGroup
type NodeGroupStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A NodeGrouo represents a logical grouping of infrastructure resources managed by a single
// administrative entity or organization. NodeGroups serve as administrative boundaries within the environment,
// providing a structured framework for organizing and managing resources based on administrative ownership
// or responsibility. E.g. A NodeGroup on one hand, can be used to represent a topology that spans multiple
// sites and regions, but a NodeGroup can also be used to group all nodes of a NodeGroup together.
// +k8s:openapi-gen=true
type NodeGroup struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   NodeGroupSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status NodeGroupStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// NodeGroupList contains a list of NodeGroups
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type NodeGroupList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []NodeGroup `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	NodeGroupKind = reflect.TypeOf(NodeGroup{}).Name()
	NodeGroupKindList = reflect.TypeOf(NodeGroupList{}).Name()
)
