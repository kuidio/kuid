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

// AdaptorSpec defines the desired state of Adaptor
type AdaptorSpec struct {
	// NodeGroupAdaptorID identifies the Adaptor identity this resource belongs to
	idv1alpha1.PartitionAdaptorID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=nodeGroupAdaptorID"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	commonv1alpha1.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,3,opt,name=userDefinedLabels"`
}

// AdaptorStatus defines the observed state of Adaptor
type AdaptorStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// An Adaptor represents a communication interface or connection point within a Node,
// facilitating network communication and data transfer between different components
// or systems within the environment. `Adaptors` serve as gateways for transmitting and
// receiving data, enabling seamless communication between Nodes.
// +k8s:openapi-gen=true
type Adaptor struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   AdaptorSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status AdaptorStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// AdaptorList contains a list of Adaptors
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AdaptorList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Adaptor `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	AdaptorKind     = reflect.TypeOf(Adaptor{}).Name()
	AdaptorKindList = reflect.TypeOf(AdaptorList{}).Name()
)
