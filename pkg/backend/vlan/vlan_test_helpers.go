package vlan

import (
	"fmt"

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
	vlanSize          uint32
	labels            map[string]string
	expectedError     bool
	expectedVLAN      *uint32
	expectedVLANRange *string
	expectedVLANSize  *uint32
}

// alias
const (
	vlanStatic  = vlanbev1alpha1.VLANClaimType_StaticID
	vlanDynamic = vlanbev1alpha1.VLANClaimType_DynamicID
	vlanRange   = vlanbev1alpha1.VLANClaimType_Range
	vlanSize    = vlanbev1alpha1.VLANClaimType_Size
)

func getIndex(index string) *vlanbev1alpha1.VLANIndex {
	return vlanbev1alpha1.BuildVLANIndex(
		metav1.ObjectMeta{Namespace: "dummy", Name: index},
		nil,
		nil,
	)
}

func (r testvlan) getDynamicVLANClaim(index string) (*vlanbev1alpha1.VLANClaim, error) {
	claim := vlanbev1alpha1.BuildVLANClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: r.name},
		&vlanbev1alpha1.VLANClaimSpec{
			Index: index,
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

func (r testvlan) getStaticVLANClaim(index string) (*vlanbev1alpha1.VLANClaim, error) {
	claim := vlanbev1alpha1.BuildVLANClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: r.name},
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

func (r testvlan) getSizeVLANClaim(index string) (*vlanbev1alpha1.VLANClaim, error) {
	claim := vlanbev1alpha1.BuildVLANClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: r.name},
		&vlanbev1alpha1.VLANClaimSpec{
			Index:    index,
			VLANSize: ptr.To[uint32](r.vlanSize),
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
	//for id, l := range cachectx.table.GetAll() {
	//	id := id
	//	//fmt.Println("entry", id, l)
	//}
}
