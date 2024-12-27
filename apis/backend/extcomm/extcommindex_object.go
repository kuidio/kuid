/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package extcomm

import (
	"errors"
	"fmt"

	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/idxtable/pkg/tree/tree16"
	"github.com/henderiw/idxtable/pkg/tree/tree32"
	"github.com/henderiw/idxtable/pkg/tree/tree64"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

var _ backend.IndexObject = &EXTCOMMIndex{}

func (r *EXTCOMMIndex) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *EXTCOMMIndex) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *EXTCOMMIndex) GetTree() gtree.GTree {
	switch GetEXTCOMMType(r.Spec.Type) {
	case ExtendedCommunityType_IPv4Address, ExtendedCommunityType_4byteAS:
		tree, err := tree16.New(fmt.Sprintf("extcommindex.%s", r.Name), 16)
		if err != nil {
			return nil
		}
		return tree
	case ExtendedCommunityType_2byteAS:
		tree, err := tree32.New(fmt.Sprintf("extcommindex.%s", r.Name), 32)
		if err != nil {
			return nil
		}
		return tree
	case ExtendedCommunityType_Opaque:
		tree, err := tree64.New(fmt.Sprintf("extcommindex.%s", r.Name), 48)
		if err != nil {
			return nil
		}
		return tree
	}
	return nil
}

func (r *EXTCOMMIndex) GetType() string {
	return r.Spec.Type
}

func (r *EXTCOMMIndex) GetMinID() *uint64 {
	if r.Spec.MinID == nil {
		return nil
	}
	return ptr.To(uint64(*r.Spec.MinID))
}

func (r *EXTCOMMIndex) GetMaxID() *uint64 {
	if r.Spec.MaxID == nil {
		return nil
	}
	return ptr.To(uint64(*r.Spec.MaxID))
}

func (r *EXTCOMMIndex) GetMax() uint64 {
	return EXTCOMMID_MaxValue[GetEXTCOMMType(r.GetType())]
}

func GetMinClaimRange(id uint64) string {
	return fmt.Sprintf("%d-%d", EXTCOMMID_Min, id-1)
}

func GetMaxClaimRange(extCommType ExtendedCommunityType, id uint64) string {
	return fmt.Sprintf("%d-%d", id+1, EXTCOMMID_MaxValue[extCommType])
}

func (r *EXTCOMMIndex) GetMinClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, backend.IndexReservedMinName),
	}
}

func (r *EXTCOMMIndex) GetMaxClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, backend.IndexReservedMaxName),
	}
}

func (r *EXTCOMMIndex) GetClaims() []backend.ClaimObject {
	claims := []backend.ClaimObject{}
	if r.GetMinID() != nil && *r.GetMinID() != 0 {
		claims = append(claims, r.GetMinClaim())
	}
	if r.GetMaxID() != nil && *r.GetMaxID() != r.GetMax() {
		claims = append(claims, r.GetMaxClaim())
	}
	for _, claim := range r.Spec.Claims {
		claims = append(claims, r.GetClaim(claim))
	}
	return claims
}

func (r *EXTCOMMIndex) GetMinClaim() backend.ClaimObject {
	return BuildEXTCOMMClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMinClaimNSN().Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       EXTCOMMIndexKind,
					Name:       r.Name,
					UID:        r.UID,
				},
			},
		},
		&EXTCOMMClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMinClaimRange(*r.Spec.MinID)),
		},
		nil,
	)
}

func (r *EXTCOMMIndex) GetMaxClaim() backend.ClaimObject {
	return BuildEXTCOMMClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMaxClaimNSN().Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       EXTCOMMIndexKind,
					Name:       r.Name,
					UID:        r.UID,
				},
			},
		},
		&EXTCOMMClaimSpec{
			Index: r.Name,
			Range: ptr.To(GetMaxClaimRange(GetEXTCOMMType(r.Spec.Type), *r.Spec.MaxID)),
		},
		nil,
	)
}

func (r *EXTCOMMIndex) GetClaim(claim EXTCOMMIndexClaim) backend.ClaimObject {
	spec := &EXTCOMMClaimSpec{
		Index: r.Name,
		ClaimLabels: common.ClaimLabels{
			UserDefinedLabels: claim.UserDefinedLabels,
		},
	}
	if claim.ID != nil {
		spec.ID = claim.ID
	} else {
		spec.Range = claim.Range
	}
	return BuildEXTCOMMClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      fmt.Sprintf("%s.%s", r.Name, claim.Name),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       EXTCOMMIndexKind,
					Name:       r.Name,
					UID:        r.UID,
				},
			},
		},
		spec,
		nil,
	)
}

func EXTCOMMIndexFromRuntime(ru runtime.Object) (backend.IndexObject, error) {
	index, ok := ru.(*EXTCOMMIndex)
	if !ok {
		return nil, errors.New("runtime object not EXTCOMMIndex")
	}
	return index, nil
}
