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

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/iputil"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const IPIndexPlural = "ipindices"
const IPIndexSingular = "ipindex"

// +k8s:deepcopy-gen=false
var _ resource.Object = &IPIndex{}
var _ resource.ObjectList = &IPIndexList{}

// GetListMeta returns the ListMeta
func (r *IPIndexList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *IPIndex) GetSingularName() string {
	return IPIndexSingular
}

func (IPIndex) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: IPIndexPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (IPIndex) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *IPIndex) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (IPIndex) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (IPIndex) New() runtime.Object {
	return &IPIndex{}
}

// NewList implements resource.Object
func (IPIndex) NewList() runtime.Object {
	return &IPIndexList{}
}

// GetCondition returns the condition based on the condition kind
func (r *IPIndex) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *IPIndex) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// IPIndexConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for IPIndex
func IPIndexConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *IPIndexList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *IPIndex) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *IPIndex) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

// GetTree satisfies the interface but should not be used
func (r *IPIndex) GetTree() gtree.GTree {
	return nil
}

func (r *IPIndex) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *IPIndex) GetType() string {
	return ""
}

func (r *IPIndex) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      IPIndexKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *IPIndex) ValidateSyntax() field.ErrorList {
	var allErrs field.ErrorList

	if len(r.Spec.Prefixes) == 0 {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.minID"),
			r,
			fmt.Errorf("a ipindex needs a prefix").Error(),
		))

	}

	return allErrs
}

func (r *IPIndex) GetClaim(prefix Prefix) (*IPClaim, error) {
	pi, err := iputil.New(prefix.Prefix)
	if err != nil {
		return nil, err
	}

	return BuildIPClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      pi.GetSubnetName(),
		},
		&IPClaimSpec{
			PrefixType:   ptr.To[IPPrefixType](IPPrefixType_Aggregate),
			Index:        r.Name,
			Prefix:       ptr.To[string](prefix.Prefix),
			PrefixLength: ptr.To[uint32](uint32(pi.GetPrefixLength())),
			CreatePrefix: ptr.To[bool](true),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: prefix.UserDefinedLabels,
			},
			Owner: commonv1alpha1.GetOwnerReference(r),
		},
		nil,
	), nil
}

// BuildIPIndex returns a reource from a client Object a Spec/Status
func BuildIPIndex(meta metav1.ObjectMeta, spec *IPIndexSpec, status *IPIndexStatus) *IPIndex {
	aspec := IPIndexSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := IPIndexStatus{}
	if status != nil {
		astatus = *status
	}
	return &IPIndex{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       IPIndexKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

// GetMinID satisfies the interface but should not be used
func (r *IPIndex) GetMinID() *uint64 { return nil }

// GetMaxID satisfies the interface but should not be used
func (r *IPIndex) GetMaxID() *uint64 { return nil }

// GetMinClaim satisfies the interface but should not be used
func (r *IPIndex) GetMinClaim() backend.ClaimObject { return nil }

// GetMaxClaim satisfies the interface but should not be used
func (r *IPIndex) GetMaxClaim() backend.ClaimObject { return nil }

func IPIndexTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			index, ok := obj.(*IPIndex)
			if !ok {
				return nil
			}

			prefixes := make([]string, 5)
			for i, prefix := range index.Spec.Prefixes {
				if i >= 5 {
					break
				}
				prefixes[i] = prefix.Prefix
			}

			return []interface{}{
				index.Name,
				index.GetCondition(conditionv1alpha1.ConditionTypeReady).Status,
				prefixes[0],
				prefixes[1],
				prefixes[2],
				prefixes[3],
				prefixes[4],
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			{Name: "Prefix0", Type: "string"},
			{Name: "Prefix1", Type: "string"},
			{Name: "Prefix2", Type: "string"},
			{Name: "Prefix3", Type: "string"},
			{Name: "Prefix4", Type: "string"},
		},
	}
}

func GetIPClaimFromPrefix(obj client.Object, prefix Prefix) *IPClaim {
	// prefix validation should have happened before
	pi, _ := iputil.New(prefix.Prefix)

	var prefixType *IPPrefixType
	if pt, ok := prefix.Labels[backend.KuidIPAMIPPrefixTypeKey]; ok {
		prefixType = GetIPPrefixTypeFromString(pt)
	}

	// topology.vpc-name -> vpc-name default is the default router
	index := obj.GetName()

	return BuildIPClaim(
		metav1.ObjectMeta{
			Namespace: obj.GetNamespace(),
			Name:      fmt.Sprintf("%s.%s", index, pi.GetSubnetName()),
		},
		&IPClaimSpec{
			PrefixType:   prefixType,
			Index:        index,
			Prefix:       ptr.To[string](prefix.Prefix),
			PrefixLength: ptr.To[uint32](uint32(pi.GetPrefixLength())),
			CreatePrefix: ptr.To[bool](true),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: prefix.UserDefinedLabels,
			},
			Owner: commonv1alpha1.GetOwnerReference(obj),
		},
		nil,
	)
}
