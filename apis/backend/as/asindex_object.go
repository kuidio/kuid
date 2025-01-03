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

package as

import (
	"errors"
	"fmt"

	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/idxtable/pkg/tree/tree32"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

var _ backend.IndexObject = &ASIndex{}

func (r *ASIndex) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *ASIndex) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *ASIndex) GetTree() gtree.GTree {
	tree, err := tree32.New(fmt.Sprintf("asindex.%s", r.Name), 32)
	if err != nil {
		panic(err)
	}
	return tree
}

func (r *ASIndex) GetType() string {
	return ""
}

func (r *ASIndex) GetMinID() *uint64 {
	if r.Spec.MinID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.MinID))
}

func (r *ASIndex) GetMaxID() *uint64 {
	if r.Spec.MaxID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.MaxID))
}

func (r *ASIndex) GetMax() uint64 {
	return ASID_Max
}

func GetMinClaimRange(id uint32) string {
	return fmt.Sprintf("%d-%d", ASID_Min, id-1)
}

func GetMaxClaimRange(id uint32) string {
	return fmt.Sprintf("%d-%d", id+1, ASID_Max)
}

func (r *ASIndex) GetMinClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, backend.IndexReservedMinName),
	}
}

func (r *ASIndex) GetMaxClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, backend.IndexReservedMaxName),
	}
}

func (r *ASIndex) GetClaims() []backend.ClaimObject {
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

func (r *ASIndex) GetMinClaim() backend.ClaimObject {
	return BuildASClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMinClaimNSN().Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       ASIndexKind,
					Name:       r.Name,
					UID:        r.UID,
				},
			},
		},
		&ASClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMinClaimRange(*r.Spec.MinID)),
		},
		nil,
	)
}

func (r *ASIndex) GetMaxClaim() backend.ClaimObject {
	return BuildASClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMaxClaimNSN().Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       ASIndexKind,
					Name:       r.Name,
					UID:        r.UID,
				},
			},
		},
		&ASClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMaxClaimRange(*r.Spec.MaxID)),
		},
		nil,
	)
}

func (r *ASIndex) GetClaim(claim ASIndexClaim) backend.ClaimObject {
	spec := &ASClaimSpec{
		Index:       r.Name,
		ClaimLabels: common.ClaimLabels{
			UserDefinedLabels: claim.UserDefinedLabels,
		},
	}
	if claim.ID != nil {
		spec.ID = claim.ID
	} else {
		spec.Range = claim.Range
	}
	return BuildASClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      fmt.Sprintf("%s.%s", r.Name, claim.Name),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       ASIndexKind,
					Name:       r.Name,
					UID:        r.UID,
				},
			},
		},
		spec,
		nil,
	)
}

func ASIndexFromRuntime(ru runtime.Object) (backend.IndexObject, error) {
	index, ok := ru.(*ASIndex)
	if !ok {
		return nil, errors.New("runtime object not ASIndex")
	}
	return index, nil
}
