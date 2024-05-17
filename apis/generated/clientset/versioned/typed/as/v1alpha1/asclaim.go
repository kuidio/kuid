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

	v1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	scheme "github.com/kuidio/kuid/apis/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ASClaimsGetter has a method to return a ASClaimInterface.
// A group's client should implement this interface.
type ASClaimsGetter interface {
	ASClaims(namespace string) ASClaimInterface
}

// ASClaimInterface has methods to work with ASClaim resources.
type ASClaimInterface interface {
	Create(ctx context.Context, aSClaim *v1alpha1.ASClaim, opts v1.CreateOptions) (*v1alpha1.ASClaim, error)
	Update(ctx context.Context, aSClaim *v1alpha1.ASClaim, opts v1.UpdateOptions) (*v1alpha1.ASClaim, error)
	UpdateStatus(ctx context.Context, aSClaim *v1alpha1.ASClaim, opts v1.UpdateOptions) (*v1alpha1.ASClaim, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ASClaim, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ASClaimList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ASClaim, err error)
	ASClaimExpansion
}

// aSClaims implements ASClaimInterface
type aSClaims struct {
	client rest.Interface
	ns     string
}

// newASClaims returns a ASClaims
func newASClaims(c *AsV1alpha1Client, namespace string) *aSClaims {
	return &aSClaims{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the aSClaim, and returns the corresponding aSClaim object, and an error if there is any.
func (c *aSClaims) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ASClaim, err error) {
	result = &v1alpha1.ASClaim{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("asclaims").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ASClaims that match those selectors.
func (c *aSClaims) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ASClaimList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ASClaimList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("asclaims").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested aSClaims.
func (c *aSClaims) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("asclaims").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a aSClaim and creates it.  Returns the server's representation of the aSClaim, and an error, if there is any.
func (c *aSClaims) Create(ctx context.Context, aSClaim *v1alpha1.ASClaim, opts v1.CreateOptions) (result *v1alpha1.ASClaim, err error) {
	result = &v1alpha1.ASClaim{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("asclaims").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aSClaim).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a aSClaim and updates it. Returns the server's representation of the aSClaim, and an error, if there is any.
func (c *aSClaims) Update(ctx context.Context, aSClaim *v1alpha1.ASClaim, opts v1.UpdateOptions) (result *v1alpha1.ASClaim, err error) {
	result = &v1alpha1.ASClaim{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("asclaims").
		Name(aSClaim.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aSClaim).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *aSClaims) UpdateStatus(ctx context.Context, aSClaim *v1alpha1.ASClaim, opts v1.UpdateOptions) (result *v1alpha1.ASClaim, err error) {
	result = &v1alpha1.ASClaim{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("asclaims").
		Name(aSClaim.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(aSClaim).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the aSClaim and deletes it. Returns an error if one occurs.
func (c *aSClaims) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("asclaims").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *aSClaims) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("asclaims").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched aSClaim.
func (c *aSClaims) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ASClaim, err error) {
	result = &v1alpha1.ASClaim{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("asclaims").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}