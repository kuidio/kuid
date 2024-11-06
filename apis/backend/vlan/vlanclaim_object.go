/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "VLAN IS" BVLANIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vlan

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/idxtable/pkg/table/table32"
	"github.com/henderiw/idxtable/pkg/tree"
	"github.com/henderiw/idxtable/pkg/tree/id32"
	"github.com/henderiw/store"
	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/backend"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

var _ backend.ClaimObject = &VLANClaim{}

func (r *VLANClaim) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *VLANClaim) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

// GetCondition returns the condition bVLANed on the condition kind
func (r *VLANClaim) GetCondition(t condition.ConditionType) condition.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *VLANClaim) SetConditions(c ...condition.Condition) {
	r.Status.SetConditions(c...)
}

func (r *VLANClaim) ValidateSyntax(s string) field.ErrorList {
	var allErrs field.ErrorList

	if err := r.ValidateVLANClaimType(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath(""),
			r,
			err.Error(),
		))
		return allErrs
	}
	var v SyntaxValidator
	claimType := r.GetClaimType()
	switch claimType {
	case backend.ClaimType_DynamicID:
		v = &VLANDynamicIDSyntaxValidator{name: string(claimType)}
	case backend.ClaimType_StaticID:
		v = &VLANStaticIDSyntaxValidator{name: string(claimType)}
	case backend.ClaimType_Range:
		v = &VLANRangeSyntaxValidator{name: string(claimType)}
	default:
		return allErrs
	}
	return v.Validate(r)
}

func (r *VLANClaim) ValidateVLANRange() error {
	if r.Spec.Range == nil {
		return fmt.Errorf("no vlan range provided")
	}
	parts := strings.SplitN(*r.Spec.Range, "-", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid vlan range, expected <start>-<end>, got: %s", *r.Spec.Range)
	}
	var errm error
	if r.Name == r.Spec.Index {
		// to be able to check if the entry is reserved we get a parentname (rang name) equal to index
		// this is because the ownerreference uses the name of the index in its labels in the cache
		errm = errors.Join(errm, fmt.Errorf("a name of range cannot be the same as the index"))
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid vlan range start, got: %s, err: %s", *r.Spec.Range, err.Error()))
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid vlan range end, got: %s, err: %s", *r.Spec.Range, err.Error()))
	}
	if errm != nil {
		return errm
	}
	if start > end {
		errm = errors.Join(errm, fmt.Errorf("invalid vlan range start > end %s", *r.Spec.Range))
	}
	if err := validateVLANID(start); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid vlan start err %s", err.Error()))
	}
	if err := validateVLANID(end); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid vlan end err %s", err.Error()))
	}
	return errm
}

func (r *VLANClaim) ValidateVLANID() error {
	if r.Spec.ID == nil {
		return fmt.Errorf("no vlan id provided")
	}
	if err := validateVLANID(int(*r.Spec.ID)); err != nil {
		return fmt.Errorf("invalid vlan id err %s", err.Error())
	}
	return nil
}

func (r *VLANClaim) ValidateVLANClaimType() error {
	var sb strings.Builder
	count := 0
	if r.Spec.ID != nil {
		sb.WriteString(fmt.Sprintf("id: %d", *r.Spec.ID))
		count++

	}
	if r.Spec.Range != nil {
		if count > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("range: %s", *r.Spec.Range))
		count++

	}
	if count > 1 {
		return fmt.Errorf("a claim can only have 1 addressing, got %s", sb.String())
	}
	return nil
}

func (r *VLANClaim) GetIndex() string { return r.Spec.Index }

func (r *VLANClaim) GetSelector() *metav1.LabelSelector { return r.Spec.Selector }

func (r *VLANClaim) IsOwner(labels labels.Set) bool {
	ownerLabels := r.getOnwerLabels()
	for k, v := range ownerLabels {
		if val, ok := labels[k]; !ok || val != v {
			return false
		}
	}
	return true
}

func (r *VLANClaim) getOnwerLabels() map[string]string {
	claimName := r.Name
	claimKind := r.Kind
	claimUID := r.UID
	for _, owner := range r.GetOwnerReferences() {
		if owner.APIVersion == SchemeGroupVersion.Identifier() &&
			owner.Kind == VLANIndexKind {
			claimName = owner.Name
			claimKind = owner.Kind
			claimUID = owner.UID
		}
	}

	return map[string]string{
		backend.KuidClaimNameKey: claimName,
		backend.KuidClaimUIDKey:  string(claimUID),
		backend.KuidOwnerKindKey: claimKind,
	}
}

// GetOwnerSelector selects the route bVLANed on the name of the claim
func (r *VLANClaim) GetOwnerSelector() (labels.Selector, error) {
	l := r.getOnwerLabels()

	fullselector := labels.NewSelector()
	for k, v := range l {
		req, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			return nil, err
		}
		fullselector = fullselector.Add(*req)
	}
	return fullselector, nil
}

func (r *VLANClaim) GetLabelSelector() (labels.Selector, error) { return r.Spec.GetLabelSelector() }

