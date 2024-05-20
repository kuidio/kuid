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

const GENIDIndexPlural = "genidindices"
const GENIDIndexSingular = "genidindex"

// +k8s:deepcopy-gen=false
var _ resource.Object = &GENIDIndex{}
var _ resource.ObjectList = &GENIDIndexList{}

// GetListMeta returns the ListMeta
func (r *GENIDIndexList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *GENIDIndex) GetSingularName() string {
	return GENIDIndexSingular
}

func (GENIDIndex) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: GENIDIndexPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (GENIDIndex) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *GENIDIndex) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (GENIDIndex) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (GENIDIndex) New() runtime.Object {
	return &GENIDIndex{}
}

// NewList implements resource.Object
func (GENIDIndex) NewList() runtime.Object {
	return &GENIDIndexList{}
}

// GetCondition returns the condition based on the condition kind
func (r *GENIDIndex) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *GENIDIndex) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// GENIDIndexConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for GENIDIndex
func GENIDIndexConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *GENIDIndexList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *GENIDIndex) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
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

func (r *GENIDIndex) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *GENIDIndex) GetType() string {
	return r.Spec.Type
}

func (r *GENIDIndex) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      GENIDIndexKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *GENIDIndex) ValidateSyntax(_ string) field.ErrorList {
	var allErrs field.ErrorList

	if GetGenIDType(r.Spec.Type) == GENIDType_Invalid {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.type"),
			r,
			fmt.Errorf("invalid GENID Type %s", r.Spec.Type).Error(),
		))
	}

	if r.Spec.MinID != nil {
		if err := validateGENIDID(GetGenIDType(r.Spec.Type), *r.Spec.MinID); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.minID"),
				r,
				fmt.Errorf("invalid GENID ID %d", *r.Spec.MinID).Error(),
			))
		}
	}
	if r.Spec.MaxID != nil {
		if err := validateGENIDID(GetGenIDType(r.Spec.Type), *r.Spec.MaxID); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.maxID"),
				r,
				fmt.Errorf("invalid GENID ID %d", *r.Spec.MaxID).Error(),
			))
		}
	}
	if r.Spec.MinID != nil && r.Spec.MaxID != nil {
		if *r.Spec.MinID > *r.Spec.MaxID {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.maxID"),
				r,
				fmt.Errorf("min GENID ID %d cannot be bigger than max GENID ID %d", *r.Spec.MinID, *r.Spec.MaxID).Error(),
			))
		}
	}
	return allErrs
}

func GetMinClaimRange(id int64) string {
	return fmt.Sprintf("%d-%d", GENIDID_Min, id-1)
}

func GetMaxClaimRange(genidType GENIDType, id int64) string {
	return fmt.Sprintf("%d-%d", id+1, GENIDID_MaxValue[genidType])
}

func (r *GENIDIndex) GetMinID() *uint64 {
	if r.Spec.MinID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.MinID))
}

func (r *GENIDIndex) GetMaxID() *uint64 {
	if r.Spec.MaxID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.MaxID))
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

func (r *GENIDIndex) GetMinClaim() backend.ClaimObject {
	return BuildGENIDClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMinClaimNSN().Name,
		},
		&GENIDClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMinClaimRange(*r.Spec.MinID)),
			Owner: commonv1alpha1.GetOwnerReference(r),
		},
		nil,
	)
}

func (r *GENIDIndex) GetMaxClaim() backend.ClaimObject {
	return BuildGENIDClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMaxClaimNSN().Name,
		},
		&GENIDClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMaxClaimRange(GetGenIDType(r.Spec.Type), *r.Spec.MaxID)),
			Owner: commonv1alpha1.GetOwnerReference(r),
		},
		nil,
	)
}

// BuildGENIDIndex returns a reource from a client Object a Spec/Status
func BuildGENIDIndex(meta metav1.ObjectMeta, spec *GENIDIndexSpec, status *GENIDIndexStatus) *GENIDIndex {
	aspec := GENIDIndexSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := GENIDIndexStatus{}
	if status != nil {
		astatus = *status
	}
	return &GENIDIndex{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       GENIDIndexKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func GENIDIndexTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			index, ok := obj.(*GENIDIndex)
			if !ok {
				return nil
			}
			return []interface{}{
				index.GetName(),
				index.GetCondition(conditionv1alpha1.ConditionTypeReady).Status,
				index.Spec.Type,
				index.GetMinID(),
				index.GetMaxID(),
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			{Name: "Type", Type: "string"},
			{Name: "MinID", Type: "integer"},
			{Name: "MaxID", Type: "integer"},
		},
	}
}

func (r *GENIDIndex) GetSpec() any {
	return r.Spec
}

func (r *GENIDIndex) SetSpec(s any) {
	if spec, ok := s.(GENIDIndexSpec); ok {
		r.Spec = spec
	}
}

func (r *GENIDIndex) NewObjList() backend.GenericObjectList {
	return &GENIDIndexList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       GENIDIndexListKind},
	}
}

func (r *GENIDIndexList) GetObjects() []backend.GenericObject {
	objs := make([]backend.GenericObject, 0, len(r.Items))
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}
