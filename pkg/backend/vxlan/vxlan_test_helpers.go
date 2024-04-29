package vxlan

import (
	"context"
	"fmt"

	"github.com/henderiw/idxtable/pkg/table32"
	"github.com/henderiw/store"
	vxlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

type testvxlan struct {
	name               string
	claimType          vxlanbev1alpha1.VXLANClaimType
	id                 uint32
	vxlanRange         string
	labels             map[string]string
	selector           *metav1.LabelSelector
	expectedError      bool
	expectedVXLAN      *uint32
	expectedVXLANRange *string
}

// alias
const (
	namespace    = "dummy"
	vxlanStatic  = vxlanbev1alpha1.VXLANClaimType_StaticID
	vxlanDynamic = vxlanbev1alpha1.VXLANClaimType_DynamicID
	vxlanRange   = vxlanbev1alpha1.VXLANClaimType_Range
)

func getIndex(index string) *vxlanbev1alpha1.VXLANIndex {
	return vxlanbev1alpha1.BuildVXLANIndex(
		metav1.ObjectMeta{Namespace: namespace, Name: index},
		nil,
		nil,
	)
}

func (r testvxlan) getDynamicVXLANClaim(index string) (*vxlanbev1alpha1.VXLANClaim, error) {
	claim := vxlanbev1alpha1.BuildVXLANClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&vxlanbev1alpha1.VXLANClaimSpec{
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

func (r testvxlan) getStaticVXLANClaim(index string) (*vxlanbev1alpha1.VXLANClaim, error) {
	claim := vxlanbev1alpha1.BuildVXLANClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&vxlanbev1alpha1.VXLANClaimSpec{
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

func (r testvxlan) getRangeVXLANClaim(index string) (*vxlanbev1alpha1.VXLANClaim, error) {
	claim := vxlanbev1alpha1.BuildVXLANClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: r.name},
		&vxlanbev1alpha1.VXLANClaimSpec{
			Index: index,
			Range: ptr.To[string](r.vxlanRange),
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
