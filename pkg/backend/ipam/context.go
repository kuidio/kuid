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
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
)

func initClaimContext(ctx context.Context, op string, claim *ipambev1alpha1.IPClaim) context.Context {
	var l *slog.Logger

	ipClaimType, err := claim.GetIPClaimType()
	if err != nil {
		return ctx
	}
	switch ipClaimType {
	case ipambev1alpha1.IPClaimType_DynamicAddress:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s dynamic address claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Index,
			)
	case ipambev1alpha1.IPClaimType_DynamicPrefix:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s dynamic prefix claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Index,
				"prefixType", claim.GetIPPrefixType(),
			)
	case ipambev1alpha1.IPClaimType_StaticAddress:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s static address claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Idx,
				"address", *claim.Spec.Address,
			)
	case ipambev1alpha1.IPClaimType_StaticPrefix:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s static prefix claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Index,
				"prefix", *claim.Spec.Prefix,
				"prefixType", claim.GetIPPrefixType(),
			)
	case ipambev1alpha1.IPClaimType_StaticRange:
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
