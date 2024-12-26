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

package extcomm

import (
	fmt "fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// +kubebuilder:object:generate=false
// +k8s:deepcopy-gen:false
type SyntaxValidator interface {
	Validate(claim *EXTCOMMClaim, extCommType ExtendedCommunityType) field.ErrorList
}

type EXTCOMMRangeSyntaxValidator struct {
	name string
}

func (r *EXTCOMMRangeSyntaxValidator) Validate(claim *EXTCOMMClaim, extCommType ExtendedCommunityType) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateEXTCOMMRange(extCommType); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.range"),
			claim,
			fmt.Errorf("invalid EXTCOMM range %s", r.name).Error(),
		))
	}
	return allErrs
}

type EXTCOMMDynamicIDSyntaxValidator struct {
	name string
}

func (r *EXTCOMMDynamicIDSyntaxValidator) Validate(claim *EXTCOMMClaim, extCommType ExtendedCommunityType) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

type EXTCOMMStaticIDSyntaxValidator struct {
	name string
}

func (r *EXTCOMMStaticIDSyntaxValidator) Validate(claim *EXTCOMMClaim, extCommType ExtendedCommunityType) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateEXTCOMMID(extCommType); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.id"),
			claim,
			fmt.Errorf("invalid EXTCOMM id %s", r.name).Error(),
		))
	}
	return allErrs
}
