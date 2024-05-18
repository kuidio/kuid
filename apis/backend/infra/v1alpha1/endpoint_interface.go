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
	"github.com/henderiw/store"
	"github.com/kuidio/kuid/apis/backend"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const EndpointPlural = "endpoints"
const EndpointSingular = "endpoint"

// +k8s:deepcopy-gen=false
var _ resource.Object = &Endpoint{}
var _ resource.ObjectList = &EndpointList{}

// GetListMeta returns the ListMeta
func (r *EndpointList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *Endpoint) GetSingularName() string {
	return EndpointSingular
}

func (Endpoint) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: EndpointPlural,
	}
}

// IsStorageVersion returns true -- v1alpha1.Config is used as the internal version.
// IsStorageVersion implements resource.Object.
func (Endpoint) IsStorageVersion() bool {
	return true
}

// GetObjectMeta implements resource.Object
func (r *Endpoint) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object.
func (Endpoint) NamespaceScoped() bool {
	return true
}

// New implements resource.Object
func (Endpoint) New() runtime.Object {
	return &Endpoint{}
}

// NewList implements resource.Object
func (Endpoint) NewList() runtime.Object {
	return &EndpointList{}
}

// GetCondition returns the condition based on the condition kind
func (r *Endpoint) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *Endpoint) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

// EndpointConvertFieldSelector is the schema conversion function for normalizing the FieldSelector for Endpoint
func EndpointConvertFieldSelector(label, value string) (internalLabel, internalValue string, err error) {
	switch label {
	case "metadata.name":
		return label, value, nil
	case "metadata.namespace":
		return label, value, nil
	default:
		return "", "", fmt.Errorf("%q is not a known field selector", label)
	}
}

func (r *EndpointList) GetItems() []backend.Object {
	objs := []backend.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *Endpoint) CalculateHash() ([sha1.Size]byte, error) {
	// Convert the struct to JSON
	jsonData, err := json.Marshal(r)
	if err != nil {
		return [sha1.Size]byte{}, err
	}

	// Calculate SHA-1 hash
	return sha1.Sum(jsonData), nil
}

func (r *Endpoint) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *Endpoint) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *Endpoint) GetRegion() string {
	return r.Spec.Region
}

func (r *Endpoint) GetSite() string {
	return r.Spec.Site
}

func (r *Endpoint) GetNodeGroupID() NodeGroupID {
	return NodeGroupID{
		SiteID:    r.Spec.SiteID,
		NodeGroup: r.Spec.NodeGroup,
	}
}

func (r *Endpoint) GetNodeID() NodeID {
	return NodeID{
		Node:        r.Spec.Node,
		NodeGroupID: r.GetNodeGroupID(),
	}
}

func (r *Endpoint) GetEndpointID() EndpointID {
	return EndpointID{
		NodeID:   r.GetNodeID(),
		Endpoint: r.Name,
	}
}

func (r *Endpoint) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      EndpointKind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *Endpoint) ValidateSyntax() field.ErrorList {
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

// BuildEndpoint returns a reource from a client Object a Spec/Status
func BuildEndpoint(meta metav1.ObjectMeta, spec *EndpointSpec, status *EndpointStatus) *Endpoint {
	aspec := EndpointSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := EndpointStatus{}
	if status != nil {
		astatus = *status
	}
	return &Endpoint{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       EndpointKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}

func EndpointTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			r, ok := obj.(*Endpoint)
			if !ok {
				return nil
			}
			return []interface{}{
				r.GetName(),
				r.GetCondition(conditionv1alpha1.ConditionTypeReady).Status,
				r.Spec.Region,
				r.Spec.Site,
				r.Spec.NodeGroup,
				r.Spec.Node,
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			{Name: "Region", Type: "string"},
			{Name: "Site", Type: "string"},
			{Name: "Topology", Type: "string"},
			{Name: "Node", Type: "string"},
		},
	}
}
