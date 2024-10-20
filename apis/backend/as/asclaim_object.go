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

package as

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/henderiw/idxtable/pkg/table"
	"github.com/henderiw/idxtable/pkg/table/table32"
	"github.com/henderiw/idxtable/pkg/tree"
	"github.com/henderiw/idxtable/pkg/tree/id32"
	"github.com/henderiw/store"
	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/backend"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

var _ backend.ClaimObject = &ASClaim{}

func (r *ASClaim) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *ASClaim) GetKey() store.Key {
	return store.KeyFromNSN(types.NamespacedName{Namespace: r.Namespace, Name: r.Spec.Index})
}

// GetCondition returns the condition based on the condition kind
func (r *ASClaim) GetCondition(t condition.ConditionType) condition.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *ASClaim) SetConditions(c ...condition.Condition) {
	r.Status.SetConditions(c...)
}

func (r *ASClaim) ValidateSyntax(s string) field.ErrorList {
	var allErrs field.ErrorList

	if err := r.ValidateASClaimType(); err != nil {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath(""),
			r,
			err.Error(),
		))
		return allErrs
	}
	var v SyntaxValidator
	claimType := r.GetClaimType()
	switch claimType {
	case backend.ClaimType_DynamicID:
		v = &ASDynamicIDSyntaxValidator{Name: string(claimType)}
	case backend.ClaimType_StaticID:
		v = &ASStaticIDSyntaxValidator{Name: string(claimType)}
	case backend.ClaimType_Range:
		v = &ASRangeSyntaxValidator{Name: string(claimType)}
	default:
		return allErrs
	}
	return v.Validate(r)
}

func (r *ASClaim) ValidateASRange() error {
	if r.Spec.Range == nil {
		return fmt.Errorf("no AS range provided")
	}
	parts := strings.SplitN(*r.Spec.Range, "-", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid AS range, expected <start>-<end>, got: %s", *r.Spec.Range)
	}
	var errm error
	if r.Name == r.Spec.Index {
		// to be able to check if the entry is reserved we get a parentname (rang name) equal to index
		// this is because the ownerreference uses the name of the index in its labels in the cache
		errm = errors.Join(errm, fmt.Errorf("a name of range cannot be the same as the index"))
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid AS range start, got: %s, err: %s", *r.Spec.Range, err.Error()))
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid AS range end, got: %s, err: %s", *r.Spec.Range, err.Error()))
	}
	if errm != nil {
		return errm
	}
	if start > end {
		errm = errors.Join(errm, fmt.Errorf("invalid AS range start > end %s", *r.Spec.Range))
	}
	if err := validateASID(start); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid AS start err %s", err.Error()))
	}
	if err := validateASID(end); err != nil {
		errm = errors.Join(errm, fmt.Errorf("invalid AS end err %s", err.Error()))
	}
	return errm
}

func (r *ASClaim) ValidateASID() error {
	if r.Spec.ID == nil {
		return fmt.Errorf("no id provided")
	}
	if err := validateASID(int(*r.Spec.ID)); err != nil {
		return fmt.Errorf("invalid id err %s", err.Error())
	}
	return nil
}

func (r *ASClaim) ValidateASClaimType() error {
	var sb strings.Builder
	count := 0
	if r.Spec.ID != nil {
		sb.WriteString(fmt.Sprintf("id: %d", *r.Spec.ID))
		count++

	}
	if r.Spec.Range != nil {
		if count > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("range: %s", *r.Spec.Range))
		count++

	}
	if count > 1 {
		return fmt.Errorf("a claim can only have 1 type, got %s", sb.String())
	}
	return nil
}

func (r *ASClaim) GetIndex() string { return r.Spec.Index }

func (r *ASClaim) GetSelector() *metav1.LabelSelector { return r.Spec.Selector }

// GetOwnerSelector selects the route based on the name of the claim
func (r *ASClaim) GetOwnerSelector() (labels.Selector, error) {
	claimName := r.Name
	claimKind := r.Kind
	claimUID := r.UID
	for _, owner := range r.GetOwnerReferences() {
		if owner.APIVersion == SchemeGroupVersion.Identifier() &&
			owner.Kind == ASIndexKind {
			claimName = owner.Name
			claimKind = owner.Kind
			claimUID = owner.UID
		}
	}

	l := map[string]string{
		backend.KuidClaimNameKey: claimName,
		backend.KuidClaimUIDKey:  string(claimUID),
		backend.KuidOwnerKindKey: claimKind,
	}

	fullselector := labels.NewSelector()
	for k, v := range l {
		req, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			return nil, err
		}
		fullselector = fullselector.Add(*req)
	}
	return fullselector, nil
}

func (r *ASClaim) GetLabelSelector() (labels.Selector, error) { return r.Spec.GetLabelSelector() }

