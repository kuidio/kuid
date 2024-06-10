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

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
)

const LinkPlural = "links"
const LinkSingular = "link"

// +k8s:deepcopy-gen=false
var _ resource.Object = &Link{}
var _ resource.ObjectList = &LinkList{}
var _ backend.ObjectList = &LinkList{}
var _ backend.GenericObject = &Link{}
var _ backend.GenericObjectList = &LinkList{}

// GetListMeta returns the ListMeta
func (r *LinkList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *Link) GetSingularName() string {
	return LinkSingular
}

func (Link) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: LinkPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (Link) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *Link) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (Link) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (Link) New() runtime.Object {
	return &Link{}
}

// NewList implements resource.Object
func (Link) NewList() runtime.Object {
	return &LinkList{}
}

func (r *Link) NewObjList() backend.GenericObjectList {
	return &LinkList{
		TypeMeta: metav1.TypeMeta{APIVersion: SchemeGroupVersion.Identifier(), Kind: LinkKindList},
	}
}

func (r *Link) SchemaGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind(LinkKind)
}

// GetCondition returns the condition based on the condition kind
func (r *Link) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *Link) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// LinkConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for Link
func LinkConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *LinkList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *LinkList) GetObjects() []backend.GenericObject {
	objs := []backend.GenericObject{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *LinkList) GetLinks() []*Link {
	objs := []*Link{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *Link) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *Link) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *Link) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *Link) GetEndPointIDA() *NodeGroupEndpointID {
	if len(r.Spec.Endpoints) != 2 {
		return nil
	}
	return r.Spec.Endpoints[0]
}

func (r *Link) GetEndPointIDB() *NodeGroupEndpointID {
	if len(r.Spec.Endpoints) != 2 {
		return nil
	}
	return r.Spec.Endpoints[1]
}

func (r *Link) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      LinkKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *Link) ValidateSyntax(_ string) field.ErrorList {
	var allErrs field.ErrorList

	if len(r.Spec.Endpoints) != 2 {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.endpoints"),
			r,
			fmt.Errorf("a link always need 2 endpoints got %d", len(r.Spec.Endpoints)).Error(),
		))
	}

	return allErrs
}

func (r *Link) GetSpec() any {
	return r.Spec
}

func (r *Link) SetSpec(s any) {
	if spec, ok := s.(LinkSpec); ok {
		r.Spec = spec
	}
}

// BuildLink returns a reource from a client Object a Spec/Status
func BuildLink(meta metav1.ObjectMeta, spec *LinkSpec, status *LinkStatus) *Link {
	aspec := LinkSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := LinkStatus{}
	if status != nil {
		astatus = *status
	}
	return &Link{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       LinkKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func LinkTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			r, ok := obj.(*Link)
			if !ok {
				return nil
			}
			return []interface{}{
				r.GetName(),
				r.GetCondition(conditionv1alpha1.ConditionTypeReady).Status,
				r.GetEndPointIDA().KuidString(),
				r.GetEndPointIDB().KuidString(),
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			{Name: "EPA", Type: "string"},
			{Name: "EPB", Type: "string"},
		},
	}
}

func LinkParseFieldSelector(ctx context.Context, fieldSelector fields.Selector) (backend.Filter, error) {
	var filter *LinkFilter

	// add the namespace to the list
	namespace, ok := genericapirequest.NamespaceFrom(ctx)
	if fieldSelector == nil {
		if ok {
			return &LinkFilter{Namespace: namespace}, nil
		}
		return filter, nil
	}
	requirements := fieldSelector.Requirements()
	for _, requirement := range requirements {
		filter = &LinkFilter{}
		switch requirement.Operator {
		case selection.Equals, selection.DoesNotExist:
			if requirement.Value == "" {
				return filter, apierrors.NewBadRequest(fmt.Sprintf("unsupported fieldSelector value %q for field %q with operator %q", requirement.Value, requirement.Field, requirement.Operator))
			}
		default:
			return filter, apierrors.NewBadRequest(fmt.Sprintf("unsupported fieldSelector operator %q for field %q", requirement.Operator, requirement.Field))
		}

		switch requirement.Field {
		case "metadata.name":
			filter.Name = requirement.Value
		case "metadata.namespace":
			filter.Namespace = requirement.Value
		default:
			return filter, apierrors.NewBadRequest(fmt.Sprintf("unknown fieldSelector field %q", requirement.Field))
		}
	}
	// add namespace to the filter selector if specified
	if ok {
		if filter != nil {
			filter.Namespace = namespace
		} else {
			filter = &LinkFilter{Namespace: namespace}
		}
	}

	return &LinkFilter{}, nil
}

