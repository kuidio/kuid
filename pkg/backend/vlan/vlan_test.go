package vlan

import (
	"context"
	"testing"

	"github.com/henderiw/store"
	vlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

func TestVLAN(t *testing.T) {
	tests := map[string]struct {
		index string
		vlans []testvlan
	}{
		"VLANMix": {
			index: "a",
			vlans: []testvlan{
				{claimType: vlanDynamic, name: "claim1", expectedError: false, expectedVLAN: ptr.To[uint32](2)},
				{claimType: vlanStatic, name: "claim2", id: 100, expectedError: false},
				{claimType: vlanStatic, name: "claim3", id: 4000, expectedError: false},
				{claimType: vlanRange, name: "claim4", vlanRange: "10-19", expectedError: false},
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
			cache, err := be.GetCache(ctx, store.KeyFromNSN(types.NamespacedName{Namespace: "dummy", Name: tc.index}))
			assert.NoError(t, err)

			for _, v := range tc.vlans {
				v := v
				var claim *vlanbev1alpha1.VLANClaim
				var err error

				switch v.claimType {
				case vlanStatic:
					claim, err = v.getStaticVLANClaim(tc.index)
				case vlanDynamic:
					claim, err = v.getDynamicVLANClaim(tc.index)
				case vlanRange:
					claim, err = v.getRangeVLANClaim(tc.index)
				case vlanSize:
					claim, err = v.getSizeVLANClaim(tc.index)
				}
				assert.NoError(t, err)
				if err != nil {
					return
				}

				err = be.Claim(ctx, claim)
				if v.expectedError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				switch v.claimType {
				case vlanStatic, vlanDynamic:
					if claim.Status.ID == nil {
						t.Errorf("expecting vlan status id got nil")
					} else {
						expectedVLAN := v.id
						if v.expectedVLAN != nil {
							expectedVLAN = *v.expectedVLAN
						}
						if *claim.Status.ID != expectedVLAN {
							t.Errorf("expecting vlan id got %d, want %d\n", *claim.Status.ID, expectedVLAN)
						}
					}
				case vlanRange:
					if claim.Status.Range == nil {
						t.Errorf("expecting vlan status id got nil")
					} else {
						expectedVLANRange := v.vlanRange
						if v.expectedVLANRange != nil {
							expectedVLANRange = *v.expectedVLANRange
						}
						if *claim.Status.Range != expectedVLANRange {
							t.Errorf("expecting vlan range got %s, want %s\n", *claim.Status.Range, expectedVLANRange)
						}
					}
				}
				//fmt.Println("entries after claim", v.name)
				printEntries(cache)
			}
		})
	}
}
