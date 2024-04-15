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

	"github.com/henderiw/iputil"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type prefixClaimSyntaxValidator struct{}

func (r *prefixClaimSyntaxValidator) ValidateSyntax(ctx context.Context, claim *ipambev1alpha1.IPClaim) field.ErrorList {
	//log := log.FromContext(ctx)
	var allErrs field.ErrorList
	// ipClaim with prefix
	pi, err := iputil.New(*claim.Spec.Prefix)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec.prefix"),
			claim,
			err.Error(),
		))
		return allErrs
	}
	return r.validateSyntax(ctx, claim, pi)
}

func (r *prefixClaimSyntaxValidator) validateSyntax(_ context.Context, ipClaim *ipambev1alpha1.IPClaim, pi *iputil.Prefix) field.ErrorList {
	var allErrs field.ErrorList

	if ipClaim.Spec.Kind == ipambev1alpha1.PrefixKindAggregate || ipClaim.Spec.Kind == ipambev1alpha1.PrefixKindNetwork {
		if pi.IsAddressPrefix() {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.prefix"),
				ipClaim,
				fmt.Sprintf("a prefix claim for kind %s is not allowed with /32 for ipv4 or /128 for ipv6 notation", ipClaim.Spec.Kind),
			))
		}
	}

	if ipClaim.Spec.Kind != ipambev1alpha1.PrefixKindNetwork {
		if pi.GetIPSubnet().String() != pi.GetIPPrefix().String() {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.prefix"),
				ipClaim,
				fmt.Sprintf("net <> address is not allowed for prefixkind: %s", ipClaim.Spec.Kind),
			))
		}
	}
	if ipClaim.Spec.CreatePrefix != nil {
		if pi.IsAddressPrefix() {
			// a create prefix should have an address different from /32 or /128
			if pi.GetIPSubnet().String() != pi.GetIPPrefix().String() {
				allErrs = append(allErrs, field.Invalid(
					field.NewPath("spec.prefix"),
					ipClaim,
					fmt.Sprintf("create prefix is not allowed with /32 for ipv4 or /128 for ipv6, got: %s", pi.GetIPPrefix().String()),
				))
			}
		}
	}
	return allErrs
}
