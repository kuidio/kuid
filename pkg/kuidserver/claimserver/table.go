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

package claimserver

/*

import (
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/kuidio/kuid/apis/backend"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	condv1alpha1 "github.com/kform-dev/choreo/apis/condition/v1alpha1"
)

func DefaultTableConvertor(gr schema.GroupResource) registry.TableConvertor {
	return registry.TableConvertor{
		Resource: gr,
		Cells: func(obj runtime.Object) []interface{} {
			claim, ok := obj.(backend.ClaimObject)
			if !ok {
				return nil
			}
			return []interface{}{
				claim.GetName(),
				claim.GetCondition(condv1alpha1.ConditionTypeReady).Status,
				claim.GetIndex(),
				string(claim.GetClaimType()),
				claim.GetClaimRequest(),
				claim.GetClaimResponse(),
			}
		},
		Columns: []metav1.TableColumnDefinition{
			{Name: "Name", Type: "string"},
			{Name: "Ready", Type: "string"},
			{Name: "Index", Type: "string"},
			{Name: "ClaimType", Type: "string"},
			{Name: "ClaimReq", Type: "string"},
			{Name: "ClaimRsp", Type: "string"},
		},
	}
}
*/