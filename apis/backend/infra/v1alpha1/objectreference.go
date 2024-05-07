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

type ObjectReference struct {
	// APIVersion of the target resources
	APIVersion string `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty" protobuf:"bytes,1,opt,name=apiVersion"`

	// Kind of the target resources
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" protobuf:"bytes,2,opt,name=kind"`

	// Name of the target resource
	// +optional
	Name *string `yaml:"name" json:"name" protobuf:"bytes,3,opt,name=name"`

	// Note: Namespace is not allowed; the namespace
	// must match the namespace of the PackageVariantSet resource
}
