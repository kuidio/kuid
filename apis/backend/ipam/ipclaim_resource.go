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

package ipam

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
	IPClaimPlural   = "ipclaims"
	IPClaimSingular = "ipclaim"
)

var (
	IPClaimShortNames = []string{}
	IPClaimCategories = []string{"kuid", "knet"}
)

// +k8s:deepcopy-gen=false
var _ resource.InternalObject = &IPClaim{}
var _ resource.ObjectList = &IPClaimList{}
var _ resource.ObjectWithStatusSubResource = &IPClaim{}
var _ resource.StatusSubResource = &IPClaimStatus{}

func (IPClaim) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    SchemeGroupVersion.Group,
		Version:  SchemeGroupVersion.Version,
		Resource: IPClaimPlural,
	}
}

// IsStorageVersion returns true -- Config is used VLAN the internal version.
// IsStorageVersion implements resource.Object
func (IPClaim) IsStorageVersion() bool {
	return true
}

// NamespaceScoped returns true to indicate Fortune is a namespaced resource.
// NamespaceScoped implements resource.Object
func (IPClaim) NamespaceScoped() bool {
	return true
}

// GetObjectMeta implements resource.Object
// GetObjectMeta implements resource.Object
func (r *IPClaim) GetObjectMeta() *metav1.ObjectMeta {
	return &r.ObjectMeta
}

// GetSingularName returns the singular name of the resource
// GetSingularName implements resource.Object
func (IPClaim) GetSingularName() string {
	return IPClaimSingular
}

// GetShortNames returns the shortnames for the resource
// GetShortNames implements resource.Object
func (IPClaim) GetShortNames() []string {
	return IPClaimShortNames
}

// GetCategories return the categories of the resource
// GetCategories implements resource.Object
func (IPClaim) GetCategories() []string {
	return IPClaimCategories
}

// New return an empty resource
// New implements resource.Object
func (IPClaim) New() runtime.Object {
	return &IPClaim{}
}

// NewList return an empty resourceList
// NewList implements resource.Object
func (IPClaim) NewList() runtime.Object {
	return &IPClaimList{}
}

// IsEqual returns a bool indicating if the desired state of both resources is equal or not
func (r *IPClaim) IsEqual(ctx context.Context, obj, old runtime.Object) bool {
	newobj := obj.(*IPClaim)
	oldobj := old.(*IPClaim)

	if !apiequality.Semantic.DeepEqual(oldobj.ObjectMeta, newobj.ObjectMeta) {
		return false
	}
	// if equal we also test the spec
	return apiequality.Semantic.DeepEqual(oldobj.Spec, newobj.Spec)
}

// GetStatus return the resource.StatusSubResource interface
func (r *IPClaim) GetStatus() resource.StatusSubResource {
	return r.Status
}

// IsStatusEqual returns a bool indicating if the status of both resources is equal or not
func (r *IPClaim) IsStatusEqual(ctx context.Context, obj, old runtime.Object) bool {
	newobj := obj.(*IPClaim)
	oldobj := old.(*IPClaim)
	return apiequality.Semantic.DeepEqual(oldobj.Status, newobj.Status)
}

// PrepareForStatusUpdate prepares the status update
func (r *IPClaim) PrepareForStatusUpdate(ctx context.Context, obj, old runtime.Object) {
	newObj := obj.(*IPClaim)
	oldObj := old.(*IPClaim)
	newObj.Spec = oldObj.Spec

	// Status updates are for only for updating status, not objectmeta.
	metav1.ResetObjectMetaForStatus(&newObj.ObjectMeta, &newObj.ObjectMeta)
}

// ValidateStatusUpdate validates status updates
func (r *IPClaim) ValidateStatusUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

// SubResourceName resturns the name of the subresource
// SubResourceName implements the resource.StatusSubResource
func (IPClaimStatus) SubResourceName() string {
	return fmt.Sprintf("%s/%s", IPClaimPlural, "status")
}

// CopyTo copies the content of the status subresource to a parent resource.
// CopyTo implements the resource.StatusSubResource
func (r IPClaimStatus) CopyTo(obj resource.ObjectWithStatusSubResource) {
	parent, ok := obj.(*IPClaim)
	if ok {
		parent.Status = r
	}
}

