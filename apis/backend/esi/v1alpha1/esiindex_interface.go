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
	"github.com/henderiw/idxtable/pkg/tree/tree64"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	rresource "github.com/kuidio/kuid/pkg/reconcilers/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

const ESIIndexPlural = "esiindices"
const ESIIndexSingular = "esiindex"

// +k8s:deepcopy-gen=false
var _ resource.Object = &ESIIndex{}
var _ resource.ObjectList = &ESIIndexList{}

// GetListMeta returns the ListMeta
func (r *ESIIndexList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *ESIIndex) GetSingularName() string {
	return ESIIndexSingular
}

func (ESIIndex) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: ESIIndexPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (ESIIndex) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *ESIIndex) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (ESIIndex) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (ESIIndex) New() runtime.Object {
	return &ESIIndex{}
}

// NewList implements resource.Object
func (ESIIndex) NewList() runtime.Object {
	return &ESIIndexList{}
}

// GetCondition returns the condition based on the condition kind
func (r *ESIIndex) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *ESIIndex) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// ConvertESIIndexFieldSelector is the schema conversion function for normalizing the FieldSelector for ESIIndex
func ConvertESIIndexFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *ESIIndexList) GetItems() []rresource.Object {
	objs := []rresource.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *ESIIndex) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *ESIIndex) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *ESIIndex) GetTree() gtree.GTree {
	tree, err := tree64.New(64)
	if err != nil {
		return nil
	}
	return tree
}

func (r *ESIIndex) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *ESIIndex) GetType() string {
	return ""
}

func (r *ESIIndex) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      ESIIndexKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *ESIIndex) ValidateSyntax() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.MinID != nil {
		if err := validateESIID(int(*r.Spec.MinID)); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.minID"),
				r,
				fmt.Errorf("invalid ESI ID %d", *r.Spec.MinID).Error(),
			))
		}
	}
	if r.Spec.MaxID != nil {
		if err := validateESIID(int(*r.Spec.MaxID)); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.maxID"),
				r,
				fmt.Errorf("invalid ESI ID %d", *r.Spec.MaxID).Error(),
			))
		}
	}
	if r.Spec.MinID != nil && r.Spec.MaxID != nil {
		if *r.Spec.MinID > *r.Spec.MaxID {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.maxID"),
				r,
				fmt.Errorf("min ESI ID %d cannot be bigger than max ESI ID %d", *r.Spec.MinID, *r.Spec.MaxID).Error(),
			))
		}
	}
	return allErrs
}

func GetMinClaimRange(id uint64) string {
	return fmt.Sprintf("%d-%d", ESIID_Min, id-1)
}

func GetMaxClaimRange(id uint64) string {
	return fmt.Sprintf("%d-%d", id+1, ESIID_Max)
}

func (r *ESIIndex) GetMinID() *uint64 {
	if r.Spec.MinID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.MinID))
}

func (r *ESIIndex) GetMaxID() *uint64 {
	if r.Spec.MaxID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.MaxID))
}

func (r *ESIIndex) GetMinClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, ESIIndexReservedMinName),
	}
}

func (r *ESIIndex) GetMaxClaimNSN() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.Namespace,
		Name:      fmt.Sprintf("%s.%s", r.Name, ESIIndexReservedMaxName),
	}
}

func (r *ESIIndex) GetMinClaim() backend.ClaimObject {
	return BuildESIClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMinClaimNSN().Name,
		},
		&ESIClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMinClaimRange(*r.Spec.MinID)),
			Owner: commonv1alpha1.GetOwnerReference(r),
		},
		nil,
	)
}

func (r *ESIIndex) GetMaxClaim() backend.ClaimObject {
	return BuildESIClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      r.GetMaxClaimNSN().Name,
		},
		&ESIClaimSpec{
			Index: r.Name,
			Range: ptr.To[string](GetMaxClaimRange(*r.Spec.MaxID)),
			Owner: commonv1alpha1.GetOwnerReference(r),
		},
		nil,
	)
}

// BuildESIIndex returns a reource from a client Object a Spec/Status
func BuildESIIndex(meta metav1.ObjectMeta, spec *ESIIndexSpec, status *ESIIndexStatus) *ESIIndex {
	aspec := ESIIndexSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := ESIIndexStatus{}
	if status != nil {
		astatus = *status
	}
	return &ESIIndex{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       ESIIndexKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}
