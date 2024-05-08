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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/idxtable/pkg/table/table16"
	"github.com/henderiw/idxtable/pkg/table/table32"
	"github.com/henderiw/idxtable/pkg/table/table64"
	"github.com/henderiw/idxtable/pkg/tree"
	"github.com/henderiw/idxtable/pkg/tree/id16"
	"github.com/henderiw/idxtable/pkg/tree/id32"
	"github.com/henderiw/idxtable/pkg/tree/id64"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

const EXTCOMMClaimPlural = "extcommclaims"
const EXTCOMMClaimSingular = "extcommclaim"
const EXTCOMMID_Min = 0

var EXTCOMMID_MaxBits = map[ExtendedCommunityType]int{
	ExtendedCommunityType_Invalid:     0,
	ExtendedCommunityType_2byteAS:     32,
	ExtendedCommunityType_4byteAS:     16,
	ExtendedCommunityType_IPv4Address: 16,
	ExtendedCommunityType_Opaque:      48,
}

var EXTCOMMID_MaxValue = map[ExtendedCommunityType]int64{
	ExtendedCommunityType_Invalid:     1<<EXTCOMMID_MaxBits[ExtendedCommunityType_Invalid] - 1,
	ExtendedCommunityType_2byteAS:     1<<EXTCOMMID_MaxBits[ExtendedCommunityType_2byteAS] - 1,
	ExtendedCommunityType_4byteAS:     1<<EXTCOMMID_MaxBits[ExtendedCommunityType_4byteAS] - 1,
	ExtendedCommunityType_IPv4Address: 1<<EXTCOMMID_MaxBits[ExtendedCommunityType_IPv4Address] - 1,
	ExtendedCommunityType_Opaque:      1<<EXTCOMMID_MaxBits[ExtendedCommunityType_Opaque] - 1,
}

// +k8s:deepcopy-gen=false
var _ resource.Object = &EXTCOMMClaim{}
var _ resource.ObjectList = &EXTCOMMClaimList{}

var _ resource.ObjectWithStatusSubResource = &EXTCOMMClaim{}

func (EXTCOMMClaimStatus) SubResourceName() string {
	return fmt.Sprintf("%s/%s", EXTCOMMClaimPlural, "status")
}

func (r EXTCOMMClaimStatus) CopyTo(obj resource.ObjectWithStatusSubResource) {
	cfg, ok := obj.(*EXTCOMMClaim)
	if ok {
		cfg.Status = r
	}
}

func (r *EXTCOMMClaim) GetStatus() resource.StatusSubResource {
	return r.Status
}

// GetListMeta returns the ListMeta
func (r *EXTCOMMClaimList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *EXTCOMMClaim) GetSingularName() string {
	return EXTCOMMClaimSingular
}

func (EXTCOMMClaim) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: EXTCOMMClaimPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (EXTCOMMClaim) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *EXTCOMMClaim) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (EXTCOMMClaim) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (EXTCOMMClaim) New() runtime.Object {
	return &EXTCOMMClaim{}
}

// NewList implements resource.Object
func (EXTCOMMClaim) NewList() runtime.Object {
	return &EXTCOMMClaimList{}
}

// GetCondition returns the condition based on the condition kind
func (r *EXTCOMMClaim) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *EXTCOMMClaim) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// EXTCOMMClaimConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for EXTCOMMClaim
func EXTCOMMClaimConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	case "spec.index":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *EXTCOMMClaimList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *EXTCOMMClaim) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *EXTCOMMClaim) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *EXTCOMMClaim) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

func (r *EXTCOMMClaim) GetIndex() string {
	return r.Spec.Index
}

func (r *EXTCOMMClaim) GetSelector() *metav1.LabelSelector {
	return r.Spec.Selector
}

func (r *EXTCOMMClaim) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      r.Kind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *EXTCOMMClaim) GetStaticID() *uint64 {
	if r.Spec.ID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.ID))
}

