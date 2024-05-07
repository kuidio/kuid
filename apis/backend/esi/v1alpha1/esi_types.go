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

type ESIClaimType string

const (
	ClaimType_Invalid   ESIClaimType = "invalid"
	ClaimType_StaticID  ESIClaimType = "staticESIID"
	ClaimType_DynamicID ESIClaimType = "dynamicESIID"
	ClaimType_Range     ESIClaimType = "esiRange"
)

func GetClaimTypeFromString(s string) ESIClaimType {
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
	ESIIndexReservedMinName = "rangeReservedMin"
	ESIIndexReservedMaxName = "rangeReservedMax"
)
