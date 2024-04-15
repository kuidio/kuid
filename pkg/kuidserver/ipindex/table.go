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

package ipindex

import (
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	ipambe1v1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
)

func NewTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			IPIndex, ok := obj.(*ipambe1v1alpha1.IPIndex)
			if !ok {
				return nil
			}
			return []interface{}{
				IPIndex.Name,
				IPIndex.GetCondition(conditionv1alpha1.ConditionTypeReady).Status,
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
		},
	}
}