func (r *EXTCOMMClaim) GetStaticTreeID(s string) tree.ID {
	if r.Spec.ID == nil {
		return nil
	}
	switch GetExtendedCommunityType(s) {
	case ExtendedCommunityType_4byteAS, ExtendedCommunityType_IPv4Address:
		return id16.NewID(uint16(*r.Spec.ID), id16.IDBitSize)
	case ExtendedCommunityType_2byteAS:
		return id32.NewID(uint32(*r.Spec.ID), id32.IDBitSize)
	case ExtendedCommunityType_Opaque:
		return id64.NewID(uint64(*r.Spec.ID), id64.IDBitSize)
	default:
		return nil
	}
}

func (r *EXTCOMMClaim) GetClaimID(s string, id uint64) tree.ID {
	switch GetExtendedCommunityType(s) {
	case ExtendedCommunityType_4byteAS, ExtendedCommunityType_IPv4Address:
		return id16.NewID(uint16(id), id16.IDBitSize)
	case ExtendedCommunityType_2byteAS:
		return id32.NewID(uint32(id), id32.IDBitSize)
	case ExtendedCommunityType_Opaque:
		return id64.NewID(uint64(id), id64.IDBitSize)
	default:
		return nil
	}
}

func (r *EXTCOMMClaim) GetRange() *string {
	return r.Spec.Range
}

func (r *EXTCOMMClaim) GetRangeID(s string) (tree.Range, error) {
	if r.Spec.Range == nil {
		return nil, fmt.Errorf("cannot provide a range without an id")
	}
	switch GetExtendedCommunityType(s) {
	case ExtendedCommunityType_4byteAS, ExtendedCommunityType_IPv4Address:
		return id16.ParseRange(*r.Spec.Range)
	case ExtendedCommunityType_2byteAS:
		return id32.ParseRange(*r.Spec.Range)
	case ExtendedCommunityType_Opaque:
		return id64.ParseRange(*r.Spec.Range)
	default:
		return nil, fmt.Errorf("invalid extended communitytype: got %s", s)
	}

}

func (r *EXTCOMMClaim) GetTable(s string, to, from uint64) table.Table {
	switch GetExtendedCommunityType(s) {
	case ExtendedCommunityType_4byteAS, ExtendedCommunityType_IPv4Address:
		return table16.New(uint16(to), uint16(from))
	case ExtendedCommunityType_2byteAS:
		return table32.New(uint32(to), uint32(from))
	case ExtendedCommunityType_Opaque:
		return table64.New(to, from)
	default:
		return nil
	}
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
		return fmt.Errorf("a claim can only have 1 type, got %s", sb.String())
	}
	return nil
}

func validateEXTCOMMID(extCommType ExtendedCommunityType, id int64) error {
	if id < EXTCOMMID_Min {
		return fmt.Errorf("invalid id, got %d", id)
	}
	if id > int64(EXTCOMMID_MaxValue[extCommType]) {
		return fmt.Errorf("invalid id, got %d", id)
	}
	return nil
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

func (r *EXTCOMMClaim) GetEXTCOMMRange() (int, int) {
	if r.Spec.Range == nil {
		return 0, 0
	}
	parts := strings.SplitN(*r.Spec.Range, "-", 2)
	if len(parts) != 2 {
		return 0, 0
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0
	}
	return start, end
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
	if err := validateEXTCOMMID(extCommType, int64(start)); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid EXTCOMM start err %s", err.Error()))
	}
	if err := validateEXTCOMMID(extCommType, int64(end)); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid EXTCOMM end err %s", err.Error()))
	}
	return errm
}

