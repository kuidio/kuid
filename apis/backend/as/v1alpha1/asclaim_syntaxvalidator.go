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
	Validate(claim *ASClaim) field.ErrorList
}

type ASRangeSyntaxValidator struct {
	name string
}

func (r *ASRangeSyntaxValidator) Validate(claim *ASClaim) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateASRange(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.range"),
			claim,
			fmt.Errorf("invalid AS range %s", r.name).Error(),
		))
	}
	return allErrs
}

type ASDynamicIDSyntaxValidator struct {
	name string
}

func (r *ASDynamicIDSyntaxValidator) Validate(claim *ASClaim) field.ErrorList {
	var allErrs field.ErrorList
	return allErrs
}

type ASStaticIDSyntaxValidator struct {
	name string
}

func (r *ASStaticIDSyntaxValidator) Validate(claim *ASClaim) field.ErrorList {
	var allErrs field.ErrorList
	if err := claim.ValidateASID(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.id"),
			claim,
			fmt.Errorf("invalid AS id %s", r.name).Error(),
		))
	}
	return allErrs
}
