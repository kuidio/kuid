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

type VLANClaimType string

const (
	ClaimType_Invalid   VLANClaimType = "invalid"
	ClaimType_StaticID  VLANClaimType = "staticVLANID"
	ClaimType_DynamicID VLANClaimType = "dynamicVLANID"
	ClaimType_Range     VLANClaimType = "vlanRange"
)

func GetClaimTypeFromString(s string) VLANClaimType {
	switch s {
	case string(ClaimType_StaticID):
		return ClaimType_StaticID
	case string(ClaimType_DynamicID):
		return ClaimType_DynamicID
	case string(ClaimType_Range):
		return ClaimType_Range
	default:
		return ClaimType_Invalid
	}
}

const (
	VLANIndexReservedMinName = "rangeReservedMin"
	VLANIndexReservedMaxName = "rangeReservedMax"
)
