package ipam

import (
	"context"
	"testing"

	"github.com/kuidio/kuid/apis/backend/ipam"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
)

func TestIPAMIndexNormal(t *testing.T) {
	tests := map[string]struct {
		index         string
		indexPrefixes []ipam.Prefix
	}{
		"CreateDelete": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "172.0.0.0/8"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			apiserver := apiServer()
			_, err := initBackend(ctx, apiserver)
			if err != nil {
				t.Errorf("cannot get backend, err: %v", err)
			}

			indexStorage, err := getStorage(ctx, apiserver, schema.GroupResource{
				Group:    ipam.SchemeGroupVersion.Group,
				Resource: ipam.IPIndexPlural,
			})
			if err != nil {
				t.Errorf("cannot get index storage, err: %v", err)
				return
			}

			index := getIndex(tc.index, tc.indexPrefixes)
			ctx = genericapirequest.WithNamespace(ctx, index.GetNamespace())
			_, err = indexStorage.Create(ctx, index, nil, &metav1.CreateOptions{
				FieldManager: "backend",
			})
			if err != nil {
				t.Errorf("cannot create index, err: %v", err)
				return
			}
			_, _, err = indexStorage.Delete(ctx, index.GetName(), nil, &metav1.DeleteOptions{})
			if err != nil {
				t.Errorf("cannot delete index, err: %v", err)
				return
			}
			_, _, err = indexStorage.Delete(ctx, index.GetName(), nil, &metav1.DeleteOptions{})
			if err == nil {
				t.Errorf("cannot delete non existing index, err: %v", err)
				return
			}
		})
	}
}
