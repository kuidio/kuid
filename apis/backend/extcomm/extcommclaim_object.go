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

package extcomm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/idxtable/pkg/table/table16"
	"github.com/henderiw/idxtable/pkg/tree"
	"github.com/henderiw/idxtable/pkg/tree/id16"
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

var _ backend.ClaimObject = &EXTCOMMClaim{}

func (r *EXTCOMMClaim) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *EXTCOMMClaim) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

// GetCondition returns the condition bVLANed on the condition kind
func (r *EXTCOMMClaim) GetCondition(t condition.ConditionType) condition.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *EXTCOMMClaim) SetConditions(c ...condition.Condition) {
	r.Status.SetConditions(c...)
}

func (r *EXTCOMMClaim) ValidateSyntax(s string) field.ErrorList {
	var allErrs field.ErrorList

	if GetEXTCOMMType(s) == ExtendedCommunityType_Invalid {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath(""),
			r,
			fmt.Errorf("invalid extended community type. got %s", s).Error(),
		))
		return allErrs
	}

	if err := r.ValidateEXTCOMMClaimType(); err != nil {
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
		v = &EXTCOMMDynamicIDSyntaxValidator{name: string(claimType)}
	case backend.ClaimType_StaticID:
		v = &EXTCOMMStaticIDSyntaxValidator{name: string(claimType)}
	case backend.ClaimType_Range:
		v = &EXTCOMMRangeSyntaxValidator{name: string(claimType)}
	default:
		return allErrs
	}
	return v.Validate(r, GetEXTCOMMType(s))
}

func (r *EXTCOMMClaim) ValidateEXTCOMMRange(extCommType ExtendedCommunityType) error {
	if r.Spec.Range == nil {
		return fmt.Errorf("no EXTCOMM range provided")
	}
	parts := strings.SplitN(*r.Spec.Range, "-", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid EXTCOMM range, expected <start>-<end>, got: %s", *r.Spec.Range)
	}
	var errm error
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid EXTCOMM range start, got: %s, err: %s", *r.Spec.Range, err.Error()))
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid EXTCOMM range end, got: %s, err: %s", *r.Spec.Range, err.Error()))
	}
	if errm != nil {
		return errm
	}
	if start > end {
		errm = errors.Join(errm, fmt.Errorf("invalid EXTCOMM range start > end %s", *r.Spec.Range))
	}
	if err := validateEXTCOMMID(extCommType, uint64(start)); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid EXTCOMM start err %s", err.Error()))
	}
	if err := validateEXTCOMMID(extCommType, uint64(end)); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid EXTCOMM end err %s", err.Error()))
	}
	return errm
}

func (r *EXTCOMMClaim) ValidateEXTCOMMID(extCommType ExtendedCommunityType) error {
	if r.Spec.ID == nil {
		return fmt.Errorf("no id provided")
	}
	if err := validateEXTCOMMID(extCommType, *r.Spec.ID); err != nil {
		return fmt.Errorf("invalid id err %s", err.Error())
	}
	return nil
}

func validateEXTCOMMID(extCommType ExtendedCommunityType, id uint64) error {
	if id < EXTCOMMID_Min {
		return fmt.Errorf("invalid id, got %d", id)
	}
	if id > uint64(EXTCOMMID_MaxValue[extCommType]) {
		return fmt.Errorf("invalid id, got %d", id)
	}
	return nil
}

