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

type SiteID struct {
	// Region defines the region of the resource
	Region string `json:"region" yaml:"region" protobuf:"bytes,1,opt,name=region"`
	// Site defines the site of the resource
	Site string `json:"site" yaml:"site" protobuf:"bytes,2,opt,name=site"`
}

type NodeID struct {
	SiteID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=siteID"`
	// Node defines the name of the node
	Node string `json:"node" yaml:"node" protobuf:"bytes,2,opt,name=node"`
}

type PartitionNodeID struct {
	// Partition defines the partition this resource belongs to
	Partition string `json:"partition" yaml:"partition" protobuf:"bytes,1,opt,name=partition"`
	// SiteID define the siteid of the node
	SiteID `json:",inline" yaml:",inline" protobuf:"bytes,2,opt,name=siteID"`
	// Node defines the name of the node
	Node string `json:"node" yaml:"node" protobuf:"bytes,3,opt,name=node"`
}

type PortID struct {
	NodeID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=nodeID"`
	// ModuleBay defines the moduleBay reference id
	ModuleBay *int `json:"moduleBay,omitempty" yaml:"moduleBay,omitempty" protobuf:"bytes,2,opt,name=moduleBay"`
	// Module defines the module reference id
	Module *int `json:"module,omitempty" yaml:"module,omitempty" protobuf:"bytes,3,opt,name=module"`
	// Port defines the id of the port
	Port int `json:"port" yaml:"port" protobuf:"bytes,4,opt,name=port"`
}

type AdaptorID struct {
	NodeID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=nodeID"`
	// ModuleBay defines the moduleBay reference id
	ModuleBay *int `json:"moduleBay,omitempty" yaml:"moduleBay,omitempty" protobuf:"bytes,2,opt,name=moduleBay"`
	// Module defines the module reference id
	Module *int `json:"module,omitempty" yaml:"module,omitempty" protobuf:"bytes,3,opt,name=module"`
	// Port defines the id of the port
	Port int `json:"port" yaml:"port" protobuf:"bytes,4,opt,name=port"`
	// Adaptor defines the name of the adaptor
	Adaptor string `json:"adaptor" yaml:"adaptor" protobuf:"bytes,5,opt,name=adaptor"`
}

type EndpointID struct {
	NodeID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=nodeID"`
	// ModuleBay defines the moduleBay reference id
	ModuleBay *int `json:"moduleBay,omitempty" yaml:"moduleBay,omitempty" protobuf:"bytes,2,opt,name=moduleBay"`
	// Module defines the module reference id
	Module *int `json:"module,omitempty" yaml:"module,omitempty" protobuf:"bytes,3,opt,name=module"`
	// Port defines the id of the port
	Port int `json:"port" yaml:"port" protobuf:"bytes,4,opt,name=port"`
	// Adaptor defines the name of the adaptor
	Adaptor string `json:"adaptor" yaml:"adaptor" protobuf:"bytes,5,opt,name=adaptor"`
	// Endpoint defines the name of the endpoint
	Endpoint int `json:"endpoint" yaml:"endpoint" protobuf:"bytes,6,opt,name=endpoint"`
}

type PartitionEndpointID struct {
	// Partition defines the partition this resource belongs to
	Partition string `json:"partition" yaml:"partition" protobuf:"bytes,1,opt,name=partition"`

	NodeID `json:",inline" yaml:",inline" protobuf:"bytes,2,opt,name=nodeID"`
	// ModuleBay defines the moduleBay reference id
	ModuleBay *int `json:"moduleBay,omitempty" yaml:"moduleBay,omitempty" protobuf:"bytes,3,opt,name=moduleBay"`
	// Module defines the module reference id
	Module *int `json:"module,omitempty" yaml:"module,omitempty" protobuf:"bytes,4,opt,name=module"`
	// Port defines the id of the port
	Port int `json:"port" yaml:"port" protobuf:"bytes,5,opt,name=port"`
	// Adaptor defines the name of the adaptor
	Adaptor *string `json:"adaptor,omitempty" yaml:"adaptor,omitempty" protobuf:"bytes,6,opt,name=adaptor"`
	// Endpoint defines the name of the endpoint
	Endpoint int `json:"endpoint" yaml:"endpoint" protobuf:"bytes,7,opt,name=endpoint"`
	// Name is used to refer to internal names of the node
	Name *string `json:"name,omitempty" yaml:"name,omitempty" protobuf:"bytes,8,opt,name=name"`
}

type ClusterID struct {
	SiteID `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=siteID"`
	// Cluster defines the name of the cluster
	Cluster string `json:"cluster" yaml:"cluster" protobuf:"bytes,2,opt,name=cluster"`
}

type PartitionClusterID struct {
	// Partition defines the partition this resource belongs to
	Partition string `json:"partition" yaml:"partition" protobuf:"bytes,1,opt,name=partition"`
	// SiteID define the siteid of the node
	SiteID `json:",inline" yaml:",inline" protobuf:"bytes,2,opt,name=siteID"`
	// Cluster defines the name of the cluster
	Cluster string `json:"cluster" yaml:"cluster" protobuf:"bytes,3,opt,name=cluster"`
}

type PartitionAttachmentID struct {
	// Partition defines the partition this resource belongs to
	Partition string `json:"partition" yaml:"partition" protobuf:"bytes,1,opt,name=partition"`
	// SiteID define the siteid of the node
	SiteID `json:",inline" yaml:",inline" protobuf:"bytes,2,opt,name=siteID"`
	// Cluster defines the name of the cluster
	Cluster *string `json:"cluster,omitempty" yaml:"cluster,omitempty" protobuf:"bytes,3,opt,name=cluster"`
	// Node defines the name of the node
	Node *string `json:"node,omitempty" yaml:"node,omitempty" protobuf:"bytes,4,opt,name=node"`
	// Node defines the name of the nodeset
	NodeSet *string `json:"nodeset,omitempty" yaml:"nodeset,omitempty" protobuf:"bytes,5,opt,name=nodeset"`
	// Interface defines the name of the interface
	Interface string `json:"interface" yaml:"interface" protobuf:"bytes,1,opt,name=interface"`
}

/*
Endpoint.Connector
ModuleBay.Module.Endpoint.Connector

// Satelite
Endpoint.ModuleAdaptor.ModuleEndpoint.Connector
ModuleBay.Module.Endpoint.ModuleAdaptor.ModuleEndpoint.Connector
*/