type LinkFilter struct {
	// Name filters by the name of the objects
	Name string `protobuf:"bytes,1,opt,name=name"`

	// Namespace filters by the namespace of the objects
	Namespace string `protobuf:"bytes,2,opt,name=namespace"`
}

func (r *LinkFilter) Filter(ctx context.Context, obj runtime.Object) bool {
	f := false // result of the previous filter
	o, ok := obj.(*Link)
	if !ok {
		return f
	}
	if r.Name != "" {
		if o.GetName() == r.Name {
			f = false
		} else {
			f = true
		}
	}
	if r.Namespace != "" {
		if o.GetNamespace() == r.Namespace {
			f = false
		} else {
			f = true
		}
	}
	return f
}

func (r *Link) GetUserDefinedLabels() map[string]string {
	return r.Spec.GetUserDefinedLabels()
}

func (r *Link) GetProvider() string {
	return ""
}

func (r *Link) GetISISLevel() ISISLevel {
	if r.Spec.ISIS != nil &&
		r.Spec.ISIS.Level != nil {
		return *r.Spec.ISIS.Level
	}
	return ISISLevelUnknown
}

func (r *Link) GetISISNetworkType() NetworkType {
	if r.Spec.ISIS != nil &&
		r.Spec.ISIS.NetworkType != nil {
		return *r.Spec.ISIS.NetworkType
	}
	return NetworkTypeUnknown
}

func (r *Link) GetISISPassive() bool {
	if r.Spec.ISIS != nil &&
		r.Spec.ISIS.Passive != nil {
		return *r.Spec.ISIS.Passive
	}
	return false
}

func (r *Link) GetISISBFD() bool {
	if r.Spec.ISIS != nil &&
		r.Spec.ISIS.BFD != nil {
		return *r.Spec.ISIS.BFD
	}
	return false
}

func (r *Link) GetISISMetric() uint32 {
	if r.Spec.ISIS != nil &&
		r.Spec.ISIS.Metric != nil {
		return *r.Spec.ISIS.Metric
	}
	return 0
}

func (r *Link) GetOSPFArea() string {
	if r.Spec.OSPF != nil &&
		r.Spec.OSPF.Area != nil {
		return *r.Spec.OSPF.Area
	}
	return ""
}

func (r *Link) GetOSPFNetworkType() NetworkType {
	if r.Spec.OSPF != nil &&
		r.Spec.OSPF.NetworkType != nil {
		return *r.Spec.OSPF.NetworkType
	}
	return NetworkTypeUnknown
}

func (r *Link) GetOSPFPassive() bool {
	if r.Spec.OSPF != nil &&
		r.Spec.OSPF.Passive != nil {
		return *r.Spec.OSPF.Passive
	}
	return false
}

func (r *Link) GetOSPFBFD() bool {
	if r.Spec.OSPF != nil &&
		r.Spec.OSPF.BFD != nil {
		return *r.Spec.OSPF.BFD
	}
	return false
}

func (r *Link) GetOSPFMetric() uint32 {
	if r.Spec.OSPF != nil &&
		r.Spec.OSPF.Metric != nil {
		return *r.Spec.OSPF.Metric
	}
	return 0
}

func (r *Link) GetBGPBFD() bool {
	if r.Spec.BGP != nil &&
		r.Spec.BGP.BFD != nil {
		return *r.Spec.BGP.BFD
	}
	return false
}
