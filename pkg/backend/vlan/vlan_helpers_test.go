package vlan

import (
	"fmt"

	"github.com/kuidio/kuid/apis/backend"
	vlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
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

func getIndex(index, _ string) (*vlanbev1alpha1.VLANIndex, error) {
	idx := vlanbev1alpha1.BuildVLANIndex(
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

func (r testCtx) getDynamicClaim(index, testType string) (*vlanbev1alpha1.VLANClaim, error) {
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
	fielErrList := claim.ValidateSyntax(testType) // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return claim, nil
}

func (r testCtx) getStaticClaim(index, testType string) (*vlanbev1alpha1.VLANClaim, error) {
	claim := vlanbev1alpha1.BuildVLANClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&vlanbev1alpha1.VLANClaimSpec{
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

func (r testCtx) getRangeClaim(index, testType string) (*vlanbev1alpha1.VLANClaim, error) {
	claim := vlanbev1alpha1.BuildVLANClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&vlanbev1alpha1.VLANClaimSpec{
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
