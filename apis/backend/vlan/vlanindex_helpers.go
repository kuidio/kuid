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

package vlan

import (
	"fmt"

	"github.com/kform-dev/choreo/apis/condition"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// GetCondition returns the condition based on the condition kind
func (r *VLANIndex) GetCondition(t condition.ConditionType) condition.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *VLANIndex) SetConditions(c ...condition.Condition) {
	r.Status.SetConditions(c...)
}

func (r *VLANIndex) ValidateSyntax(s string) field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.MinID != nil {
		if err := validateVLANID(int(*r.Spec.MinID)); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.minID"),
				r,
				fmt.Errorf("invalid vlan ID %d", *r.Spec.MinID).Error(),
			))
		}
	}
	if r.Spec.MaxID != nil {
		if err := validateVLANID(int(*r.Spec.MaxID)); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.maxID"),
				r,
				fmt.Errorf("invalid vlan ID %d", *r.Spec.MaxID).Error(),
			))
		}
	}
	if r.Spec.MinID != nil && r.Spec.MaxID != nil {
		if *r.Spec.MinID > *r.Spec.MaxID {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.maxID"),
				r,
				fmt.Errorf("min vlan ID %d cannot be bigger than max vlan ID %d", *r.Spec.MinID, *r.Spec.MaxID).Error(),
			))
		}
	}
	return allErrs
}

// BuildVLANIndex returns a reource from a client Object a Spec/Status
func BuildVLANIndex(meta metav1.ObjectMeta, spec *VLANIndexSpec, status *VLANIndexStatus) *VLANIndex {
	aspec := VLANIndexSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := VLANIndexStatus{}
	if status != nil {
		astatus = *status
	}
	return &VLANIndex{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       VLANIndexKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}
