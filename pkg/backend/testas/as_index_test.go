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
	"testing"

	"github.com/kuidio/kuid/apis/backend/as"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
)

func TestIndex(t *testing.T) {
	tests := map[string]struct {
		index    string
		testType string
	}{
		"CreateDelete": {
			index:    "a",
			testType: "",
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
				Group:    as.SchemeGroupVersion.Group,
				Resource: as.ASIndexPlural,
			})
			if err != nil {
				t.Errorf("cannot get index storage, err: %v", err)
			}

			index, err := getIndex(tc.index, tc.testType)
			assert.NoError(t, err)
			ctx = genericapirequest.WithNamespace(ctx, index.GetNamespace())
			_, err = indexStorage.Create(ctx, index, nil, &metav1.CreateOptions{
				FieldManager: "backend",
			})
			if err != nil {
				assert.Error(t, err)
			}
			_, _, err = indexStorage.Delete(ctx, index.GetName(), nil, &metav1.DeleteOptions{})
			if err != nil {
				assert.Error(t, err)
			}
			_, _, err = indexStorage.Delete(ctx, index.GetName(), nil, &metav1.DeleteOptions{})
			if err != nil {
				assert.Error(t, err)
			}
		})
	}
}
