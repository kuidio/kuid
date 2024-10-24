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

package genericserver

/*

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
	genObj, ok := obj.(backend.GenericObject)
	if !ok {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath(""),
			genObj,
			fmt.Errorf("unexpected new object got: %s", reflect.TypeOf(obj)).Error(),
		))
		return allErrs
	}

	return genObj.ValidateSyntax("")
}

func (r *strategy) Update(ctx context.Context, key types.NamespacedName, obj, old runtime.Object, dryrun bool) (runtime.Object, error) {
	log := log.FromContext(ctx)
	// check if there is a change
	newGenObj, ok := obj.(backend.GenericObject)
	if !ok {
		return obj, fmt.Errorf("unexpected new object, got: %s", reflect.TypeOf(obj))
	}
	oldGenObj, ok := old.(backend.GenericObject)
	if !ok {
		return obj, fmt.Errorf("unexpected old object, got: %s", reflect.TypeOf(obj))
	}

	newHash, err := newGenObj.CalculateHash()
	if err != nil {
		return obj, err
	}
	oldHash, err := oldGenObj.CalculateHash()
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
*/
