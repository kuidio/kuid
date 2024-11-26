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

package vlan

import (
	"errors"
	"fmt"

	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/idxtable/pkg/tree/tree16"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

var _ backend.IndexObject = &VLANIndex{}

func (r *VLANIndex) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *VLANIndex) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *VLANIndex) GetTree() gtree.GTree {
	//tree, err := tree32.New(32)
	tree, err := tree16.New(fmt.Sprintf("vlanidindex.%s", r.Name), 12)
	if err != nil {
		panic(err)
	}
	return tree
}

func (r *VLANIndex) GetType() string {
	return ""
}

func (r *VLANIndex) GetMinID() *uint64 {
	if r.Spec.MinID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.MinID))
}

func (r *VLANIndex) GetMaxID() *uint64 {
	if r.Spec.MaxID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.MaxID))
}

func (r *VLANIndex) GetMax() uint64 {
	return VLANID_Max
}

func GetMinClaimRange(id uint32) string {
	return fmt.Sprintf("%d-%d", VLANID_Min, id-1)
}

func GetMaxClaimRange(id uint32) string {
	return fmt.Sprintf("%d-%d", id+1, VLANID_Max)
}

func (r *VLANIndex) GetMinClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, backend.IndexReservedMinName),
	}
}

func (r *VLANIndex) GetMaxClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, backend.IndexReservedMaxName),
	}
}

func (r *VLANIndex) GetClaims() []backend.ClaimObject {
	claims := []backend.ClaimObject{}
	if r.GetMinID() != nil && *r.GetMinID() != 0 {
		claims = append(claims, r.GetMinClaim())
	}
	if r.GetMaxID() != nil && *r.GetMaxID() != r.GetMax() {
		claims = append(claims, r.GetMaxClaim())
	}
	return claims
}

func (r *VLANIndex) GetMinClaim() backend.ClaimObject {
	return BuildVLANClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMinClaimNSN().Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       VLANIndexKind,
					Name:       r.Name,
					UID:        r.UID,
				},
			},
		},
		&VLANClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMinClaimRange(*r.Spec.MinID)),
		},
		nil,
	)
}

func (r *VLANIndex) GetMaxClaim() backend.ClaimObject {
	return BuildVLANClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMaxClaimNSN().Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       VLANIndexKind,
					Name:       r.Name,
					UID:        r.UID,
				},
			},
		},
		&VLANClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMaxClaimRange(*r.Spec.MaxID)),
		},
		nil,
	)
}

func VLANIndexFromRuntime(ru runtime.Object) (backend.IndexObject, error) {
	index, ok := ru.(*VLANIndex)
	if !ok {
		return nil, errors.New("runtime object not VLANIndex")
	}
	return index, nil
}
