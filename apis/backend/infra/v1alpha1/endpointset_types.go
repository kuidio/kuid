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

// EndpointSetSpec defines the desired state of EndpointSet
// An EndpointSet can be a LAG (single Homed) or ESI (multiHomed). The EndpointSet
// can only belong to a single NodeGroup
type EndpointSetSpec struct {
	// Endpoints defines the Endpoints that are part of the EndpointSet
	// Min 1, Max 16
	Endpoints []*EndpointID `json:"endpoints" yaml:"endpoints" protobuf:"bytes,1,opt,name=endpoints"`
	// Lacp defines if the lag enabled LACP
	// +optional
	Lacp *bool `json:"lacp,omitempty" yaml:"lacp,omitempty" protobuf:"bytes,2,opt,name=lacp"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	commonv1alpha1.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,3,opt,name=userDefinedLabels"`
}

// EndpointSetStatus defines the observed state of EndpointSet
type EndpointSetStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// ESI defines the ethernet segment identifier of the logical link
	// if set this is a multi-homed logical endpoint
	// the ESI is a global unique identifier within the administrative domain
	// +optional
	ESI *uint32 `json:"esi,omitempty" yaml:"esi,omitempty" protobuf:"bytes,2,opt,name=esi"`
	// LagId defines the lag id for the logical single-homed or multi-homed
	// endpoint
	// +optional
	LagId *uint32 `json:"lagID,omitempty" yaml:"lagID,omitempty" protobuf:"bytes,3,opt,name=lagID"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A EndpointSet represents a set of endpoints that belong together within a nodeGroup.
// E.g. it can be used to model a logical Link Aggregation group within
// a node or it can be used to represent a logical multi-homing construction
// between a set of nodes belonging to a single nodeGroup.
// +k8s:openapi-gen=true
type EndpointSet struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   EndpointSetSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status EndpointSetStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// EndpointSetList contains a list of EndpointSets
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EndpointSetList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []EndpointSet `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	EndpointSetKind = reflect.TypeOf(EndpointSet{}).Name()
	EndpointSetKindList = reflect.TypeOf(EndpointSetList{}).Name()
)
