package ipam

import (
	"context"
	"fmt"
	"testing"

	"github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestIPAMStaticAddress(t *testing.T) {
	tests := map[string]struct {
		niName   string
		prefixes []testprefix
	}{/*
		"NoParent": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/32", claimType: "", claimInfo: "address", expectedError: true},
			},
		},
		"Address_AggregateParent": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/32", claimType: "", claimInfo: "address", expectedError: true},
			},
		},
		"Address_NetworkParent_OwnerClash": { // since netwotk prefixes get expanded the address is clashing
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "network", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/32", claimType: "", claimInfo: "address", expectedError: true},
			},
		},
		"Address_NetworkParent": { // 32 or /128 not possible in network Addresses
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "network", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.1/32", claimType: "", claimInfo: "address", expectedError: true},
			},
		},
		"Address_First_PoolParent": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "pool", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/32", claimType: "", claimInfo: "address", expectedError: false},
			},
		},
		"Address_PoolParent": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "pool", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.1/32", claimType: "", claimInfo: "address", expectedError: false},
			},
		},
		"Address_First_OtherParent": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "other", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/32", claimType: "", claimInfo: "address", expectedError: false},
			},
		},
		"PrefixAddress_AggregateParent": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.1/24", claimType: "", claimInfo: "address", expectedError: true},
			},
		},
		"PrefixAddress_NetworkParent_OwnerClash": { // since netwotk prefixes get expanded the address is clashing
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "network", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "", claimInfo: "address", expectedError: true},
			},
		},
		"PrefixAddress_NetworkParent": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "network", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.1/24", claimType: "", claimInfo: "address", expectedError: false},
			},
		},
		"PrefixAddress_NetworkParentWithAddress": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.1/24", claimType: "network", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.2/24", claimType: "", claimInfo: "address", expectedError: false},
			},
		},
		"PrefixAddress_First_PoolParent": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "pool", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "", claimInfo: "address", expectedError: true},
			},
		},
		"PrefixAddress_PoolParent": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "pool", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.1/24", claimType: "", claimInfo: "address", expectedError: true},
			},
		},
		"PrefixAddress_First_OtherParent": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "other", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "", claimInfo: "address", expectedError: true},
			},
		},*/
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			be := New(nil)
			ctx := context.Background()
			if tc.niName != "" {
				index := getIPIndex(tc.niName)
				err := be.CreateIndex(ctx, index)
				assert.NoError(t, err)
			}

			for _, p := range tc.prefixes {
				p := p
				var ipClaim *v1alpha1.IPClaim
				var err error
				if p.claimType == "aggregate" {
					ipClaim, err = getIPClaimFromNetworkPrefix(tc.niName, p.prefix, p.labels)

				} else {
					switch p.claimInfo {
					case "prefix":
						if p.prefix != "" {
							ipClaim, err = getIPClaimFromPrefix(tc.niName, p.prefix, v1alpha1.GetIPClaimTypeFromString(p.claimType), p.labels)
						} else {
							ipClaim, err = getIPClaimFromDynamicPrefix(p.name, tc.niName, v1alpha1.GetIPClaimTypeFromString(p.claimType), p.prefixLength, p.labels, p.selector)
						}
					case "address":
						if p.prefix != "" {
							ipClaim, err = getIPClaimFromAddress(tc.niName, p.prefix, p.labels)
						} else {
							ipClaim, err = getIPClaimFromDynamicAddress(p.name, tc.niName, p.labels, p.selector)
						}
					case "range":
						ipClaim, err = getIPClaimFromRange(p.name, tc.niName, p.prefix, p.labels)
					}
				}
				assert.NoError(t, err)

				err = be.Claim(ctx, ipClaim)
				if p.expectedError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				if ipClaim.Status.Prefix != nil {
					fmt.Println("status", *ipClaim.Status.Prefix)
				}
			}
		})
	}
}
