/*
Copyright 2023 The Nephio Authors.

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

package backend

import "k8s.io/apimachinery/pkg/util/sets"

const (
	// system defined common
	KuidOwnerGroupKey     = "be.kuid.dev/owner-group"
	KuidOwnerVersionKey   = "be.kuid.dev/owner-version"
	KuidOwnerKindKey      = "be.kuid.dev/owner-kind"
	KuidOwnerNameKey      = "be.kuid.dev/owner-name"
	KuidOwnerNamespaceKey = "be.kuid.dev/owner-namespace"
	KuidClaimNameKey      = "be.kuid.dev/claim-name"
	KuidClaimTypeKey      = "be.kuid.dev/claim-type"
	// system defined ipam
	KuidIPAMIPPrefixTypeKey     = "ipam.be.kuid.dev/ipprefix-type"
	KuidIPAMClaimSummaryTypeKey = "ipam.be.kuid.dev/claim-summary-type"
	KuidIPAMddressFamilyKey     = "ipam.be.kuid.dev/address-family"
	KuidIPAMSubnetKey           = "ipam.be.kuid.dev/subnet" // this is the subnet in prefix annotation used for GW selection
	KuidIPAMDefaultGatewayKey   = "ipam.be.kuid.dev/default-gateway"
	KuidIPAMIndexKey            = "ipam.be.kuid.dev/index"

	// DNS used keys
	KuidINVNetworkKey  = "inv.kuid.dev/network"
	KuidINVEndpointKey = "inv.kuid.dev/endpoint-name"
	KuidINVHeadlessKey = "inv.kuid.dev/headless-name"
	KuidINVAliasKey    = "inv.kuid.dev/alias-name"
	KuidINVPortNS      = "port.kuid.dev"
	// Inventory
	KuidINVRegionKey           = "infra.be.kuid.dev/region"
	KuidINVAvailabilityZoneKey = "infra.be.kuid.dev/availability-zone"
	KuidINVSiteKey             = "infra.be.kuid.dev/site"
	KuidINVLocationKey         = "infra.be.kuid.dev/location"
	KuidINVRackKey             = "infra.be.kuid.dev/rack"
	KuidINVPositionKey         = "infra.be.kuid.dev/position" // position in the rack
	KuidINVNodeTypeKey         = "infra.be.kuid.dev/node-type"
	KuidINVNodeKey             = "infra.be.kuid.dev/node"
	KuidINVLinkKey             = "infra.be.kuid.dev/link"
	KuidINVLinkTypeKey         = "infra.be.kuid.dev/link-type"
	KuidINVPurpose             = "infra.be.kuid.dev/purpose"
	KuidINVExclude             = "infra.be.kuid.dev/exclude"
	// Network Inventory
	KuidINVNetworkDeviceType   = "network.infra.be.kuid.dev/device-type"     // edge, core or maybe p and pe
	KuidINVNetworkLinkBFD      = "link.network.infra.be.kuid.dev/bfd"        // enable or disable true or false
	KuidINVNetworkLinkBFDMinTX = "link.network.infra.be.kuid.dev/bfd-min-tx" //
	KuidINVNetworkLinkBFDMinRX = "link.network.infra.be.kuid.dev/bfd-min-rx" //
)

var BackendSystemKeys = sets.New[string](KuidOwnerGroupKey,
	KuidOwnerVersionKey,
	KuidOwnerKindKey,
	KuidOwnerNameKey,
	KuidOwnerNamespaceKey,
	KuidClaimNameKey,
)

var BackendIPAMSystemKeys = sets.New[string](KuidOwnerGroupKey,
	KuidIPAMIPPrefixTypeKey,
	KuidIPAMClaimSummaryTypeKey,
	KuidIPAMddressFamilyKey,
	KuidIPAMSubnetKey,
	KuidIPAMDefaultGatewayKey,
	KuidIPAMIndexKey,
)

var BackendSystemClaimKeys = sets.New[string](KuidOwnerGroupKey,
	KuidClaimTypeKey,
)
