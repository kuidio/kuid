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

package vxlanentry

import (
	"context"
	"fmt"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-store/pkg/rest"
	"github.com/henderiw/apiserver-store/pkg/storebackend"
	"github.com/henderiw/logger/log"
	vxlanbe1v1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
	"github.com/kuidio/kuid/pkg/backend"
	"github.com/kuidio/kuid/pkg/backend/vxlan"
	watchermanager "github.com/kuidio/kuid/pkg/kuidserver/watcher-manager"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewStrategy creates and returns a fischerStrategy instance
func NewStrategy(ctx context.Context, typer runtime.ObjectTyper, client client.Client, store storebackend.Storer[runtime.Object], be backend.Backend[*vxlan.CacheContext]) *strategy {
	watcherManager := watchermanager.New(32)

	go watcherManager.Start(ctx)

	return &strategy{
		ObjectTyper:    typer,
		NameGenerator:  names.SimpleNameGenerator,
		client:         client,
		store:          store,
		be:             be,
		gr:             vxlanbe1v1alpha1.Resource(vxlanbe1v1alpha1.VXLANEntryPlural),
		resource:       &vxlanbe1v1alpha1.VXLANEntry{},
		watcherManager: watcherManager,
	}
}

// MatchPackageRevision is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func Match(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object does not match
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	api, ok := obj.(*vxlanbe1v1alpha1.VXLANEntry)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a %s", vxlanbe1v1alpha1.VXLANEntryKind)
	}
	return labels.Set(api.ObjectMeta.Labels), SelectableFields(api), nil
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *vxlanbe1v1alpha1.VXLANEntry) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

var _ rest.RESTGetStrategy = &strategy{}
var _ rest.RESTListStrategy = &strategy{}
var _ rest.RESTCreateStrategy = &strategy{}
var _ rest.RESTUpdateStrategy = &strategy{}
var _ rest.RESTDeleteStrategy = &strategy{}
var _ rest.RESTWatchStrategy = &strategy{}

type strategy struct {
	runtime.ObjectTyper
	names.NameGenerator
	client         client.Client
	store          storebackend.Storer[runtime.Object]
	be             backend.Backend[*vxlan.CacheContext]
	gr             schema.GroupResource
	resource       resource.Object
	watcherManager watchermanager.WatcherManager
}

func (r *strategy) NamespaceScoped() bool {
	return true
}

func (r *strategy) Canonicalize(obj runtime.Object) {}

func (r *strategy) notifyWatcher(ctx context.Context, event watch.Event) {
	log := log.FromContext(ctx).With("eventType", event.Type)
	log.Info("notify watcherManager")

	r.watcherManager.WatchChan() <- event
}
