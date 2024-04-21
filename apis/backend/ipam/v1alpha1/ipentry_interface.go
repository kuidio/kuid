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
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/netip"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/iputil"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

const IPEntryPlural = "ipentries"
const IPEntrySingular = "ipentry"

// +k8s:deepcopy-gen=false
var _ resource.Object = &IPEntry{}
var _ resource.ObjectList = &IPEntryList{}

// GetListMeta returns the ListMeta
func (r *IPEntryList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *IPEntry) GetSingularName() string {
	return IPEntrySingular
}

func (IPEntry) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: IPEntryPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (IPEntry) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *IPEntry) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (IPEntry) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (IPEntry) New() runtime.Object {
	return &IPEntry{}
}

// NewList implements resource.Object
func (IPEntry) NewList() runtime.Object {
	return &IPEntryList{}
}

// GetCondition returns the condition based on the condition kind
func (r *IPEntry) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *IPEntry) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// ConvertIPEntryFieldSelector is the schema conversion function for normalizing the FieldSelector for IPEntry
func ConvertIPEntryFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *IPEntry) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *IPEntry) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *IPEntry) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return r.Spec.Owner
}

func (r *IPEntry) GetClaimName() string {
	return r.Spec.IPClaim
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

func (r *IPEntry) GetIPPrefix() string {
	return r.Spec.Prefix
}

func GetIPEntry(ctx context.Context, k store.Key, prefix netip.Prefix, labels map[string]string) *IPEntry {
	//log := log.FromContext(ctx)
	pi := iputil.NewPrefixInfo(prefix)

	ni := k.Name
	ns := k.Namespace

	spec := &IPEntrySpec{
		NetworkInstance: ni,
		PrefixType:      GetIPPrefixTypeFromString(labels[backend.KuidIPAMIPPrefixTypeKey]),
		ClaimType:       GetIPClaimTypeFromString(labels[backend.KuidIPAMClaimTypeKey]),
		AddressFamily:   ptr.To[iputil.AddressFamily](pi.GetAddressFamily()),
		Prefix:          pi.String(),
		IPClaim:         labels[backend.KuidClaimNameKey],
		Owner: &commonv1alpha1.OwnerReference{
			Group:     labels[backend.KuidOwnerGroupKey],
			Version:   labels[backend.KuidOwnerVersionKey],
			Kind:      labels[backend.KuidOwnerKindKey],
			Namespace: labels[backend.KuidOwnerNamespaceKey],
			Name:      labels[backend.KuidOwnerNameKey],
		},
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

	status := &IPEntryStatus{}
	status.SetConditions(conditionv1alpha1.Ready())

	return BuildIPEntry(
		metav1.ObjectMeta{
			Name:      pi.GetSubnetName(),
			Namespace: ns,
		},
		spec,
		status,
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