func (r *ASClaim) GetClaimLabels() labels.Set {
	labels := r.Spec.GetUserDefinedLabels()

	// for claims originated from the index we need to use the ownerreferences, since these claims
	// are never stored in the apiserver, the ip entries need to reference the index instead
	claimName := r.Name
	claimKind := ASClaimKind
	claimUID := r.UID
	for _, owner := range r.GetOwnerReferences() {
		if owner.APIVersion == SchemeGroupVersion.Identifier() &&
			owner.Kind == ASIndexKind {
			claimName = owner.Name
			claimKind = owner.Kind
			claimUID = owner.UID
		}
	}
	// system defined labels
	labels[backend.KuidClaimTypeKey] = string(r.GetClaimType())
	labels[backend.KuidClaimNameKey] = claimName
	labels[backend.KuidClaimUIDKey] = string(claimUID)
	labels[backend.KuidOwnerKindKey] = claimKind
	return labels
}

func (r *ASClaim) ValidateOwner(labels labels.Set) error {
	routeClaimName := labels[backend.KuidClaimNameKey]
	routeClaimUID := labels[backend.KuidClaimUIDKey]

	if string(r.UID) != routeClaimUID && r.Name != routeClaimName {
		return fmt.Errorf("route owned by different claim got name %s/%s uid %s/%s",
			r.Name,
			routeClaimName,
			string(r.UID),
			routeClaimUID,
		)
	}
	return nil
}

func (r *ASClaim) GetClaimType() backend.ClaimType {
	claimType := backend.ClaimType_Invalid
	count := 0
	if r.Spec.ID != nil {
		claimType = backend.ClaimType_StaticID
		count++

	}
	if r.Spec.Range != nil {
		claimType = backend.ClaimType_Range
		count++

	}
	if count > 1 {
		return backend.ClaimType_Invalid
	}
	if count == 0 {
		return backend.ClaimType_DynamicID
	}
	return claimType
}
func (r *ASClaim) GetStaticID() *uint64 {
	if r.Spec.ID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Spec.ID))
}
func (r *ASClaim) GetStaticTreeID(t string) tree.ID {
	if r.Spec.ID == nil {
		return nil
	}
	return id32.NewID(*r.Spec.ID, id32.IDBitSize)
}

func (r *ASClaim) GetClaimID(t string, id uint64) tree.ID {
	return id32.NewID(uint32(id), id32.IDBitSize)
}

func (r *ASClaim) GetRange() *string {
	return r.Spec.Range
}

func (r *ASClaim) GetRangeID(t string) (tree.Range, error) {
	if r.Spec.Range == nil {
		return nil, fmt.Errorf("cannot provide a range without an id")
	}
	return id32.ParseRange(*r.Spec.Range)
}

func (r *ASClaim) GetTable(t string, to, from uint64) table.Table {
	return table32.New(uint32(to), uint32(from))
}

func (r *ASClaim) SetStatusRange(s *string) {
	r.Status.Range = s
}

func (r *ASClaim) SetStatusID(s *uint64) {
	if s == nil {
		r.Status.ID = nil
		return
	}
	r.Status.ID = ptr.To[uint32](uint32(*s))
}

func (r *ASClaim) GetStatusID() *uint64 {
	if r.Status.ID == nil {
		return nil
	}
	return ptr.To[uint64](uint64(*r.Status.ID))
}

func (r *ASClaim) GetClaimRequest() string {
	if r.Spec.ID != nil {
		return getASDot(*r.Spec.ID)
	}
	if r.Spec.Range != nil {
		range32, err := id32.ParseRange(*r.Spec.Range)
		if err != nil {
			return *r.Spec.Range
		}
		return fmt.Sprintf("%s-%s", getASDot(uint32(range32.From().ID())), getASDot(uint32(range32.To().ID())))
	}
	return ""
}

func (r *ASClaim) GetClaimResponse() string {
	if r.Status.ID != nil {
		return getASDot(*r.Status.ID)
	}
	if r.Status.Range != nil {
		range32, err := id32.ParseRange(*r.Status.Range)
		if err != nil {
			return *r.Status.Range
		}
		return fmt.Sprintf("%s-%s", getASDot(uint32(range32.From().ID())), getASDot(uint32(range32.To().ID())))
	}
	return ""
}

func ASClaimFromUnstructured(ru runtime.Unstructured) (backend.ClaimObject, error) {
	obj := &ASClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(ru.UnstructuredContent(), obj)
	if err != nil {
		return nil, fmt.Errorf("error converting unstructured to asIndex: %v", err)
	}
	return obj, nil
}

func ASClaimFromRuntime(ru runtime.Object) (backend.ClaimObject, error) {
	claim, ok := ru.(*ASClaim)
	if !ok {
		return nil, errors.New("runtime object not ASClaim")
	}
	return claim, nil
}

// BuildASClaim returns a reource from a client Object a Spec/Status
func BuildASClaim(meta metav1.ObjectMeta, spec *ASClaimSpec, status *ASClaimStatus) backend.ClaimObject {
	aspec := ASClaimSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := ASClaimStatus{}
	if status != nil {
		astatus = *status
	}
	return &ASClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       ASClaimKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}
