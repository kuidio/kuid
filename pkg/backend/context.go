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

package backend

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/apis/backend"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func InitIndexContext(ctx context.Context, op string, idx client.Object) context.Context {
	l := log.FromContext(ctx).
		With(
			"op", fmt.Sprintf("%s index", op),
			"gvk", idx.GetObjectKind().GroupVersionKind().String(),
			"nsn", types.NamespacedName{Namespace: idx.GetNamespace(), Name: idx.GetName()}.String(),
		)
	return log.IntoContext(ctx, l)
}

func InitClaimContext(ctx context.Context, op string, claim backend.ClaimObject) context.Context {
	var l *slog.Logger

	claimType := claim.GetClaimType()
	switch claimType {
	case backend.ClaimType_DynamicID:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s %s claim", op, string(claimType)),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.GetIndex(),
			)
	case backend.ClaimType_StaticID:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s %s claim", op, string(claimType)),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.GetIndex(),
				"id", *claim.GetStaticID(), // safe
			)
	case backend.ClaimType_Range:
		l = log.FromContext(ctx).
			With(
				"op", fmt.Sprintf("%s %s claim", op, string(claimType)),
				"nsn", claim.GetNamespacedName().String(),
				"index", claim.GetIndex(),
				"range", *claim.GetRange(), // safe
			)
	}
	return log.IntoContext(ctx, l)
}
