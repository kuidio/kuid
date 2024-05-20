package vxlan

import (
	"fmt"

	"github.com/kuidio/kuid/apis/backend"
	vxlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

type testCtx struct {
	name          string
	claimType     backend.ClaimType
	id            uint64
	tRange        string
	labels        map[string]string
	selector      *metav1.LabelSelector
	expectedError bool
	expectedID    *uint64
	expectedRange *string
}

// alias
const (
	namespace    = "dummy"
	staticClaim  = backend.ClaimType_StaticID
	dynamicClaim = backend.ClaimType_DynamicID
	rangeClaim   = backend.ClaimType_Range
)

func getIndex(index, _ string) (*vxlanbev1alpha1.VXLANIndex, error) {
	idx := vxlanbev1alpha1.BuildVXLANIndex(
		metav1.ObjectMeta{Namespace: namespace, Name: index},
		nil,
		nil,
	)

	fieldErrs := idx.ValidateSyntax("")
	if len(fieldErrs) != 0 {
		return nil, fmt.Errorf("syntax errors %v", fieldErrs)
	}
	return idx, nil
}

func (r testCtx) getDynamicClaim(index, testType string) (*vxlanbev1alpha1.VXLANClaim, error) {
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
	fielErrList := claim.ValidateSyntax(testType) // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return claim, nil
}

func (r testCtx) getStaticClaim(index, testType string) (*vxlanbev1alpha1.VXLANClaim, error) {
	claim := vxlanbev1alpha1.BuildVXLANClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&vxlanbev1alpha1.VXLANClaimSpec{
			Index: index,
			ID:    ptr.To[uint32](uint32(r.id)),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := claim.ValidateSyntax(testType) // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return claim, nil
}

func (r testCtx) getRangeClaim(index, testType string) (*vxlanbev1alpha1.VXLANClaim, error) {
	claim := vxlanbev1alpha1.BuildVXLANClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&vxlanbev1alpha1.VXLANClaimSpec{
			Index: index,
			Range: ptr.To[string](r.tRange),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := claim.ValidateSyntax(testType) // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return claim, nil
}
