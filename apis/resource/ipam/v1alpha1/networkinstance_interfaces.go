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

package v1alpha1

import (
	"github.com/henderiw/iputil"
	"github.com/henderiw/store"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	"github.com/kuidio/kuid/pkg/reconcilers/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

// GetCondition returns the condition based on the condition kind
func (r *NetworkInstance) GetCondition(t conditionv1alpha1.ConditionType) conditionv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *NetworkInstance) SetConditions(c ...conditionv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

func (r *NetworkInstanceList) GetItems() []resource.Object {
	objs := []resource.Object{}
	for _, r := range r.Items {
		r := r
		objs = append(objs, &r)
	}
	return objs
}

func (r *NetworkInstance) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

func (r *NetworkInstance) GetOwnerReference() *commonv1alpha1.OwnerReference {
	return &commonv1alpha1.OwnerReference{
		Group:     SchemeGroupVersion.Group,
		Version:   SchemeGroupVersion.Version,
		Kind:      r.Kind,
		Namespace: r.Namespace,
		Name:      r.Name,
	}
}

func (r *NetworkInstance) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *NetworkInstance) GetIPClaim(prefix Prefix) (*ipambev1alpha1.IPClaim, error) {
	pi, err := iputil.New(prefix.Prefix)
	if err != nil {
		return nil, err
	}

	return ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{
			Namespace: r.GetNamespace(),
			Name:      pi.GetSubnetName(),
		},
		&ipambev1alpha1.IPClaimSpec{
			PrefixType:      ptr.To[ipambev1alpha1.IPPrefixType](ipambev1alpha1.IPPrefixType_Aggregate),
			NetworkInstance: r.Name,
			Prefix:          ptr.To[string](prefix.Prefix),
			PrefixLength:    ptr.To[uint32](uint32(pi.GetPrefixLength())),
			CreatePrefix:    ptr.To[bool](true),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: prefix.UserDefinedLabels,
			},
			Owner: commonv1alpha1.GetOwnerReference(r),
		},
		nil,
	), nil

}

func BuildNetworkInstance(meta metav1.ObjectMeta, spec *NetworkInstanceSpec, status *NetworkInstanceStatus) *NetworkInstance {
	aspec := NetworkInstanceSpec{}
	if spec != nil {
		aspec = *spec
	}
	astatus := NetworkInstanceStatus{}
	if status != nil {
		astatus = *status
	}
	return &NetworkInstance{
		TypeMeta: metav1.TypeMeta{
			APIVersion: SchemeGroupVersion.Identifier(),
			Kind:       NetworkInstanceKind,
		},
		ObjectMeta: meta,
		Spec:       aspec,
		Status:     astatus,
	}
}
