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
	fmt "fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

type SyntaxValidator interface {
	Validate(claim *VLANClaim) field.ErrorList
}

type vlanRangeSyntaxValidator struct {
	name string
}

func (r *vlanRangeSyntaxValidator) Validate(claim *VLANClaim) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateVLANRange(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.range"),
			claim,
			fmt.Errorf("invalid vlan range %s", r.name).Error(),
		))
	}
	return allErrs
}

type vlanDynamicIDSyntaxValidator struct {
	name string
}

func (r *vlanDynamicIDSyntaxValidator) Validate(claim *VLANClaim) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

type vlanStaticIDSyntaxValidator struct {
	name string
}

func (r *vlanStaticIDSyntaxValidator) Validate(claim *VLANClaim) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateVLANID(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.id"),
			claim,
			fmt.Errorf("invalid vlan id %s", r.name).Error(),
		))
	}
	return allErrs
}

type vlanSizeSyntaxValidator struct {
	name string
}

func (r *vlanSizeSyntaxValidator) Validate(claim *VLANClaim) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateVLANSize(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.size"),
			claim,
			fmt.Errorf("invalid vlan id %s", r.name).Error(),
		))
	}
	return allErrs
}
