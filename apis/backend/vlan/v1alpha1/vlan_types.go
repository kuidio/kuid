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
	VLANClaimType_Invalid   VLANClaimType = "invalid"
	VLANClaimType_StaticID  VLANClaimType = "staticVLANID"
	VLANClaimType_DynamicID VLANClaimType = "dynamicVLANID"
	VLANClaimType_Range     VLANClaimType = "vlanRange"
	VLANClaimType_Size      VLANClaimType = "vlanSize"
)

func GetIPClaimTypeFromString(s string) VLANClaimType {
	switch s {
	case string(VLANClaimType_StaticID):
		return VLANClaimType_StaticID
	case string(VLANClaimType_DynamicID):
		return VLANClaimType_DynamicID
	case string(VLANClaimType_Range):
		return VLANClaimType_Range
	case string(VLANClaimType_Size):
		return VLANClaimType_Size
	default:
		return VLANClaimType_Invalid
	}
}
