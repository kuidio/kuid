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
		/*
		"Ranges": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "10.0.0.0/24", prefixType: aggregate, expectedError: false},
				{claimType: staticRange, name: "range1", ip: "10.0.0.10-10.0.0.19", expectedError: false},
				{claimType: staticAddress, ip: "10.0.0.10/32", expectedError: false},
				{claimType: staticAddress, ip: "10.0.0.11/32", expectedError: false},
				{claimType: staticAddress, ip: "10.0.0.19/32", expectedError: false},
				{claimType: staticRange, name: "range2", ip: "10.0.0.20-10.0.0.29", expectedError: false},
				{claimType: staticAddress, ip: "10.0.0.20/32", expectedError: false},
				{claimType: staticAddress, ip: "10.0.0.21/32", expectedError: false},
				{claimType: staticAddress, ip: "10.0.0.29/32", expectedError: false},
			},
		},
		"Overlap": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "10.0.0.0/24", prefixType: aggregate, expectedError: false},
				{claimType: staticRange, name: "range1", ip: "10.0.0.10-10.0.0.20", expectedError: false},
				{claimType: staticRange, name: "range2", ip: "10.0.0.20-10.0.0.29", expectedError: true},
			},
		},
		"Range2Network": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "10.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticRange, name: "range1", ip: "10.0.0.10-10.0.0.100", expectedError: false},
			},
		},
		"Range2Pool": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "10.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticRange, name: "range1", ip: "10.0.0.10-10.0.0.100", expectedError: false},
			},
		},
		*/
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
