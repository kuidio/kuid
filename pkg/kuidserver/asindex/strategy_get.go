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

package asindex

import (
	"context"

	"github.com/henderiw/apiserver-store/pkg/storebackend"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func (r *strategy) Get(ctx context.Context, key types.NamespacedName) (runtime.Object, error) {
	obj, err := r.store.Get(ctx, storebackend.KeyFromNSN(key))
	if err != nil {
		return nil, apierrors.NewNotFound(r.gr, key.Name)
	}
	return obj, nil
}