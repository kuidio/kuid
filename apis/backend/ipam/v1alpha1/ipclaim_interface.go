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

	"github.com/hansthienpondt/nipam/pkg/table"
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
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func (r *IPClaim) ValidateOwner(labels labels.Set) error {
	fmt.Println("ValidateOwner claim", r)
	fmt.Println("ValidateOwner labels", labels)
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

func (r *IPClaim) ValidateExistingChildren(route table.Route) error {
	routeLabels := route.Labels()
	switch r.Spec.Kind {
	case PrefixKindNetwork:
		if routeLabels[backend.KuidIPAMKindKey] != string(PrefixKindNetwork) {
			return fmt.Errorf("child prefix kind mismatch got %s/%s", string(PrefixKindNetwork), routeLabels[backend.KuidIPAMKindKey])
		}
		if route.Prefix().Addr().Is4() && route.Prefix().Bits() != 32 {
			return fmt.Errorf("a child prefix of kind %s, can only be an address prefix (/32), got: %v", string(PrefixKindNetwork), route.Prefix())
		}
		if route.Prefix().Addr().Is6() && route.Prefix().Bits() != 128 {
			return fmt.Errorf("a child prefix of kind %s, can only be an address prefix (/128), got: %v", string(PrefixKindNetwork), route.Prefix())
		}
		return nil
	case PrefixKindAggregate:
		// nesting is possible in aggregate
		return nil
	default:
		return fmt.Errorf("a more specific prefix was already claimed %s/%s, nesting not allowed for %s",
			routeLabels[backend.KuidIPAMKindKey],
			routeLabels[backend.KuidClaimNameKey],
			string(r.Spec.Kind))
	}
}

// GetDummyLabelsFromPrefix returns a map with the labels from the spec
// augmented with the prefixkind and the subnet from the prefixInfo
func (r *IPClaim) GetDummyLabelsFromPrefix(pi iputil.Prefix) map[string]string {
	labels := map[string]string{}
	for k, v := range r.Spec.GetUserDefinedLabels() {
		labels[k] = v
	}
	labels[backend.KuidIPAMKindKey] = string(r.Spec.Kind)
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
func (r *IPClaim) GetGatewayLabelSelector(subnetString string) (labels.Selector, error) {
	l := map[string]string{
		backend.KuidIPAMGatewayKey: "true",
		backend.KuidIPAMSubnetKey:  subnetString,
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

func (r *IPClaim) ValidateExistingParent(pi *iputil.Prefix, route table.Route) error {
	routeLabels := route.Labels()
	switch r.Spec.Kind {
	case PrefixKindAggregate:
		if routeLabels[backend.KuidIPAMKindKey] != string(PrefixKindAggregate) {
			return fmt.Errorf("nesting aggregate prefixes with anything other than an aggregate prefix is not allowed, prefix %s/%s kind %s/%s",
				route.Prefix().String(),
				pi.String(),
				routeLabels[backend.KuidIPAMKindKey],
				string(PrefixKindAggregate),
			)
		}
		return nil
	case PrefixKindLoopback:
		if routeLabels[backend.KuidIPAMKindKey] != string(PrefixKindAggregate) &&
			routeLabels[backend.KuidIPAMKindKey] != string(PrefixKindLoopback) {
			return fmt.Errorf("nesting loopback prefixes with anything other than an aggregate/loopback prefix is not allowed, prefix %s/%s kind %s/%s",
				route.Prefix().String(),
				pi.String(),
				routeLabels[backend.KuidIPAMKindKey],
				string(PrefixKindAggregate),
			)
		}
		if pi.IsAddressPrefix() {
			// address (/32 or /128) can parant with aggregate or loopback
			switch routeLabels[backend.KuidIPAMKindKey] {
			case string(PrefixKindAggregate), string(PrefixKindLoopback):
				// /32 or /128 can be parented with aggregates or loopbacks
			default:
				return fmt.Errorf("nesting loopback prefixes only possible with address (/32, /128) based prefixes, got %s", pi.GetIPPrefix().String())
			}
		}

		if !pi.IsAddressPrefix() {
			switch routeLabels[backend.KuidIPAMKindKey] {
			case string(PrefixKindAggregate):
				// none /32 or /128 can only be parented with aggregates
			default:
				return fmt.Errorf("nesting (none /32, /128)loopback prefixes only possible with aggregate prefixes, got %s", route.String())
			}
		}
		return nil
	case PrefixKindNetwork:
		if r.Spec.CreatePrefix != nil {
			if routeLabels[backend.KuidIPAMKindKey] != string(PrefixKindAggregate) {
				return fmt.Errorf("nesting network prefixes with anything other than an aggregate prefix is not allowed, prefix %s/%s kind %s/%s",
					route.Prefix().String(),
					pi.String(),
					routeLabels[backend.KuidIPAMKindKey],
					string(PrefixKindAggregate),
				)
			}
		} else {
			if routeLabels[backend.KuidIPAMKindKey] != string(PrefixKindNetwork) {
				return fmt.Errorf("nesting network address prefixes with anything other than an network prefix is not allowed, prefix %s/%s kind %s/%s",
					route.Prefix().String(),
					pi.String(),
					routeLabels[backend.KuidIPAMKindKey],
					string(PrefixKindAggregate),
				)
			}
		}
		return nil
	case PrefixKindPool:
		// if the parent is not an aggregate we dont allow the prefix to be created
		if routeLabels[backend.KuidIPAMKindKey] != string(PrefixKindAggregate) &&
			routeLabels[backend.KuidIPAMKindKey] != string(PrefixKindPool) {
			return fmt.Errorf("nesting loopback prefixes with anything other than an aggregate/pool prefix is not allowed, prefix %s/%s kind %s/%s",
				route.Prefix().String(),
				pi.String(),
				routeLabels[backend.KuidIPAMKindKey],
				string(PrefixKindAggregate),
			)
		}
		return nil
	default:
		return fmt.Errorf("unknown prefix kind %s", r.Spec.Kind)
	}
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

func GetIPClaimFromPrefix(kind PrefixKind, ni, prefix string, udLabels commonv1alpha1.UserDefinedLabels, obj client.Object) (*IPClaim, error) {
	pi, err := iputil.New(prefix)
	if err != nil {
		return nil, fmt.Errorf("cannot build ip claim bad prefix %s", prefix)
	}
	return BuildIPClaim(
		metav1.ObjectMeta{
			Namespace: obj.GetNamespace(),
			Name:      pi.GetSubnetName(),
		},
		GetIPClaimSpec(kind, ni, pi, udLabels, obj),
		nil,
	), nil
}

func GetIPClaimSpec(kind PrefixKind, ni string, pi *iputil.Prefix, udLabels commonv1alpha1.UserDefinedLabels, obj client.Object) *IPClaimSpec {
	return &IPClaimSpec{
		Kind:            kind,
		NetworkInstance: ni,
		AddressFamily:   ptr.To[iputil.AddressFamily](pi.GetAddressFamily()),
		Prefix:          ptr.To[string](pi.Prefix.String()),
		PrefixLength:    ptr.To[uint32](uint32(pi.GetPrefixLength())),
		CreatePrefix:    ptr.To[bool](!pi.IsAddressPrefix()),
		ClaimLabels: commonv1alpha1.ClaimLabels{
			UserDefinedLabels: udLabels,
		},
		Owner: commonv1alpha1.GetOwnerReference(obj),
	}
}

func (r *IPClaim) UpdateSpecFromPrefix(kind PrefixKind, ni, prefix string, udLabels commonv1alpha1.UserDefinedLabels, obj client.Object) error {
	pi, err := iputil.New(prefix)
	if err != nil {
		return fmt.Errorf("cannot build ip claim bad prefix %s", prefix)
	}
	r.Spec = *(GetIPClaimSpec(kind, ni, pi, udLabels, obj))
	return nil
}
