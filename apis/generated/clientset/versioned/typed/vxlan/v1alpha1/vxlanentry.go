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
	"time"

	v1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
	scheme "github.com/kuidio/kuid/apis/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VXLANEntriesGetter has a method to return a VXLANEntryInterface.
// A group's client should implement this interface.
type VXLANEntriesGetter interface {
	VXLANEntries(namespace string) VXLANEntryInterface
}

// VXLANEntryInterface has methods to work with VXLANEntry resources.
type VXLANEntryInterface interface {
	Create(ctx context.Context, vXLANEntry *v1alpha1.VXLANEntry, opts v1.CreateOptions) (*v1alpha1.VXLANEntry, error)
	Update(ctx context.Context, vXLANEntry *v1alpha1.VXLANEntry, opts v1.UpdateOptions) (*v1alpha1.VXLANEntry, error)
	UpdateStatus(ctx context.Context, vXLANEntry *v1alpha1.VXLANEntry, opts v1.UpdateOptions) (*v1alpha1.VXLANEntry, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.VXLANEntry, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.VXLANEntryList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.VXLANEntry, err error)
	VXLANEntryExpansion
}

// vXLANEntries implements VXLANEntryInterface
type vXLANEntries struct {
	client rest.Interface
	ns     string
}

// newVXLANEntries returns a VXLANEntries
func newVXLANEntries(c *VxlanV1alpha1Client, namespace string) *vXLANEntries {
	return &vXLANEntries{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the vXLANEntry, and returns the corresponding vXLANEntry object, and an error if there is any.
func (c *vXLANEntries) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.VXLANEntry, err error) {
	result = &v1alpha1.VXLANEntry{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("vxlanentries").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VXLANEntries that match those selectors.
func (c *vXLANEntries) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.VXLANEntryList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.VXLANEntryList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("vxlanentries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested vXLANEntries.
func (c *vXLANEntries) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("vxlanentries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a vXLANEntry and creates it.  Returns the server's representation of the vXLANEntry, and an error, if there is any.
func (c *vXLANEntries) Create(ctx context.Context, vXLANEntry *v1alpha1.VXLANEntry, opts v1.CreateOptions) (result *v1alpha1.VXLANEntry, err error) {
	result = &v1alpha1.VXLANEntry{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("vxlanentries").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(vXLANEntry).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a vXLANEntry and updates it. Returns the server's representation of the vXLANEntry, and an error, if there is any.
func (c *vXLANEntries) Update(ctx context.Context, vXLANEntry *v1alpha1.VXLANEntry, opts v1.UpdateOptions) (result *v1alpha1.VXLANEntry, err error) {
	result = &v1alpha1.VXLANEntry{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("vxlanentries").
		Name(vXLANEntry.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(vXLANEntry).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *vXLANEntries) UpdateStatus(ctx context.Context, vXLANEntry *v1alpha1.VXLANEntry, opts v1.UpdateOptions) (result *v1alpha1.VXLANEntry, err error) {
	result = &v1alpha1.VXLANEntry{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("vxlanentries").
		Name(vXLANEntry.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(vXLANEntry).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the vXLANEntry and deletes it. Returns an error if one occurs.
func (c *vXLANEntries) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("vxlanentries").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *vXLANEntries) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("vxlanentries").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched vXLANEntry.
func (c *vXLANEntries) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.VXLANEntry, err error) {
	result = &v1alpha1.VXLANEntry{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("vxlanentries").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
