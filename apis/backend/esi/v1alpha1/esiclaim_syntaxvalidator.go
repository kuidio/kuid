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
	Validate(claim *ESIClaim) field.ErrorList
}

type ESIRangeSyntaxValidator struct {
	name string
}

func (r *ESIRangeSyntaxValidator) Validate(claim *ESIClaim) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateESIRange(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.range"),
			claim,
			fmt.Errorf("invalid ESI range %s", r.name).Error(),
		))
	}
	return allErrs
}

type ESIDynamicIDSyntaxValidator struct {
	name string
}

func (r *ESIDynamicIDSyntaxValidator) Validate(claim *ESIClaim) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

type ESIStaticIDSyntaxValidator struct {
	name string
}

func (r *ESIStaticIDSyntaxValidator) Validate(claim *ESIClaim) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateESIID(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.id"),
			claim,
			fmt.Errorf("invalid ESI id %s", r.name).Error(),
		))
	}
	return allErrs
}
