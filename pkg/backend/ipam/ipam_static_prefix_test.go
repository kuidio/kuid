package ipam

import (
	"context"
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
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: true}, // rib not ready
			},
		},
		"AggregateNormal": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "10.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "2000::/48", prefixType: aggregate, expectedError: false},
			},
		},
		"AggregateNestingIPv4": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", prefixType: aggregate, expectedError: true}, // nesting aggregate
				{claimType: staticPrefix, ip: "2000::/48", prefixType: aggregate, expectedError: false},
			},
		},
		"AggregateNestingIPv6": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "2000::/48", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "2000::/32", prefixType: aggregate, expectedError: true}, // nesting aggregate
			},
		},
		"NoAggregate": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", expectedError: true}, // no aggregate
			},
		},
		"NormalNesting": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: false},
			},
		},
		"Normal2Pool": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: true},
			},
		},
		"Normal2Network": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: true},
			},
		},
		"Network2Network": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: network, expectedError: true},
			},
		},
		"Network2Pool": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: network, expectedError: true},
			},
		},
		"Pool2Network": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: pool, expectedError: true},
			},
		},
		"Pool2Pool": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: pool, expectedError: true},
			},
		},
		"Pool2Aggregate": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: false},
			},
		},
		"Network2Aggregate": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: false},
			},
		},
		"InsertNormal2Normal": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", expectedError: false},
			},
		},
		"InsertNormal2Pool": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: pool, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", expectedError: false},
			},
		},
		"InsertNormal2Network": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", prefixType: network, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", expectedError: false},
			},
		},
		"InsertPool": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: pool, expectedError: true},
			},
		},
		"InsertNetwork": {
			niName: "a",
			prefixes: []testprefix{
				{claimType: staticPrefix, ip: "172.0.0.0/8", prefixType: aggregate, expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/16", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/27", expectedError: false},
				{claimType: staticPrefix, ip: "172.0.0.0/24", prefixType: network, expectedError: true},
			},
		},
		"InsertAggregate": {
			niName: "a",
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
			if tc.niName != "" {
				index := getNI(tc.niName)
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
						ipClaim, err = p.getIPClaimFromNetworkPrefix(tc.niName)
					} else {
						ipClaim, err = p.getStaticPrefixIPClaim(tc.niName)
					}
				case staticRange:
					ipClaim, err = p.getStaticRangeIPClaim(tc.niName)
				case staticAddress:
					ipClaim, err = p.getStaticAddressIPClaim(tc.niName)
				case dynamicPrefix:
					ipClaim, err = p.getDynamicPrefixIPClaim(tc.niName)
				case dynamicAddress:
					ipClaim, err = p.getDynamicAddressIPClaim(tc.niName)
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
