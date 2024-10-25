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
	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/kuidio/kuid/apis/backend/vlan"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +k8s:deepcopy-gen=false
var _ resource.Object = &VLANIndex{}
var _ resource.ObjectList = &VLANIndexList{}
var _ resource.MultiVersionObject = &VLANIndex{}

func (VLANIndex) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: vlan.VLANIndexPlural,
	}
}

// IsStorageVersion returns true -- Config is used as the internal version.
// IsStorageVersion implements resource.Object
func (VLANIndex) IsStorageVersion() bool {
	return false
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object
func (VLANIndex) NamespaceScoped() bool {
	return true
}

// GetObjectMeta implements resource.Object
// GetObjectMeta implements resource.Object
func (r *VLANIndex) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// New return an empty resource
// New implements resource.Object
func (VLANIndex) New() runtime.Object {
	return &VLANIndex{}
}

// NewList return an empty resourceList
// NewList implements resource.Object
func (VLANIndex) NewList() runtime.Object {
	return &VLANIndexList{}
}

// GetListMeta returns the ListMeta
// GetListMeta implements resource.ObjectList
func (r *VLANIndexList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

// RegisterConversions registers the conversions.
// RegisterConversions implements resource.MultiVersionObject
func (VLANIndex) RegisterConversions() func(s *runtime.Scheme) error {
	return RegisterConversions
}
