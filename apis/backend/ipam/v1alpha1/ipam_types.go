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
	"k8s.io/utils/ptr"
)

type IPClaimType string

const (
	IPClaimType_Invalid        IPClaimType = "invalid"
	IPClaimType_StaticAddress  IPClaimType = "staticAddress"
	IPClaimType_StaticPrefix   IPClaimType = "staticPrefix"
	IPClaimType_StaticRange    IPClaimType = "staticRange"
	IPClaimType_DynamicAddress IPClaimType = "dynamicAddress"
	IPClaimType_DynamicPrefix  IPClaimType = "dynamicPrefix"
)

func GetIPClaimTypeFromString(s string) IPClaimType {
	switch s {
	case string(IPClaimType_StaticAddress):
		return IPClaimType_StaticAddress
	case string(IPClaimType_StaticPrefix):
		return IPClaimType_StaticPrefix
	case string(IPClaimType_StaticRange):
		return IPClaimType_StaticRange
	case string(IPClaimType_DynamicAddress):
		return IPClaimType_DynamicAddress
	case string(IPClaimType_DynamicPrefix):
		return IPClaimType_DynamicPrefix
	default:
		return IPClaimType_Invalid
	}
}

type IPClaimSummaryType string

const (
	IPClaimSummaryType_Prefix  IPClaimSummaryType = "prefix"
	IPClaimSummaryType_Address IPClaimSummaryType = "address"
	IPClaimSummaryType_Range   IPClaimSummaryType = "range"
	IPClaimSummaryType_Invalid IPClaimSummaryType = "invalid"
)

func GetIPClaimSummaryTypeFromString(s string) IPClaimSummaryType {
	switch s {
	case string(IPClaimSummaryType_Prefix):
		return IPClaimSummaryType_Prefix
	case string(IPClaimSummaryType_Address):
		return IPClaimSummaryType_Address
	case string(IPClaimSummaryType_Range):
		return IPClaimSummaryType_Range
	default:
		return IPClaimSummaryType_Invalid
	}
}

type IPPrefixType string

const (
	IPPrefixType_Invalid   IPPrefixType = "invalid"
	IPPrefixType_Other     IPPrefixType = "other"
	IPPrefixType_Pool      IPPrefixType = "pool"
	IPPrefixType_Network   IPPrefixType = "network"
	IPPrefixType_Aggregate IPPrefixType = "aggregate"
)

func GetIPPrefixTypeFromString(s string) *IPPrefixType {
	switch s {
	case string(IPPrefixType_Pool):
		return ptr.To[IPPrefixType](IPPrefixType_Pool)
	case string(IPPrefixType_Network):
		return ptr.To[IPPrefixType](IPPrefixType_Network)
	case string(IPPrefixType_Aggregate):
		return ptr.To[IPPrefixType](IPPrefixType_Aggregate)
	default:
		return nil
	}
}
