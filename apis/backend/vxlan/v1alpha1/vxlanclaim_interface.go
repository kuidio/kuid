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
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	rresource "github.com/kuidio/kuid/pkg/reconcilers/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const VXLANClaimPlural = "vxlanclaims"
const VXLANClaimSingular = "vxlanclaim"

// +k8s:deepcopy-gen=false
var _ resource.Object = &VXLANClaim{}
var _ resource.ObjectList = &VXLANClaimList{}

var _ resource.ObjectWithStatusSubResource = &VXLANClaim{}

func (VXLANClaimStatus) SubResourceName() string {
	return fmt.Sprintf("%s/%s", VXLANClaimPlural, "status")
}

func (r VXLANClaimStatus) CopyTo(obj resource.ObjectWithStatusSubResource) {
	cfg, ok := obj.(*VXLANClaim)
	if ok {
		cfg.Status = r
	}
}

func (r *VXLANClaim) GetStatus() resource.StatusSubResource {
	return r.Status
}

// GetListMeta returns the ListMeta
func (r *VXLANClaimList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *VXLANClaim) GetSingularName() string {
	return VXLANClaimSingular
}

func (VXLANClaim) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: VXLANClaimPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (VXLANClaim) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *VXLANClaim) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (VXLANClaim) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (VXLANClaim) New() runtime.Object {
	return &VXLANClaim{}
}

// NewList implements resource.Object
func (VXLANClaim) NewList() runtime.Object {
	return &VXLANClaimList{}
}

// GetCondition returns the condition based on the condition kind
func (r *VXLANClaim) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *VXLANClaim) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// ConvertVXLANClaimFieldSelector is the schema conversion function for normalizing the FieldSelector for VXLANClaim
func ConvertVXLANClaimFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
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

func (r *VXLANClaimList) GetItems() []rresource.Object {
	objs := []rresource.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *VXLANClaim) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *VXLANClaim) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *VXLANClaim) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

func (r *VXLANClaim) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      r.Kind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *VXLANClaim) GetClaimRequest() string {
	// we assume validation is already done when calling this
	if r.Spec.ID != nil {
		return strconv.Itoa(int(*r.Spec.ID))
	}
	if r.Spec.Range != nil {
		return *r.Spec.Range
	}
	return ""
}

func (r *VXLANClaim) GetClaimResponse() string {
	// we assume validation is already done when calling this
	if r.Status.ID != nil {
		return strconv.Itoa(int(*r.Status.ID))
	}
	if r.Status.Range != nil {
		return *r.Status.Range
	}
	return ""
}

func (r *VXLANClaim) GetClaimType() VXLANClaimType {
	claimType := VXLANClaimType_Invalid
	count := 0
	if r.Spec.ID != nil {
		claimType = VXLANClaimType_StaticID
		count++

	}
	if r.Spec.Range != nil {
		claimType = VXLANClaimType_Range
		count++

	}
	if count > 1 {
		return VXLANClaimType_Invalid
	}
	if count == 0 {
		return VXLANClaimType_DynamicID
	}
	return claimType
}

func (r *VXLANClaim) ValidateVXLANClaimType() error {
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
		return fmt.Errorf("an ipclaim can only have 1 addressing, got %s", sb.String())
	}
	return nil
}

func validateVXLANID(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid VXLAN id, got %d", id)
	}
	if id >= 4095 {
		return fmt.Errorf("invalid VXLAN id, got %d", id)
	}
	return nil
}

func (r *VXLANClaim) ValidateVXLANID() error {
	if r.Spec.ID == nil {
		return fmt.Errorf("no VXLAN id provided")
	}
	if err := validateVXLANID(int(*r.Spec.ID)); err != nil {
		return fmt.Errorf("invalid VXLAN id err %s", err.Error())
	}
	return nil
}

func (r *VXLANClaim) GetVXLANRange() (int, int) {
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

func (r *VXLANClaim) ValidateVXLANRange() error {
	if r.Spec.Range == nil {
		return fmt.Errorf("no VXLAN range provided")
	}
	parts := strings.SplitN(*r.Spec.Range, "-", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid VXLAN range, expected <start>-<end>, got: %s", *r.Spec.Range)
	}
	var errm error
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid VXLAN range start, got: %s, err: %s", *r.Spec.Range, err.Error()))
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid VXLAN range end, got: %s, err: %s", *r.Spec.Range, err.Error()))
	}
	if errm != nil {
		return errm
	}
	if start >= end {
		errm = errors.Join(errm, fmt.Errorf("invalid VXLAN range start >= end %s", *r.Spec.Range))
	}
	if err := validateVXLANID(start); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid VXLAN start err %s", err.Error()))
	}
	if err := validateVXLANID(end); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid VXLAN end err %s", err.Error()))
	}
	return errm
}

func (r *VXLANClaim) ValidateSyntax() field.ErrorList {
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

	if err := r.ValidateVXLANClaimType(); err != nil {
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
	case VXLANClaimType_DynamicID:
		v = &VXLANDynamicIDSyntaxValidator{name: string(claimType)}
	case VXLANClaimType_StaticID:
		v = &VXLANStaticIDSyntaxValidator{name: string(claimType)}
	case VXLANClaimType_Range:
		v = &VXLANRangeSyntaxValidator{name: string(claimType)}
	default:
		return allErrs
	}
	return v.Validate(r)
}

func (r *VXLANClaim) ValidateOwner(labels labels.Set) error {
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
func (r *VXLANClaim) GetLabelSelector() (labels.Selector, error) {
	return r.Spec.GetLabelSelector()
}

func (r *VXLANClaim) GetClaimLabels() labels.Set {
	labels := r.Spec.GetUserDefinedLabels()
	// system defined labels
	labels[backend.KuidVXLANClaimTypeKey] = string(r.GetClaimType())
	labels[backend.KuidClaimNameKey] = r.Name
	labels[backend.KuidOwnerGroupKey] = r.Spec.Owner.Group
	labels[backend.KuidOwnerVersionKey] = r.Spec.Owner.Version
	labels[backend.KuidOwnerKindKey] = r.Spec.Owner.Kind
	labels[backend.KuidOwnerNamespaceKey] = r.Spec.Owner.Namespace
	labels[backend.KuidOwnerNameKey] = r.Spec.Owner.Name
	return labels
}

// GetOwnerSelector returns a label selector to select the owner of the claim in the backend
func (r *VXLANClaim) GetOwnerSelector() (labels.Selector, error) {
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

// BuildVXLANClaim returns a reource from a client Object a Spec/Status
func BuildVXLANClaim(meta metav1.ObjectMeta, spec *VXLANClaimSpec, status *VXLANClaimStatus) *VXLANClaim {
	aspec := VXLANClaimSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := VXLANClaimStatus{}
	if status != nil {
		astatus = *status
	}
	return &VXLANClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       VXLANClaimKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}
