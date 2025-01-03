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
	"github.com/henderiw/store"
	condv1alpha1 "github.com/kform-dev/choreo/apis/condition/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *VLANClaim) GetKey() store.Key {
	return store.KeyFromNSN(r.GetNamespacedName())
}

func (r *VLANClaim) GetNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: r.GetNamespace(),
		Name:      r.GetName(),
	}
}

// GetCondition returns the condition based on the condition kind
func (r *VLANClaim) GetCondition(t condv1alpha1.ConditionType) condv1alpha1.Condition {
	return r.Status.GetCondition(t)
}

// SetConditions sets the conditions on the resource. it allows for 0, 1 or more conditions
// to be set at once
func (r *VLANClaim) SetConditions(c ...condv1alpha1.Condition) {
	r.Status.SetConditions(c...)
}
