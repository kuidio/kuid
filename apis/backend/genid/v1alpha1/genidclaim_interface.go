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

const GENIDClaimPlural = "genidclaims"
const GENIDClaimSingular = "genidclaim"
const GENIDID_Min = 0

var GENIDID_MaxBits = map[GENIDType]int{
	GENIDType_Invalid: 0,
	GENIDType_16bit:   16,
	GENIDType_32bit:   32,
	GENIDType_48bit:   48,
	GENIDType_64bit:   63, // workaround for apiserver issue with uint64
}

var GENIDID_MaxValue = map[GENIDType]int64{
	GENIDType_Invalid: 0,
	GENIDType_16bit:   1<<GENIDID_MaxBits[GENIDType_16bit] - 1,
	GENIDType_32bit:   1<<GENIDID_MaxBits[GENIDType_32bit] - 1,
	GENIDType_48bit:   1<<GENIDID_MaxBits[GENIDType_48bit] - 1,
	GENIDType_64bit:   1<<GENIDID_MaxBits[GENIDType_64bit] - 1,
}

// +k8s:deepcopy-gen=false
var _ resource.Object = &GENIDClaim{}
var _ resource.ObjectList = &GENIDClaimList{}

var _ resource.ObjectWithStatusSubResource = &GENIDClaim{}

func (GENIDClaimStatus) SubResourceName() string {
	return fmt.Sprintf("%s/%s", GENIDClaimPlural, "status")
}

func (r GENIDClaimStatus) CopyTo(obj resource.ObjectWithStatusSubResource) {
	cfg, ok := obj.(*GENIDClaim)
	if ok {
		cfg.Status = r
	}
}

func (r *GENIDClaim) GetStatus() resource.StatusSubResource {
	return r.Status
}

// GetListMeta returns the ListMeta
func (r *GENIDClaimList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *GENIDClaim) GetSingularName() string {
	return GENIDClaimSingular
}

func (GENIDClaim) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: GENIDClaimPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (GENIDClaim) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *GENIDClaim) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (GENIDClaim) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (GENIDClaim) New() runtime.Object {
	return &GENIDClaim{}
}

// NewList implements resource.Object
func (GENIDClaim) NewList() runtime.Object {
	return &GENIDClaimList{}
}

// GetCondition returns the condition based on the condition kind
func (r *GENIDClaim) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *GENIDClaim) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// GENIDConvertClaimFieldSelector is the schema conversion function for normalizing the FieldSelector for GENIDClaim
func GENIDClaimConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
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

func (r *GENIDClaimList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *GENIDClaim) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *GENIDClaim) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *GENIDClaim) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

func (r *GENIDClaim) GetIndex() string {
	return r.Spec.Index
}

func (r *GENIDClaim) GetSelector() *metav1.LabelSelector {
	return r.Spec.Selector
}

func (r *GENIDClaim) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      r.Kind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *GENIDClaim) GetStaticID() *uint64 {
	if r.Spec.ID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.ID))
}

func (r *GENIDClaim) GetStaticTreeID(s string) tree.ID {
	if r.Spec.ID == nil {
		return nil
	}
	switch GetGenIDType(s) {
	case GENIDType_16bit:
		return id16.NewID(uint16(*r.Spec.ID), id16.IDBitSize)
	case GENIDType_32bit:
		return id32.NewID(uint32(*r.Spec.ID), id32.IDBitSize)
	case GENIDType_48bit:
		return id64.NewID(uint64(*r.Spec.ID), 48)
	case GENIDType_64bit:
		return id64.NewID(uint64(*r.Spec.ID), 48)
	default:
		return nil
	}
}

func (r *GENIDClaim) GetClaimID(s string, id uint64) tree.ID {
	switch GetGenIDType(s) {
	case GENIDType_16bit:
		return id16.NewID(uint16(*r.Spec.ID), id16.IDBitSize)
	case GENIDType_32bit:
		return id32.NewID(uint32(*r.Spec.ID), id32.IDBitSize)
	case GENIDType_48bit:
		return id64.NewID(uint64(*r.Spec.ID), 48)
	case GENIDType_64bit:
		return id64.NewID(uint64(*r.Spec.ID), 48)
	default:
		return nil
	}
}

func (r *GENIDClaim) GetRange() *string {
	return r.Spec.Range
}

func (r *GENIDClaim) GetRangeID(s string) (tree.Range, error) {
	if r.Spec.Range == nil {
		return nil, fmt.Errorf("cannot provide a range without an id")
	}
	switch GetGenIDType(s) {
	case GENIDType_16bit:
		return id16.ParseRange(*r.Spec.Range)
	case GENIDType_32bit:
		return id32.ParseRange(*r.Spec.Range)
	case GENIDType_48bit:
		return id64.ParseRange(*r.Spec.Range)
	case GENIDType_64bit:
		return id64.ParseRange(*r.Spec.Range)
	default:
		return nil, fmt.Errorf("invalid type: got %s", s)
	}

}

