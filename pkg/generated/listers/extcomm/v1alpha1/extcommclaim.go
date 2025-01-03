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
	v1alpha1 "github.com/kuidio/kuid/apis/backend/extcomm/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// EXTCOMMClaimLister helps list EXTCOMMClaims.
// All objects returned here must be treated as read-only.
type EXTCOMMClaimLister interface {
	// List lists all EXTCOMMClaims in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.EXTCOMMClaim, err error)
	// EXTCOMMClaims returns an object that can list and get EXTCOMMClaims.
	EXTCOMMClaims(namespace string) EXTCOMMClaimNamespaceLister
	EXTCOMMClaimListerExpansion
}

// eXTCOMMClaimLister implements the EXTCOMMClaimLister interface.
type eXTCOMMClaimLister struct {
	listers.ResourceIndexer[*v1alpha1.EXTCOMMClaim]
}

// NewEXTCOMMClaimLister returns a new EXTCOMMClaimLister.
func NewEXTCOMMClaimLister(indexer cache.Indexer) EXTCOMMClaimLister {
	return &eXTCOMMClaimLister{listers.New[*v1alpha1.EXTCOMMClaim](indexer, v1alpha1.Resource("extcommclaim"))}
}

// EXTCOMMClaims returns an object that can list and get EXTCOMMClaims.
func (s *eXTCOMMClaimLister) EXTCOMMClaims(namespace string) EXTCOMMClaimNamespaceLister {
	return eXTCOMMClaimNamespaceLister{listers.NewNamespaced[*v1alpha1.EXTCOMMClaim](s.ResourceIndexer, namespace)}
}

// EXTCOMMClaimNamespaceLister helps list and get EXTCOMMClaims.
// All objects returned here must be treated as read-only.
type EXTCOMMClaimNamespaceLister interface {
	// List lists all EXTCOMMClaims in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.EXTCOMMClaim, err error)
	// Get retrieves the EXTCOMMClaim from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.EXTCOMMClaim, error)
	EXTCOMMClaimNamespaceListerExpansion
}

// eXTCOMMClaimNamespaceLister implements the EXTCOMMClaimNamespaceLister
// interface.
type eXTCOMMClaimNamespaceLister struct {
	listers.ResourceIndexer[*v1alpha1.EXTCOMMClaim]
}
