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
	"context"
	"fmt"
	"testing"

	"github.com/henderiw/iputil"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
)

func TestGetIPEntry(t *testing.T) {
	tests := map[string]struct {
		name      string
		namespace string
		prefix    string
		labels    map[string]string
	}{
		"Normal": {
			name:      "a",
			namespace: "b",
			prefix:    "10.0.0.0/24",
			labels: map[string]string{
				backend.KuidOwnerVersionKey:   "v1alpha1",
				backend.KuidOwnerKindKey:      "ipam.res.kuid.dev",
				backend.KuidOwnerNameKey:      "vpc1.10.0.0.0-24",
				backend.KuidOwnerNamespaceKey: "default",
				backend.KuidClaimNameKey:      "10.0.0.0-24",
				backend.KuidIPAMSubnetKey:     "10.0.0.0-24",
				backend.KuidIPAMIndexKey:      "1",
				backend.KuidIPAMGatewayKey:    "true",
				"x":                           "y",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			prefix, err := iputil.New(tc.prefix)
			if err != nil {
				assert.Error(t, err)
			}

			ipEntry := GetIPEntry(context.TODO(), store.KeyFromNSN(types.NamespacedName{
				Name:      tc.name,
				Namespace: tc.namespace,
			}), prefix.Prefix, tc.labels)

			fmt.Println("ipEntry", ipEntry)
		})
	}
}
