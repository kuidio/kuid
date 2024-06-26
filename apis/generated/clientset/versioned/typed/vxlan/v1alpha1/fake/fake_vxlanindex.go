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

package fake

import (
	"context"

	v1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVXLANIndexes implements VXLANIndexInterface
type FakeVXLANIndexes struct {
	Fake *FakeVxlanV1alpha1
	ns   string
}

var vxlanindexesResource = v1alpha1.SchemeGroupVersion.WithResource("vxlanindexes")

var vxlanindexesKind = v1alpha1.SchemeGroupVersion.WithKind("VXLANIndex")

// Get takes name of the vXLANIndex, and returns the corresponding vXLANIndex object, and an error if there is any.
func (c *FakeVXLANIndexes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.VXLANIndex, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(vxlanindexesResource, c.ns, name), &v1alpha1.VXLANIndex{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VXLANIndex), err
}

// List takes label and field selectors, and returns the list of VXLANIndexes that match those selectors.
func (c *FakeVXLANIndexes) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.VXLANIndexList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(vxlanindexesResource, vxlanindexesKind, c.ns, opts), &v1alpha1.VXLANIndexList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.VXLANIndexList{ListMeta: obj.(*v1alpha1.VXLANIndexList).ListMeta}
	for _, item := range obj.(*v1alpha1.VXLANIndexList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested vXLANIndexes.
func (c *FakeVXLANIndexes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(vxlanindexesResource, c.ns, opts))

}

// Create takes the representation of a vXLANIndex and creates it.  Returns the server's representation of the vXLANIndex, and an error, if there is any.
func (c *FakeVXLANIndexes) Create(ctx context.Context, vXLANIndex *v1alpha1.VXLANIndex, opts v1.CreateOptions) (result *v1alpha1.VXLANIndex, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(vxlanindexesResource, c.ns, vXLANIndex), &v1alpha1.VXLANIndex{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VXLANIndex), err
}

// Update takes the representation of a vXLANIndex and updates it. Returns the server's representation of the vXLANIndex, and an error, if there is any.
func (c *FakeVXLANIndexes) Update(ctx context.Context, vXLANIndex *v1alpha1.VXLANIndex, opts v1.UpdateOptions) (result *v1alpha1.VXLANIndex, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(vxlanindexesResource, c.ns, vXLANIndex), &v1alpha1.VXLANIndex{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VXLANIndex), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVXLANIndexes) UpdateStatus(ctx context.Context, vXLANIndex *v1alpha1.VXLANIndex, opts v1.UpdateOptions) (*v1alpha1.VXLANIndex, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(vxlanindexesResource, "status", c.ns, vXLANIndex), &v1alpha1.VXLANIndex{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VXLANIndex), err
}

// Delete takes name of the vXLANIndex and deletes it. Returns an error if one occurs.
func (c *FakeVXLANIndexes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(vxlanindexesResource, c.ns, name, opts), &v1alpha1.VXLANIndex{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVXLANIndexes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(vxlanindexesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.VXLANIndexList{})
	return err
}

// Patch applies the patch and returns the patched vXLANIndex.
func (c *FakeVXLANIndexes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.VXLANIndex, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(vxlanindexesResource, c.ns, name, pt, data, subresources...), &v1alpha1.VXLANIndex{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VXLANIndex), err
}
