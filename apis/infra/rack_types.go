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

// RackSpec defines the desired state of Rack
type RackSpec struct {
	// SiteID identifies the siteID this resource belongs to
	id.SiteID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=siteID"`
	// Location defines the location information where this resource is located
	// in lon/lat coordinates
	Location *Location `json:"location,omitempty" yaml:"location,omitempty" protobuf:"bytes,2,opt,name=location"`
	// The height of the rack, measured in units.
	Height string `json:"height,omitempty" yaml:"height,omitempty" protobuf:"bytes,3,opt,name=height"`
	// The canonical distance between the two vertical rails on a face. In inch
	Width string `json:"width,omitempty" yaml:"width,omitempty" protobuf:"bytes,4,opt,name=width"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined label
	common.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,5,opt,name=userDefinedLabels"`
}

// RackStatus defines the observed state of Rack
type RackStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A rack represents a physical equipment rack within your environment. Each rack is designed to accommodate
// the installation of devices and equipment.
type Rack struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   RackSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status RackStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// RackList contains a list of Racks
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RackList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Rack `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	RackKind     = reflect.TypeOf(Rack{}).Name()
	RackKindList = reflect.TypeOf(RackList{}).Name()
)
