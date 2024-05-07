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

package v1alpha1

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/idxtable/pkg/tree/tree16"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

const VLANIndexPlural = "vlanindices"
const VLANIndexSingular = "vlanindex"

// +k8s:deepcopy-gen=false
var _ resource.Object = &VLANIndex{}
var _ resource.ObjectList = &VLANIndexList{}

// GetListMeta returns the ListMeta
func (r *VLANIndexList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *VLANIndex) GetSingularName() string {
	return VLANIndexSingular
}

func (VLANIndex) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: VLANIndexPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (VLANIndex) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *VLANIndex) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (VLANIndex) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (VLANIndex) New() runtime.Object {
	return &VLANIndex{}
}

// NewList implements resource.Object
func (VLANIndex) NewList() runtime.Object {
	return &VLANIndexList{}
}

// GetCondition returns the condition based on the condition kind
func (r *VLANIndex) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *VLANIndex) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// ConvertVLANIndexFieldSelector is the schema conversion function for normalizing the FieldSelector for VLANIndex
func ConvertVLANIndexFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *VLANIndexList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *VLANIndex) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *VLANIndex) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *VLANIndex) GetTree() gtree.GTree {
	tree, err := tree16.New(12)
	if err != nil {
		return nil
	}
	return tree
}

func (r *VLANIndex) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *VLANIndex) GetType() string {
	return ""
}

func (r *VLANIndex) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      VLANIndexKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *VLANIndex) ValidateSyntax() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.MinID != nil {
		if err := validateVLANID(int(*r.Spec.MinID)); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.minID"),
				r,
				fmt.Errorf("invalid vlan ID %d", *r.Spec.MinID).Error(),
			))
		}
	}
	if r.Spec.MaxID != nil {
		if err := validateVLANID(int(*r.Spec.MaxID)); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.maxID"),
				r,
				fmt.Errorf("invalid vlan ID %d", *r.Spec.MaxID).Error(),
			))
		}
	}
	if r.Spec.MinID != nil && r.Spec.MaxID != nil {
		if *r.Spec.MinID > *r.Spec.MaxID {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.maxID"),
				r,
				fmt.Errorf("min vlan ID %d cannot be bigger than max vlan ID %d", *r.Spec.MinID, *r.Spec.MaxID).Error(),
			))
		}
	}
	return allErrs
}

func GetMinClaimRange(id uint32) string {
	return fmt.Sprintf("%d-%d", VLANID_Min, id-1)
}

func GetMaxClaimRange(id uint32) string {
	return fmt.Sprintf("%d-%d", id+1, VLANID_Max)
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

func (r *VLANIndex) GetMinClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, VLANIndexReservedMinName),
	}
}

func (r *VLANIndex) GetMaxClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, VLANIndexReservedMaxName),
	}
}

func (r *VLANIndex) GetMinClaim() backend.ClaimObject {
	return BuildVLANClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMinClaimNSN().Name,
		},
		&VLANClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMinClaimRange(*r.Spec.MinID)),
			Owner: commonv1alpha1.GetOwnerReference(r),
		},
		nil,
	)
}

func (r *VLANIndex) GetMaxClaim() backend.ClaimObject {
	return BuildVLANClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMaxClaimNSN().Name,
		},
		&VLANClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMaxClaimRange(*r.Spec.MaxID)),
			Owner: commonv1alpha1.GetOwnerReference(r),
		},
		nil,
	)
}

// BuildVLANIndex returns a reource from a client Object a Spec/Status
func BuildVLANIndex(meta metav1.ObjectMeta, spec *VLANIndexSpec, status *VLANIndexStatus) *VLANIndex {
	aspec := VLANIndexSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := VLANIndexStatus{}
	if status != nil {
		astatus = *status
	}
	return &VLANIndex{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       VLANIndexKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}
