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

package vlan

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/henderiw/logger/log"
	vlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
)

func initIndexContext(ctx context.Context, op string, idx *vlanbev1alpha1.VLANIndex) context.Context {
	l := log.FromContext(ctx).
		With(
			"op", fmt.Sprintf("%s index", op),
			"nsn", idx.GetNamespacedName().String(),
		)
	return log.IntoContext(ctx, l)
}

func initClaimContext(ctx context.Context, op string, claim *vlanbev1alpha1.VLANClaim) context.Context {
	var l *slog.Logger

	claimType := claim.GetClaimType()
	switch claimType {
	case vlanbev1alpha1.VLANClaimType_DynamicID:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s %s claim", op, string(claimType)),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Index,
			)
	case vlanbev1alpha1.VLANClaimType_StaticID:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s %s claim", op, string(claimType)),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Index,
				"id", *claim.Spec.ID,
			)
	case vlanbev1alpha1.VLANClaimType_Range:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s %s claim", op, string(claimType)),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.Spec.Index,
				"range", *claim.Spec.Range,
			)
	}
	return log.IntoContext(ctx, l)
}
