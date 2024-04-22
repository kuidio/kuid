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
	"github.com/henderiw/logger/log"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
)

func (r *staticAddressApplicator) Apply(ctx context.Context, claim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx).With("name", claim.GetName())
	log.Info("static address claim")
	pi, err := iputil.New(*claim.Spec.Address)
	if err != nil {
		return err
	}
	fmt.Println("applyAddress", *claim.Spec.Address, r.parentClaimSummaryType, r.parentRangeName, r.parentNetwork, r.parentLabels)
	if r.parentClaimSummaryType == ipambev1alpha1.IPClaimSummaryType_Range {
		if err := r.applyAddressInRange(ctx, claim, pi, r.parentRangeName, r.parentLabels); err != nil {
			return err
		}
	} else {
		if err := r.apply(ctx, claim, []*iputil.Prefix{pi}, false, r.parentLabels); err != nil {
			return err
		}
	}

	r.updateClaimAddressStatus(ctx, claim, pi, r.parentNetwork)
	return nil
}
