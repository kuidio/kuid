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

package genid

import (
	"fmt"
)

const GENIDID_Min = 0

var GENIDID_MaxBits = map[GENIDType]int{
	GENIDType_Invalid: 0,
	GENIDType_16bit:   16,
	GENIDType_32bit:   32,
	GENIDType_48bit:   48,
	GENIDType_64bit:   64,
}

var GENIDID_MaxValue = map[GENIDType]uint64{
	GENIDType_Invalid: 0,
	GENIDType_16bit:   1<<GENIDID_MaxBits[GENIDType_16bit] - 1,
	GENIDType_32bit:   1<<GENIDID_MaxBits[GENIDType_32bit] - 1,
	GENIDType_48bit:   1<<GENIDID_MaxBits[GENIDType_48bit] - 1,
	GENIDType_64bit:   1<<GENIDID_MaxBits[GENIDType_64bit] - 1,
}

func validateGENIDID(genidType GENIDType, id uint64) error {
	if id < GENIDID_Min {
		return fmt.Errorf("invalid id, got %d", id)
	}
	if id > GENIDID_MaxValue[genidType] {
		return fmt.Errorf("invalid id, got %d", id)
	}
	return nil
}
