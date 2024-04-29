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

type ASClaimType string

const (
	ASClaimType_Invalid   ASClaimType = "invalid"
	ASClaimType_StaticID  ASClaimType = "staticASID"
	ASClaimType_DynamicID ASClaimType = "dynamicASID"
	ASClaimType_Range     ASClaimType = "asRange"
)

func GetIPClaimTypeFromString(s string) ASClaimType {
	switch s {
	case string(ASClaimType_StaticID):
		return ASClaimType_StaticID
	case string(ASClaimType_DynamicID):
		return ASClaimType_DynamicID
	case string(ASClaimType_Range):
		return ASClaimType_Range
	default:
		return ASClaimType_Invalid
	}
}
