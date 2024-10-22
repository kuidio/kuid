package astest

import (
	"context"
	"testing"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/kuidio/kuid/apis/backend/as"
	asbev1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
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
			asStorageProviders := as.NewStorageProviders(ctx, true, &options.Options{
				Type: options.StorageType_Memory,
			})
			apiserver := builder.APIServer.
				WithServerName("kuid-api-server").
				WithOpenAPIDefinitions("Config", "v1alpha1", openapi.GetOpenAPIDefinitions).
				WithResourceAndHandler(&as.ASIndex{}, asStorageProviders.GetIndexStorageProvider()).
				WithResourceAndHandler(&as.ASClaim{}, asStorageProviders.GetClaimStorageProvider()).
				WithResourceAndHandler(&as.ASEntry{}, asStorageProviders.GetEntryStorageProvider()).
				WithResourceAndHandler(&asbev1alpha1.ASIndex{}, asStorageProviders.GetIndexStorageProvider()).
				WithResourceAndHandler(&asbev1alpha1.ASClaim{}, asStorageProviders.GetClaimStorageProvider()).
				WithResourceAndHandler(&asbev1alpha1.ASEntry{}, asStorageProviders.GetEntryStorageProvider()).
				WithoutEtcd()

			if _, err := apiserver.Build(ctx); err != nil {
				panic(err)
			}
			if err := asStorageProviders.ApplyStorageToBackend(ctx, apiserver); err != nil {
				panic(err)
			}
			be := asStorageProviders.GetBackend()
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
