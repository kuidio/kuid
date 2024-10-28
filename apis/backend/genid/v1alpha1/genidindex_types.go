/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "VLAN IS" BVLANIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"reflect"

	condv1alpha1 "github.com/kform-dev/choreo/apis/condition/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GENIDIndexSpec defines the desired state of GENIDIndex
type GENIDIndexSpec struct {
	// MinID defines the min ID the index supports
	// +optional
	MinID *uint64 `json:"minID,omitempty" yaml:"minID,omitempty" protobuf:"bytes,1,opt,name=minID"`
	// MaxID defines the max ID the index supports
	// +optional
	MaxID *uint64 `json:"maxID,omitempty" yaml:"maxID,omitempty" protobuf:"bytes,2,opt,name=maxID"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	commonv1alpha1.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,3,opt,name=userDefinedLabels"`
	// Type defines the type of the GENID
	// 16bit, 32bit, 48bit, 64bit
	Type string `json:"type,omitempty" yaml:"type,omitempty" protobuf:"bytes,4,opt,name=type"`
}

// GENIDIndexStatus defines the observed state of GENIDIndex
type GENIDIndexStatus struct {
	// MinID defines the min ID the index supports
	// +optional
	MinID *uint64 `json:"minID,omitempty" yaml:"minID,omitempty" protobuf:"bytes,1,opt,name=minID"`
	// MaxID defines the max ID the index supports
	// +optional
	MaxID *uint64 `json:"maxID,omitempty" yaml:"maxID,omitempty" protobuf:"bytes,2,opt,name=maxID"`
	// ConditionedStatus provides the status of the GENIDIndex using conditions
	// - a ready condition indicates the overall status of the resource
	condv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,3,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GENIDIndex is the Schema for the GENIDIndex API
type GENIDIndex struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   GENIDIndexSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status GENIDIndexStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// GENIDIndexList contains a list of GENIDIndexs
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type GENIDIndexList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []GENIDIndex `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	GENIDIndexKind     = reflect.TypeOf(GENIDIndex{}).Name()
	GENIDIndexListKind = reflect.TypeOf(GENIDIndexList{}).Name()
)
