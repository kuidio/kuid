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
	"testing"

	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/apis/backend/as"
	genericbe "github.com/kuidio/kuid/pkg/backend/generic"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/utils/ptr"
)

func Test(t *testing.T) {
	tests := map[string]struct {
		index string
		ctxs  []testCtx
	}{
		"Mix": {
			index: "a",
			ctxs: []testCtx{
				{claimType: dynamicClaim, name: "claim1", expectedError: false, expectedID: ptr.To[uint64](0)},
				{claimType: staticClaim, name: "claim2", id: 100, expectedError: false},
				{claimType: staticClaim, name: "claim3", id: 4000, expectedError: false},
				{claimType: rangeClaim, name: "claim4", tRange: "10-19", expectedError: false},
				{claimType: rangeClaim, name: "claim4", tRange: "11-19", expectedError: false}, // reclaim
				{claimType: rangeClaim, name: "claim5", tRange: "5-10", expectedError: false},  // claim a new range
				{claimType: rangeClaim, name: "claim6", tRange: "19-100", expectedError: true}, // overlap with a static id of claim2
				{claimType: staticClaim, name: "claim7", id: 12, expectedError: false},
				{claimType: staticClaim, name: "claim7", id: 13, expectedError: false}, // reclaim an existing id
				{claimType: dynamicClaim, name: "claim8", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedID: ptr.To[uint64](11)}, // a dynamic claim from a range
				{claimType: dynamicClaim, name: "claim9", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedID: ptr.To[uint64](12)}, // a dynamic claim from a range that was claimed before
				{claimType: dynamicClaim, name: "claim10", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"},
				}, expectedError: false, expectedID: ptr.To[uint64](14)},
				{claimType: dynamicClaim, name: "claim10", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "claim4"}, // update
				}, expectedError: false, expectedID: ptr.To[uint64](14)},
				{claimType: rangeClaim, name: "claim4", tRange: "11-19", expectedError: false}, // update
			},
		},
	}

	for name, tc := range tests {
		tc := tc

		testTypes := []string{""}

		for _, testType := range testTypes {
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

				claimStorage, err := getStorage(ctx, apiserver, schema.GroupResource{
					Group:    as.SchemeGroupVersion.Group,
					Resource: as.ASClaimPlural,
				})
				if err != nil {
					t.Errorf("cannot get claim storage, err: %v", err)
				}

				index, err := getIndex(tc.index, testType)
				assert.NoError(t, err)
				ctx = genericapirequest.WithNamespace(ctx, index.GetNamespace())
				_, err = indexStorage.Create(ctx, index, nil, &metav1.CreateOptions{
					FieldManager: "backend",
				})
				assert.NoError(t, err)

				for _, v := range tc.ctxs {
					v := v
					var claim backend.ClaimObject
					var err error

					switch v.claimType {
					case staticClaim:
						claim, err = v.getStaticClaim(tc.index, testType)
					case dynamicClaim:
						claim, err = v.getDynamicClaim(tc.index, testType)
					case rangeClaim:
						claim, err = v.getRangeClaim(tc.index, testType)
						fmt.Println("claim range", *claim.GetRange())
					}
					assert.NoError(t, err)
					if err != nil {
						return
					}

					exists := true
					if _, err := claimStorage.Get(ctx, claim.GetName(), &metav1.GetOptions{}); err != nil {
						exists = false
					}
					var newClaim runtime.Object
					if !exists {
						newClaim, err = claimStorage.Create(ctx, claim, nil, &metav1.CreateOptions{FieldManager: "test"})
					} else {

						defaultObjInfo := rest.DefaultUpdatedObjectInfo(claim, genericbe.ClaimTransformer)
						newClaim, _, err = claimStorage.Update(ctx, claim.GetName(), defaultObjInfo, nil, nil, false, &metav1.UpdateOptions{
							FieldManager: "backend",
						})
					}

					if v.expectedError {
						assert.Error(t, err)
						continue
					}
					if err != nil {
						assert.NoError(t, err)
						continue
					}

					asClaim, ok := newClaim.(*as.ASClaim)
					if !ok {
						t.Errorf("expecting ipClaim, got: %v", reflect.TypeOf(newClaim).Name())
						continue
					}

					//uobj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(newClaim)
					//if err != nil {
					//	assert.Error(t, err)
					//	continue
					//}
					//u := &unstructured.Unstructured{
					//	Object: uobj,
					//}
					//status := uobj["status"]
					//statusObj, ok := status.(map[string]any)
					//if !ok {
					//	t.Errorf("expecting status id got nil")
					//}

					switch v.claimType {
					case staticClaim, dynamicClaim:
						expectedID := v.id
						if v.expectedID != nil {
							expectedID = *v.expectedID
						}
						fmt.Printf("expected/received %v/%v\n", expectedID, asClaim.Status.ID)
					/*
						id, ok := statusObj["id"]
						if !ok {
							t.Errorf("expecting status id got nil")
							continue
						} else {
							expectedID := v.id
							if v.expectedID != nil {
								expectedID = *v.expectedID
							}
							if uint64(id) != expectedID {
								t.Errorf("expecting id got %d, want %d\n", *claim.Status.ID, expectedID)
							}
						}
					*/
					case rangeClaim:
						expectedRange := v.tRange
						if v.expectedRange != nil {
							expectedRange = *v.expectedRange
						}
						fmt.Printf("expected/received %v/%v\n", expectedRange, asClaim.Status.Range)
						/*
							if claim.Status.Range == nil {
								t.Errorf("expecting status id got nil")
							} else {
								expectedRange := v.tRange
								if v.expectedRange != nil {
									expectedRange = *v.expectedRange
								}
								if *claim.Status.Range != expectedRange {
									t.Errorf("expecting range got %s, want %s\n", *claim.Status.Range, expectedRange)
								}
							}
						*/
					}
					//fmt.Println("entries after claim", v.name)
					//key := store.KeyFromNSN(types.NamespacedName{Namespace: namespace, Name: tc.index})
					//be.PrintEntries(ctx, tc.index)
					//assert.NoError(t, err)
				}
			})
		}
	}
}
