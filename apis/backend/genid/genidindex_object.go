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

package genid

import (
	"errors"
	"fmt"

	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/idxtable/pkg/tree/tree16"
	"github.com/henderiw/idxtable/pkg/tree/tree32"
	"github.com/henderiw/idxtable/pkg/tree/tree64"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

var _ backend.IndexObject = &GENIDIndex{}

func (r *GENIDIndex) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *GENIDIndex) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *GENIDIndex) GetTree() gtree.GTree {
	switch GetGenIDType(r.Spec.Type) {
	case GENIDType_16bit:
		tree, err := tree16.New(16)
		if err != nil {
			return nil
		}
		return tree
	case GENIDType_32bit:
		tree, err := tree32.New(32)
		if err != nil {
			return nil
		}
		return tree
	case GENIDType_48bit:
		tree, err := tree64.New(48)
		if err != nil {
			return nil
		}
		return tree
	case GENIDType_64bit:
		tree, err := tree64.New(64)
		if err != nil {
			return nil
		}
		return tree
	default:
		return nil
	}
}

func (r *GENIDIndex) GetType() string {
	return r.Spec.Type
}

func (r *GENIDIndex) GetMinID() *uint64 {
	if r.Spec.MinID == nil {
		return nil
	}
	return ptr.To(uint64(*r.Spec.MinID))
}

func (r *GENIDIndex) GetMaxID() *uint64 {
	if r.Spec.MaxID == nil {
		return nil
	}
	return ptr.To(uint64(*r.Spec.MaxID))
}

func (r *GENIDIndex) GetMax() uint64 {
	return GENIDID_MaxValue[GetGenIDType(r.Spec.Type)]
}

func GetMinClaimRange(id uint64) string {
	return fmt.Sprintf("%d-%d", GENIDID_Min, id-1)
}

func GetMaxClaimRange(genidType GENIDType, id uint64) string {
	return fmt.Sprintf("%d-%d", id+1, GENIDID_MaxValue[genidType])
}

func (r *GENIDIndex) GetMinClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, backend.IndexReservedMinName),
	}
}

func (r *GENIDIndex) GetMaxClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, backend.IndexReservedMaxName),
	}
}

func (r *GENIDIndex) GetClaims() []backend.ClaimObject {
	claims := []backend.ClaimObject{}
	if r.GetMinID() != nil && *r.GetMinID() != 0 {
		claims = append(claims, r.GetMinClaim())
	}
	if r.GetMaxID() != nil && *r.GetMaxID() != r.GetMax() {
		claims = append(claims, r.GetMaxClaim())
	}
	return claims
}

func (r *GENIDIndex) GetMinClaim() backend.ClaimObject {
	return BuildGENIDClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMinClaimNSN().Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       GENIDIndexKind,
					Name:       r.Name,
					UID:        r.UID,
				},
			},
		},
		&GENIDClaimSpec{
			Index: r.Name,
			Range: ptr.To(GetMinClaimRange(*r.Spec.MinID)),
		},
		nil,
	)
}

func (r *GENIDIndex) GetMaxClaim() backend.ClaimObject {
	return BuildGENIDClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMaxClaimNSN().Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       r.Kind,
					Name:       r.Name,
					UID:        r.UID,
				},
			},
		},
		&GENIDClaimSpec{
			Index: r.Name,
			Range: ptr.To(GetMaxClaimRange(GetGenIDType(r.Spec.Type), *r.Spec.MaxID)),
		},
		nil,
	)
}

func GENIDIndexFromRuntime(ru runtime.Object) (backend.IndexObject, error) {
	index, ok := ru.(*GENIDIndex)
	if !ok {
		return nil, errors.New("runtime object not GENIDIndex")
	}
	return index, nil
}
