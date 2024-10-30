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
		"AggregateNoParents": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "10.0.0.0/8", prefixType: aggregate, expectedError: true},
				{claimType: staticPrefix, ip: "2000::/48", prefixType: aggregate, expectedError: true},
			},
		},
		"AggregateNoParentIPv6": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/16", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "2000::/48", prefixType: aggregate, expectedError: true},
			},
		},
		"AggregateNormalIPv6": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "2000::/32", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "2000::/48", prefixType: aggregate, expectedError: false},
			},
		},
		"Normal2Pool": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/16", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: pool, expectedError: false},
			},
		},
		"NestingNetwork": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/16", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: network, expectedError: true},
			},
		},
		"Network2Pool": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: network, expectedError: true},
			},
		},
		"Pool2Network": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: pool, expectedError: true},
			},
		},
		"Pool2Pool": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: pool, expectedError: true},
			},
		},
		"Pool2Aggregate": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
			},
		},
		"Network2Aggregate": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
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
