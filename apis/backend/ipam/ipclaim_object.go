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

package ipam

import (
	"fmt"
	"strings"

	"github.com/henderiw/store"
	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/backend"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func (r *IPClaim) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *IPClaim) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

// GetCondition returns the condition based on the condition kind
func (r *IPClaim) GetCondition(t condition.ConditionType) condition.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *IPClaim) SetConditions(c ...condition.Condition) {
	r.Status.SetConditions(c...)
}

func (r *IPClaim) ValidateSyntax(s string) field.ErrorList {
	var allErrs field.ErrorList

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

func (r *IPClaim) GetIndex() string { return r.Spec.Index }

func (r *IPClaim) GetSelector() *metav1.LabelSelector { return r.Spec.Selector }

// GetOwnerSelector selects the route based on the name of the claim
func (r *IPClaim) GetOwnerSelector() (labels.Selector, error) {
	l := map[string]string{
		backend.KuidClaimNameKey: r.Name,
		backend.KuidClaimUIDKey:  string(r.UID),
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

func (r *IPClaim) GetLabelSelector() (labels.Selector, error) { return r.Spec.GetLabelSelector() }

func (r *IPClaim) GetClaimLabels() labels.Set {
	labels := r.Spec.GetUserDefinedLabels()

	// system defined labels
	labels[backend.KuidClaimTypeKey] = string(r.GetIPClaimSummaryType())
	labels[backend.KuidIPAMIPPrefixTypeKey] = string(r.GetIPPrefixType())
	labels[backend.KuidClaimNameKey] = r.Name
	labels[backend.KuidClaimUIDKey] = string(r.UID)
	return labels
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

func (r *IPClaim) ValidateOwner(labels labels.Set) error {
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

// IsOwnedByIPIndex returns true if the owner is the IPIndex
func (r *IPClaim) IsOwnedByIPIndex() bool {
	for _, ownerRef := range r.OwnerReferences {
		if strings.HasPrefix(ownerRef.APIVersion, SchemeGroupVersion.Group) &&
			ownerRef.Kind == IPIndexKind {
			return true
		}
	}
	return false
}
