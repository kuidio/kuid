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

package genid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/henderiw/store"
	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/backend"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var _ backend.EntryObject = &GENIDEntry{}

func (r *GENIDEntry) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}
func (r *GENIDEntry) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

// GetCondition returns the condition based on the condition kind
func (r *GENIDEntry) GetCondition(t condition.ConditionType) condition.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *GENIDEntry) SetConditions(c ...condition.Condition) {
	r.Status.SetConditions(c...)
}

func (r *GENIDEntry) ValidateSyntax(s string) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

func (r *GENIDEntry) GetIndex() string                { return r.Spec.Index }
func (r *GENIDEntry) IsIndexEntry() bool              { return r.Spec.IndexEntry }
func (r *GENIDEntry) GetClaimType() backend.ClaimType { return r.Spec.ClaimType }
func (r *GENIDEntry) GetSpecID() string               { return r.Spec.ID }

func (r *GENIDEntry) GetChoreoAPIVersion() string {
	return schema.GroupVersion{Group: GroupName, Version: "genid"}.String()
}

func GENIDEntryFromRuntime(ru runtime.Object) (backend.EntryObject, error) {
	entry, ok := ru.(*GENIDEntry)
	if !ok {
		return nil, errors.New("runtime object not ASIndex")
	}
	return entry, nil
}

func GENIDEntryFromUnstructured(ru runtime.Unstructured) (backend.EntryObject, error) {
	obj := &GENIDEntry{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(ru.UnstructuredContent(), obj)
	if err != nil {
		return nil, fmt.Errorf("error converting unstructured: %v", err)
	}
	return obj, nil
}

func GetGENIDEntry(k store.Key, vrange, id string, labels map[string]string) backend.EntryObject {
	index := k.Name
	ns := k.Namespace

	spec := &GENIDEntrySpec{
		Index:     index,
		ClaimType: backend.GetClaimTypeFromString(labels[backend.KuidClaimTypeKey]),
		ID:        id,
	}
	// filter the system defined labels from the labels to prepare for the user defined labels
	udLabels := map[string]string{}
	for k, v := range labels {
		if !backend.BackendSystemKeys.Has(k) {
			udLabels[k] = v
		}
	}
	spec.UserDefinedLabels.Labels = udLabels

	id = strings.ReplaceAll(id, "/", "-")
	name := fmt.Sprintf("%s.%s", index, id)
	if vrange != "" {
		name = fmt.Sprintf("%s.%s", vrange, id)
	}

	return BuildGENIDEntry(
		metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			OwnerReferences: []metav1.OwnerReference{
				{
					// this is a bit of a hack for choreo to ensure we point to the proper external reference
					APIVersion: schema.GroupVersion{Group: SchemeGroupVersion.Group, Version: "v1alpha1"}.Identifier(),
					Kind: labels[backend.KuidOwnerKindKey],
					Name: labels[backend.KuidClaimNameKey],
					UID:  types.UID(labels[backend.KuidClaimUIDKey]),
				},
			},
		},
		spec,
		nil,
	)
}

func BuildGENIDEntry(meta metav1.ObjectMeta, spec *GENIDEntrySpec, status *GENIDEntryStatus) backend.EntryObject {
	vlanspec := GENIDEntrySpec{}
	if spec != nil {
		vlanspec = *spec
	}
	vlanstatus := GENIDEntryStatus{}
	if status != nil {
		vlanstatus = *status
	}
	return &GENIDEntry{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       GENIDEntryKind,
		},
		ObjectMeta: meta,
		Spec:       vlanspec,
		Status:     vlanstatus,
	}
}
