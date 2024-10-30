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
	"log/slog"

	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/apis/backend/ipam"
)

// this is implemented locally since the ip claim does not follow the backend.ClaimObject structure
func initClaimContext(ctx context.Context, op string, claim *ipam.IPClaim) context.Context {
	var l *slog.Logger

	ipClaimType, err := claim.GetIPClaimType()
	if err != nil {
		return ctx
	}
	switch ipClaimType {
	case ipam.IPClaimType_DynamicAddress:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s dynamic address claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Index,
			)
	case ipam.IPClaimType_DynamicPrefix:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s dynamic prefix claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Index,
				"prefixType", claim.GetIPPrefixType(),
			)
	case ipam.IPClaimType_StaticAddress:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s static address claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Idx,
				"address", *claim.Spec.Address,
			)
	case ipam.IPClaimType_StaticPrefix:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s static prefix claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Index,
				"prefix", *claim.Spec.Prefix,
				"prefixType", claim.GetIPPrefixType(),
			)
	case ipam.IPClaimType_StaticRange:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s static range claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Index,
				"range", *claim.Spec.Range,
			)
	}
	return log.IntoContext(ctx, l)
}
