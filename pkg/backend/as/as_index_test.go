package as

import (
	"context"
	"testing"

	"github.com/kuidio/kuid/pkg/backend/backend"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestIndex(t *testing.T) {
	tests := map[string]struct {
		index    string
		testType string
	}{
		"CreateDelete": {
			index: "a",
			testType: "",
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
