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

type dynamicAddressSyntaxValidator struct {
	name string
}

func (r *dynamicAddressSyntaxValidator) Validate(claim *IPClaim) field.ErrorList {
	var allErrs field.ErrorList

	ipPrefixType := claim.GetIPPrefixType()
	if ipPrefixType == IPPrefixType_Invalid {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.type"),
			claim,
			fmt.Errorf("%s, invalid claim type, got %s", r.name, string(ipPrefixType)).Error(),
		))
	}

	if claim.Spec.DefaultGateway != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.defaultGateway"),
			claim,
			fmt.Errorf("%s cannot have a defaultGateway", r.name).Error(),
		))

	}
	if claim.Spec.CreatePrefix != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.createPrefix"),
			claim,
			fmt.Errorf("%s cannot have a createPrefix", r.name).Error(),
		))
	}
	if claim.Spec.PrefixLength != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.prefixLength"),
			claim,
			fmt.Errorf("%s cannot have a prefixLength", r.name).Error(),
		))
	}

	return allErrs
}
