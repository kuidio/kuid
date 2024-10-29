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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// ASClaimLister helps list ASClaims.
// All objects returned here must be treated as read-only.
type ASClaimLister interface {
	// List lists all ASClaims in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ASClaim, err error)
	// ASClaims returns an object that can list and get ASClaims.
	ASClaims(namespace string) ASClaimNamespaceLister
	ASClaimListerExpansion
}

// aSClaimLister implements the ASClaimLister interface.
type aSClaimLister struct {
	listers.ResourceIndexer[*v1alpha1.ASClaim]
}

// NewASClaimLister returns a new ASClaimLister.
func NewASClaimLister(indexer cache.Indexer) ASClaimLister {
	return &aSClaimLister{listers.New[*v1alpha1.ASClaim](indexer, v1alpha1.Resource("asclaim"))}
}

// ASClaims returns an object that can list and get ASClaims.
func (s *aSClaimLister) ASClaims(namespace string) ASClaimNamespaceLister {
	return aSClaimNamespaceLister{listers.NewNamespaced[*v1alpha1.ASClaim](s.ResourceIndexer, namespace)}
}

// ASClaimNamespaceLister helps list and get ASClaims.
// All objects returned here must be treated as read-only.
type ASClaimNamespaceLister interface {
	// List lists all ASClaims in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ASClaim, err error)
	// Get retrieves the ASClaim from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ASClaim, error)
	ASClaimNamespaceListerExpansion
}

// aSClaimNamespaceLister implements the ASClaimNamespaceLister
// interface.
type aSClaimNamespaceLister struct {
	listers.ResourceIndexer[*v1alpha1.ASClaim]
}