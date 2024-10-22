package astest

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/apis/backend/as"
	"github.com/kuidio/kuid/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

func getIndex(index, _ string) (*as.ASIndex, error) {
	idx := as.BuildASIndex(
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

func (r testCtx) getDynamicClaim(index, testType string) (backend.ClaimObject, error) {
	claim := as.BuildASClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&as.ASClaimSpec{
			Index: index,
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
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

func (r testCtx) getStaticClaim(index, testType string) (backend.ClaimObject, error) {
	claim := as.BuildASClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&as.ASClaimSpec{
			Index: index,
			ID:    ptr.To[uint32](uint32(r.id)),
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
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

func (r testCtx) getRangeClaim(index, testType string) (backend.ClaimObject, error) {
	fmt.Println("getRangeClaim", r.name, r.tRange)
	claim := as.BuildASClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&as.ASClaimSpec{
			Index: index,
			Range: ptr.To[string](r.tRange),
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := claim.ValidateSyntax(testType) // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	fmt.Println("getRangeClaim", *claim.GetRange())
	return claim, nil
}

func transformer(_ context.Context, newObj runtime.Object, oldObj runtime.Object) (runtime.Object, error) {
	// Type assertion to specific object types, assuming we are working with a type that has Spec and Status fields
	new, ok := newObj.(backend.ClaimObject)
	if !ok {
		return nil, fmt.Errorf("newObj is not of type ClaimObject %s", reflect.TypeOf(newObj).Name())
	}
	old, ok := oldObj.(backend.ClaimObject)
	if !ok {
		return nil, fmt.Errorf("oldObj is not of type ClaimObject %s", reflect.TypeOf(newObj).Name())
	}

	new.SetResourceVersion(old.GetResourceVersion())
	new.SetUID(old.GetUID())

	if new.GetRange() != nil {
		fmt.Println("transformer", "new", *new.GetRange() )
	}

	if old.GetRange() != nil {
		fmt.Println("transformer", "old", *old.GetRange() )
	}
	

	return new, nil
}
