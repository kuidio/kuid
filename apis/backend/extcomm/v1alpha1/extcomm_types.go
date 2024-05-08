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

type ExtendedCommunityType string

const (
	ExtendedCommunityType_Invalid     ExtendedCommunityType = "invalid"
	ExtendedCommunityType_2byteAS     ExtendedCommunityType = "2byteAS"     // 0x00, 0x40 -> 4 byte local admin
	ExtendedCommunityType_IPv4Address ExtendedCommunityType = "ipv4Address" // 0x01, 0x41 -> 2 byte local admin
	ExtendedCommunityType_4byteAS     ExtendedCommunityType = "4byteAS"     // 0x02, 0x42 -> 2 byte local admin
	ExtendedCommunityType_Opaque      ExtendedCommunityType = "opaque"      // 0x03, 0x43 -> 6 byte local admin
)

func GetExtendedCommunityType(s string) ExtendedCommunityType {
	switch s {
	case string(ExtendedCommunityType_2byteAS):
		return ExtendedCommunityType_2byteAS
	case string(ExtendedCommunityType_IPv4Address):
		return ExtendedCommunityType_IPv4Address
	case string(ExtendedCommunityType_4byteAS):
		return ExtendedCommunityType_4byteAS
	case string(ExtendedCommunityType_Opaque):
		return ExtendedCommunityType_Opaque
	default:
		return ExtendedCommunityType_Invalid
	}
}

type ExtendedCommunitySubType string

const (
	ExtendedCommunitySubType_Invalid     ExtendedCommunitySubType = "invalid"
	ExtendedCommunitySubType_RouteTarget ExtendedCommunitySubType = "target" // 0x02
	ExtendedCommunitySubType_RouteOrigin ExtendedCommunitySubType = "origin" // 0x03
)

func GetExtendedCommunitySubType(s string) ExtendedCommunitySubType {
	switch s {
	case string(ExtendedCommunitySubType_RouteTarget):
		return ExtendedCommunitySubType_RouteTarget
	case string(ExtendedCommunitySubType_RouteOrigin):
		return ExtendedCommunitySubType_RouteOrigin
	default:
		return ExtendedCommunitySubType_Invalid
	}
}
