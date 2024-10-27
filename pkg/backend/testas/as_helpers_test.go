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

package testas

import (
	"context"
	"fmt"
	"reflect"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/apis/backend/as"
	"github.com/kuidio/kuid/apis/backend/as/register"
	asbev1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	"github.com/kuidio/kuid/apis/common"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	"github.com/kuidio/kuid/pkg/config"
	"github.com/kuidio/kuid/pkg/generated/openapi"
	"github.com/kuidio/kuid/pkg/registry/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/utils/ptr"
)

type testCtx struct {
	name          string
	claimType     backend.ClaimType
	id            uint64
	tRange        string
	labels        map[string]string
	selector      *metav1.LabelSelector
	expectedError bool
	expectedID    *uint64
	expectedRange *string
}

// alias
const (
	namespace    = "dummy"
	staticClaim  = backend.ClaimType_StaticID
	dynamicClaim = backend.ClaimType_DynamicID
	rangeClaim   = backend.ClaimType_Range
)

func apiServer() *builder.Server {
	return builder.NewAPIServer().
		WithServerName("kuid-api-server").
		WithOpenAPIDefinitions("Config", "v1alpha1", openapi.GetOpenAPIDefinitions).
		WithoutEtcd()
}

func initBackend(ctx context.Context, apiserver *builder.Server) (bebackend.Backend, error) {
	groupConfig := config.GroupConfig{
		BackendFn:               register.NewBackend,
		ApplyStorageToBackendFn: register.ApplyStorageToBackend,
		Resources: []*config.ResourceConfig{
			{StorageProviderFn: register.NewIndexStorageProvider, Internal: &as.ASIndex{}, ResourceVersions: []resource.Object{&as.ASIndex{}, &asbev1alpha1.ASIndex{}}},
			{StorageProviderFn: register.NewClaimStorageProvider, Internal: &as.ASClaim{}, ResourceVersions: []resource.Object{&as.ASClaim{}, &asbev1alpha1.ASClaim{}}},
			{StorageProviderFn: register.NewStorageProvider, Internal: &as.ASEntry{}, ResourceVersions: []resource.Object{&as.ASEntry{}, &asbev1alpha1.ASEntry{}}},
		},
	}

	be := groupConfig.BackendFn()
	for _, resource := range groupConfig.Resources {
		storageProvider := resource.StorageProviderFn(ctx, resource.Internal, be, true, &options.Options{
			Type: options.StorageType_Memory,
		})
		for _, resourceVersion := range resource.ResourceVersions {
			apiserver.WithResourceAndHandler(resourceVersion, storageProvider)
		}
	}

	if _, err := apiserver.Build(ctx); err != nil {
		return nil, err
	}
	if err := groupConfig.ApplyStorageToBackendFn(ctx, be, apiserver); err != nil {
		return nil, err
	}
	return be, nil
}

func getStorage(ctx context.Context, apiServer *builder.Server, gr schema.GroupResource) (*registry.Store, error) {
	storageProvider := apiServer.StorageProvider[gr]
	storage, err := storageProvider.Get(ctx, apiServer.Schemes[0], &Getter{})
	if err != nil {
		return nil, err
	}
	registryStore, ok := storage.(*registry.Store)
	if !ok {
		return nil, fmt.Errorf("index store is not a *registry.Store, got: %v", reflect.TypeOf(storage).Name())
	}
	return registryStore, nil
}

var _ generic.RESTOptionsGetter = &Getter{}

type Getter struct{}

func (r *Getter) GetRESTOptions(resource schema.GroupResource, example runtime.Object) (generic.RESTOptions, error) {
	return generic.RESTOptions{}, nil
}

func getIndex(index, _ string) (*as.ASIndex, error) {
	idx := as.BuildASIndex(
		metav1.ObjectMeta{Namespace: namespace, Name: index},
		nil,
		nil,
	)

	fieldErrs := idx.ValidateSyntax("")
	if len(fieldErrs) != 0 {
		return nil, fmt.Errorf("syntax errors %v", fieldErrs)
	}
	return idx, nil
}

func (r testCtx) getDynamicClaim(index, testType string) (backend.ClaimObject, error) {
	claim := as.BuildASClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&as.ASClaimSpec{
			Index: index,
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
				Selector:          r.selector,
			},
		},
		nil,
	)
	fielErrList := claim.ValidateSyntax(testType) // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return claim, nil
}

func (r testCtx) getStaticClaim(index, testType string) (backend.ClaimObject, error) {
	claim := as.BuildASClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&as.ASClaimSpec{
			Index: index,
			ID:    ptr.To[uint32](uint32(r.id)),
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := claim.ValidateSyntax(testType) // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return claim, nil
}

func (r testCtx) getRangeClaim(index, testType string) (backend.ClaimObject, error) {
	fmt.Println("getRangeClaim", r.name, r.tRange)
	claim := as.BuildASClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&as.ASClaimSpec{
			Index: index,
			Range: ptr.To[string](r.tRange),
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := claim.ValidateSyntax(testType) // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	fmt.Println("getRangeClaim", *claim.GetRange())
	return claim, nil
}

/*
func transformer(_ context.Context, newObj runtime.Object, oldObj runtime.Object) (runtime.Object, error) {
	// Type assertion to specific object types, assuming we are working with a type that has Spec and Status fields
	new, ok := newObj.(backend.ClaimObject)
	if !ok {
		return nil, fmt.Errorf("newObj is not of type ClaimObject %s", reflect.TypeOf(newObj).Name())
	}
	old, ok := oldObj.(backend.ClaimObject)
	if !ok {
		return nil, fmt.Errorf("oldObj is not of type ClaimObject %s", reflect.TypeOf(newObj).Name())
	}

	new.SetResourceVersion(old.GetResourceVersion())
	new.SetUID(old.GetUID())

	if new.GetRange() != nil {
		fmt.Println("transformer", "new", *new.GetRange())
	}

	if old.GetRange() != nil {
		fmt.Println("transformer", "old", *old.GetRange())
	}

	return new, nil
}
*/
