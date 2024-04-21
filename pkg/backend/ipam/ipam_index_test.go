package ipam

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPAMIndexNormal(t *testing.T) {
	tests := map[string]struct {
		niName string
	}{
		"CreateDelete": {
			niName: "a",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			be := New(nil)
			ctx := context.Background()
			index := getNI(tc.niName)
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
