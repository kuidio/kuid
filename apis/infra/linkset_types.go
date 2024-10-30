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

// LinkSetSpec defines the desired state of LinkSet
type LinkSetSpec struct {
	// Endpoints define the endpoint identifiers of the LinkSet
	Endpoints []*id.PartitionEndpointID `json:"endpoints" yaml:"endpoints" protobuf:"bytes,1,opt,name=endpoints"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	common.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,2,opt,name=userDefinedLabels"`
}

// LinkSetStatus defines the observed state of LinkSet
type LinkSetStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// ESI defines the ethernet segment identifier of the logical link
	// if set this is a multi-homed linkset
	// the ESI is a global unique identifier within the administrative domain/topology
	ESI *uint32 `json:"esi,omitempty" yaml:"esi,omitempty" protobuf:"varint,2,opt,name=esi"`
	// LagId defines the lag id for the logical single-homed or multi-homed
	// endpoint
	LagId *uint32 `json:"lagId,omitempty" yaml:"lagId,omitempty" protobuf:"varint,3,opt,name=lagId"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A linkSet represents a set of links that belong together within a node group or accross nodeGroups.
// E.g. it can be used to model a logical Link Aggregation group between 2 nodes or
// it can be used to represent a logical multi-homing construction between a set of nodes
// belonging to 1 or multiple nodeGroups/Topologies.
type LinkSet struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   LinkSetSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status LinkSetStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// LinkSetList contains a list of LinkSets
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type LinkSetList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []LinkSet `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	LinkSetKind     = reflect.TypeOf(LinkSet{}).Name()
	LinkSetKindList = reflect.TypeOf(LinkSetList{}).Name()
)
