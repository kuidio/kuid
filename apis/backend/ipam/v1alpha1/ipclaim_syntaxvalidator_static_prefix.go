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
	"k8s.io/utils/ptr"
)

type staticPrefixSyntaxValidator struct {
	name string
}

func (r *staticPrefixSyntaxValidator) Validate(claim *IPClaim) field.ErrorList {
	var allErrs field.ErrorList
	if claim.Spec.Prefix == nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.prefix"),
			claim,
			fmt.Errorf("%s requires a prefix", r.name).Error(),
		))
	}
	pi, err := iputil.New(*claim.Spec.Prefix)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.prefix"),
			claim,
			fmt.Errorf("%s invalid prefix; err %s", r.name, err.Error()).Error(),
		))
	}
	// this is for user convenience
	if claim.Spec.CreatePrefix == nil {
		claim.Spec.CreatePrefix = ptr.To[bool](true)
	}
	if claim.Spec.PrefixLength == nil {
		claim.Spec.PrefixLength = ptr.To[uint32](uint32(pi.GetPrefixLength()))
	}
	if pi.IsAddressPrefix() {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.prefix"),
			claim,
			fmt.Sprintf("%s invalid prefix, no /32 for ipv4 or /128 for ipv6 notationallowed, got %s", r.name, *claim.Spec.Prefix),
		))
	}
	ipPrefixType := claim.GetIPPrefixType()
	if ipPrefixType == IPPrefixType_Invalid {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.type"),
			claim,
			fmt.Errorf("%s, invalid claim type, got %s", r.name, string(ipPrefixType)).Error(),
		))
	}
	if ipPrefixType != IPPrefixType_Network {
		if pi.GetIPSubnet().String() != pi.GetIPPrefix().String() {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.prefix"),
				claim,
				fmt.Sprintf("%s, invalid prefix net <> address is not allowed for claimType: %s", r.name, string(ipPrefixType)),
			))
		}
	}

	if claim.Spec.DefaultGateway != nil {
		if ipPrefixType != IPPrefixType_Network {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.type"),
				claim,
				fmt.Errorf("%s cannot have a defaultGateway on a prefix which is not a network type", r.name).Error(),
			))
		}
		if !pi.IsNorLastNorFirst() {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.prefix"),
				claim,
				fmt.Errorf("%s cannot have a defaultGateway on a prefix which is the first or last ip in the prefix", r.name).Error(),
			))
		}
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
