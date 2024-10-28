/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "VLAN IS" BVLANIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package genid

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
	GENIDClaimPlural   = "genidclaims"
	GENIDClaimSingular = "genidclaim"
)

var (
	GENIDClaimShortNames = []string{}
	GENIDClaimCategories = []string{"kuid", "knet"}
)

// +k8s:deepcopy-gen=false
var _ resource.InternalObject = &GENIDClaim{}
var _ resource.ObjectList = &GENIDClaimList{}
var _ resource.ObjectWithStatusSubResource = &GENIDClaim{}
var _ resource.StatusSubResource = &GENIDClaimStatus{}

func (GENIDClaim) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: GENIDClaimPlural,
	}
}

// IsStorageVersion returns true -- Config is used VLAN the internal version.
// IsStorageVersion implements resource.Object
func (GENIDClaim) IsStorageVersion() bool {
	return true
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object
func (GENIDClaim) NamespaceScoped() bool {
	return true
}

// GetObjectMeta implements resource.Object
// GetObjectMeta implements resource.Object
func (r *GENIDClaim) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// GetSingularName returns the singular name of the resource
// GetSingularName implements resource.Object
func (GENIDClaim) GetSingularName() string {
	return GENIDClaimSingular
}

// GetShortNames returns the shortnames for the resource
// GetShortNames implements resource.Object
func (GENIDClaim) GetShortNames() []string {
	return GENIDClaimShortNames
}

// GetCategories return the categories of the resource
// GetCategories implements resource.Object
func (GENIDClaim) GetCategories() []string {
	return GENIDClaimCategories
}

// New return an empty resource
// New implements resource.Object
func (GENIDClaim) New() runtime.Object {
	return &GENIDClaim{}
}

// NewList return an empty resourceList
// NewList implements resource.Object
func (GENIDClaim) NewList() runtime.Object {
	return &GENIDClaimList{}
}

// IsEqual returns a bool indicating if the desired state of both resources is equal or not
func (r *GENIDClaim) IsEqual(ctx context.Context, obj, old runtime.Object) bool {
	newobj := obj.(*GENIDClaim)
	oldobj := old.(*GENIDClaim)

	if !apiequality.Semantic.DeepEqual(oldobj.ObjectMeta, newobj.ObjectMeta) {
		return false
	}
	// if equal we also test the spec
	return apiequality.Semantic.DeepEqual(oldobj.Spec, newobj.Spec)
}

// GetStatus return the resource.StatusSubResource interface
func (r *GENIDClaim) GetStatus() resource.StatusSubResource {
	return r.Status
}

// IsStatusEqual returns a bool indicating if the status of both resources is equal or not
func (r *GENIDClaim) IsStatusEqual(ctx context.Context, obj, old runtime.Object) bool {
	newobj := obj.(*GENIDClaim)
	oldobj := old.(*GENIDClaim)
	return apiequality.Semantic.DeepEqual(oldobj.Status, newobj.Status)
}

// PrepareForStatusUpdate prepares the status update
func (r *GENIDClaim) PrepareForStatusUpdate(ctx context.Context, obj, old runtime.Object) {
	newObj := obj.(*GENIDClaim)
	oldObj := old.(*GENIDClaim)
	newObj.Spec = oldObj.Spec

	// Status updates are for only for updating status, not objectmeta.
	metav1.ResetObjectMetaForStatus(&newObj.ObjectMeta, &newObj.ObjectMeta)
}

// ValidateStatusUpdate validates status updates
func (r *GENIDClaim) ValidateStatusUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

// SubResourceName resturns the name of the subresource
// SubResourceName implements the resource.StatusSubResource
func (GENIDClaimStatus) SubResourceName() string {
	return fmt.Sprintf("%s/%s", GENIDClaimPlural, "status")
}

// CopyTo copies the content of the status subresource to a parent resource.
// CopyTo implements the resource.StatusSubResource
func (r GENIDClaimStatus) CopyTo(obj resource.ObjectWithStatusSubResource) {
	parent, ok := obj.(*GENIDClaim)
	if ok {
		parent.Status = r
	}
}

// GetListMeta returns the ListMeta
// GetListMeta implements the resource.ObjectList
func (r *GENIDClaimList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

// TableConvertor return the table format of the resource
func (r *GENIDClaim) TableConvertor() func(gr schema.GroupResource) rest.TableConvertor {
	return func(gr schema.GroupResource) rest.TableConvertor {
		return registry.NewTableConverter(
			gr,
			func(obj runtime.Object) []interface{} {
				claim, ok := obj.(*GENIDClaim)
				if !ok {
					return nil
				}
				return []interface{}{
					claim.GetName(),
					claim.GetCondition(condition.ConditionTypeReady).Status,
					claim.GetIndex(),
					string(claim.GetClaimType()),
					claim.GetClaimRequest(),
					claim.GetClaimResponse(),
				}
			},
			[]metav1.TableColumnDefinition{
				{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			{Name: "Index", Type: "string"},
			{Name: "ClaimType", Type: "string"},
			{Name: "ClaimReq", Type: "string"},
			{Name: "ClaimRsp", Type: "string"},
			},
		)
	}
}

// FieldLabelConversion is the schema conversion function for normalizing the FieldSelector for the resource
func (r *GENIDClaim) FieldLabelConversion() runtime.FieldLabelConversionFunc {
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

func (r *GENIDClaim) FieldSelector() func(ctx context.Context, fieldSelector fields.Selector) (resource.Filter, error) {
	return func(ctx context.Context, fieldSelector fields.Selector) (resource.Filter, error) {
		var filter *GENIDClaimFilter

		// add the namespace to the list
		namespace, ok := genericapirequest.NamespaceFrom(ctx)
		if fieldSelector == nil {
			if ok {
				return &GENIDClaimFilter{Namespace: namespace}, nil
			}
			return filter, nil
		}
		requirements := fieldSelector.Requirements()
		for _, requirement := range requirements {
			filter = &GENIDClaimFilter{}
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
				filter = &GENIDClaimFilter{Namespace: namespace}
			}
			return filter, nil
		}

		return &GENIDClaimFilter{}, nil
	}

}

type GENIDClaimFilter struct {
	// Name filters by the name of the objects
	Name string `protobuf:"bytes,1,opt,name=name"`

	// Namespace filters by the namespace of the objects
	Namespace string `protobuf:"bytes,2,opt,name=namespace"`
}

func (r *GENIDClaimFilter) Filter(ctx context.Context, obj runtime.Object) bool {
	f := false // result of the previous filter
	o, ok := obj.(*GENIDClaim)
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

func (r *GENIDClaim) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	// status cannot be set upon create -> reset it
	newobj := obj.(*GENIDClaim)
	newobj.Status = GENIDClaimStatus{}
}

// ValidateCreate statically validates
func (r *GENIDClaim) ValidateCreate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return r.ValidateSyntax("")
}

func (r *GENIDClaim) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	// ensure the status dont get updated
	newobj := obj.(*GENIDClaim)
	oldObj := old.(*GENIDClaim)
	newobj.Status = oldObj.Status
}

func (r *GENIDClaim) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return r.ValidateSyntax("")
}
