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

// +kubebuilder:object:generate=true
// +groupName=infra.kuid.dev
package infra

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	GroupName = "infra.kuid.dev"
	Version   = runtime.APIVersionInternal
)

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: Version}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Cluster{},
		&ClusterList{},
		&Endpoint{},
		&EndpointList{},
		&EndpointSet{},
		&EndpointSetList{},
		&Link{},
		&LinkList{},
		&LinkSet{},
		&LinkSetList{},
		&Module{},
		&ModuleList{},
		&ModuleBay{},
		&ModuleBayList{},
		&Node{},
		&NodeList{},
		&NodeItem{},
		&NodeItemList{},
		&NodeSet{},
		&NodeSetList{},
		&Partition{},
		&PartitionList{},
		&Rack{},
		&RackList{},
		&Region{},
		&RegionList{},
		&Site{},
		&SiteList{},
	)
	return nil
}
