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

// ModuleBaySpec defines the desired state of ModuleBay
type ModuleBaySpec struct {
	// NodeID identifies the node identity this resource belongs to
	NodeID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=nodeID"`
	// Position defines the position in the node the moduleBay is deployed
	Position string `json:"psoition" yaml:"psoition" protobuf:"bytes,2,opt,name=psoition"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	commonv1alpha1.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,3,opt,name=userDefinedLabels"`
}

// ModuleBayStatus defines the observed state of ModuleBay
type ModuleBayStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A ModuleBay serves as a modular slot or enclosure within a Node, designed to accommodate additional modules.
// ModuleBays provide a flexible and scalable approach to extending the capabilities of Nodes,
// allowing users to customize and enhance their infrastructure deployments according to specific requirements.
// +k8s:openapi-gen=true
type ModuleBay struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ModuleBaySpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status ModuleBayStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ModuleBayList contains a list of ModuleBays
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ModuleBayList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []ModuleBay `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	ModuleBayKind = reflect.TypeOf(ModuleBay{}).Name()
)
