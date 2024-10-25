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
	ModulePlural   = "modules"
	ModuleSingular = "module"
)

var (
	ModuleShortNames = []string{}
	ModuleCategories = []string{"kuid", "knet"}
)

// +k8s:deepcopy-gen=false
var _ resource.InternalObject = &Module{}
var _ resource.ObjectList = &ModuleList{}
var _ resource.ObjectWithStatusSubResource = &Module{}
var _ resource.StatusSubResource = &ModuleStatus{}

func (Module) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: ModulePlural,
	}
}

// IsStorageVersion returns true -- Config is used as the internal version.
// IsStorageVersion implements resource.Object
func (Module) IsStorageVersion() bool {
	return true
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object
func (Module) NamespaceScoped() bool {
	return true
}

// GetObjectMeta implements resource.Object
// GetObjectMeta implements resource.Object
func (r *Module) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// GetSingularName returns the singular name of the resource
// GetSingularName implements resource.Object
func (Module) GetSingularName() string {
	return ModuleSingular
}

// GetShortNames returns the shortnames for the resource
// GetShortNames implements resource.Object
func (Module) GetShortNames() []string {
	return ModuleShortNames
}

// GetCategories return the categories of the resource
// GetCategories implements resource.Object
func (Module) GetCategories() []string {
	return ModuleCategories
}

// New return an empty resource
// New implements resource.Object
func (Module) New() runtime.Object {
	return &Module{}
}

// NewList return an empty resourceList
// NewList implements resource.Object
func (Module) NewList() runtime.Object {
	return &ModuleList{}
}

// IsEqual returns a bool indicating if the desired state of both resources is equal or not
func (r *Module) IsEqual(ctx context.Context, obj, old runtime.Object) bool {
	newobj := obj.(*Module)
	oldobj := old.(*Module)

	if !apiequality.Semantic.DeepEqual(oldobj.ObjectMeta, newobj.ObjectMeta) {
		return false
	}
	// if equal we also test the spec
	return apiequality.Semantic.DeepEqual(oldobj.Spec, newobj.Spec)
}

// GetStatus return the resource.StatusSubResource interface
func (r *Module) GetStatus() resource.StatusSubResource {
	return r.Status
}

// IsStatusEqual returns a bool indicating if the status of both resources is equal or not
func (r *Module) IsStatusEqual(ctx context.Context, obj, old runtime.Object) bool {
	newobj := obj.(*Module)
	oldobj := old.(*Module)
	return apiequality.Semantic.DeepEqual(oldobj.Status, newobj.Status)
}

// PrepareForStatusUpdate prepares the status update
func (r *Module) PrepareForStatusUpdate(ctx context.Context, obj, old runtime.Object) {
	newObj := obj.(*Module)
	oldObj := old.(*Module)
	newObj.Spec = oldObj.Spec

	// Status updates are for only for updating status, not objectmeta.
	metav1.ResetObjectMetaForStatus(&newObj.ObjectMeta, &newObj.ObjectMeta)
}

// ValidateStatusUpdate validates status updates
func (r *Module) ValidateStatusUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

// SubResourceName resturns the name of the subresource
// SubResourceName implements the resource.StatusSubResource
func (ModuleStatus) SubResourceName() string {
	return fmt.Sprintf("%s/%s", ModulePlural, "status")
}

// CopyTo copies the content of the status subresource to a parent resource.
// CopyTo implements the resource.StatusSubResource
func (r ModuleStatus) CopyTo(obj resource.ObjectWithStatusSubResource) {
	parent, ok := obj.(*Module)
	if ok {
		parent.Status = r
	}
}

// GetListMeta returns the ListMeta
// GetListMeta implements the resource.ObjectList
func (r *ModuleList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

// TableConvertor return the table format of the resource
func (r *Module) TableConvertor() func(gr schema.GroupResource) rest.TableConvertor {
	return func(gr schema.GroupResource) rest.TableConvertor {
		return registry.NewTableConverter(
			gr,
			func(obj runtime.Object) []interface{} {
				claim, ok := obj.(*Module)
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
func (r *Module) FieldLabelConversion() runtime.FieldLabelConversionFunc {
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

func (r *Module) FieldSelector() func(ctx context.Context, fieldSelector fields.Selector) (resource.Filter, error) {
	return func(ctx context.Context, fieldSelector fields.Selector) (resource.Filter, error) {
		var filter *ModuleFilter

		// add the namespace to the list
		namespace, ok := genericapirequest.NamespaceFrom(ctx)
		if fieldSelector == nil {
			if ok {
				return &ModuleFilter{Namespace: namespace}, nil
			}
			return filter, nil
		}
		requirements := fieldSelector.Requirements()
		for _, requirement := range requirements {
			filter = &ModuleFilter{}
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
				filter = &ModuleFilter{Namespace: namespace}
			}
			return filter, nil
		}

		return &ModuleFilter{}, nil
	}

}

type ModuleFilter struct {
	// Name filters by the name of the objects
	Name string `protobuf:"bytes,1,opt,name=name"`

	// Namespace filters by the namespace of the objects
	Namespace string `protobuf:"bytes,2,opt,name=namespace"`
}

func (r *ModuleFilter) Filter(ctx context.Context, obj runtime.Object) bool {
	f := false // result of the previous filter
	o, ok := obj.(*Module)
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

func (r *Module) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	// status cannot be set upon create -> reset it
	newobj := obj.(*Module)
	newobj.Status = ModuleStatus{}
}

// ValidateCreate statically validates
func (r *Module) ValidateCreate(ctx context.Context, obj runtime.Object) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

func (r *Module) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	// ensure the status dont get updated
	newobj := obj.(*Module)
	oldObj := old.(*Module)
	newobj.Status = oldObj.Status
}

func (r *Module) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}
