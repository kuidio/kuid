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

// ASEntryLister helps list ASEntries.
// All objects returned here must be treated as read-only.
type ASEntryLister interface {
	// List lists all ASEntries in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ASEntry, err error)
	// ASEntries returns an object that can list and get ASEntries.
	ASEntries(namespace string) ASEntryNamespaceLister
	ASEntryListerExpansion
}

// aSEntryLister implements the ASEntryLister interface.
type aSEntryLister struct {
	listers.ResourceIndexer[*v1alpha1.ASEntry]
}

// NewASEntryLister returns a new ASEntryLister.
func NewASEntryLister(indexer cache.Indexer) ASEntryLister {
	return &aSEntryLister{listers.New[*v1alpha1.ASEntry](indexer, v1alpha1.Resource("asentry"))}
}

// ASEntries returns an object that can list and get ASEntries.
func (s *aSEntryLister) ASEntries(namespace string) ASEntryNamespaceLister {
	return aSEntryNamespaceLister{listers.NewNamespaced[*v1alpha1.ASEntry](s.ResourceIndexer, namespace)}
}

// ASEntryNamespaceLister helps list and get ASEntries.
// All objects returned here must be treated as read-only.
type ASEntryNamespaceLister interface {
	// List lists all ASEntries in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ASEntry, err error)
	// Get retrieves the ASEntry from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ASEntry, error)
	ASEntryNamespaceListerExpansion
}

// aSEntryNamespaceLister implements the ASEntryNamespaceLister
// interface.
type aSEntryNamespaceLister struct {
	listers.ResourceIndexer[*v1alpha1.ASEntry]
}