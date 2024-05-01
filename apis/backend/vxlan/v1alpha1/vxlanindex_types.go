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

// VXLANIndexSpec defines the desired state of VXLANIndex
type VXLANIndexSpec struct {
	// MinID defines the min ID the index supports
	// +optional
	MinID *uint32 `json:"minID,omitempty" yaml:"minID,omitempty" protobuf:"bytes,1,opt,name=minID"`
	// MaxID defines the max ID the index supports
	// +optional
	MaxID *uint32 `json:"maxID,omitempty" yaml:"maxID,omitempty" protobuf:"bytes,2,opt,name=maxID"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	commonv1alpha1.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,3,opt,name=userDefinedLabels"`
}

// VXLANIndexStatus defines the observed state of VXLANIndex
type VXLANIndexStatus struct {
	// MinID defines the min ID the index is configured with
	// +optional
	MinID *uint32 `json:"minID,omitempty" yaml:"minID,omitempty" protobuf:"bytes,1,opt,name=minID"`
	// MaxID defines the max ID the index is configured with
	// +optional
	MaxID *uint32 `json:"maxID,omitempty" yaml:"maxID,omitempty" protobuf:"bytes,2,opt,name=maxID"`
	// ConditionedStatus provides the status of the VLANIndex using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,3,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VXLANIndex is the Schema for the VXLANIndex API
//
// +k8s:openapi-gen=true
type VXLANIndex struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   VXLANIndexSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status VXLANIndexStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// VXLANIndexList contains a list of VXLANIndexs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VXLANIndexList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []VXLANIndex `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	VXLANIndexKind = reflect.TypeOf(VXLANIndex{}).Name()
)
