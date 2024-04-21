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
	"fmt"
	"strings"

	"github.com/henderiw/iputil"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	"github.com/kuidio/kuid/pkg/reconcilers/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

// GetCondition returns the condition based on the condition kind
func (r *IP) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *IP) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

func (r *IPList) GetItems() []resource.Object {
	objs := []resource.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *IP) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *IP) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      r.Kind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *IP) GetAddressing() (ipambev1alpha1.IPClaimType, error) {
	addressing := ipambev1alpha1.IPClaimType_Invalid
	var sb strings.Builder
	count := 0
	if r.Spec.Address != nil {
		sb.WriteString(fmt.Sprintf("address: %s", *r.Spec.Address))
		addressing = ipambev1alpha1.IPClaimType_StaticAddress
		count++

	}
	if r.Spec.Prefix != nil {
		if count > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("prefix: %s", *r.Spec.Prefix))
		addressing = ipambev1alpha1.IPClaimType_StaticPrefix
		count++

	}
	if r.Spec.Range != nil {
		if count > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("range: %s", *r.Spec.Range))
		addressing = ipambev1alpha1.IPClaimType_StaticRange
		count++
	}
	if count > 1 {
		return ipambev1alpha1.IPClaimType_Invalid, fmt.Errorf("an ipclaim can only have 1 addressing, got %s", sb.String())
	}
	return addressing, nil
}

func (r *IP) GetIPClaimSummaryType() ipambev1alpha1.IPClaimSummaryType {
	addressing, err := r.GetAddressing()
	if err != nil {
		return ipambev1alpha1.IPClaimSummaryType_Invalid
	}
	switch addressing {
	case ipambev1alpha1.IPClaimType_DynamicAddress, ipambev1alpha1.IPClaimType_StaticAddress:
		return ipambev1alpha1.IPClaimSummaryType_Address
	case ipambev1alpha1.IPClaimType_DynamicPrefix, ipambev1alpha1.IPClaimType_StaticPrefix:
		return ipambev1alpha1.IPClaimSummaryType_Prefix
	case ipambev1alpha1.IPClaimType_StaticRange:
		return ipambev1alpha1.IPClaimSummaryType_Range
	default:
		return ipambev1alpha1.IPClaimSummaryType_Invalid
	}
}

func (r *IP) GetIPPrefixType() *ipambev1alpha1.IPPrefixType {
	if r.Spec.PrefixType == nil {
		return ptr.To[ipambev1alpha1.IPPrefixType](ipambev1alpha1.IPPrefixType_Other)
	}
	switch *r.Spec.PrefixType {
	case ipambev1alpha1.IPPrefixType_Network, ipambev1alpha1.IPPrefixType_Pool:
		return r.Spec.PrefixType
	default:
		return ptr.To[ipambev1alpha1.IPPrefixType](ipambev1alpha1.IPPrefixType_Invalid)
	}
}

func BuildIP(meta metav1.ObjectMeta, spec *IPSpec, status *IPStatus) *IP {
	aspec := IPSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := IPStatus{}
	if status != nil {
		astatus = *status
	}
	return &IP{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       IPKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func (r *IP) GetIPClaim() (*ipambev1alpha1.IPClaim, error) {
	var name string
	var spec *ipambev1alpha1.IPClaimSpec
	switch r.GetIPClaimSummaryType() {
	case ipambev1alpha1.IPClaimSummaryType_Address:
		pi, err := iputil.New(*r.Spec.Address)
		if err != nil {
			return nil, fmt.Errorf("invalid IP Address, err: %s", err.Error())
		}
		name = pi.GetSubnetName()
		spec = r.getIPClaimAddressSpec()
	case ipambev1alpha1.IPClaimSummaryType_Prefix:
		pi, err := iputil.New(*r.Spec.Prefix)
		if err != nil {
			return nil, fmt.Errorf("invalid IP Address, err: %s", err.Error())
		}
		name = pi.GetSubnetName()
		spec = r.getIPClaimPrefixSpec(pi)
	case ipambev1alpha1.IPClaimSummaryType_Range:
		name = r.Name
		spec = r.getIPClaimRangeSpec()
	default:
		return nil, fmt.Errorf("invalid IP")
	}

	return ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      name,
		},
		spec,
		nil,
	), nil
}

func (r *IP) getIPClaimPrefixSpec(pi *iputil.Prefix) *ipambev1alpha1.IPClaimSpec {
	return &ipambev1alpha1.IPClaimSpec{
		PrefixType:      r.GetIPPrefixType(),
		NetworkInstance: r.Spec.NetworkInstance,
		Prefix:          r.Spec.Prefix,
		PrefixLength:    ptr.To[uint32](uint32(pi.GetPrefixLength())),
		CreatePrefix:    ptr.To[bool](true),
		DefaultGateway:  r.Spec.DefaultGateway,
		ClaimLabels: commonv1alpha1.ClaimLabels{
			UserDefinedLabels: r.Spec.UserDefinedLabels,
		},
		Owner: commonv1alpha1.GetOwnerReference(r),
	}
}

func (r *IP) getIPClaimAddressSpec() *ipambev1alpha1.IPClaimSpec {
	return &ipambev1alpha1.IPClaimSpec{
		PrefixType:      nil,
		NetworkInstance: r.Spec.NetworkInstance,
		Address:         r.Spec.Address,
		ClaimLabels: commonv1alpha1.ClaimLabels{
			UserDefinedLabels: r.Spec.UserDefinedLabels,
		},
		Owner: commonv1alpha1.GetOwnerReference(r),
	}
}

func (r *IP) getIPClaimRangeSpec() *ipambev1alpha1.IPClaimSpec {
	return &ipambev1alpha1.IPClaimSpec{
		PrefixType:      nil,
		NetworkInstance: r.Spec.NetworkInstance,
		Range:           r.Spec.Range,
		ClaimLabels: commonv1alpha1.ClaimLabels{
			UserDefinedLabels: r.Spec.UserDefinedLabels,
		},
		Owner: commonv1alpha1.GetOwnerReference(r),
	}
}
