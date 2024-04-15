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

	"github.com/henderiw/logger/log"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
)

type dynamicApplicator struct {
	applicator
}

func (r *dynamicApplicator) Apply(ctx context.Context, claim *ipambev1alpha1.IPClaim) error {
	log := log.FromContext(ctx).With("name", claim.GetName(), "kind", claim.Spec.Kind)
	log.Info("dynamic claim")
	
	// claim a prefix
	pi, err := r.claimPrefix(ctx, claim)
	if err != nil {
		return err
	}
	return r.apply(ctx, claim, pi)
}