func (r *EXTCOMMClaim) ValidateSyntax(s string) field.ErrorList {
	var allErrs field.ErrorList

	gv, err := schema.ParseGroupVersion(r.APIVersion)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("apiVersion"),
			r,
			fmt.Errorf("invalid apiVersion: err: %s", err.Error()).Error(),
		))
		return allErrs
	}

	// this is for user convenience
	if r.Spec.Owner == nil {
		r.Spec.Owner = &commonv1alpha1.OwnerReference{
			Group:     gv.Group,
			Version:   gv.Version,
			Kind:      r.Kind,
			Namespace: r.Namespace,
			Name:      r.Name,
		}
	}

	if GetExtendedCommunityType(s) == ExtendedCommunityType_Invalid {
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
	return v.Validate(r, GetExtendedCommunityType(s))
}

func (r *EXTCOMMClaim) ValidateOwner(labels labels.Set) error {
	routeClaimName := labels[backend.KuidClaimNameKey]
	routeOwner := commonv1alpha1.OwnerReference{
		Group:     labels[backend.KuidOwnerGroupKey],
		Version:   labels[backend.KuidOwnerVersionKey],
		Kind:      labels[backend.KuidOwnerKindKey],
		Namespace: labels[backend.KuidOwnerNamespaceKey],
		Name:      labels[backend.KuidOwnerNameKey],
	}
	if (r.Spec.Owner != nil && *r.Spec.Owner != routeOwner) || r.Name != routeClaimName {
		return fmt.Errorf("route owned by different claim got name %s/%s owner %s/%s",
			r.Name,
			routeClaimName,
			r.Spec.Owner.String(),
			routeOwner.String(),
		)
	}
	return nil
}

// GetLabelSelector returns a labels selector based on the label selector
func (r *EXTCOMMClaim) GetLabelSelector() (labels.Selector, error) {
	return r.Spec.GetLabelSelector()
}

func (r *EXTCOMMClaim) GetClaimLabels() labels.Set {
	labels := r.Spec.GetUserDefinedLabels()
	// system defined labels
	labels[backend.KuidClaimTypeKey] = string(r.GetClaimType())
	labels[backend.KuidClaimNameKey] = r.Name
	labels[backend.KuidOwnerGroupKey] = r.Spec.Owner.Group
	labels[backend.KuidOwnerVersionKey] = r.Spec.Owner.Version
	labels[backend.KuidOwnerKindKey] = r.Spec.Owner.Kind
	labels[backend.KuidOwnerNamespaceKey] = r.Spec.Owner.Namespace
	labels[backend.KuidOwnerNameKey] = r.Spec.Owner.Name
	return labels
}

// GetOwnerSelector returns a label selector to select the owner of the claim in the backend
func (r *EXTCOMMClaim) GetOwnerSelector() (labels.Selector, error) {
	l := map[string]string{
		backend.KuidOwnerGroupKey:     r.Spec.Owner.Group,
		backend.KuidOwnerVersionKey:   r.Spec.Owner.Version,
		backend.KuidOwnerKindKey:      r.Spec.Owner.Kind,
		backend.KuidOwnerNamespaceKey: r.Spec.Owner.Namespace,
		backend.KuidOwnerNameKey:      r.Spec.Owner.Name,
		backend.KuidClaimNameKey:      r.Name,
	}

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

// BuildEXTCOMMClaim returns a reource from a client Object a Spec/Status
func BuildEXTCOMMClaim(meta metav1.ObjectMeta, spec *EXTCOMMClaimSpec, status *EXTCOMMClaimStatus) *EXTCOMMClaim {
	aspec := EXTCOMMClaimSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := EXTCOMMClaimStatus{}
	if status != nil {
		astatus = *status
	}
	return &EXTCOMMClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       EXTCOMMClaimKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func (r *EXTCOMMClaim) SetStatusRange(s *string) {
	r.Status.Range = s
}

func (r *EXTCOMMClaim) SetStatusID(s *uint64) {
	if s == nil {
		r.Status.ID = nil
		return
	}
	r.Status.ID = ptr.To[int64](int64(*s))
}

func (r *EXTCOMMClaim) GetStatusID() *uint64 {
	if r.Status.ID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Status.ID))
}
