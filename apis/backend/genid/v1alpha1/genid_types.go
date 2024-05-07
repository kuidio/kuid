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

type GENIDType string

const (
	GENIDType_Invalid GENIDType = "invalid"
	GENIDType_16bit   GENIDType = "16bit"
	GENIDType_32bit   GENIDType = "32bit"
	GENIDType_48bit   GENIDType = "48bit"
	GENIDType_64bit   GENIDType = "64bit"
)

func GetGenIDType(s string) GENIDType {
	switch s {
	case string(GENIDType_16bit):
		return GENIDType_16bit
	case string(GENIDType_32bit):
		return GENIDType_32bit
	case string(GENIDType_48bit):
		return GENIDType_48bit
	case string(GENIDType_64bit):
		return GENIDType_64bit
	default:
		return GENIDType_Invalid
	}
}
