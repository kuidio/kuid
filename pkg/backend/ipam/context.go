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

func initIndexContext(ctx context.Context, op string, idx *ipambev1alpha1.IPIndex) context.Context {
	l := log.FromContext(ctx).
		With(
			"op", fmt.Sprintf("%s index", op),
			"nsn", idx.GetNamespacedName().String(),
		)
	return log.IntoContext(ctx, l)
}

func initClaimContext(ctx context.Context, op string, claim *ipambev1alpha1.IPClaim) context.Context {
	var l *slog.Logger

	addressing, err := claim.GetAddressing()
	if err != nil {
		return ctx
	}
	switch addressing {
	case ipambev1alpha1.IPClaimAddressing_DynamicAddress:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s dynamic address claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"ni", claim.Spec.NetworkInstance,
			)
	case ipambev1alpha1.IPClaimAddressing_DynamicPrefix:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s dynamic prefix claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"ni", claim.Spec.NetworkInstance,
			)
	case ipambev1alpha1.IPClaimAddressing_StaticAddress:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s static address claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"ni", claim.Spec.NetworkInstance,
				"address", *claim.Spec.Address,
			)
	case ipambev1alpha1.IPClaimAddressing_StaticPrefix:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s static prefix claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"ni", claim.Spec.NetworkInstance,
				"prefix", *claim.Spec.Prefix,
			)
	case ipambev1alpha1.IPClaimAddressing_StaticRange:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s static range claim", op),
				"nsn", claim.GetNamespacedName().String(),
				"ni", claim.Spec.NetworkInstance,
				"range", *claim.Spec.Range,
			)
	}

	return log.IntoContext(ctx, l)
}
