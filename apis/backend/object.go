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
	"context"
	"crypto/sha1"

	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/idxtable/pkg/tree"
	"github.com/henderiw/idxtable/pkg/tree/gtree"
	"github.com/henderiw/store"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IndexObject interface {
	Object
	GetTree() gtree.GTree
	GetType() string
	GetMinID() *uint64
	GetMaxID() *uint64
	GetMinClaim() ClaimObject
	GetMaxClaim() ClaimObject
	ValidateSyntax() field.ErrorList
}

type ClaimObject interface {
	Object
	GetIndex() string                   // implement
	GetSelector() *metav1.LabelSelector // implement
	GetOwnerSelector() (labels.Selector, error)
	GetClaimType() ClaimType
	GetLabelSelector() (labels.Selector, error)
	GetClaimLabels() labels.Set
	ValidateOwner(labels labels.Set) error
	GetStaticID() *uint64
	GetStaticTreeID(t string) tree.ID
	GetClaimID(t string, id uint64) tree.ID
	GetRange() *string
	GetRangeID(t string) (tree.Range, error)
	GetTable(t string, to, from uint64) table.Table
	SetStatusRange(*string)
	SetStatusID(*uint64)
	GetStatusID() *uint64
	ValidateSyntax(s string) field.ErrorList
	GetClaimRequest() string
	GetClaimResponse() string
}

type EntryObject interface {
	Object
	GetIndex() string
	GetClaimType() ClaimType
	GetOwnerGVK() schema.GroupVersionKind
	GetOwnerNSN() types.NamespacedName
	SetSpec(x any)
	GetSpec() any
	GetSpecID() string
}

type GenericObject interface {
	Object
	ValidateSyntax() field.ErrorList
	SetSpec(x any)
	GetSpec() any
	NewObjList() GenericObjectList
	GroupVersionKind() schema.GroupVersionKind
}

type Object interface {
	client.Object
	GetNamespacedName() types.NamespacedName
	GetKey() store.Key
	GetOwnerReference() *commonv1alpha1.OwnerReference
	GetObjectMeta() *metav1.ObjectMeta
	GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition
	SetConditions(c ...conditionv1alpha1.Condition)
	CalculateHash() ([sha1.Size]byte, error)
}

type ObjectList interface {
	GetItems() []Object
	client.ObjectList
}

type GenericObjectList interface {
	GetObjects() []GenericObject
	client.ObjectList
}

type Filter interface {
	Filter(ctx context.Context, obj runtime.Object) bool
}
