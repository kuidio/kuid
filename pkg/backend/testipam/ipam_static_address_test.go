package ipam

import (
	"testing"

	"github.com/kuidio/kuid/apis/backend/ipam"
	"k8s.io/utils/ptr"
)

func TestIPAMStaticAddress(t *testing.T) {
	tests := map[string]prefixTest{
		"NoParent": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticAddress, ip: "172.0.0.0/32", expectedError: true},
			},
		},
		"Address_AggregateParent": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticAddress, ip: "172.0.0.0/32", expectedError: false},
			},
		},
		"Address_NetworkParent_OwnerClash": { // since netwotk prefixes get expanded the address is clashing
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticAddress, ip: "172.0.0.0/32", expectedError: true},
			},
		},
		"Address_NetworkParent": { // 32 or /128 not possible in network Addresses
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticAddress, ip: "172.0.0.1/32", expectedError: true},
			},
		},
		"Address_First_PoolParent": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticAddress, ip: "172.0.0.0/32", expectedError: false},
			},
		},
		"Address_PoolParent": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticAddress, ip: "172.0.0.1/32", expectedError: false},
			},
		},
		"Address_First_OtherParent": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", expectedError: true},
				{claimType: staticAddress, ip: "172.0.0.0/32", expectedError: false},
			},
		},
		"PrefixAddress_AggregateParent": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticAddress, ip: "172.0.0.1/24", expectedError: true},
			},
		},
		"PrefixAddress_NetworkParent_OwnerClash": { // since netwotk prefixes get expanded the address is clashing
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: ptr.To(ipam.IPPrefixType_Aggregate)},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticAddress, ip: "172.0.0.0/32", expectedError: true},
			},
		},
		"PrefixAddress_NetworkParent": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticAddress, ip: "172.0.0.1/24", expectedError: false},
			},
		},
		"PrefixAddress_NetworkParentWithAddress": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.1/24", prefixType: network, expectedError: false},
				{claimType: staticAddress, ip: "172.0.0.2/24", expectedError: false},
			},
		},
		"PrefixAddress_First_PoolParent": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticAddress, ip: "172.0.0.0/32", expectedError: false},
			},
		},
		"PrefixAddress_PoolParent": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticAddress, ip: "172.0.0.1/32", expectedError: false},
			},
		},
		"PrefixAddress_First_OtherParent": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8", PrefixType: aggregate},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: aggregate, expectedError: false},
				{claimType: staticAddress, ip: "172.0.0.1/24", prefixType: network, expectedError: true},
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
