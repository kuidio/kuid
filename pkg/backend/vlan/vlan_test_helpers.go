package vlan

import (
	"context"
	"fmt"

	"github.com/henderiw/idxtable/pkg/table12"
	"github.com/henderiw/store"
	vlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

type testvlan struct {
	name              string
	claimType         vlanbev1alpha1.VLANClaimType
	id                uint32
	vlanRange         string
	labels            map[string]string
	selector          *metav1.LabelSelector
	expectedError     bool
	expectedVLAN      *uint32
	expectedVLANRange *string
}

// alias
const (
	namespace   = "dummy"
	vlanStatic  = vlanbev1alpha1.VLANClaimType_StaticID
	vlanDynamic = vlanbev1alpha1.VLANClaimType_DynamicID
	vlanRange   = vlanbev1alpha1.VLANClaimType_Range
)

func getIndex(index string) *vlanbev1alpha1.VLANIndex {
	return vlanbev1alpha1.BuildVLANIndex(
		metav1.ObjectMeta{Namespace: namespace, Name: index},
		nil,
		nil,
	)
}

func (r testvlan) getDynamicVLANClaim(index string) (*vlanbev1alpha1.VLANClaim, error) {
	claim := vlanbev1alpha1.BuildVLANClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&vlanbev1alpha1.VLANClaimSpec{
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

func (r testvlan) getStaticVLANClaim(index string) (*vlanbev1alpha1.VLANClaim, error) {
	claim := vlanbev1alpha1.BuildVLANClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&vlanbev1alpha1.VLANClaimSpec{
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

func (r testvlan) getRangeVLANClaim(index string) (*vlanbev1alpha1.VLANClaim, error) {
	claim := vlanbev1alpha1.BuildVLANClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: r.name},
		&vlanbev1alpha1.VLANClaimSpec{
			Index: index,
			Range: ptr.To[string](r.vlanRange),
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
	cachectx.ranges.List(ctx, func(ctx context.Context, k store.Key, t *table12.Table12) {
		fmt.Println("range", k.Name)
		for _, entry := range t.GetAll() {
			entry := entry
			fmt.Println("entry", entry.ID().String())
		}
	})
}
