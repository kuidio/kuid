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

package config

import (
	"context"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-builder/pkg/builder/rest"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	"github.com/kuidio/kuid/pkg/registry/options"
	"k8s.io/apimachinery/pkg/runtime"
)

var Groups = map[string]*GroupConfig{}

type BackendFn func() bebackend.Backend

type StorageProviderFn func(ctx context.Context, obj resource.InternalObject, be bebackend.Backend, sync bool, options *options.Options) *rest.StorageProvider

type ApplyStorageToBackendFn func(ctx context.Context, be bebackend.Backend, apiServer *builder.Server) error

type GroupConfig struct {
	AddToScheme             func(s *runtime.Scheme) error
	BackendFn               BackendFn
	ApplyStorageToBackendFn ApplyStorageToBackendFn
	Resources               []*ResourceConfig
}

type ResourceConfig struct {
	StorageProviderFn StorageProviderFn
	Internal          resource.InternalObject
	ResourceVersions  []resource.Object
}

func Register(groupName string, addToScheme func(s *runtime.Scheme) error, befn BackendFn, applybefn ApplyStorageToBackendFn, resources []*ResourceConfig) {
	Groups[groupName] = &GroupConfig{
		AddToScheme:             addToScheme,
		BackendFn:               befn,
		ApplyStorageToBackendFn: applybefn,
		Resources:               resources,
	}
}
