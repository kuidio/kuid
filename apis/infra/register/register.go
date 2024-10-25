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

package register

import (
	"context"

	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-builder/pkg/builder/rest"
	"github.com/kuidio/kuid/apis/infra"
	infrav1alpha1 "github.com/kuidio/kuid/apis/infra/v1alpha1"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	"github.com/kuidio/kuid/pkg/config"
	genericregistry "github.com/kuidio/kuid/pkg/registry/generic"
	"github.com/kuidio/kuid/pkg/registry/options"
)

func init() {
	config.Register(
		infra.SchemeGroupVersion.Group,
		infrav1alpha1.AddToScheme,
		nil,
		nil,
		[]*config.ResourceConfig{
			{StorageProviderFn: NewStorageProvider, Internal: &infra.Cluster{}, ResourceVersions: []resource.Object{&infra.Cluster{}, &infrav1alpha1.Cluster{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.Endpoint{}, ResourceVersions: []resource.Object{&infra.Endpoint{}, &infrav1alpha1.Endpoint{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.EndpointSet{}, ResourceVersions: []resource.Object{&infra.EndpointSet{}, &infrav1alpha1.EndpointSet{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.Link{}, ResourceVersions: []resource.Object{&infra.Link{}, &infrav1alpha1.Link{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.LinkSet{}, ResourceVersions: []resource.Object{&infra.LinkSet{}, &infrav1alpha1.LinkSet{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.Module{}, ResourceVersions: []resource.Object{&infra.Module{}, &infrav1alpha1.Module{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.ModuleBay{}, ResourceVersions: []resource.Object{&infra.ModuleBay{}, &infrav1alpha1.ModuleBay{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.Node{}, ResourceVersions: []resource.Object{&infra.Node{}, &infrav1alpha1.Node{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.NodeItem{}, ResourceVersions: []resource.Object{&infra.NodeItem{}, &infrav1alpha1.NodeItem{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.NodeSet{}, ResourceVersions: []resource.Object{&infra.NodeSet{}, &infrav1alpha1.NodeSet{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.Partition{}, ResourceVersions: []resource.Object{&infra.Partition{}, &infrav1alpha1.Partition{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.Rack{}, ResourceVersions: []resource.Object{&infra.Rack{}, &infrav1alpha1.Rack{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.Region{}, ResourceVersions: []resource.Object{&infra.Region{}, &infrav1alpha1.Region{}}},
			{StorageProviderFn: NewStorageProvider, Internal: &infra.Site{}, ResourceVersions: []resource.Object{&infra.Site{}, &infrav1alpha1.Site{}}},
		},
	)
}

func NewStorageProvider(ctx context.Context, obj resource.InternalObject, be bebackend.Backend, sync bool, options *options.Options) *rest.StorageProvider {
	return genericregistry.NewStorageProvider(ctx, obj, options)
}
