package ipam

import (
	"context"
	"fmt"
	"testing"

	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIPAMDynamicAddress(t *testing.T) {
	tests := map[string]struct {
		niName   string
		prefixes []testprefix
	}{
		"FromAggregate": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "10.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{name: "addrClaim1", claimInfo: "address", expectedError: true},
			},
		},
		/*
			"FromNetwork": {
				niName: "a",
				prefixes: []testprefix{
					{prefix: "10.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
					{prefix: "10.0.0.0/24", claimType: "network", claimInfo: "prefix", expectedError: false},
					{name: "addrClaim1", claimInfo: "address", expectedError: false},
					{name: "addrClaim1", claimInfo: "address", expectedError: false},
					{name: "addrClaim2", claimInfo: "address", expectedError: false},
					{name: "addrClaim3", claimInfo: "address", expectedError: false},
					{name: "addrClaim4", claimInfo: "address", expectedError: false},
					{name: "addrClaim5", claimInfo: "address", expectedError: false},
				},
			},
			"FromPool": {
				niName: "a",
				prefixes: []testprefix{
					{prefix: "10.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
					{prefix: "10.0.0.0/24", claimType: "pool", claimInfo: "prefix", expectedError: false},
					{name: "addrClaim1", claimInfo: "address", expectedError: false},
					{name: "addrClaim1", claimInfo: "address", expectedError: false},
					{name: "addrClaim2", claimInfo: "address", expectedError: false},
					{name: "addrClaim3", claimInfo: "address", expectedError: false},
					{name: "addrClaim4", claimInfo: "address", expectedError: false},
					{name: "addrClaim5", claimInfo: "address", expectedError: false},
				},
			},
			"FromOther": {
				niName: "a",
				prefixes: []testprefix{
					{prefix: "10.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
					{prefix: "10.0.0.0/24", claimType: "other", claimInfo: "prefix", expectedError: false},
					{name: "addrClaim1", claimInfo: "address", expectedError: true},
				},
			},
		*/
		"FromRange": {
			niName: "a",
			prefixes: []testprefix{
				{prefix: "10.0.0.0/8", claimType: "aggregate", claimInfo: "prefix", expectedError: false},
				{prefix: "10.0.0.0/24", claimType: "other", claimInfo: "prefix", expectedError: false},
				{name: "range1", prefix: "10.0.0.10-10.0.0.100", claimInfo: "range", expectedError: false},
				{name: "addrClaim1", claimInfo: "address", expectedError: false, selector: &v1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "range1"},
				}},
				{name: "addrClaim1", claimInfo: "address", expectedError: false, selector: &v1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "range1"},
				}},
				{name: "addrClaim2", claimInfo: "address", expectedError: false, selector: &v1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "range1"},
				}},
			},
		},
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
					fmt.Println("prefix status", *ipClaim.Status.Prefix)
				}
				if ipClaim.Status.Address != nil {
					fmt.Println("address status", *ipClaim.Status.Address)
				}
			}
		})
	}
}
