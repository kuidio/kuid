package as

import (
	"context"
	"testing"

	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	asbev1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

func Test(t *testing.T) {
	tests := map[string]struct {
		index string
		ctxs  []testCtx
	}{
		"Mix": {
			index: "a",
			ctxs: []testCtx{
				{claimType: dynamicClaim, name: "claim1", expectedError: false, expectedID: ptr.To[uint32](0)},
				{claimType: staticClaim, name: "claim2", id: 100, expectedError: false},
				{claimType: staticClaim, name: "claim3", id: 4000, expectedError: false},
				{claimType: rangeClaim, name: "claim4", asRange: "10-19", expectedError: false},
				{claimType: rangeClaim, name: "claim4", asRange: "11-19", expectedError: false}, // reclaim
				{claimType: rangeClaim, name: "claim5", asRange: "5-10", expectedError: false},  // claim a new entry
				{claimType: rangeClaim, name: "claim6", asRange: "19-100", expectedError: true}, // overlap
				{claimType: staticClaim, name: "claim7", id: 12, expectedError: false},
				{claimType: staticClaim, name: "claim7", id: 13, expectedError: false}, // reclaim an existing id
				{claimType: dynamicClaim, name: "claim8", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedID: ptr.To[uint32](11)}, // a dynamic claim from a range
				{claimType: dynamicClaim, name: "claim9", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedID: ptr.To[uint32](12)}, // a dynamic claim from a range that was claimed before
				{claimType: dynamicClaim, name: "claim10", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedID: ptr.To[uint32](14)},
				{claimType: dynamicClaim, name: "claim10", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"}, // update
				}, expectedError: false, expectedID: ptr.To[uint32](14)},
				{claimType: rangeClaim, name: "claim4", asRange: "11-19", expectedError: false}, // update
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

			for _, v := range tc.ctxs {
				v := v
				var claim *asbev1alpha1.ASClaim
				var err error

				switch v.claimType {
				case staticClaim:
					claim, err = v.getStaticClaim(tc.index)
				case dynamicClaim:
					claim, err = v.getDynamicClaim(tc.index)
				case rangeClaim:
					claim, err = v.getRangeClaim(tc.index)
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
				case staticClaim, dynamicClaim:
					if claim.Status.ID == nil {
						t.Errorf("expecting status id got nil")
					} else {
						expectedID := v.id
						if v.expectedID != nil {
							expectedID = *v.expectedID
						}
						if *claim.Status.ID != expectedID {
							t.Errorf("expecting id got %d, want %d\n", *claim.Status.ID, expectedID)
						}
					}
				case rangeClaim:
					if claim.Status.Range == nil {
						t.Errorf("expecting status id got nil")
					} else {
						expectedRange := v.asRange
						if v.expectedRange != nil {
							expectedRange = *v.expectedRange
						}
						if *claim.Status.Range != expectedRange {
							t.Errorf("expecting range got %s, want %s\n", *claim.Status.Range, expectedRange)
						}
					}
				}
				//fmt.Println("entries after claim", v.name)
				printEntries(cache)
			}
		})
	}
}
