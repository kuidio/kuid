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

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"

	"github.com/henderiw/apiserver-store/pkg/storebackend"
	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/apis/backend"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/watch"
)

func (r *strategy) BeginUpdate(ctx context.Context) error { return nil }

func (r *strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (r *strategy) AllowCreateOnUpdate() bool { return false }

func (r *strategy) AllowUnconditionalUpdate() bool { return false }

func (r *strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	var allErrs field.ErrorList

	claim, ok := obj.(backend.ClaimObject)
	if !ok {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath(""),
			claim,
			fmt.Errorf("unexpected new object, got: %s", reflect.TypeOf(obj)).Error(),
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

func (r *strategy) Update(ctx context.Context, key types.NamespacedName, obj, old runtime.Object, dryrun bool) (runtime.Object, error) {
	log := log.FromContext(ctx)
	// check if there is a change
	newObj, ok := obj.(backend.ClaimObject)
	if !ok {
		return obj, fmt.Errorf("unexpected new object, got: %s", reflect.TypeOf(obj))
	}
	oldObj, ok := old.(backend.ClaimObject)
	if !ok {
		return obj, fmt.Errorf("unexpected old object, got: %s", reflect.TypeOf(obj))
	}

	newHash, err := newObj.CalculateHash()
	if err != nil {
		return obj, err
	}
	oldHash, err := oldObj.CalculateHash()
	if err != nil {
		return obj, err
	}

	if oldHash == newHash {
		log.Debug("update nothing to do", "oldHash", hex.EncodeToString(oldHash[:]), "newHash", hex.EncodeToString(newHash[:]))
		return obj, nil
	}
	log.Debug("updating", "oldHash", hex.EncodeToString(oldHash[:]), "newHash", hex.EncodeToString(newHash[:]))
	if dryrun {
		return obj, nil
	}
	if err := updateResourceVersion(ctx, obj, old); err != nil {
		return obj, apierrors.NewInternalError(err)
	}
	if err := r.store.Update(ctx, storebackend.KeyFromNSN(key), obj); err != nil {
		return obj, apierrors.NewInternalError(err)
	}
	r.notifyWatcher(ctx, watch.Event{
		Type:   watch.Modified,
		Object: obj,
	})
	return obj, nil
}

func (r *strategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

func updateResourceVersion(_ context.Context, obj, old runtime.Object) error {
	accessorNew, err := meta.Accessor(obj)
	if err != nil {
		return nil
	}
	accessorOld, err := meta.Accessor(old)
	if err != nil {
		return nil
	}
	resourceVersion, err := strconv.Atoi(accessorOld.GetResourceVersion())
	if err != nil {
		return err
	}
	resourceVersion++
	accessorNew.SetResourceVersion(strconv.Itoa(resourceVersion))
	return nil
}
