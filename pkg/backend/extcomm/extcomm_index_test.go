package extcomm

import (
	"context"
	"testing"

	"github.com/kuidio/kuid/pkg/backend/backend"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	extcommbev1alpha1 "github.com/kuidio/kuid/apis/backend/extcomm/v1alpha1"
)

func TestIndex(t *testing.T) {
	tests := map[string]struct {
		index    string
		testType string
	}{
		"CreateDelete-2ByteAS": {
			index: "a",
			testType: string(extcommbev1alpha1.ExtendedCommunityType_2byteAS),
		},
		"CreateDelete-4ByteAS": {
			index: "a",
			testType: string(extcommbev1alpha1.ExtendedCommunityType_4byteAS),
		},
		"CreateDelete-IPv4Address": {
			index: "a",
			testType: string(extcommbev1alpha1.ExtendedCommunityType_IPv4Address),
		},
		"CreateDelete-Opaque": {
			index: "a",
			testType: string(extcommbev1alpha1.ExtendedCommunityType_Opaque),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			be := backend.New(nil, nil, nil, nil, nil, schema.GroupVersionKind{}, schema.GroupVersionKind{})
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
