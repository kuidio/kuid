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

package infra

import (
	"context"
	"fmt"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/kform-dev/choreo/apis/condition"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

const (
	LinkPlural   = "links"
	LinkSingular = "links"
)

var (
	LinkShortNames = []string{}
	LinkCategories = []string{"kuid", "knet"}
)

// +k8s:deepcopy-gen=false
var _ resource.InternalObject = &Link{}
var _ resource.ObjectList = &LinkList{}
var _ resource.ObjectWithStatusSubResource = &Link{}
var _ resource.StatusSubResource = &LinkStatus{}

func (Link) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: LinkPlural,
	}
}

// IsStorageVersion returns true -- Config is used as the internal version.
// IsStorageVersion implements resource.Object
func (Link) IsStorageVersion() bool {
	return true
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object
func (Link) NamespaceScoped() bool {
	return true
}

// GetObjectMeta implements resource.Object
// GetObjectMeta implements resource.Object
func (r *Link) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// GetSingularName returns the singular name of the resource
// GetSingularName implements resource.Object
func (Link) GetSingularName() string {
	return LinkSingular
}

// GetShortNames returns the shortnames for the resource
// GetShortNames implements resource.Object
func (Link) GetShortNames() []string {
	return LinkShortNames
}

// GetCategories return the categories of the resource
// GetCategories implements resource.Object
func (Link) GetCategories() []string {
	return LinkCategories
}

// New return an empty resource
// New implements resource.Object
func (Link) New() runtime.Object {
	return &Link{}
}

// NewList return an empty resourceList
// NewList implements resource.Object
func (Link) NewList() runtime.Object {
	return &LinkList{}
}

// IsEqual returns a bool indicating if the desired state of both resources is equal or not
func (r *Link) IsEqual(ctx context.Context, obj, old runtime.Object) bool {
	newobj := obj.(*Link)
	oldobj := old.(*Link)

	if !apiequality.Semantic.DeepEqual(oldobj.ObjectMeta, newobj.ObjectMeta) {
		return false
	}
	// if equal we also test the spec
	return apiequality.Semantic.DeepEqual(oldobj.Spec, newobj.Spec)
}

// GetStatus return the resource.StatusSubResource interface
func (r *Link) GetStatus() resource.StatusSubResource {
	return r.Status
}

// IsStatusEqual returns a bool indicating if the status of both resources is equal or not
func (r *Link) IsStatusEqual(ctx context.Context, obj, old runtime.Object) bool {
	newobj := obj.(*Link)
	oldobj := old.(*Link)
	return apiequality.Semantic.DeepEqual(oldobj.Status, newobj.Status)
}

// PrepareForStatusUpdate prepares the status update
func (r *Link) PrepareForStatusUpdate(ctx context.Context, obj, old runtime.Object) {
	newObj := obj.(*Link)
	oldObj := old.(*Link)
	newObj.Spec = oldObj.Spec

	// Status updates are for only for updating status, not objectmeta.
	metav1.ResetObjectMetaForStatus(&newObj.ObjectMeta, &newObj.ObjectMeta)
}

// ValidateStatusUpdate validates status updates
func (r *Link) ValidateStatusUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

// SubResourceName resturns the name of the subresource
// SubResourceName implements the resource.StatusSubResource
func (LinkStatus) SubResourceName() string {
	return fmt.Sprintf("%s/%s", LinkPlural, "status")
}

// CopyTo copies the content of the status subresource to a parent resource.
// CopyTo implements the resource.StatusSubResource
func (r LinkStatus) CopyTo(obj resource.ObjectWithStatusSubResource) {
	parent, ok := obj.(*Link)
	if ok {
		parent.Status = r
	}
}

// GetListMeta returns the ListMeta
// GetListMeta implements the resource.ObjectList
func (r *LinkList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

// TableConvertor return the table format of the resource
func (r *Link) TableConvertor() func(gr schema.GroupResource) rest.TableConvertor {
	return func(gr schema.GroupResource) rest.TableConvertor {
		return registry.NewTableConverter(
			gr,
			func(obj runtime.Object) []interface{} {
				claim, ok := obj.(*Link)
				if !ok {
					return nil
				}
				return []interface{}{
					claim.GetName(),
					claim.GetCondition(condition.ConditionTypeReady).Status,
				}
			},
			[]metav1.TableColumnDefinition{
				{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			},
		)
	}
}

// FieldLabelConversion is the schema conversion function for normalizing the FieldSelector for the resource
func (r *Link) FieldLabelConversion() runtime.FieldLabelConversionFunc {
	return func(label, value string) (internalLabel, internalValue string, err error) {
		switch label {
		case "metadata.name":
			return label, value, nil
		case "metadata.namespace":
			return label, value, nil
		default:
			return "", "", fmt.Errorf("%q is not a known field selector", label)
		}
	}
}

func (r *Link) FieldSelector() func(ctx context.Context, fieldSelector fields.Selector) (resource.Filter, error) {
	return func(ctx context.Context, fieldSelector fields.Selector) (resource.Filter, error) {
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
			return filter, nil
		}

		return &LinkFilter{}, nil
	}

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
	if r == nil {
		return false
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

func (r *Link) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	// status cannot be set upon create -> reset it
	newobj := obj.(*Link)
	newobj.Status = LinkStatus{}
}

// ValidateCreate statically validates
func (r *Link) ValidateCreate(ctx context.Context, obj runtime.Object) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

func (r *Link) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	// ensure the status dont get updated
	newobj := obj.(*Link)
	oldObj := old.(*Link)
	newobj.Status = oldObj.Status
}

func (r *Link) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}
