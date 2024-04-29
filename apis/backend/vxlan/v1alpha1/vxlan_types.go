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

type VXLANClaimType string

const (
	VXLANClaimType_Invalid   VXLANClaimType = "invalid"
	VXLANClaimType_StaticID  VXLANClaimType = "staticVXLANID"
	VXLANClaimType_DynamicID VXLANClaimType = "dynamicVXLANID"
	VXLANClaimType_Range     VXLANClaimType = "vxlanRange"
)

func GetClaimTypeFromString(s string) VXLANClaimType {
	switch s {
	case string(VXLANClaimType_StaticID):
		return VXLANClaimType_StaticID
	case string(VXLANClaimType_DynamicID):
		return VXLANClaimType_DynamicID
	case string(VXLANClaimType_Range):
		return VXLANClaimType_Range
	default:
		return VXLANClaimType_Invalid
	}
}
