package ipam

import (
	"testing"

	"github.com/kuidio/kuid/apis/backend/ipam"
)

func TestIPAMDynamicPrefix(t *testing.T) {
	tests := map[string]prefixTest{
		"Normal": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/8"},
			},
			prefixes: []testprefix{
				{claimType: dynamicPrefix, name: "prefix1", prefixLength: 24, expectedError: false, expectedIP: "10.0.0.0/24"},
				{claimType: dynamicPrefix, name: "prefix2", prefixLength: 24, expectedError: false, expectedIP: "10.0.1.0/24"},
				{claimType: dynamicPrefix, name: "prefix3", prefixLength: 24, prefixType: network, expectedError: false, expectedIP: "10.0.2.0/24"},
				{claimType: dynamicPrefix, name: "prefix4", prefixLength: 24, prefixType: network, expectedError: false, expectedIP: "10.0.3.0/24"},
				{claimType: dynamicPrefix, name: "prefix5", prefixLength: 24, expectedError: false, expectedIP: "10.0.4.0/24"},
			},
		},
		"NoAvailable": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/24"},
			},
			prefixes: []testprefix{
				{claimType: dynamicPrefix, name: "prefix1", prefixLength: 24, expectedError: true},
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			if err := prefixTestRun(name, tc); err != nil {
				t.Errorf("test %s failed err: %v", name, err)
			}
		})
	}
}
