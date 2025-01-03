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

	v1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeIPClaims implements IPClaimInterface
type FakeIPClaims struct {
	Fake *FakeIpamV1alpha1
	ns   string
}

var ipclaimsResource = v1alpha1.SchemeGroupVersion.WithResource("ipclaims")

var ipclaimsKind = v1alpha1.SchemeGroupVersion.WithKind("IPClaim")

// Get takes name of the iPClaim, and returns the corresponding iPClaim object, and an error if there is any.
func (c *FakeIPClaims) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.IPClaim, err error) {
	emptyResult := &v1alpha1.IPClaim{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(ipclaimsResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.IPClaim), err
}

// List takes label and field selectors, and returns the list of IPClaims that match those selectors.
func (c *FakeIPClaims) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.IPClaimList, err error) {
	emptyResult := &v1alpha1.IPClaimList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(ipclaimsResource, ipclaimsKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.IPClaimList{ListMeta: obj.(*v1alpha1.IPClaimList).ListMeta}
	for _, item := range obj.(*v1alpha1.IPClaimList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested iPClaims.
func (c *FakeIPClaims) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(ipclaimsResource, c.ns, opts))

}

// Create takes the representation of a iPClaim and creates it.  Returns the server's representation of the iPClaim, and an error, if there is any.
func (c *FakeIPClaims) Create(ctx context.Context, iPClaim *v1alpha1.IPClaim, opts v1.CreateOptions) (result *v1alpha1.IPClaim, err error) {
	emptyResult := &v1alpha1.IPClaim{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(ipclaimsResource, c.ns, iPClaim, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.IPClaim), err
}

// Update takes the representation of a iPClaim and updates it. Returns the server's representation of the iPClaim, and an error, if there is any.
func (c *FakeIPClaims) Update(ctx context.Context, iPClaim *v1alpha1.IPClaim, opts v1.UpdateOptions) (result *v1alpha1.IPClaim, err error) {
	emptyResult := &v1alpha1.IPClaim{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(ipclaimsResource, c.ns, iPClaim, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.IPClaim), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeIPClaims) UpdateStatus(ctx context.Context, iPClaim *v1alpha1.IPClaim, opts v1.UpdateOptions) (result *v1alpha1.IPClaim, err error) {
	emptyResult := &v1alpha1.IPClaim{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(ipclaimsResource, "status", c.ns, iPClaim, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.IPClaim), err
}

// Delete takes name of the iPClaim and deletes it. Returns an error if one occurs.
func (c *FakeIPClaims) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(ipclaimsResource, c.ns, name, opts), &v1alpha1.IPClaim{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeIPClaims) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(ipclaimsResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.IPClaimList{})
	return err
}

// Patch applies the patch and returns the patched iPClaim.
func (c *FakeIPClaims) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.IPClaim, err error) {
	emptyResult := &v1alpha1.IPClaim{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(ipclaimsResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.IPClaim), err
}
