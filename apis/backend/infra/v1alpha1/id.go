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
	"fmt"
	"strings"
)

type SiteID struct {
	// Region defines the region this sites belongs to
	Region string `json:"region" yaml:"region" protobuf:"bytes,1,opt,name=region"`
	// Site defines the site in which the node is deployed
	Site string `json:"site" yaml:"site" protobuf:"bytes,2,opt,name=site"`
}

type NodeGroupID struct {
	SiteID `json:",inline" yaml:",inline" protobuf:"bytes,4,opt,name=siteID"`
	// NodeGroup defines the node group the resource belongs to.
	NodeGroup string `json:"nodeGroup" yaml:"nodeGroup" protobuf:"bytes,3,opt,name=nodeGroup"`
}

type NodeID struct {
	NodeGroupID `json:",inline" yaml:",inline" protobuf:"bytes,5,opt,name=nodeGroupID"`
	// Node defines the node the resource belongs to.
	Node string `json:"node" yaml:"node" protobuf:"bytes,4,opt,name=node"`
}

type EndpointID struct {
	NodeID `json:",inline" yaml:",inline" protobuf:"bytes,6,opt,name=nodeID"`
	// Endpoint defines the name of the endpoint
	Endpoint string `json:"endpoint" yaml:"endpoint" protobuf:"bytes,5,opt,name=endpoint"`
}

func (r SiteID) KuidString() string {
	return fmt.Sprintf(
		"%s.%s",
		r.Region,
		r.Site,
	)
}

func (r NodeGroupID) KuidString() string {
	return fmt.Sprintf(
		"%s.%s",
		r.NodeGroup,
		r.SiteID.KuidString(),
	)
}

func (r NodeID) KuidString() string {
	return fmt.Sprintf(
		"%s.%s",
		r.NodeGroupID.KuidString(),
		r.Node,
	)
}

func String2NodeID(s string) *NodeID {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return nil
	}
	return &NodeID{
		Node: parts[3],
		NodeGroupID: NodeGroupID{
			NodeGroup: parts[0],
			SiteID: SiteID{
				Region: parts[1],
				Site:   parts[2],
			},
		},
	}
}

func (r EndpointID) KuidString() string {
	return fmt.Sprintf(
		"%s.%s",
		r.NodeID.KuidString(),
		r.Endpoint,
	)
}

func String2EndpointID(s string) *EndpointID {
	parts := strings.Split(s, ".")
	if len(parts) != 5 {
		return nil
	}
	return &EndpointID{
		Endpoint: parts[4],
		NodeID: NodeID{
			Node: parts[3],
			NodeGroupID: NodeGroupID{
				NodeGroup: parts[0],
				SiteID: SiteID{
					Region: parts[1],
					Site:   parts[2],
				},
			},
		},
	}
}
