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

	"github.com/henderiw/iputil"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type staticAddressSyntaxValidator struct {
	name string
}

func (r *staticAddressSyntaxValidator) Validate(claim *IPClaim) field.ErrorList {
	var allErrs field.ErrorList
	if claim.Spec.Address == nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.address"),
			claim,
			fmt.Errorf("%s requires a address", r.name).Error(),
		))
	}
	pi, err := iputil.New(*claim.Spec.Address)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.address"),
			claim,
			fmt.Errorf("%s invalid prefix; err %s", r.name, err.Error()).Error(),
		))
	}
	if pi.IsAddressPrefix() { //
		if claim.Spec.DefaultGateway != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.defaultGateway"),
				claim,
				fmt.Errorf("%s cannot have a defaultGateway on a prefix which is not a network type", r.name).Error(),
			))
		}

	} else {
		if !((pi.GetAddressFamily() == iputil.AddressFamilyIpv4 && pi.GetPrefixLength().Int() == 31) ||
			(pi.GetAddressFamily() == iputil.AddressFamilyIpv6 && pi.GetPrefixLength().Int() == 127)) {
			if !pi.IsNorLastNorFirst() {
				allErrs = append(allErrs, field.Invalid(
					field.NewPath("spec.prefix"),
					claim,
					fmt.Errorf("%s cannot have a address that is the first or last address in a prefix based address", r.name).Error(),
				))
			}
		}
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
	if claim.Spec.AddressFamily != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.addressFamily"),
			claim,
			fmt.Errorf("%s cannot have a addressFamily", r.name).Error(),
		))
	}
	if claim.Spec.Idx != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.index"),
			claim,
			fmt.Errorf("%s cannot have a index", r.name).Error(),
		))
	}
	return allErrs
}