// GetListMeta returns the ListMeta
// GetListMeta implements the resource.ObjectList
func (r *IPClaimList) GetListMeta() *metav1.ListMeta {
	return &r.ListMeta
}

func (r *IPClaim) GetClaimRequest() string {
	// we assume validation is already done when calling this
	if r.Spec.Address != nil {
		return *r.Spec.Address
	}
	if r.Spec.Prefix != nil {
		return *r.Spec.Prefix
	}
	if r.Spec.Range != nil {
		return *r.Spec.Range
	}
	return ""
}

func (r *IPClaim) GetClaimResponse() string {
	// we assume validation is already done when calling this
	if r.Status.Address != nil {
		return *r.Status.Address
	}
	if r.Status.Prefix != nil {
		return *r.Status.Prefix
	}
	if r.Status.Range != nil {
		return *r.Status.Range
	}
	return ""
}

// TableConvertor return the table format of the resource
func (r *IPClaim) TableConvertor() func(gr schema.GroupResource) rest.TableConvertor {
	return func(gr schema.GroupResource) rest.TableConvertor {
		return registry.NewTableConverter(
			gr,
			func(obj runtime.Object) []interface{} {
				ipclaim, ok := obj.(*IPClaim)
				if !ok {
					return nil
				}
				ipClaimType, _ := ipclaim.GetIPClaimType()
				return []interface{}{
					ipclaim.Name,
					ipclaim.GetCondition(condition.ConditionTypeReady).Status,
					ipclaim.Spec.Index,
					string(ipClaimType),
					string(ipclaim.GetIPPrefixType()),
					ipclaim.GetClaimRequest(),
					ipclaim.GetClaimResponse(),
					ipclaim.Status.DefaultGateway,
				}
			},
			[]metav1.TableColumnDefinition{
				{Name: "Name", Type: "string"},
				{Name: "Ready", Type: "string"},
				{Name: "Index", Type: "string"},
				{Name: "ClaimType", Type: "string"},
				{Name: "PrefixType", Type: "string"},
				{Name: "ClaimReq", Type: "string"},
				{Name: "ClaimRsp", Type: "string"},
				{Name: "DefaultGateway", Type: "string"},
			},
		)
	}
}

// FieldLabelConversion is the schema conversion function for normalizing the FieldSelector for the resource
func (r *IPClaim) FieldLabelConversion() runtime.FieldLabelConversionFunc {
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

func (r *IPClaim) FieldSelector() func(ctx context.Context, fieldSelector fields.Selector) (resource.Filter, error) {
	return func(ctx context.Context, fieldSelector fields.Selector) (resource.Filter, error) {
		var filter *IPClaimFilter

		// add the namespace to the list
		namespace, ok := genericapirequest.NamespaceFrom(ctx)
		if fieldSelector == nil {
			if ok {
				return &IPClaimFilter{Namespace: namespace}, nil
			}
			return filter, nil
		}
		requirements := fieldSelector.Requirements()
		for _, requirement := range requirements {
			filter = &IPClaimFilter{}
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
				filter = &IPClaimFilter{Namespace: namespace}
			}
			return filter, nil
		}

		return &IPClaimFilter{}, nil
	}

}

type IPClaimFilter struct {
	// Name filters by the name of the objects
	Name string `protobuf:"bytes,1,opt,name=name"`

	// Namespace filters by the namespace of the objects
	Namespace string `protobuf:"bytes,2,opt,name=namespace"`
}

func (r *IPClaimFilter) Filter(ctx context.Context, obj runtime.Object) bool {
	f := false // result of the previous filter
	o, ok := obj.(*IPClaim)
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

func (r *IPClaim) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	// status cannot be set upon create -> reset it
	newobj := obj.(*IPClaim)
	newobj.Status = IPClaimStatus{}
}

// ValidateCreate statically validates
func (r *IPClaim) ValidateCreate(ctx context.Context, obj runtime.Object) field.ErrorList {
	newobj := obj.(*IPClaim)
	return newobj.ValidateSyntax("")
}

func (r *IPClaim) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	// ensure the status dont get updated
	newobj := obj.(*IPClaim)
	oldObj := old.(*IPClaim)
	newobj.Status = oldObj.Status
}

func (r *IPClaim) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newobj := obj.(*IPClaim)
	return newobj.ValidateSyntax("")
}
