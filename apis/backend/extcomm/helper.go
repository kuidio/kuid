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

package extcomm

const EXTCOMMID_Min = 0

var EXTCOMMID_MaxBits = map[ExtendedCommunityType]int{
	ExtendedCommunityType_Invalid:     0,
	ExtendedCommunityType_2byteAS:     32,
	ExtendedCommunityType_4byteAS:     16,
	ExtendedCommunityType_IPv4Address: 16,
	ExtendedCommunityType_Opaque:      48,
}

var EXTCOMMID_MaxValue = map[ExtendedCommunityType]uint64{
	ExtendedCommunityType_Invalid:     1<<EXTCOMMID_MaxBits[ExtendedCommunityType_Invalid] - 1,
	ExtendedCommunityType_2byteAS:     1<<EXTCOMMID_MaxBits[ExtendedCommunityType_2byteAS] - 1,
	ExtendedCommunityType_4byteAS:     1<<EXTCOMMID_MaxBits[ExtendedCommunityType_4byteAS] - 1,
	ExtendedCommunityType_IPv4Address: 1<<EXTCOMMID_MaxBits[ExtendedCommunityType_IPv4Address] - 1,
	ExtendedCommunityType_Opaque:      1<<EXTCOMMID_MaxBits[ExtendedCommunityType_Opaque] - 1,
}