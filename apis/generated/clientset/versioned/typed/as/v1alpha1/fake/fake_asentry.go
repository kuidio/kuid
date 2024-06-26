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

	v1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeASEntries implements ASEntryInterface
type FakeASEntries struct {
	Fake *FakeAsV1alpha1
	ns   string
}

var asentriesResource = v1alpha1.SchemeGroupVersion.WithResource("asentries")

var asentriesKind = v1alpha1.SchemeGroupVersion.WithKind("ASEntry")

// Get takes name of the aSEntry, and returns the corresponding aSEntry object, and an error if there is any.
func (c *FakeASEntries) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ASEntry, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(asentriesResource, c.ns, name), &v1alpha1.ASEntry{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ASEntry), err
}

// List takes label and field selectors, and returns the list of ASEntries that match those selectors.
func (c *FakeASEntries) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ASEntryList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(asentriesResource, asentriesKind, c.ns, opts), &v1alpha1.ASEntryList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ASEntryList{ListMeta: obj.(*v1alpha1.ASEntryList).ListMeta}
	for _, item := range obj.(*v1alpha1.ASEntryList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested aSEntries.
func (c *FakeASEntries) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(asentriesResource, c.ns, opts))

}

// Create takes the representation of a aSEntry and creates it.  Returns the server's representation of the aSEntry, and an error, if there is any.
func (c *FakeASEntries) Create(ctx context.Context, aSEntry *v1alpha1.ASEntry, opts v1.CreateOptions) (result *v1alpha1.ASEntry, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(asentriesResource, c.ns, aSEntry), &v1alpha1.ASEntry{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ASEntry), err
}

// Update takes the representation of a aSEntry and updates it. Returns the server's representation of the aSEntry, and an error, if there is any.
func (c *FakeASEntries) Update(ctx context.Context, aSEntry *v1alpha1.ASEntry, opts v1.UpdateOptions) (result *v1alpha1.ASEntry, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(asentriesResource, c.ns, aSEntry), &v1alpha1.ASEntry{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ASEntry), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeASEntries) UpdateStatus(ctx context.Context, aSEntry *v1alpha1.ASEntry, opts v1.UpdateOptions) (*v1alpha1.ASEntry, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(asentriesResource, "status", c.ns, aSEntry), &v1alpha1.ASEntry{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ASEntry), err
}

// Delete takes name of the aSEntry and deletes it. Returns an error if one occurs.
func (c *FakeASEntries) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(asentriesResource, c.ns, name, opts), &v1alpha1.ASEntry{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeASEntries) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(asentriesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ASEntryList{})
	return err
}

// Patch applies the patch and returns the patched aSEntry.
func (c *FakeASEntries) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ASEntry, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(asentriesResource, c.ns, name, pt, data, subresources...), &v1alpha1.ASEntry{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ASEntry), err
}
