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

	v1alpha1 "github.com/kuidio/kuid/apis/backend/infra/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRegions implements RegionInterface
type FakeRegions struct {
	Fake *FakeInfraV1alpha1
	ns   string
}

var regionsResource = v1alpha1.SchemeGroupVersion.WithResource("regions")

var regionsKind = v1alpha1.SchemeGroupVersion.WithKind("Region")

// Get takes name of the region, and returns the corresponding region object, and an error if there is any.
func (c *FakeRegions) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Region, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(regionsResource, c.ns, name), &v1alpha1.Region{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Region), err
}

// List takes label and field selectors, and returns the list of Regions that match those selectors.
func (c *FakeRegions) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.RegionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(regionsResource, regionsKind, c.ns, opts), &v1alpha1.RegionList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RegionList{ListMeta: obj.(*v1alpha1.RegionList).ListMeta}
	for _, item := range obj.(*v1alpha1.RegionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested regions.
func (c *FakeRegions) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(regionsResource, c.ns, opts))

}

// Create takes the representation of a region and creates it.  Returns the server's representation of the region, and an error, if there is any.
func (c *FakeRegions) Create(ctx context.Context, region *v1alpha1.Region, opts v1.CreateOptions) (result *v1alpha1.Region, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(regionsResource, c.ns, region), &v1alpha1.Region{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Region), err
}

// Update takes the representation of a region and updates it. Returns the server's representation of the region, and an error, if there is any.
func (c *FakeRegions) Update(ctx context.Context, region *v1alpha1.Region, opts v1.UpdateOptions) (result *v1alpha1.Region, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(regionsResource, c.ns, region), &v1alpha1.Region{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Region), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRegions) UpdateStatus(ctx context.Context, region *v1alpha1.Region, opts v1.UpdateOptions) (*v1alpha1.Region, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(regionsResource, "status", c.ns, region), &v1alpha1.Region{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Region), err
}

// Delete takes name of the region and deletes it. Returns an error if one occurs.
func (c *FakeRegions) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(regionsResource, c.ns, name, opts), &v1alpha1.Region{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRegions) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(regionsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.RegionList{})
	return err
}

// Patch applies the patch and returns the patched region.
func (c *FakeRegions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Region, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(regionsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Region{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Region), err
}