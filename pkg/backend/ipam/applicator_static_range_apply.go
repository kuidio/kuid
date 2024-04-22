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

	"github.com/henderiw/iputil"
	"github.com/henderiw/logger/log"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	"go4.org/netipx"
)

func (r *staticRangeApplicator) Apply(ctx context.Context, claim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx).With("name", claim.GetName())
	log.Info("static range claim")
	ipRange, err := netipx.ParseIPRange(*claim.Spec.Range)
	if err != nil {
		return err
	}

	// add each prefix in the routing table -> we convey them all together
	pis := make([]*iputil.Prefix, 0, len(ipRange.Prefixes()))
	for _, prefix := range ipRange.Prefixes() {
		pis = append(pis, iputil.NewPrefixInfo(prefix))
	}
	if err := r.apply(ctx, claim, pis, false, map[string]string{}); err != nil {
		return err
	}
	if err := r.applyRange(ctx, claim, ipRange); err != nil {
		return err
	}

	r.updateClaimRangeStatus(ctx, claim)
	return nil
}
