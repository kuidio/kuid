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
	v1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// IPClaimLister helps list IPClaims.
// All objects returned here must be treated as read-only.
type IPClaimLister interface {
	// List lists all IPClaims in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.IPClaim, err error)
	// IPClaims returns an object that can list and get IPClaims.
	IPClaims(namespace string) IPClaimNamespaceLister
	IPClaimListerExpansion
}

// iPClaimLister implements the IPClaimLister interface.
type iPClaimLister struct {
	indexer cache.Indexer
}

// NewIPClaimLister returns a new IPClaimLister.
func NewIPClaimLister(indexer cache.Indexer) IPClaimLister {
	return &iPClaimLister{indexer: indexer}
}

// List lists all IPClaims in the indexer.
func (s *iPClaimLister) List(selector labels.Selector) (ret []*v1alpha1.IPClaim, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.IPClaim))
	})
	return ret, err
}

// IPClaims returns an object that can list and get IPClaims.
func (s *iPClaimLister) IPClaims(namespace string) IPClaimNamespaceLister {
	return iPClaimNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// IPClaimNamespaceLister helps list and get IPClaims.
// All objects returned here must be treated as read-only.
type IPClaimNamespaceLister interface {
	// List lists all IPClaims in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.IPClaim, err error)
	// Get retrieves the IPClaim from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.IPClaim, error)
	IPClaimNamespaceListerExpansion
}

// iPClaimNamespaceLister implements the IPClaimNamespaceLister
// interface.
type iPClaimNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all IPClaims in the indexer for a given namespace.
func (s iPClaimNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.IPClaim, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.IPClaim))
	})
	return ret, err
}

// Get retrieves the IPClaim from the indexer for a given namespace and name.
func (s iPClaimNamespaceLister) Get(name string) (*v1alpha1.IPClaim, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("ipclaim"), name)
	}
	return obj.(*v1alpha1.IPClaim), nil
}
