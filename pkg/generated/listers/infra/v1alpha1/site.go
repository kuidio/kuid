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
	v1alpha1 "github.com/kuidio/kuid/apis/infra/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// SiteLister helps list Sites.
// All objects returned here must be treated as read-only.
type SiteLister interface {
	// List lists all Sites in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Site, err error)
	// Sites returns an object that can list and get Sites.
	Sites(namespace string) SiteNamespaceLister
	SiteListerExpansion
}

// siteLister implements the SiteLister interface.
type siteLister struct {
	listers.ResourceIndexer[*v1alpha1.Site]
}

// NewSiteLister returns a new SiteLister.
func NewSiteLister(indexer cache.Indexer) SiteLister {
	return &siteLister{listers.New[*v1alpha1.Site](indexer, v1alpha1.Resource("site"))}
}

// Sites returns an object that can list and get Sites.
func (s *siteLister) Sites(namespace string) SiteNamespaceLister {
	return siteNamespaceLister{listers.NewNamespaced[*v1alpha1.Site](s.ResourceIndexer, namespace)}
}

// SiteNamespaceLister helps list and get Sites.
// All objects returned here must be treated as read-only.
type SiteNamespaceLister interface {
	// List lists all Sites in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Site, err error)
	// Get retrieves the Site from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Site, error)
	SiteNamespaceListerExpansion
}

// siteNamespaceLister implements the SiteNamespaceLister
// interface.
type siteNamespaceLister struct {
	listers.ResourceIndexer[*v1alpha1.Site]
}