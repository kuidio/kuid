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

// EXTCOMMEntryLister helps list EXTCOMMEntries.
// All objects returned here must be treated as read-only.
type EXTCOMMEntryLister interface {
	// List lists all EXTCOMMEntries in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.EXTCOMMEntry, err error)
	// EXTCOMMEntries returns an object that can list and get EXTCOMMEntries.
	EXTCOMMEntries(namespace string) EXTCOMMEntryNamespaceLister
	EXTCOMMEntryListerExpansion
}

// eXTCOMMEntryLister implements the EXTCOMMEntryLister interface.
type eXTCOMMEntryLister struct {
	listers.ResourceIndexer[*v1alpha1.EXTCOMMEntry]
}

// NewEXTCOMMEntryLister returns a new EXTCOMMEntryLister.
func NewEXTCOMMEntryLister(indexer cache.Indexer) EXTCOMMEntryLister {
	return &eXTCOMMEntryLister{listers.New[*v1alpha1.EXTCOMMEntry](indexer, v1alpha1.Resource("extcommentry"))}
}

// EXTCOMMEntries returns an object that can list and get EXTCOMMEntries.
func (s *eXTCOMMEntryLister) EXTCOMMEntries(namespace string) EXTCOMMEntryNamespaceLister {
	return eXTCOMMEntryNamespaceLister{listers.NewNamespaced[*v1alpha1.EXTCOMMEntry](s.ResourceIndexer, namespace)}
}

// EXTCOMMEntryNamespaceLister helps list and get EXTCOMMEntries.
// All objects returned here must be treated as read-only.
type EXTCOMMEntryNamespaceLister interface {
	// List lists all EXTCOMMEntries in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.EXTCOMMEntry, err error)
	// Get retrieves the EXTCOMMEntry from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.EXTCOMMEntry, error)
	EXTCOMMEntryNamespaceListerExpansion
}

// eXTCOMMEntryNamespaceLister implements the EXTCOMMEntryNamespaceLister
// interface.
type eXTCOMMEntryNamespaceLister struct {
	listers.ResourceIndexer[*v1alpha1.EXTCOMMEntry]
}