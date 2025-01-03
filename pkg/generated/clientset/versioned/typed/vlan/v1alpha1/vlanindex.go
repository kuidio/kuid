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

	v1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	scheme "github.com/kuidio/kuid/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// VLANIndexesGetter has a method to return a VLANIndexInterface.
// A group's client should implement this interface.
type VLANIndexesGetter interface {
	VLANIndexes(namespace string) VLANIndexInterface
}

// VLANIndexInterface has methods to work with VLANIndex resources.
type VLANIndexInterface interface {
	Create(ctx context.Context, vLANIndex *v1alpha1.VLANIndex, opts v1.CreateOptions) (*v1alpha1.VLANIndex, error)
	Update(ctx context.Context, vLANIndex *v1alpha1.VLANIndex, opts v1.UpdateOptions) (*v1alpha1.VLANIndex, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, vLANIndex *v1alpha1.VLANIndex, opts v1.UpdateOptions) (*v1alpha1.VLANIndex, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.VLANIndex, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.VLANIndexList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.VLANIndex, err error)
	VLANIndexExpansion
}

// vLANIndexes implements VLANIndexInterface
type vLANIndexes struct {
	*gentype.ClientWithList[*v1alpha1.VLANIndex, *v1alpha1.VLANIndexList]
}

// newVLANIndexes returns a VLANIndexes
func newVLANIndexes(c *VlanV1alpha1Client, namespace string) *vLANIndexes {
	return &vLANIndexes{
		gentype.NewClientWithList[*v1alpha1.VLANIndex, *v1alpha1.VLANIndexList](
			"vlanindexes",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *v1alpha1.VLANIndex { return &v1alpha1.VLANIndex{} },
			func() *v1alpha1.VLANIndexList { return &v1alpha1.VLANIndexList{} }),
	}
}
