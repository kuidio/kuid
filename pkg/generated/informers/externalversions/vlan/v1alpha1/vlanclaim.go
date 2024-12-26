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
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	vlanv1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	versioned "github.com/kuidio/kuid/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/kuidio/kuid/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kuidio/kuid/pkg/generated/listers/vlan/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// VLANClaimInformer provides access to a shared informer and lister for
// VLANClaims.
type VLANClaimInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.VLANClaimLister
}

type vLANClaimInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewVLANClaimInformer constructs a new informer for VLANClaim type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewVLANClaimInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredVLANClaimInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredVLANClaimInformer constructs a new informer for VLANClaim type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredVLANClaimInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.VlanV1alpha1().VLANClaims(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.VlanV1alpha1().VLANClaims(namespace).Watch(context.TODO(), options)
			},
		},
		&vlanv1alpha1.VLANClaim{},
		resyncPeriod,
		indexers,
	)
}

func (f *vLANClaimInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredVLANClaimInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *vLANClaimInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&vlanv1alpha1.VLANClaim{}, f.defaultInformer)
}

func (f *vLANClaimInformer) Lister() v1alpha1.VLANClaimLister {
	return v1alpha1.NewVLANClaimLister(f.Informer().GetIndexer())
}
