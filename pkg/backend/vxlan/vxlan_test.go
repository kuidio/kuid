package vxlan

import (
	"context"
	"testing"

	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	vxlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

func TestVXLAN(t *testing.T) {
	tests := map[string]struct {
		index string
		vxlans []testvxlan
	}{
		"VXLANMix": {
			index: "a",
			vxlans: []testvxlan{
				{claimType: vxlanDynamic, name: "claim1", expectedError: false, expectedVXLAN: ptr.To[uint32](0)},
				{claimType: vxlanStatic, name: "claim2", id: 100, expectedError: false},
				{claimType: vxlanStatic, name: "claim3", id: 4000, expectedError: false},
				{claimType: vxlanRange, name: "claim4", vxlanRange: "10-19", expectedError: false},
				{claimType: vxlanRange, name: "claim4", vxlanRange: "11-19", expectedError: false}, // reclaim
				{claimType: vxlanRange, name: "claim5", vxlanRange: "5-10", expectedError: false},  // claim a new entry
				{claimType: vxlanRange, name: "claim6", vxlanRange: "19-100", expectedError: true}, // overlap
				{claimType: vxlanStatic, name: "claim7", id: 12, expectedError: false},
				{claimType: vxlanStatic, name: "claim7", id: 13, expectedError: false}, // reclaim an existing id
				{claimType: vxlanDynamic, name: "claim8", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedVXLAN: ptr.To[uint32](11)}, // a dynamic claim from a range
				{claimType: vxlanDynamic, name: "claim9", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedVXLAN: ptr.To[uint32](12)}, // a dynamic claim from a range that was claimed before
				{claimType: vxlanDynamic, name: "claim10", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedVXLAN: ptr.To[uint32](14)},
				{claimType: vxlanDynamic, name: "claim10", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"}, // update
				}, expectedError: false, expectedVXLAN: ptr.To[uint32](14)},
				{claimType: vxlanRange, name: "claim4", vxlanRange: "11-19", expectedError: false}, // update
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

			for _, v := range tc.vxlans {
				v := v
				var claim *vxlanbev1alpha1.VXLANClaim
				var err error

				switch v.claimType {
				case vxlanStatic:
					claim, err = v.getStaticVXLANClaim(tc.index)
				case vxlanDynamic:
					claim, err = v.getDynamicVXLANClaim(tc.index)
				case vxlanRange:
					claim, err = v.getRangeVXLANClaim(tc.index)
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
				case vxlanStatic, vxlanDynamic:
					if claim.Status.ID == nil {
						t.Errorf("expecting vxlan status id got nil")
					} else {
						expectedVXLAN := v.id
						if v.expectedVXLAN != nil {
							expectedVXLAN = *v.expectedVXLAN
						}
						if *claim.Status.ID != expectedVXLAN {
							t.Errorf("expecting vxlan id got %d, want %d\n", *claim.Status.ID, expectedVXLAN)
						}
					}
				case vxlanRange:
					if claim.Status.Range == nil {
						t.Errorf("expecting vxlan status id got nil")
					} else {
						expectedVXLANRange := v.vxlanRange
						if v.expectedVXLANRange != nil {
							expectedVXLANRange = *v.expectedVXLANRange
						}
						if *claim.Status.Range != expectedVXLANRange {
							t.Errorf("expecting vxlan range got %s, want %s\n", *claim.Status.Range, expectedVXLANRange)
						}
					}
				}
				//fmt.Println("entries after claim", v.name)
				printEntries(cache)
			}
		})
	}
}
