package ipam

import (
	"testing"

	"github.com/kuidio/kuid/apis/backend/ipam"
)

func TestIPAMStaticRange(t *testing.T) {
	tests := map[string]prefixTest{
		"Normal": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/8"},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "10.0.0.0/24", expectedError: false},
				{claimType: staticRange, name: "range1", ip: "10.0.0.10-10.0.0.100", expectedError: false},
				{claimType: staticAddress, ip: "10.0.0.10/32", expectedError: false},
				{claimType: staticAddress, ip: "10.0.0.11/32", expectedError: false},
				{claimType: staticAddress, ip: "10.0.0.100/32", expectedError: false},
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
