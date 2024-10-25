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
	fmt "fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// +kubebuilder:object:generate=false
type SyntaxValidator interface {
	Validate(claim *VLANClaim) field.ErrorList
}

type VLANDynamicIDSyntaxValidator struct {
	name string
}

func (r *VLANDynamicIDSyntaxValidator) Validate(claim *VLANClaim) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

type VLANStaticIDSyntaxValidator struct {
	name string
}

func (r *VLANStaticIDSyntaxValidator) Validate(claim *VLANClaim) field.ErrorList {
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

type VLANRangeSyntaxValidator struct {
	name string
}

func (r *VLANRangeSyntaxValidator) Validate(claim *VLANClaim) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateVLANRange(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.range"),
			claim,
			fmt.Errorf("invalid vlan range %s, err: %s", r.name, err.Error()).Error(),
		))
	}
	return allErrs
}
