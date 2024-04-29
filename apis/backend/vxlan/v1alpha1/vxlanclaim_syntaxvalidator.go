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
	Validate(claim *VXLANClaim) field.ErrorList
}

type VXLANRangeSyntaxValidator struct {
	name string
}

func (r *VXLANRangeSyntaxValidator) Validate(claim *VXLANClaim) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateVXLANRange(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.range"),
			claim,
			fmt.Errorf("invalid VXLAN range %s", r.name).Error(),
		))
	}
	return allErrs
}

type VXLANDynamicIDSyntaxValidator struct {
	name string
}

func (r *VXLANDynamicIDSyntaxValidator) Validate(claim *VXLANClaim) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

type VXLANStaticIDSyntaxValidator struct {
	name string
}

func (r *VXLANStaticIDSyntaxValidator) Validate(claim *VXLANClaim) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateVXLANID(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.id"),
			claim,
			fmt.Errorf("invalid VXLAN id %s", r.name).Error(),
		))
	}
	return allErrs
}
