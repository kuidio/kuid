package genid

import (
	"context"
	"testing"

	"github.com/kuidio/kuid/pkg/backend/backend"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genidbev1alpha1 "github.com/kuidio/kuid/apis/backend/genid/v1alpha1"
)

func TestIndex(t *testing.T) {
	tests := map[string]struct {
		index    string
		testType string
	}{
		"CreateDelete-16bit": {
			index: "a",
			testType: string(genidbev1alpha1.GENIDType_16bit),
		},
		"CreateDelete-32bit": {
			index: "a",
			testType: string(genidbev1alpha1.GENIDType_32bit),
		},
		"CreateDelete-48bit": {
			index: "a",
			testType: string(genidbev1alpha1.GENIDType_48bit),
		},
		"CreateDelete-64bit": {
			index: "a",
			testType: string(genidbev1alpha1.GENIDType_64bit),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			be := backend.New(nil, nil, nil, nil, nil, schema.GroupVersionKind{},schema.GroupVersionKind{})
			ctx := context.Background()
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
