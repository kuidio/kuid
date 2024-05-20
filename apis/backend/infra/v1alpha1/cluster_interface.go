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

const ClusterPlural = "clusters"
const ClusterSingular = "cluster"

// +k8s:deepcopy-gen=false
var _ resource.Object = &Cluster{}
var _ resource.ObjectList = &ClusterList{}
var _ backend.ObjectList = &ClusterList{}
var _ backend.GenericObject = &Cluster{}
var _ backend.GenericObjectList = &ClusterList{}

// GetListMeta returns the ListMeta
func (r *ClusterList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *Cluster) GetSingularName() string {
	return ClusterSingular
}

func (Cluster) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: ClusterPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (Cluster) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *Cluster) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (Cluster) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (Cluster) New() runtime.Object {
	return &Cluster{}
}

// NewList implements resource.Object
func (Cluster) NewList() runtime.Object {
	return &ClusterList{}
}

func (r *Cluster) NewObjList() backend.GenericObjectList {
	return &ClusterList{
		TypeMeta: metav1.TypeMeta{APIVersion: SchemeGroupVersion.Identifier(), Kind: ClusterKindList},
	}
}

func (r *Cluster) SchemaGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind(ClusterKind)
}

// GetCondition returns the condition based on the condition kind
func (r *Cluster) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *Cluster) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// ClusterConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for Cluster
func ClusterConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *ClusterList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *ClusterList) GetObjects() []backend.GenericObject {
	objs := []backend.GenericObject{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *Cluster) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *Cluster) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *Cluster) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *Cluster) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      ClusterKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *Cluster) ValidateSyntax(_ string) field.ErrorList {
	var allErrs field.ErrorList

	/*
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.type"),
				r,
				fmt.Errorf("invalid GENID Type %s", r.Spec.Type).Error(),
			))
		}
	*/
	return allErrs
}

func (r *Cluster) GetSpec() any {
	return r.Spec
}

func (r *Cluster) SetSpec(s any) {
	if spec, ok := s.(ClusterSpec); ok {
		r.Spec = spec
	}
}

// BuildCluster returns a reource from a client Object a Spec/Status
func BuildCluster(meta metav1.ObjectMeta, spec *ClusterSpec, status *ClusterStatus) *Cluster {
	aspec := ClusterSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := ClusterStatus{}
	if status != nil {
		astatus = *status
	}
	return &Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       ClusterKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func ClusterTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			r, ok := obj.(*Cluster)
			if !ok {
				return nil
			}
			return []interface{}{
				r.GetName(),
				r.GetCondition(conditionv1alpha1.ConditionTypeReady).Status,
				r.Spec.NodeGroup,
				r.Spec.Provider,
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			{Name: "Domain", Type: "string"},
			{Name: "Provider", Type: "string"},
		},
	}
}

func ClusterParseFieldSelector(ctx context.Context, fieldSelector fields.Selector) (backend.Filter, error) {
	var filter *ClusterFilter

	// add the namespace to the list
	namespace, ok := genericapirequest.NamespaceFrom(ctx)
	if fieldSelector == nil {
		if ok {
			return &ClusterFilter{Namespace: namespace}, nil
		}
		return filter, nil
	}
	requirements := fieldSelector.Requirements()
	for _, requirement := range requirements {
		filter = &ClusterFilter{}
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
			filter = &ClusterFilter{Namespace: namespace}
		}
	}

	return &ClusterFilter{}, nil
}

type ClusterFilter struct {
	// Name filters by the name of the objects
	Name string `protobuf:"bytes,1,opt,name=name"`

	// Namespace filters by the namespace of the objects
	Namespace string `protobuf:"bytes,2,opt,name=namespace"`
}

func (r *ClusterFilter) Filter(ctx context.Context, obj runtime.Object) bool {
	f := false // result of the previous filter
	o, ok := obj.(*Cluster)
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

func (r *Cluster) GetUserDefinedLabels() map[string]string {
	return r.Spec.GetUserDefinedLabels()
}

func (r *Cluster) GetProvider() string {
	return r.Spec.Provider
}
