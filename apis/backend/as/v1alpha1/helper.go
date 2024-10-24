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

import (
	"fmt"
	"strconv"
)

const ASID_Min = 0
const ASID_Max = 4294967295

func validateASID(id int) error {
	if id < ASID_Min {
		return fmt.Errorf("invalid id, got %d", id)
	}
	if id > ASID_Max {
		return fmt.Errorf("invalid id, got %d", id)
	}
	return nil
}

func getASDot(asn uint32) string {
	if asn > 65536 {
		a := asn / 65536
		b := asn - (a * 65536)
		return fmt.Sprintf("%d.%d", a, b)
	}
	return strconv.Itoa(int(asn))
}
