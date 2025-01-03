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

package claimserver

/*

import (
	"context"
	"fmt"
	"reflect"

	"github.com/henderiw/apiserver-store/pkg/storebackend"
	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/apis/backend"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/watch"
)

func (r *strategy) BeginCreate(ctx context.Context) error { return nil }

func (r *strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (r *strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	var allErrs field.ErrorList

	claim, ok := obj.(backend.ClaimObject)
	if !ok {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath(""),
			claim,
			fmt.Errorf("unexpected new object got: %s", reflect.TypeOf(obj)).Error(),
		))
		return allErrs
	}

	if r.serverObjContext.NewIndexFn != nil {
		index := r.serverObjContext.NewIndexFn()
		if err := r.client.Get(ctx, types.NamespacedName{Namespace: claim.GetNamespace(), Name: claim.GetIndex()}, index); err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("spec.index"),
				claim,
				fmt.Errorf("index does not exist cannot validate syntax").Error(),
			))
			return allErrs
		}

		return claim.ValidateSyntax(index.GetType())
	}
	return claim.ValidateSyntax("")
}

func (r *strategy) Create(ctx context.Context, key types.NamespacedName, obj runtime.Object, dryrun bool) (runtime.Object, error) {
	log := log.FromContext(ctx)
	if dryrun {
		return obj, nil
	}

	log.Info("create claim in storage", "key", key, "obj", obj)

	if err := r.store.Create(ctx, storebackend.KeyFromNSN(key), obj); err != nil {
		return obj, apierrors.NewInternalError(err)
	}
	r.notifyWatcher(ctx, watch.Event{
		Type:   watch.Added,
		Object: obj,
	})
	return obj, nil
}

func (r *strategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}
*/