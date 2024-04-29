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
	"strings"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/iputil"
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

const IPClaimPlural = "ipclaims"
const IPClaimSingular = "ipclaim"

// +k8s:deepcopy-gen=false
var _ resource.Object = &IPClaim{}
var _ resource.ObjectList = &IPClaimList{}

var _ resource.ObjectWithStatusSubResource = &IPClaim{}

func (IPClaimStatus) SubResourceName() string {
	return fmt.Sprintf("%s/%s", IPClaimPlural, "status")
}

func (r IPClaimStatus) CopyTo(obj resource.ObjectWithStatusSubResource) {
	cfg, ok := obj.(*IPClaim)
	if ok {
		cfg.Status = r
	}
}

func (r *IPClaim) GetStatus() resource.StatusSubResource {
	return r.Status
}

// GetListMeta returns the ListMeta
func (r *IPClaimList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *IPClaim) GetSingularName() string {
	return IPClaimSingular
}

func (IPClaim) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: IPClaimPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (IPClaim) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *IPClaim) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (IPClaim) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (IPClaim) New() runtime.Object {
	return &IPClaim{}
}

// NewList implements resource.Object
func (IPClaim) NewList() runtime.Object {
	return &IPClaimList{}
}

// GetCondition returns the condition based on the condition kind
func (r *IPClaim) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *IPClaim) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// ConvertIPClaimFieldSelector is the schema conversion function for normalizing the FieldSelector for IPClaim
func ConvertIPClaimFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	case "spec.networkInstance":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *IPClaimList) GetItems() []rresource.Object {
	objs := []rresource.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *IPClaim) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *IPClaim) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *IPClaim) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.NetworkInstance})
}

func (r *IPClaim) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      r.Kind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *IPClaim) GetIPPrefixType() IPPrefixType {
	if r.Spec.PrefixType == nil {
		return IPPrefixType_Other
	}
	switch *r.Spec.PrefixType {
	case IPPrefixType_Aggregate, IPPrefixType_Network, IPPrefixType_Pool:
		return *r.Spec.PrefixType
	default:
		return IPPrefixType_Invalid
	}
}

func (r *IPClaim) GetClaimRequest() string {
	// we assume validation is already done when calling this
	if r.Spec.Address != nil {
		return *r.Spec.Address
	}
	if r.Spec.Prefix != nil {
		return *r.Spec.Prefix
	}
	if r.Spec.Range != nil {
		return *r.Spec.Range
	}
	return ""
}

func (r *IPClaim) GetClaimResponse() string {
	// we assume validation is already done when calling this
	if r.Status.Address != nil {
		return *r.Status.Address
	}
	if r.Status.Prefix != nil {
		return *r.Status.Prefix
	}
	if r.Status.Range != nil {
		return *r.Status.Range
	}
	return ""
}

func (r *IPClaim) GetIPClaimType() (IPClaimType, error) {
	claimType := IPClaimType_Invalid
	var sb strings.Builder
	count := 0
	if r.Spec.Address != nil {
		sb.WriteString(fmt.Sprintf("address: %s", *r.Spec.Address))
		claimType = IPClaimType_StaticAddress
		count++

	}
	if r.Spec.Prefix != nil {
		if count > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("prefix: %s", *r.Spec.Prefix))
		claimType = IPClaimType_StaticPrefix
		count++

	}
	if r.Spec.Range != nil {
		if count > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("range: %s", *r.Spec.Range))
		claimType = IPClaimType_StaticRange
		count++
	}
	if count > 1 {
		return IPClaimType_Invalid, fmt.Errorf("an ipclaim can only have 1 claimType, got %s", sb.String())
	}
	if count == 0 {
		if r.Spec.CreatePrefix != nil {
			return IPClaimType_DynamicPrefix, nil
		} else {
			return IPClaimType_DynamicAddress, nil
		}
	}
	return claimType, nil
}

func (r *IPClaim) GetIPClaimSummaryType() IPClaimSummaryType {
	ipClaimType, err := r.GetIPClaimType()
	if err != nil {
		return IPClaimSummaryType_Invalid
	}
	switch ipClaimType {
	case IPClaimType_DynamicAddress, IPClaimType_StaticAddress:
		return IPClaimSummaryType_Address
	case IPClaimType_DynamicPrefix, IPClaimType_StaticPrefix:
		return IPClaimSummaryType_Prefix
	case IPClaimType_StaticRange:
		return IPClaimSummaryType_Range
	default:
		return IPClaimSummaryType_Invalid
	}
}

func (r *IPClaim) ValidateSyntax() field.ErrorList {
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

	ipClaimType, err := r.GetIPClaimType()
	if err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath(""),
			r,
			err.Error(),
		))
		return allErrs
	}
	var v SyntaxValidator
	switch ipClaimType {
	case IPClaimType_StaticAddress:
		v = &staticAddressSyntaxValidator{name: "staticIPAddress"}
	case IPClaimType_StaticPrefix:
		v = &staticPrefixSyntaxValidator{name: "staticIPprefix"}
	case IPClaimType_StaticRange:
		v = &staticRangeSyntaxValidator{name: "staticIPRange"}
	case IPClaimType_DynamicAddress:
		v = &dynamicAddressSyntaxValidator{name: "dynamicIPRange"}
	case IPClaimType_DynamicPrefix:
		v = &dynamicPrefixSyntaxValidator{name: "dynamicIPprefix"}
	default:
		return allErrs
	}
	return v.Validate(r)
}

func (r *IPClaim) ValidateOwner(labels labels.Set) error {
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

// GetDummyLabelsFromPrefix returns a map with the labels from the spec
// augmented with the prefixkind and the subnet from the prefixInfo
func (r *IPClaim) GetDummyLabelsFromPrefix(pi iputil.Prefix) map[string]string {
	labels := map[string]string{}
	for k, v := range r.Spec.GetUserDefinedLabels() {
		labels[k] = v
	}
	labels[backend.KuidIPAMIPPrefixTypeKey] = string(r.GetIPPrefixType())
	labels[backend.KuidIPAMClaimSummaryTypeKey] = string(r.GetIPClaimSummaryType())
	labels[backend.KuidIPAMSubnetKey] = string(pi.GetSubnetName())

	return labels
}

// GetLabelSelector returns a labels selector based on the label selector
func (r *IPClaim) GetLabelSelector() (labels.Selector, error) {
	return r.Spec.GetLabelSelector()
}

// GetOwnerSelector returns a label selector to select the owner of the claim in the backend
func (r *IPClaim) GetOwnerSelector() (labels.Selector, error) {
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

// GetGatewayLabelSelector returns a label selector to select the gateway of the claim in the backend
func (r *IPClaim) GetDefaultGatewayLabelSelector(subnetString string) (labels.Selector, error) {
	l := map[string]string{
		backend.KuidIPAMDefaultGatewayKey: "true",
		backend.KuidIPAMSubnetKey:         subnetString,
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

// BuildIPClaim returns a reource from a client Object a Spec/Status
func BuildIPClaim(meta metav1.ObjectMeta, spec *IPClaimSpec, status *IPClaimStatus) *IPClaim {
	aspec := IPClaimSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := IPClaimStatus{}
	if status != nil {
		astatus = *status
	}
	return &IPClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       IPClaimKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}
