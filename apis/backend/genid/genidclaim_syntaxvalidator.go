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
	fmt "fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// +kubebuilder:object:generate=false
type SyntaxValidator interface {
	Validate(claim *GENIDClaim, genidType GENIDType) field.ErrorList
}

type GENIDRangeSyntaxValidator struct {
	name string
}

func (r *GENIDRangeSyntaxValidator) Validate(claim *GENIDClaim, genidType GENIDType) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateGENIDRange(genidType); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.range"),
			claim,
			fmt.Errorf("invalid GENID range %s", r.name).Error(),
		))
	}
	return allErrs
}

type GENIDDynamicIDSyntaxValidator struct {
	name string
}

func (r *GENIDDynamicIDSyntaxValidator) Validate(claim *GENIDClaim, genidType GENIDType) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

type GENIDStaticIDSyntaxValidator struct {
	name string
}

func (r *GENIDStaticIDSyntaxValidator) Validate(claim *GENIDClaim, genidType GENIDType) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateGENIDID(genidType); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.id"),
			claim,
			fmt.Errorf("invalid GENID id %s", r.name).Error(),
		))
	}
	return allErrs
}