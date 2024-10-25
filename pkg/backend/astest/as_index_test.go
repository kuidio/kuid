package astest

import (
	"context"
	"testing"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/kuidio/kuid/apis/backend/as"
	"github.com/kuidio/kuid/apis/backend/as/register"
	asbev1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	"github.com/kuidio/kuid/pkg/config"
	"github.com/kuidio/kuid/pkg/generated/openapi"
	"github.com/kuidio/kuid/pkg/registry/options"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	tests := map[string]struct {
		index    string
		testType string
	}{
		//"CreateDelete": {
		//	index:    "a",
		//	testType: "",
		//},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			ctx := context.Background()
			apiserver := builder.APIServer.
				WithServerName("kuid-api-server").
				WithOpenAPIDefinitions("Config", "v1alpha1", openapi.GetOpenAPIDefinitions).
				WithoutEtcd()

			groupConfig := config.GroupConfig{
				BackendFn:               register.NewBackend,
				ApplyStorageToBackendFn: register.ApplyStorageToBackend,
				Resources: []*config.ResourceConfig{
					{StorageProviderFn: register.NewIndexStorageProvider, Internal: &as.ASIndex{},ResourceVersions: []resource.Object{&as.ASIndex{}, &asbev1alpha1.ASIndex{}}},
					{StorageProviderFn: register.NewClaimStorageProvider, Internal: &as.ASClaim{}, ResourceVersions: []resource.Object{&as.ASClaim{}, &asbev1alpha1.ASClaim{}}},
					{StorageProviderFn: register.NewStorageProvider, Internal: &as.ASEntry{}, ResourceVersions: []resource.Object{&as.ASEntry{}, &asbev1alpha1.ASEntry{}}},
				},
			}

			be := groupConfig.BackendFn()
			for _, resource := range groupConfig.Resources {
				storageProvider := resource.StorageProviderFn(ctx,resource.Internal, be, true, &options.Options{
					Type: options.StorageType_Memory,
				})
				for _, resourceVersion := range resource.ResourceVersions {
					apiserver.WithResourceAndHandler(resourceVersion, storageProvider)
				}
			}

			if _, err := apiserver.Build(ctx); err != nil {
				panic(err)
			}
			if err := groupConfig.ApplyStorageToBackendFn(ctx, be, apiserver); err != nil {
				panic(err)
			}

			index, err := getIndex(tc.index, tc.testType)
			assert.NoError(t, err)
			if err := be.CreateIndex(ctx, index); err != nil {
				assert.Error(t, err)
			}
			if err := be.DeleteIndex(ctx, index); err != nil {
				assert.Error(t, err)
			}
			if err := be.DeleteIndex(ctx, index); err != nil {
				assert.Error(t, err)
			}
		})
	}
}
