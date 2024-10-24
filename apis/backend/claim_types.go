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

// +k8s:openapi-gen=true
// ClaimType define the type of the claim
type ClaimType string

const (
	ClaimType_Invalid   ClaimType = "invalid"
	ClaimType_StaticID  ClaimType = "staticID"
	ClaimType_DynamicID ClaimType = "dynamicID"
	ClaimType_Range     ClaimType = "range"
)

func GetClaimTypeFromString(s string) ClaimType {
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
	IndexReservedMinName = "rangeReservedMin"
	IndexReservedMaxName = "rangeReservedMax"
)
