package ipam

import (
	"testing"

	"github.com/kuidio/kuid/apis/backend/ipam"
)

func TestIPAMStaticPrefix(t *testing.T) {
	tests := map[string]prefixTest{
		/*
			"NotReady": {
				index: "",
				prefixes: []testprefix{
					{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: true}, // rib not ready
				},
			},
		*/
		"NoParents": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8"},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "10.0.0.0/8", expectedError: true},
				{claimType: staticPrefix, ip: "2000::/48", expectedError: true},
			},
		},
		"NoParentIPv6": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8"},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "2000::/48", expectedError: true},
			},
		},
		"ParentIPv6": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "2000::/32"},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "2000::/48", expectedError: false},
			},
		},
		"Nesting": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8"},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: false},
			},
		},
		"NestingNetwork": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8"},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: network, expectedError: true},
			},
		},
		"Other2Network": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8"},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: true},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := prefixTestRun(name, tc); err != nil {
				t.Errorf("test %s failed err: %v", name, err)
			}
		})
	}
}
