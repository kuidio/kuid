/*
Copyright 2023 The Nephio Authors.

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

package backend

import (
	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/idxtable/pkg/tree"
	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/store"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IndexObject interface {
	Object
	GetNamespacedName() types.NamespacedName
	GetKey() store.Key
	GetTree() gtree.GTree
	GetType() string
	GetMinID() *uint64
	GetMaxID() *uint64
	GetMinClaim() ClaimObject
	GetMaxClaim() ClaimObject
}

type ClaimObject interface {
	Object
	GetNamespacedName() types.NamespacedName
	GetKey() store.Key
	GetIndex() string                   // implement
	GetSelector() *metav1.LabelSelector // implement
	GetOwnerSelector() (labels.Selector, error)
	GetClaimType() ClaimType
	GetLabelSelector() (labels.Selector, error)
	GetClaimLabels() labels.Set
	ValidateOwner(labels labels.Set) error
	GetStaticID() *uint64                           // implement
	GetStaticTreeID(t string) tree.ID               // implement
	GetClaimID(t string, id uint64) tree.ID         // implement
	GetRange() *string                              // implement
	GetRangeID(t string) (tree.Range, error)        // implement
	GetTable(t string, to, from uint64) table.Table // implement
	SetStatusRange(*string)                         // implement
	SetStatusID(*uint64)                            // implement
	GetStatusID() *uint64                           // implement
	ValidateSyntax(s string) field.ErrorList
	
}

type EntryObject interface {
	Object
	GetNamespacedName() types.NamespacedName
	GetKey() store.Key
	GetClaimType() ClaimType
	GetOwnerGVK() schema.GroupVersionKind
	GetOwnerNSN() types.NamespacedName
	SetSpec(x any)
	GetSpec() any
}

type Object interface {
	client.Object
	GetOwnerReference() *commonv1alpha1.OwnerReference
}

type ObjectList interface {
	GetItems() []Object
	client.ObjectList
}
