package ipam

import (
	"context"
	"fmt"
	"testing"

	"github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestIPAMStaticPrefix(t *testing.T) {
	tests := map[string]struct {
		niName   string
		prefixes []testprefix
	}{
		"NotReady": {
			niName: "",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: true}, // rib not ready
			},
		},
		/*
		"AggregateNormal": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "10.0.0.8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "2000::/48", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
			},
		},
		"AggregateNestingIPv4": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "aggregate", claimInfo: "prefix", expectedError: true}, // nesting aggregate
				{prefix: "2000::/48", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
			},
		},
		"AggregateNestingIPv6": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "2000::/48", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "2000::/32", claimType: "aggregate", claimInfo: "prefix", expectedError: true}, // nesting aggregate
			},
		},
		"NoAggregate": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "normal", claimInfo: "prefix", expectedError: true}, // no aggregate
			},
		},
		"NormalNesting": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "normal", claimInfo: "prefix", expectedError: false},
			},
		},
		"Normal2Pool": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "pool", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "normal", claimInfo: "prefix", expectedError: true},
			},
		},
		"Normal2Network": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "network", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "normal", claimInfo: "prefix", expectedError: true},
			},
		},
		"Network2Network": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "network", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "network", claimInfo: "prefix", expectedError: true},
			},
		},
		"Network2Pool": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "pool", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "network", claimInfo: "prefix", expectedError: true},
			},
		},
		"Pool2Network": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "network", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "pool", claimInfo: "prefix", expectedError: true},
			},
		},
		"Pool2Pool": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "pool", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "pool", claimInfo: "prefix", expectedError: true},
			},
		},
		"Pool2Aggregate": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "pool", claimInfo: "prefix", expectedError: false},
			},
		},
		"Network2Aggregate": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "network", claimInfo: "prefix", expectedError: false},
			},
		},
		"InsertNormal2Normal": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "normal", claimInfo: "prefix", expectedError: false},
			},
		},
		"InsertNormal2Pool": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "pool", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "normal", claimInfo: "prefix", expectedError: false},
			},
		},
		"InsertNormal2Network": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "network", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "normal", claimInfo: "prefix", expectedError: false},
			},
		},
		"InsertPool": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "pool", claimInfo: "prefix", expectedError: true},
			},
		},
		"InsertNetwork": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "network", claimInfo: "prefix", expectedError: true},
			},
		},
		"InsertAggregate": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "172.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/16", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/27", claimType: "normal", claimInfo: "prefix", expectedError: false},
				{prefix: "172.0.0.0/24", claimType: "aggregate", claimInfo: "prefix", expectedError: true},
			},
			
		},
		*/
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