func (r *GENIDClaim) GetTable(s string, to, from uint64) table.Table {
	switch GetGenIDType(s) {
	case GENIDType_16bit:
		return table16.New(uint16(to), uint16(from))
	case GENIDType_32bit:
		return table32.New(uint32(to), uint32(from))
	case GENIDType_48bit:
		return table64.New(uint64(to), uint64(from))
	case GENIDType_64bit:
		return table64.New(uint64(to), uint64(from))
	default:
		return nil
	}
}

func (r *GENIDClaim) GetClaimRequest() string {
	// we assume validation is already done when calling this
	if r.Spec.ID != nil {
		return strconv.Itoa(int(*r.Spec.ID))
	}
	if r.Spec.Range != nil {
		return *r.Spec.Range
	}
	return ""
}

func (r *GENIDClaim) GetClaimResponse() string {
	// we assume validation is already done when calling this
	if r.Status.ID != nil {
		return strconv.Itoa(int(*r.Status.ID))
	}
	if r.Status.Range != nil {
		return *r.Status.Range
	}
	return ""
}

func (r *GENIDClaim) GetClaimType() backend.ClaimType {
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

func (r *GENIDClaim) ValidateGENIDClaimType() error {
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

func validateGENIDID(genidType GENIDType, id int64) error {
	if id < GENIDID_Min {
		return fmt.Errorf("invalid id, got %d", id)
	}
	if id > GENIDID_MaxValue[genidType] {
		return fmt.Errorf("invalid id, got %d", id)
	}
	return nil
}

func (r *GENIDClaim) ValidateGENIDID(genidType GENIDType) error {
	if r.Spec.ID == nil {
		return fmt.Errorf("no id provided")
	}
	if err := validateGENIDID(genidType, *r.Spec.ID); err != nil {
		return fmt.Errorf("invalid id err %s", err.Error())
	}
	return nil
}

func (r *GENIDClaim) GetGENIDRange() (int, int) {
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

func (r *GENIDClaim) ValidateGENIDRange(genidType GENIDType) error {
	if r.Spec.Range == nil {
		return fmt.Errorf("no GENID range provided")
	}
	parts := strings.SplitN(*r.Spec.Range, "-", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid GENID range, expected <start>-<end>, got: %s", *r.Spec.Range)
	}
	var errm error
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid GENID range start, got: %s, err: %s", *r.Spec.Range, err.Error()))
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid GENID range end, got: %s, err: %s", *r.Spec.Range, err.Error()))
	}
	if errm != nil {
		return errm
	}
	if start > end {
		errm = errors.Join(errm, fmt.Errorf("invalid GENID range start > end %s", *r.Spec.Range))
	}
	if err := validateGENIDID(genidType, int64(start)); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid GENID start err %s", err.Error()))
	}
	if err := validateGENIDID(genidType, int64(end)); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid GENID end err %s", err.Error()))
	}
	return errm
}

func (r *GENIDClaim) ValidateSyntax(s string) field.ErrorList {
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

	if err := r.ValidateGENIDClaimType(); err != nil {
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
		v = &GENIDDynamicIDSyntaxValidator{name: string(claimType)}
	case backend.ClaimType_StaticID:
		v = &GENIDStaticIDSyntaxValidator{name: string(claimType)}
	case backend.ClaimType_Range:
		v = &GENIDRangeSyntaxValidator{name: string(claimType)}
	default:
		return allErrs
	}
	return v.Validate(r, GetGenIDType(s))
}

func (r *GENIDClaim) ValidateOwner(labels labels.Set) error {
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
func (r *GENIDClaim) GetLabelSelector() (labels.Selector, error) {
	return r.Spec.GetLabelSelector()
}

func (r *GENIDClaim) GetClaimLabels() labels.Set {
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
func (r *GENIDClaim) GetOwnerSelector() (labels.Selector, error) {
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

// BuildGENIDClaim returns a reource from a client Object a Spec/Status
func BuildGENIDClaim(meta metav1.ObjectMeta, spec *GENIDClaimSpec, status *GENIDClaimStatus) *GENIDClaim {
	aspec := GENIDClaimSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := GENIDClaimStatus{}
	if status != nil {
		astatus = *status
	}
	return &GENIDClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       GENIDClaimKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func (r *GENIDClaim) SetStatusRange(s *string) {
	r.Status.Range = s
}

func (r *GENIDClaim) SetStatusID(s *uint64) {
	if s == nil {
		r.Status.ID = nil
		return
	}
	r.Status.ID = ptr.To[int64](int64(*s))
}

func (r *GENIDClaim) GetStatusID() *uint64 {
	if r.Status.ID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Status.ID))
}

func (r *GENIDClaim) GetSpec() any {
	return r.Spec
}

func (r *GENIDClaim) SetSpec(s any) {
	if spec, ok := s.(GENIDClaimSpec); ok {
		r.Spec = spec
	}
}

func (r *GENIDClaim) NewObjList() backend.GenericObjectList {
	return &GENIDClaimList{
		TypeMeta: metav1.TypeMeta{APIVersion: SchemeGroupVersion.Identifier(), Kind: GENIDClaimListKind},
	}
}

func (r *GENIDClaimList) GetObjects() []backend.GenericObject {
	objs := []backend.GenericObject{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}