func (r *VLANClaim) GetClaimLabels() labels.Set {
	labels := r.Spec.GetUserDefinedLabels()

	// for claims originated from the index we need to use the ownerreferences, since these claims
	// are never stored in the apiserver, the ip entries need to reference the index instead
	claimName := r.Name
	claimKind := VLANClaimKind
	claimUID := r.UID
	for _, owner := range r.GetOwnerReferences() {
		if owner.APIVersion == SchemeGroupVersion.Identifier() &&
			owner.Kind == VLANIndexKind {
			claimName = owner.Name
			claimKind = owner.Kind
			claimUID = owner.UID
		}
	}
	// system defined labels
	labels[backend.KuidClaimTypeKey] = string(r.GetClaimType())
	labels[backend.KuidClaimNameKey] = claimName
	labels[backend.KuidClaimUIDKey] = string(claimUID)
	labels[backend.KuidOwnerKindKey] = claimKind
	return labels
}

func (r *VLANClaim) ValidateOwner(labels labels.Set) error {
	routeClaimName := labels[backend.KuidClaimNameKey]
	routeClaimUID := labels[backend.KuidClaimUIDKey]

	if string(r.UID) != routeClaimUID && r.Name != routeClaimName {
		return fmt.Errorf("route owned by different claim got name %s/%s uid %s/%s",
			r.Name,
			routeClaimName,
			string(r.UID),
			routeClaimUID,
		)
	}
	return nil
}

func (r *VLANClaim) GetClaimType() backend.ClaimType {
	claimType := backend.ClaimType_Invalid
	count := 0
	if r.Spec.ID != nil {
		claimType = backend.ClaimType_StaticID
		count++

	}
	if r.Spec.Range != nil {
		claimType = backend.ClaimType_Range
		count++

	}
	if count > 1 {
		return backend.ClaimType_Invalid
	}
	if count == 0 {
		return backend.ClaimType_DynamicID
	}
	return claimType
}
func (r *VLANClaim) GetStaticID() *uint64 {
	if r.Spec.ID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.ID))
}
func (r *VLANClaim) GetStaticTreeID(t string) tree.ID {
	if r.Spec.ID == nil {
		return nil
	}
	return id32.NewID(*r.Spec.ID, id32.IDBitSize)
}

func (r *VLANClaim) GetClaimID(t string, id uint64) tree.ID {
	return id32.NewID(uint32(id), id32.IDBitSize)
}

func (r *VLANClaim) GetRange() *string {
	return r.Spec.Range
}

func (r *VLANClaim) GetRangeID(t string) (tree.Range, error) {
	if r.Spec.Range == nil {
		return nil, fmt.Errorf("cannot provide a range without an id")
	}
	return id32.ParseRange(*r.Spec.Range)
}

func (r *VLANClaim) GetTable(t string, to, from uint64) table.Table {
	return table32.New(uint32(to), uint32(from))
}

func (r *VLANClaim) SetStatusRange(s *string) {
	r.Status.Range = s
}

func (r *VLANClaim) SetStatusID(s *uint64) {
	if s == nil {
		r.Status.ID = nil
		return
	}
	r.Status.ID = ptr.To[uint32](uint32(*s))
}

func (r *VLANClaim) GetStatusID() *uint64 {
	if r.Status.ID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Status.ID))
}

func (r *VLANClaim) GetClaimRequest() string {
	// we assume validation is already done when calling this
	if r.Spec.ID != nil {
		return strconv.Itoa(int(*r.Spec.ID))
	}
	if r.Spec.Range != nil {
		return *r.Spec.Range
	}
	return ""
}

func (r *VLANClaim) GetClaimResponse() string {
	// we assume validation is already done when calling this
	if r.Status.ID != nil {
		return strconv.Itoa(int(*r.Status.ID))
	}
	if r.Status.Range != nil {
		return *r.Status.Range
	}
	return ""
}

func (r *VLANClaim) GetChoreoAPIVersion() string {
	return schema.GroupVersion{Group: GroupName, Version: "vlan"}.String()
}

func (r *VLANClaim) GetClaimSet(typ string) (sets.Set[tree.ID], error) {
	arange, err := r.GetRangeID(typ)
	if err != nil {
		return nil, fmt.Errorf("cannot get range from claim: %v", err)
	}
	// claim set represents the new entries
	newClaimSet := sets.New[tree.ID]()
	for _, rangeID := range arange.IDs() {
		newClaimSet.Insert(rangeID)
	}
	return newClaimSet, nil
}

func VLANClaimFromUnstructured(ru runtime.Unstructured) (backend.ClaimObject, error) {
	obj := &VLANClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(ru.UnstructuredContent(), obj)
	if err != nil {
		return nil, fmt.Errorf("error converting unstructured to asIndex: %v", err)
	}
	return obj, nil
}

func VLANClaimFromRuntime(ru runtime.Object) (backend.ClaimObject, error) {
	claim, ok := ru.(*VLANClaim)
	if !ok {
		return nil, errors.New("runtime object not VLANClaim")
	}
	return claim, nil
}

// BuildVLANClaim returns a reource from a client Object a Spec/Status
func BuildVLANClaim(meta metav1.ObjectMeta, spec *VLANClaimSpec, status *VLANClaimStatus) backend.ClaimObject {
	vlanspec := VLANClaimSpec{}
	if spec != nil {
		vlanspec = *spec
	}
	vlanstatus := VLANClaimStatus{}
	if status != nil {
		vlanstatus = *status
	}
	return &VLANClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       VLANClaimKind,
		},
		ObjectMeta: meta,
		Spec:       vlanspec,
		Status:     vlanstatus,
	}
}
