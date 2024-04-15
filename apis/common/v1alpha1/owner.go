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

import "sigs.k8s.io/controller-runtime/pkg/client"

// +k8s:openapi-gen=true
type OwnerReference struct {
	Group     string `json:"group" yaml:"group" protobuf:"bytes,1,opt,name=group"`
	Version   string `json:"version" yaml:"version" protobuf:"bytes,2,opt,name=version"`
	Kind      string `json:"kind" yaml:"kind" protobuf:"bytes,3,opt,name=kind"`
	Namespace string `json:"namespace" yaml:"namespace" protobuf:"bytes,4,opt,name=namespace"`
	Name      string `json:"name" yaml:"name" protobuf:"bytes,5,opt,name=name"`
}

func GetOwnerReference(obj client.Object) *OwnerReference {
	return &OwnerReference{
		Group:     obj.GetObjectKind().GroupVersionKind().Group,
		Version:   obj.GetObjectKind().GroupVersionKind().Version,
		Kind:      obj.GetObjectKind().GroupVersionKind().Kind,
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}
