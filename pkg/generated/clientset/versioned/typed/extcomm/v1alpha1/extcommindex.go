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
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"

	v1alpha1 "github.com/kuidio/kuid/apis/backend/extcomm/v1alpha1"
	scheme "github.com/kuidio/kuid/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// EXTCOMMIndexesGetter has a method to return a EXTCOMMIndexInterface.
// A group's client should implement this interface.
type EXTCOMMIndexesGetter interface {
	EXTCOMMIndexes(namespace string) EXTCOMMIndexInterface
}

// EXTCOMMIndexInterface has methods to work with EXTCOMMIndex resources.
type EXTCOMMIndexInterface interface {
	Create(ctx context.Context, eXTCOMMIndex *v1alpha1.EXTCOMMIndex, opts v1.CreateOptions) (*v1alpha1.EXTCOMMIndex, error)
	Update(ctx context.Context, eXTCOMMIndex *v1alpha1.EXTCOMMIndex, opts v1.UpdateOptions) (*v1alpha1.EXTCOMMIndex, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, eXTCOMMIndex *v1alpha1.EXTCOMMIndex, opts v1.UpdateOptions) (*v1alpha1.EXTCOMMIndex, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.EXTCOMMIndex, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.EXTCOMMIndexList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.EXTCOMMIndex, err error)
	EXTCOMMIndexExpansion
}

// eXTCOMMIndexes implements EXTCOMMIndexInterface
type eXTCOMMIndexes struct {
	*gentype.ClientWithList[*v1alpha1.EXTCOMMIndex, *v1alpha1.EXTCOMMIndexList]
}

// newEXTCOMMIndexes returns a EXTCOMMIndexes
func newEXTCOMMIndexes(c *ExtcommV1alpha1Client, namespace string) *eXTCOMMIndexes {
	return &eXTCOMMIndexes{
		gentype.NewClientWithList[*v1alpha1.EXTCOMMIndex, *v1alpha1.EXTCOMMIndexList](
			"extcommindexes",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *v1alpha1.EXTCOMMIndex { return &v1alpha1.EXTCOMMIndex{} },
			func() *v1alpha1.EXTCOMMIndexList { return &v1alpha1.EXTCOMMIndexList{} }),
	}
}
