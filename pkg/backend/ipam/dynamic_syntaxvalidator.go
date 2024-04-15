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

package ipam

import (
	"context"
	"fmt"

	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type dynamicClaimSyntaxValidator struct{}

func (r *dynamicClaimSyntaxValidator) ValidateSyntax(_ context.Context, claim *ipambev1alpha1.IPClaim) field.ErrorList {
	var allErrs field.ErrorList
	// dynamic entries with aggregate prefix kind not supported
	if claim.Spec.Kind == ipambev1alpha1.PrefixKindAggregate {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath(""),
			claim,
			fmt.Sprintf("a dynamic prefix claim is not supported for: %s", claim.Spec.Kind),
		))
		return allErrs
	}
	// a dynamic prefix claim has to set the prefixLength
	if claim.Spec.CreatePrefix != nil && claim.Spec.PrefixLength == nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.prefixLength"),
			claim,
			"a dynamic prefix claim has to specify the prefixLength",
		))
		return allErrs
	}
	if claim.Spec.PrefixLength != nil && claim.Spec.CreatePrefix == nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.createPrefix"),
			claim,
			"a dynamic prefix with prefixLength set has to also set a createPrefix",
		))
		return allErrs
	}
	// TODO Pool
	return allErrs
}
