package extcomm

import (
	"fmt"

	"github.com/kuidio/kuid/apis/backend"
	extcommbev1alpha1 "github.com/kuidio/kuid/apis/backend/extcomm/v1alpha1"
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

// aliEXTCOMM
const (
	namespace    = "dummy"
	staticClaim  = backend.ClaimType_StaticID
	dynamicClaim = backend.ClaimType_DynamicID
	rangeClaim   = backend.ClaimType_Range
)

func getIndex(index, testType string) (*extcommbev1alpha1.EXTCOMMIndex, error) {
	idx := extcommbev1alpha1.BuildEXTCOMMIndex(
		metav1.ObjectMeta{Namespace: namespace, Name: index},
		&extcommbev1alpha1.EXTCOMMIndexSpec{
			Type: testType,
			SubType: string(extcommbev1alpha1.ExtendedCommunitySubType_RouteTarget),
		},
		nil,
	)

	fieldErrs := idx.ValidateSyntax()
	if len(fieldErrs) != 0 {
		return nil, fmt.Errorf("syntax errors %v", fieldErrs)
	}
	return idx, nil
}

func (r testCtx) getDynamicClaim(index, testType string) (*extcommbev1alpha1.EXTCOMMClaim, error) {
	claim := extcommbev1alpha1.BuildEXTCOMMClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&extcommbev1alpha1.EXTCOMMClaimSpec{
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

func (r testCtx) getStaticClaim(index, testType string) (*extcommbev1alpha1.EXTCOMMClaim, error) {
	claim := extcommbev1alpha1.BuildEXTCOMMClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&extcommbev1alpha1.EXTCOMMClaimSpec{
			Index: index,
			ID:    ptr.To[uint64](r.id),
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

func (r testCtx) getRangeClaim(index, testType string) (*extcommbev1alpha1.EXTCOMMClaim, error) {
	claim := extcommbev1alpha1.BuildEXTCOMMClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&extcommbev1alpha1.EXTCOMMClaimSpec{
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
