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

package vlanentry

import (
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	vlanbe1v1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			entry, ok := obj.(*vlanbe1v1alpha1.VLANEntry)
			if !ok {
				return nil
			}
			return []interface{}{
				entry.Name,
				entry.GetCondition(conditionv1alpha1.ConditionTypeReady).Status,
				entry.Spec.Index,
				entry.Spec.ClaimType,
				entry.Spec.ID,
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			{Name: "Index", Type: "string"},
			{Name: "ClaimType", Type: "string"},
			{Name: "VLANID", Type: "integer"},
		},
	}
}
