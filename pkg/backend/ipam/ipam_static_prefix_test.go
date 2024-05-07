package ipam

import (
	"context"
	"testing"

	"github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestIPAMStaticPrefix(t *testing.T) {
	tests := map[string]struct {
		index   string
		prefixes []testprefix
	}{
		"NotReady": {
			index: "",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: true}, // rib not ready
			},
		},
		"AggregateNormal": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "10.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "2000::/48", prefixType: aggregate, expectedError: false},
			},
		},
		"AggregateNestingIPv4": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", prefixType: aggregate, expectedError: true}, // nesting aggregate
				{claimType: staticPrefix, ip: "2000::/48", prefixType: aggregate, expectedError: false},
			},
		},
		"AggregateNestingIPv6": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "2000::/48", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "2000::/32", prefixType: aggregate, expectedError: true}, // nesting aggregate
			},
		},
		"NoAggregate": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", expectedError: true}, // no aggregate
			},
		},
		"NormalNesting": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: false},
			},
		},
		"Normal2Pool": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: true},
			},
		},
		"Normal2Network": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: true},
			},
		},
		"Network2Network": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: network, expectedError: true},
			},
		},
		"Network2Pool": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: network, expectedError: true},
			},
		},
		"Pool2Network": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: pool, expectedError: true},
			},
		},
		"Pool2Pool": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: pool, expectedError: true},
			},
		},
		"Pool2Aggregate": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
			},
		},
		"Network2Aggregate": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
			},
		},
		"InsertNormal2Normal": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", expectedError: false},
			},
		},
		"InsertNormal2Pool": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: pool, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", expectedError: false},
			},
		},
		"InsertNormal2Network": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", expectedError: false},
			},
		},
		"InsertPool": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: true},
			},
		},
		"InsertNetwork": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: true},
			},
		},
		"InsertAggregate": {
			index: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: aggregate, expectedError: true},
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			be := New(nil)
			ctx := context.Background()
			if tc.index != "" {
				index := getIndex(tc.index)
				err := be.CreateIndex(ctx, index)
				assert.NoError(t, err)
			}

			for _, p := range tc.prefixes {
				p := p
				var ipClaim *v1alpha1.IPClaim
				var err error

				switch p.claimType {
				case staticPrefix:
					if p.prefixType != nil && *p.prefixType == *aggregate {
						ipClaim, err = p.getIPClaimFromNetworkPrefix(tc.index)
					} else {
						ipClaim, err = p.getStaticPrefixIPClaim(tc.index)
					}
				case staticRange:
					ipClaim, err = p.getStaticRangeIPClaim(tc.index)
				case staticAddress:
					ipClaim, err = p.getStaticAddressIPClaim(tc.index)
				case dynamicPrefix:
					ipClaim, err = p.getDynamicPrefixIPClaim(tc.index)
				case dynamicAddress:
					ipClaim, err = p.getDynamicAddressIPClaim(tc.index)
				}
				assert.NoError(t, err)

				err = be.Claim(ctx, ipClaim)
				if p.expectedError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					switch p.claimType {
					case staticPrefix, dynamicPrefix:
						if ipClaim.Status.Prefix == nil {
							t.Errorf("expecting prefix status got nil")
						} else {
							expectedIP := p.ip
							if p.expectedIP != "" {
								expectedIP = p.expectedIP
							}
							if *ipClaim.Status.Prefix != expectedIP {
								t.Errorf("expecting prefix got %s, want %s\n", *ipClaim.Status.Prefix, expectedIP)
							}
						}
					case staticAddress, dynamicAddress:
						if ipClaim.Status.Address == nil {
							t.Errorf("expecting address status got nil")
						} else {
							expectedIP := p.ip
							if p.expectedIP != "" {
								expectedIP = p.expectedIP
							}
							if *ipClaim.Status.Address != expectedIP {
								t.Errorf("expecting address got %s, want %s\n", *ipClaim.Status.Address, expectedIP)
							}
						}
						if ipClaim.Status.DefaultGateway == nil {
							if p.expectedDG != "" {
								t.Errorf("expecting defaultGateway %s got nil", p.expectedDG)
							}
						} else {
							if p.expectedDG == "" {
								t.Errorf("unexpected defaultGateway got %s", *ipClaim.Status.DefaultGateway)
							}
							if *ipClaim.Status.DefaultGateway != p.expectedDG {
								t.Errorf("expecting defaultGateway got %s, want %s\n", *ipClaim.Status.DefaultGateway, p.expectedDG)
							}
						}
					case staticRange:
						if ipClaim.Status.Range == nil {
							t.Errorf("expecting range status got nil")
						} else {
							if *ipClaim.Status.Range != p.ip {
								t.Errorf("expecting prefix got %s, want %s\n", *ipClaim.Status.Range, p.ip)
							}
						}
					}
				}
			}
		})
	}
}
