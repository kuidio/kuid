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
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/idxtable/pkg/tree/tree16"
	"github.com/henderiw/idxtable/pkg/tree/tree32"
	"github.com/henderiw/idxtable/pkg/tree/tree64"
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

const EXTCOMMIndexPlural = "extcommindices"
const EXTCOMMIndexSingular = "extcommindex"

// +k8s:deepcopy-gen=false
var _ resource.Object = &EXTCOMMIndex{}
var _ resource.ObjectList = &EXTCOMMIndexList{}

// GetListMeta returns the ListMeta
func (r *EXTCOMMIndexList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *EXTCOMMIndex) GetSingularName() string {
	return EXTCOMMIndexSingular
}

func (EXTCOMMIndex) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: EXTCOMMIndexPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (EXTCOMMIndex) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *EXTCOMMIndex) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (EXTCOMMIndex) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (EXTCOMMIndex) New() runtime.Object {
	return &EXTCOMMIndex{}
}

// NewList implements resource.Object
func (EXTCOMMIndex) NewList() runtime.Object {
	return &EXTCOMMIndexList{}
}

// GetCondition returns the condition based on the condition kind
func (r *EXTCOMMIndex) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *EXTCOMMIndex) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// EXTCOMMIndexConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for EXTCOMMIndex
func EXTCOMMIndexConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *EXTCOMMIndexList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *EXTCOMMIndex) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *EXTCOMMIndex) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *EXTCOMMIndex) GetTree() gtree.GTree {
	switch GetExtendedCommunityType(r.Spec.Type) {
	case ExtendedCommunityType_IPv4Address, ExtendedCommunityType_4byteAS:
		tree, err := tree16.New(16)
		if err != nil {
			return nil
		}
		return tree
	case ExtendedCommunityType_2byteAS:
		tree, err := tree32.New(32)
		if err != nil {
			return nil
		}
		return tree
	case ExtendedCommunityType_Opaque:
		tree, err := tree64.New(48)
		if err != nil {
			return nil
		}
		return tree
	}
	return nil
}

func (r *EXTCOMMIndex) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *EXTCOMMIndex) GetType() string {
	return r.Spec.Type
}

func (r *EXTCOMMIndex) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      EXTCOMMIndexKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *EXTCOMMIndex) ValidateSyntax() field.ErrorList {
	var allErrs field.ErrorList

	if GetExtendedCommunityType(r.Spec.Type) == ExtendedCommunityType_Invalid {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.type"),
			r,
			fmt.Errorf("invalid EXTCOMM Type %s", r.Spec.Type).Error(),
		))
	}

	if GetExtendedCommunitySubType(r.Spec.SubType) == ExtendedCommunitySubType_Invalid {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.subType"),
			r,
			fmt.Errorf("invalid EXTCOMM SubType %s", r.Spec.SubType).Error(),
		))
	}

	if r.Spec.MinID != nil {
		if err := validateEXTCOMMID(GetExtendedCommunityType(r.Spec.Type), *r.Spec.MinID); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.minID"),
				r,
				fmt.Errorf("invalid EXTCOMM ID %d", *r.Spec.MinID).Error(),
			))
		}
	}
	if r.Spec.MaxID != nil {
		if err := validateEXTCOMMID(GetExtendedCommunityType(r.Spec.Type), *r.Spec.MaxID); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.maxID"),
				r,
				fmt.Errorf("invalid EXTCOMM ID %d", *r.Spec.MaxID).Error(),
			))
		}
	}
	if r.Spec.MinID != nil && r.Spec.MaxID != nil {
		if *r.Spec.MinID > *r.Spec.MaxID {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.maxID"),
				r,
				fmt.Errorf("min EXTCOMM ID %d cannot be bigger than max EXTCOMM ID %d", *r.Spec.MinID, *r.Spec.MaxID).Error(),
			))
		}
	}
	return allErrs
}

func GetMinClaimRange(id int64) string {
	return fmt.Sprintf("%d-%d", EXTCOMMID_Min, id-1)
}

func GetMaxClaimRange(extCommType ExtendedCommunityType, id int64) string {
	return fmt.Sprintf("%d-%d", id+1, EXTCOMMID_MaxValue[extCommType])
}

func (r *EXTCOMMIndex) GetMinID() *uint64 {
	if r.Spec.MinID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.MinID))
}

func (r *EXTCOMMIndex) GetMaxID() *uint64 {
	if r.Spec.MaxID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.MaxID))
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

func (r *EXTCOMMIndex) GetMinClaim() backend.ClaimObject {
	return BuildEXTCOMMClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMinClaimNSN().Name,
		},
		&EXTCOMMClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMinClaimRange(*r.Spec.MinID)),
			Owner: commonv1alpha1.GetOwnerReference(r),
		},
		nil,
	)
}

func (r *EXTCOMMIndex) GetMaxClaim() backend.ClaimObject {
	return BuildEXTCOMMClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMaxClaimNSN().Name,
		},
		&EXTCOMMClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMaxClaimRange(GetExtendedCommunityType(r.Spec.Type), *r.Spec.MaxID)),
			Owner: commonv1alpha1.GetOwnerReference(r),
		},
		nil,
	)
}

// BuildEXTCOMMIndex returns a reource from a client Object a Spec/Status
func BuildEXTCOMMIndex(meta metav1.ObjectMeta, spec *EXTCOMMIndexSpec, status *EXTCOMMIndexStatus) *EXTCOMMIndex {
	aspec := EXTCOMMIndexSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := EXTCOMMIndexStatus{}
	if status != nil {
		astatus = *status
	}
	return &EXTCOMMIndex{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       EXTCOMMIndexKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func EXTCOMMIndexTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			index, ok := obj.(*EXTCOMMIndex)
			if !ok {
				return nil
			}
			return []interface{}{
				index.GetName(),
				index.GetCondition(conditionv1alpha1.ConditionTypeReady).Status,
				index.Spec.Transitive,
				index.Spec.Type,
				index.Spec.SubType,
				index.Spec.GlobalID,
				index.GetMinID(),
				index.GetMaxID(),
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			{Name: "Transitive", Type: "boolean"},
			{Name: "Type", Type: "string"},
			{Name: "SubType", Type: "string"},
			{Name: "GlobalID", Type: "string"},
			{Name: "MinID", Type: "integer"},
			{Name: "MaxID", Type: "integer"},
		},
	}
}