func (r *EXTCOMMClaim) ValidateEXTCOMMClaimType() error {
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

func (r *EXTCOMMClaim) GetIndex() string { return r.Spec.Index }

func (r *EXTCOMMClaim) GetSelector() *metav1.LabelSelector { return r.Spec.Selector }

func (r *EXTCOMMClaim) IsOwner(labels labels.Set) bool {
	for k, v := range r.getOwnerLabels() {
		if val, ok := labels[k]; !ok || val != v {
			return false
		}
	}
	return true
}

func (r *EXTCOMMClaim) getOwnerLabels() map[string]string {
	return map[string]string{
		backend.KuidClaimNameKey: r.Name,
		backend.KuidClaimUIDKey:  string(r.UID),
	}
}

// GetOwnerSelector selects the route bVLANed on the name of the claim
func (r *EXTCOMMClaim) GetOwnerSelector() (labels.Selector, error) {
	l := r.getOwnerLabels()

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

func (r *EXTCOMMClaim) GetLabelSelector() (labels.Selector, error) { return r.Spec.GetLabelSelector() }

func (r *EXTCOMMClaim) GetClaimLabels() labels.Set {
	labels := r.Spec.GetUserDefinedLabels()

	// system defined labels
	labels[backend.KuidClaimTypeKey] = string(r.GetClaimType())
	labels[backend.KuidClaimNameKey] = r.Name
	labels[backend.KuidClaimUIDKey] = string(r.UID)
	labels[backend.KuidOwnerKindKey] = r.Kind
	return labels
}

func (r *EXTCOMMClaim) ValidateOwner(labels labels.Set) error {
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

func (r *EXTCOMMClaim) GetClaimType() backend.ClaimType {
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
func (r *EXTCOMMClaim) GetStaticID() *uint64 {
	if r.Spec.ID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.ID))
}
func (r *EXTCOMMClaim) GetStaticTreeID(t string) tree.ID {
	if r.Spec.ID == nil {
		return nil
	}
	return id16.NewID(uint16(*r.Spec.ID), id16.IDBitSize)
}

func (r *EXTCOMMClaim) GetClaimID(t string, id uint64) tree.ID {
	return id16.NewID(uint16(id), id16.IDBitSize)
}

func (r *EXTCOMMClaim) GetStatusClaimID() tree.ID {
	if r.Status.ID == nil {
		return nil
	}
	return id16.NewID(uint16(*r.Status.ID), id16.IDBitSize) 
}

func (r *EXTCOMMClaim) GetRange() *string {
	return r.Spec.Range
}

func (r *EXTCOMMClaim) GetRangeID(t string) (tree.Range, error) {
	if r.Spec.Range == nil {
		return nil, fmt.Errorf("cannot provide a range without an id")
	}
	return id32.ParseRange(*r.Spec.Range)
}

func (r *EXTCOMMClaim) GetTable(t string, to, from uint64) table.Table {
	return table16.New(uint16(to), uint16(from))
}

func (r *EXTCOMMClaim) SetStatusRange(s *string) {
	r.Status.Range = s
}

func (r *EXTCOMMClaim) SetStatusID(s *uint64) {
	if s == nil {
		r.Status.ID = nil
		return
	}
	r.Status.ID = ptr.To[uint64](uint64(*s))
}

func (r *EXTCOMMClaim) GetStatusID() *uint64 {
	if r.Status.ID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Status.ID))
}

func (r *EXTCOMMClaim) GetClaimRequest() string {
	// we assume validation is already done when calling this
	if r.Spec.ID != nil {
		return strconv.Itoa(int(*r.Spec.ID))
	}
	if r.Spec.Range != nil {
		return *r.Spec.Range
	}
	return ""
}

func (r *EXTCOMMClaim) GetClaimResponse() string {
	// we assume validation is already done when calling this
	if r.Status.ID != nil {
		return strconv.Itoa(int(*r.Status.ID))
	}
	if r.Status.Range != nil {
		return *r.Status.Range
	}
	return ""
}

func (r *EXTCOMMClaim) GetClaimSet(typ string) (map[string]tree.ID, sets.Set[string], error) {
	arange, err := r.GetRangeID(typ)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get range from claim: %v", err)
	}
	// claim set represents the new entries
	newClaimSet := sets.New[string]()
	newClaimMap := map[string]tree.ID{}
	for _, rangeID := range arange.IDs() {
		newClaimSet.Insert(rangeID.String())
		newClaimMap[rangeID.String()] = rangeID
	}
	return newClaimMap, newClaimSet, nil
}

func (r *EXTCOMMClaim) GetChoreoAPIVersion() string {
	return schema.GroupVersion{Group: GroupName, Version: "extcomm"}.String()
}

func EXTCOMMClaimFromUnstructured(ru runtime.Unstructured) (backend.ClaimObject, error) {
	obj := &EXTCOMMClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(ru.UnstructuredContent(), obj)
	if err != nil {
		return nil, fmt.Errorf("error converting unstructured to asIndex: %v", err)
	}
	return obj, nil
}

func EXTCOMMClaimFromRuntime(ru runtime.Object) (backend.ClaimObject, error) {
	claim, ok := ru.(*EXTCOMMClaim)
	if !ok {
		return nil, errors.New("runtime object not EXTCOMMClaim")
	}
	return claim, nil
}

// BuildEXTCOMMClaim returns a reource from a client Object a Spec/Status
func BuildEXTCOMMClaim(meta metav1.ObjectMeta, spec *EXTCOMMClaimSpec, status *EXTCOMMClaimStatus) backend.ClaimObject {
	vlanspec := EXTCOMMClaimSpec{}
	if spec != nil {
		vlanspec = *spec
	}
	vlanstatus := EXTCOMMClaimStatus{}
	if status != nil {
		vlanstatus = *status
	}
	return &EXTCOMMClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       EXTCOMMClaimKind,
		},
		ObjectMeta: meta,
		Spec:       vlanspec,
		Status:     vlanstatus,
	}
}
