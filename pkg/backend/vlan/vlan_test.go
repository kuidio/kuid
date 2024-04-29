package vlan

import (
	"context"
	"testing"

	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	vlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
				{claimType: vlanDynamic, name: "claim1", expectedError: false, expectedVLAN: ptr.To[uint32](0)},
				{claimType: vlanStatic, name: "claim2", id: 100, expectedError: false},
				{claimType: vlanStatic, name: "claim3", id: 4000, expectedError: false},
				{claimType: vlanRange, name: "claim4", vlanRange: "10-19", expectedError: false},
				{claimType: vlanRange, name: "claim4", vlanRange: "11-19", expectedError: false}, // reclaim
				{claimType: vlanRange, name: "claim5", vlanRange: "5-10", expectedError: false},  // claim a new entry
				{claimType: vlanRange, name: "claim6", vlanRange: "19-100", expectedError: true}, // overlap
				{claimType: vlanStatic, name: "claim7", id: 12, expectedError: false},
				{claimType: vlanStatic, name: "claim7", id: 13, expectedError: false}, // reclaim an existing id
				{claimType: vlanDynamic, name: "claim8", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedVLAN: ptr.To[uint32](11)}, // a dynamic claim from a range
				{claimType: vlanDynamic, name: "claim9", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedVLAN: ptr.To[uint32](12)}, // a dynamic claim from a range that was claimed before
				{claimType: vlanDynamic, name: "claim10", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedVLAN: ptr.To[uint32](14)}, 
				{claimType: vlanDynamic, name: "claim10", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"}, // update
				}, expectedError: false, expectedVLAN: ptr.To[uint32](14)}, 
				{claimType: vlanRange, name: "claim4", vlanRange: "11-19", expectedError: false}, // update
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
				}
				assert.NoError(t, err)
				if err != nil {
					return
				}

				err = be.Claim(ctx, claim)
				if v.expectedError {
					assert.Error(t, err)
					continue
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
