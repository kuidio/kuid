package as

import (
	"context"
	"fmt"

	"github.com/henderiw/idxtable/pkg/table32"
	"github.com/henderiw/store"
	asbev1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

type testCtx struct {
	name          string
	claimType     asbev1alpha1.ASClaimType
	id            uint32
	asRange       string
	labels        map[string]string
	selector      *metav1.LabelSelector
	expectedError bool
	expectedID    *uint32
	expectedRange *string
}

// alias
const (
	namespace    = "dummy"
	staticClaim  = asbev1alpha1.ASClaimType_StaticID
	dynamicClaim = asbev1alpha1.ASClaimType_DynamicID
	rangeClaim   = asbev1alpha1.ASClaimType_Range
)

func getIndex(index string) *asbev1alpha1.ASIndex {
	return asbev1alpha1.BuildASIndex(
		metav1.ObjectMeta{Namespace: namespace, Name: index},
		nil,
		nil,
	)
}

func (r testCtx) getDynamicClaim(index string) (*asbev1alpha1.ASClaim, error) {
	claim := asbev1alpha1.BuildASClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&asbev1alpha1.ASClaimSpec{
			Index: index,
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
				Selector:          r.selector,
			},
		},
		nil,
	)
	fielErrList := claim.ValidateSyntax() // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return claim, nil
}

func (r testCtx) getStaticClaim(index string) (*asbev1alpha1.ASClaim, error) {
	claim := asbev1alpha1.BuildASClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&asbev1alpha1.ASClaimSpec{
			Index: index,
			ID:    ptr.To[uint32](r.id),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := claim.ValidateSyntax() // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return claim, nil
}

func (r testCtx) getRangeClaim(index string) (*asbev1alpha1.ASClaim, error) {
	claim := asbev1alpha1.BuildASClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: r.name},
		&asbev1alpha1.ASClaimSpec{
			Index: index,
			Range: ptr.To[string](r.asRange),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := claim.ValidateSyntax() // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return claim, nil
}

func printEntries(cachectx *CacheContext) {
	fmt.Println("---------")
	ctx := context.Background()
	for _, entry := range cachectx.tree.GetAll() {
		entry := entry
		fmt.Println("entry", entry.ID().String())
	}
	cachectx.ranges.List(ctx, func(ctx context.Context, k store.Key, t *table32.Table32) {
		fmt.Println("range", k.Name)
		for _, entry := range t.GetAll() {
			entry := entry
			fmt.Println("entry", entry.ID().String())
		}
	})
}
