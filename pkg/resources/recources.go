/*
Copyright 2023 The Nephio Authors.

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

package resources

import (
	"context"
	"fmt"
	"sync"

	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/pkg/reconcilers/resource"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Config struct {
	Owns []backend.GenericObject
}

func New(c client.Client, cfg Config) *Resources {
	return &Resources{
		Client:            c,
		cfg:               cfg,
		newResources:      map[corev1.ObjectReference]backend.GenericObject{},
		existingResources: map[corev1.ObjectReference]backend.GenericObject{},
	}
}

type Resources struct {
	client.Client
	cfg               Config
	m                 sync.RWMutex
	newResources      map[corev1.ObjectReference]backend.GenericObject
	existingResources map[corev1.ObjectReference]backend.GenericObject
	matchLabels       client.MatchingLabels
}

// Init initializes the new and exisiting resource inventory list
/*
func (r *Resources) Init(ml client.MatchingLabels) {
	r.matchLabels = ml
	r.newResources = map[corev1.ObjectReference]client.Object{}
	r.existingResources = map[corev1.ObjectReference]client.Object{}
}
*/

// AddNewResource adds a new resource to the inventoru
func (r *Resources) AddNewResource(ctx context.Context, cr client.Object, o backend.GenericObject) error {
	r.m.Lock()
	defer r.m.Unlock()

	log := log.FromContext(ctx)

	o.SetOwnerReferences([]metav1.OwnerReference{
		{
			APIVersion: cr.GetObjectKind().GroupVersionKind().GroupVersion().String(),
			Kind:       cr.GetObjectKind().GroupVersionKind().Kind,
			Name:       cr.GetName(),
			UID:        cr.GetUID(),
			Controller: ptr.To(true),
		},
	})

	ref := corev1.ObjectReference{
		APIVersion: o.GetObjectKind().GroupVersionKind().GroupVersion().String(),
		Kind:       o.GetObjectKind().GroupVersionKind().Kind,
		Namespace:  o.GetNamespace(),
		Name:       o.GetName(),
	}

	log.Info("add newresource", "ref", ref.String())

	r.newResources[ref] = o
	return nil
}

// GetExistingResources retrieves the exisiting resource that match the label selector and the owner reference
// and puts the results in the resource inventory
func (r *Resources) getExistingResources(ctx context.Context, cr client.Object) error {
	log := log.FromContext(ctx)
	for _, ownObj := range r.cfg.Owns {
		ownObj := ownObj

		opts := []client.ListOption{}
		if len(r.matchLabels) > 0 {
			opts = append(opts, r.matchLabels)
		}

		ownObjList := ownObj.NewObjList()
		if err := r.List(ctx, ownObjList, opts...); err != nil {
			return err
		}
		for _, o := range ownObjList.GetObjects() {
			o := o
			for _, ref := range o.GetOwnerReferences() {
				log.Info("ownerref", "refs", fmt.Sprintf("%s/%s", ref.UID, cr.GetUID()))
				if ref.UID == cr.GetUID() {
					apiVersion, kind := ownObj.SchemaGroupVersionKind().ToAPIVersionAndKind()
					log.Info("gvk", "apiVersion", apiVersion, "kind", kind)
					r.existingResources[corev1.ObjectReference{APIVersion: apiVersion, Kind: kind, Name: o.GetName(), Namespace: o.GetNamespace()}] = o
				}
			}
		}
	}
	return nil
}

// APIDelete is used to delete the existing resources that are owned by this cr
// the implementation retrieves the existing resources and deletes them
func (r *Resources) APIDelete(ctx context.Context, cr client.Object) error {
	r.m.Lock()
	defer r.m.Unlock()

	// step 0: get existing resources
	if err := r.getExistingResources(ctx, cr); err != nil {
		return err
	}
	return r.apiDelete(ctx)
}

func (r *Resources) apiDelete(ctx context.Context) error {
	// delete in priority
	for ref, o := range r.existingResources {
		if ref.Kind == "Namespace" {
			continue
		}
		if err := r.delete(ctx, ref, o); err != nil {
			return err
		}
	}
	for ref, o := range r.existingResources {
		if err := r.delete(ctx, ref, o); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resources) delete(ctx context.Context, ref corev1.ObjectReference, o backend.GenericObject) error {
	log := log.FromContext(ctx)
	log.Info("api delete existing resource", "referernce", ref.String())
	if err := r.Delete(ctx, o); err != nil {
		if resource.IgnoreNotFound(err) != nil {
			log.Error("api delete", "error", err, "object", o)
			return err
		}
		delete(r.existingResources, ref)
	}
	return nil
}

// APIApply
// step 0: get existing resources
// step 1: remove the exisiting resources from the internal resource list that overlap with the new resources
// step 2: delete the exisiting resources that are no longer needed
// step 3: apply the new resources to the api server
func (r *Resources) APIApply(ctx context.Context, cr client.Object) error {
	r.m.Lock()
	defer r.m.Unlock()

	log := log.FromContext(ctx)
	// step 0: get existing resources
	if err := r.getExistingResources(ctx, cr); err != nil {
		return err
	}

	// step 1: remove the exisiting resources that overlap with the new resources
	// since apply will change the content.
	for ref := range r.newResources {
		delete(r.existingResources, ref)
	}

	log.Info("api apply existing resources to be deleted", "existing resources", r.getExistingRefs())
	// step2b delete the exisiting resource that are no longer needed
	if err := r.apiDelete(ctx); err != nil {
		return err
	}

	// step3b apply the new resources to the api server
	return r.apiApply(ctx)
}

func (r *Resources) apiApply(ctx context.Context) error {
	// apply in priority
	for ref, o := range r.newResources {
		ref := ref
		o := o
		if ref.Kind == "Namespace" { // apply in priority
			r.apply(ctx, ref, o)
		} else {
			continue
		}
	}
	for ref, o := range r.newResources {
		r.apply(ctx, ref, o)
	}
	return nil
}

func (r *Resources) apply(ctx context.Context, ref corev1.ObjectReference, o backend.GenericObject) error {
	log := log.FromContext(ctx)
	key := types.NamespacedName{Namespace: o.GetNamespace(), Name: o.GetName()}
	log.Info("api apply object", "ref", ref.String(), "key", key.String())
	spec := o.GetSpec()
	if err := r.Client.Get(ctx, key, o); err != nil {
		if resource.IgnoreNotFound(err) != nil {
			return err
		}
		if err := r.Create(ctx, o); err != nil {
			return err
		}
		return nil
	}
	o.SetSpec(spec)
	return r.Update(ctx, o)
}

func (r *Resources) GetNewResources() map[corev1.ObjectReference]client.Object {
	r.m.RLock()
	defer r.m.RUnlock()

	res := make(map[corev1.ObjectReference]client.Object, len(r.newResources))
	for ref, o := range r.newResources {
		ref := ref
		o := o
		res[ref] = o
	}
	return res
}

func (r *Resources) getExistingRefs() []corev1.ObjectReference {
	l := []corev1.ObjectReference{}
	for ref := range r.existingResources {
		ref := ref
		l = append(l, ref)
	}
	return l
}
