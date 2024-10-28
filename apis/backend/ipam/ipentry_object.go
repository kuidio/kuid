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

package ipam

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	"github.com/henderiw/iputil"
	"github.com/henderiw/logger/log"
	"github.com/henderiw/store"
	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/backend"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

func (r *IPEntry) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}
func (r *IPEntry) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

// GetCondition returns the condition based on the condition kind
func (r *IPEntry) GetCondition(t condition.ConditionType) condition.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *IPEntry) SetConditions(c ...condition.Condition) {
	r.Status.SetConditions(c...)
}

func (r *IPEntry) ValidateSyntax(s string) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

func (r *IPEntry) GetClaimSummaryType() IPClaimSummaryType {
	switch r.Spec.ClaimType {
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

func (r *IPEntry) GetIPPrefixType() IPPrefixType {
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

func (r *IPEntry) GetIndex() string {
	return r.Spec.Index
}

func (r *IPEntry) GetIPPrefix() string {
	return r.Spec.Prefix
}

func (r *IPEntry) GetSpecID() string {
	return r.Spec.Prefix
}

func IPEntryFromRuntime(ru runtime.Object) (*IPEntry, error) {
	entry, ok := ru.(*IPEntry)
	if !ok {
		return nil, errors.New("runtime object not ASIndex")
	}
	return entry, nil
}

func GetIPEntry(ctx context.Context, k store.Key, rangeName string, prefix netip.Prefix, labels map[string]string) *IPEntry {
	log := log.FromContext(ctx)
	log.Debug("get ipEntry from iptables", "key", k.String(), "range", rangeName, "prefix", prefix, "labels", labels)

	pi := iputil.NewPrefixInfo(prefix)

	index := k.Name
	ns := k.Namespace

	// indicates if the IP entry is originated from the ip index
	indexEntry := false
	if labels[backend.KuidIndexEntryKey] == "true" {
		indexEntry = true
	}

	spec := &IPEntrySpec{
		Index:         index,
		IndexEntry:    indexEntry,
		PrefixType:    GetIPPrefixTypeFromString(labels[backend.KuidIPAMIPPrefixTypeKey]),
		ClaimType:     GetIPClaimTypeFromString(labels[backend.KuidClaimTypeKey]),
		Prefix:        pi.String(),
		AddressFamily: ptr.To(pi.GetAddressFamily()),
	}
	if _, ok := labels[backend.KuidIPAMDefaultGatewayKey]; ok {
		spec.DefaultGateway = ptr.To[bool](true)
	}
	// filter the system defined labels from the labels to prepare for the user defined labels
	udLabels := map[string]string{}
	for k, v := range labels {
		if !backend.BackendSystemKeys.Has(k) && !backend.BackendIPAMSystemKeys.Has(k) {
			udLabels[k] = v
		}
	}
	spec.UserDefinedLabels.Labels = udLabels

	name := fmt.Sprintf("%s.%s.%s", ns, index, pi.GetSubnetName())
	if rangeName != "" {
		name = fmt.Sprintf("%s.%s.%s.%s", ns, index, rangeName, pi.GetSubnetName())
	}
	return BuildIPEntry(
		metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind:       labels[backend.KuidOwnerKindKey],
					Name:       labels[backend.KuidClaimNameKey],
					UID:        types.UID(labels[backend.KuidClaimUIDKey]),
				},
			},
		},
		spec,
		nil,
	)
}

// BuildIPEntry returns a reource from a client Object a Spec/Status
func BuildIPEntry(meta metav1.ObjectMeta, spec *IPEntrySpec, status *IPEntryStatus) *IPEntry {
	aspec := IPEntrySpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := IPEntryStatus{}
	if status != nil {
		astatus = *status
	}
	return &IPEntry{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       IPEntryKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}
