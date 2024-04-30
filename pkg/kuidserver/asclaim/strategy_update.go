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

package asclaim

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"

	"github.com/henderiw/apiserver-store/pkg/storebackend"
	"github.com/henderiw/logger/log"
	asbe1v1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
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

	claim, ok := obj.(*asbe1v1alpha1.ASClaim)
	if !ok {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath(""),
			claim,
			fmt.Errorf("unexpected new object, expecting: %s, got: %s", asbe1v1alpha1.ASClaimKind, reflect.TypeOf(obj)).Error(),
		))
		return allErrs
	}
	allErrs = claim.ValidateSyntax()
	if len(allErrs) != 0 {
		return allErrs
	}
	// TODO validate mutations
	return claim.ValidateSyntax()
}

func (r *strategy) Update(ctx context.Context, key types.NamespacedName, obj, old runtime.Object, dryrun bool) (runtime.Object, error) {
	log := log.FromContext(ctx)
	// check if there is a change
	newObj, ok := obj.(*asbe1v1alpha1.ASClaim)
	if !ok {
		return obj, fmt.Errorf("unexpected new object, expecting: %s, got: %s", asbe1v1alpha1.ASClaimKind, reflect.TypeOf(obj))
	}
	oldObj, ok := old.(*asbe1v1alpha1.ASClaim)
	if !ok {
		return obj, fmt.Errorf("unexpected old object, expecting: %s, got: %s", asbe1v1alpha1.ASClaimKind, reflect.TypeOf(obj))
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

	/*
		if err := r.be.ValidateClaim(ctx, obj); err != nil {
			return obj, apierrors.NewInternalError(err)
		}

		if err := r.be.Claim(ctx, obj); err != nil {
			return obj, apierrors.NewInternalError(err)
		}
	*/

	log.Info("update ipclaim storage", "key", key, "obj", obj)

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
