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
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			index, ok := obj.(*ipambe1v1alpha1.IPIndex)
			if !ok {
				return nil
			}

			prefixes := make([]string, 5, 5)
			for i, prefix := range index.Spec.Prefixes {
				if i >= 5 {
					break
				}
				prefixes[i] = prefix.Prefix
			}

			return []interface{}{
				index.Name,
				index.GetCondition(conditionv1alpha1.ConditionTypeReady).Status,
				prefixes[0],
				prefixes[1],
				prefixes[2],
				prefixes[3],
				prefixes[4],
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			{Name: "Prefix0", Type: "string"},
			{Name: "Prefix1", Type: "string"},
			{Name: "Prefix2", Type: "string"},
			{Name: "Prefix3", Type: "string"},
			{Name: "Prefix4", Type: "string"},
		},
	}
}
