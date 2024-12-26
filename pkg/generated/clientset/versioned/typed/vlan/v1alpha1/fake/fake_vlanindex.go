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

	v1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVLANIndexes implements VLANIndexInterface
type FakeVLANIndexes struct {
	Fake *FakeVlanV1alpha1
	ns   string
}

var vlanindexesResource = v1alpha1.SchemeGroupVersion.WithResource("vlanindexes")

var vlanindexesKind = v1alpha1.SchemeGroupVersion.WithKind("VLANIndex")

// Get takes name of the vLANIndex, and returns the corresponding vLANIndex object, and an error if there is any.
func (c *FakeVLANIndexes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.VLANIndex, err error) {
	emptyResult := &v1alpha1.VLANIndex{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(vlanindexesResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.VLANIndex), err
}

// List takes label and field selectors, and returns the list of VLANIndexes that match those selectors.
func (c *FakeVLANIndexes) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.VLANIndexList, err error) {
	emptyResult := &v1alpha1.VLANIndexList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(vlanindexesResource, vlanindexesKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.VLANIndexList{ListMeta: obj.(*v1alpha1.VLANIndexList).ListMeta}
	for _, item := range obj.(*v1alpha1.VLANIndexList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested vLANIndexes.
func (c *FakeVLANIndexes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(vlanindexesResource, c.ns, opts))

}

// Create takes the representation of a vLANIndex and creates it.  Returns the server's representation of the vLANIndex, and an error, if there is any.
func (c *FakeVLANIndexes) Create(ctx context.Context, vLANIndex *v1alpha1.VLANIndex, opts v1.CreateOptions) (result *v1alpha1.VLANIndex, err error) {
	emptyResult := &v1alpha1.VLANIndex{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(vlanindexesResource, c.ns, vLANIndex, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.VLANIndex), err
}

// Update takes the representation of a vLANIndex and updates it. Returns the server's representation of the vLANIndex, and an error, if there is any.
func (c *FakeVLANIndexes) Update(ctx context.Context, vLANIndex *v1alpha1.VLANIndex, opts v1.UpdateOptions) (result *v1alpha1.VLANIndex, err error) {
	emptyResult := &v1alpha1.VLANIndex{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(vlanindexesResource, c.ns, vLANIndex, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.VLANIndex), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVLANIndexes) UpdateStatus(ctx context.Context, vLANIndex *v1alpha1.VLANIndex, opts v1.UpdateOptions) (result *v1alpha1.VLANIndex, err error) {
	emptyResult := &v1alpha1.VLANIndex{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(vlanindexesResource, "status", c.ns, vLANIndex, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.VLANIndex), err
}

// Delete takes name of the vLANIndex and deletes it. Returns an error if one occurs.
func (c *FakeVLANIndexes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(vlanindexesResource, c.ns, name, opts), &v1alpha1.VLANIndex{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVLANIndexes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(vlanindexesResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.VLANIndexList{})
	return err
}

// Patch applies the patch and returns the patched vLANIndex.
func (c *FakeVLANIndexes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.VLANIndex, err error) {
	emptyResult := &v1alpha1.VLANIndex{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(vlanindexesResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.VLANIndex), err
}
