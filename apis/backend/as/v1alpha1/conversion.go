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
	"errors"
	unsafe "unsafe"

	"github.com/kform-dev/choreo/apis/condition"
	conditionv1alpha1 "github.com/kform-dev/choreo/apis/condition/v1alpha1"
	"github.com/kuidio/kuid/apis/common"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
)

// Convert_v1alpha1_ConditionedStatus_To_condition_ConditionedStatus is hand made conversion function.
func Convert_v1alpha1_ConditionedStatus_To_condition_ConditionedStatus(in *conditionv1alpha1.ConditionedStatus, out *condition.ConditionedStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_ConditionedStatus_To_condition_ConditionedStatus(in, out, s)
}

func autoConvert_v1alpha1_ConditionedStatus_To_condition_ConditionedStatus(in *conditionv1alpha1.ConditionedStatus, out *condition.ConditionedStatus, _ conversion.Scope) error {
	out.Conditions = *(*[]condition.Condition)(unsafe.Pointer(&in.Conditions))
	return nil
}

// Convert_v1alpha1_ConditionedStatus_To_condition_ConditionedStatus is hand made conversion function.
func Convert_condition_ConditionedStatus_To_v1alpha1_ConditionedStatus(in *condition.ConditionedStatus, out *conditionv1alpha1.ConditionedStatus, s conversion.Scope) error {
	return autoConvert_condition_ConditionedStatus_To_v1alpha1_ConditionedStatus(in, out, s)
}

func autoConvert_condition_ConditionedStatus_To_v1alpha1_ConditionedStatus(in *condition.ConditionedStatus, out *conditionv1alpha1.ConditionedStatus, _ conversion.Scope) error {
	out.Conditions = *(*[]conditionv1alpha1.Condition)(unsafe.Pointer(&in.Conditions))
	return nil
}

// Convert_condition_Condition_To_v1alpha1_Condition is hand made conversion function.
func Convert_condition_Condition_To_v1alpha1_Condition(in *condition.Condition, out *conditionv1alpha1.Condition, s conversion.Scope) error {
	return autoConvert_condition_Condition_To_v1alpha1_Condition(in, out, s)
}

func autoConvert_condition_Condition_To_v1alpha1_Condition(in *condition.Condition, out *conditionv1alpha1.Condition, _ conversion.Scope) error {
	out.Condition = in.Condition
	return nil
}

// Convert_TargetStatus_To_config_TargetStatus is hand made conversion function.
func Convert_v1alpha1_Condition_To_condition_Condition(in *conditionv1alpha1.Condition, out *condition.Condition, s conversion.Scope) error {
	return autoConvert_v1alpha1_Condition_To_condition_Condition(in, out, s)
}

func autoConvert_v1alpha1_Condition_To_condition_Condition(in *conditionv1alpha1.Condition, out *condition.Condition, _ conversion.Scope) error {
	out.Condition = in.Condition
	return nil
}

// Convert_v1alpha1_ConditionedStatus_To_condition_ConditionedStatus is hand made conversion function.
func Convert_common_ClaimLabels_To_v1alpha1_ClaimLabels(in *common.ClaimLabels, out *commonv1alpha1.ClaimLabels, s conversion.Scope) error {
	return autoConvert_common_ClaimLabels_To_v1alpha1_ClaimLabels(in, out, s)
}

func autoConvert_common_ClaimLabels_To_v1alpha1_ClaimLabels(in *common.ClaimLabels, out *commonv1alpha1.ClaimLabels, _ conversion.Scope) error {
	if in == nil {
		return errors.New("input ClaimLabels is nil")
	}
	if out == nil {
		out = &commonv1alpha1.ClaimLabels{} // Allocate new structure if out is nil, depending on the use case this might be handled differently
	}

	// Assuming UserDefinedLabels can be directly copied
	out.UserDefinedLabels = commonv1alpha1.UserDefinedLabels(in.UserDefinedLabels)

	// Manually handle the conversion of the LabelSelector
	if in.Selector != nil {
		out.Selector = &metav1.LabelSelector{}
		if in.Selector.MatchLabels != nil {
			out.Selector.MatchLabels = make(map[string]string)
			for key, value := range in.Selector.MatchLabels {
				out.Selector.MatchLabels[key] = value
			}
		}
		if in.Selector.MatchExpressions != nil {
			out.Selector.MatchExpressions = make([]metav1.LabelSelectorRequirement, len(in.Selector.MatchExpressions))
			for i, expr := range in.Selector.MatchExpressions {
				out.Selector.MatchExpressions[i] = metav1.LabelSelectorRequirement{
					Key:      expr.Key,
					Operator: expr.Operator,
					Values:   append([]string{}, expr.Values...), // Copy slice to avoid reference issues
				}
			}
		}
	} else {
		out.Selector = nil // Explicitly setting to nil if the input is nil
	}

	return nil
}

func Convert_v1alpha1_ClaimLabels_To_common_ClaimLabels(in *commonv1alpha1.ClaimLabels, out *common.ClaimLabels, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClaimLabels_To_common_ClaimLabels(in, out, s)
}

func autoConvert_v1alpha1_ClaimLabels_To_common_ClaimLabels(in *commonv1alpha1.ClaimLabels, out *common.ClaimLabels, _ conversion.Scope) error {
	if in == nil {
		return errors.New("input v1alpha1.ClaimLabels is nil")
	}
	if out == nil {
		out = &common.ClaimLabels{} // Allocate new structure if out is nil
	}

	// Directly copy UserDefinedLabels assuming direct compatibility
	out.UserDefinedLabels = common.UserDefinedLabels(in.UserDefinedLabels)

	// Handle conversion of LabelSelector
	if in.Selector != nil {
		out.Selector = &metav1.LabelSelector{}
		if in.Selector.MatchLabels != nil {
			out.Selector.MatchLabels = make(map[string]string)
			for key, value := range in.Selector.MatchLabels {
				out.Selector.MatchLabels[key] = value
			}
		}
		if in.Selector.MatchExpressions != nil {
			out.Selector.MatchExpressions = make([]metav1.LabelSelectorRequirement, len(in.Selector.MatchExpressions))
			for i, expr := range in.Selector.MatchExpressions {
				out.Selector.MatchExpressions[i] = metav1.LabelSelectorRequirement{
					Key:      expr.Key,
					Operator: metav1.LabelSelectorOperator(expr.Operator),
					Values:   append([]string{}, expr.Values...), // Copy slice to avoid reference issues
				}
			}
		}
	} else {
		out.Selector = nil // Set to nil if the source is nil
	}

	return nil
}

func Convert_common_UserDefinedLabels_To_v1alpha1_UserDefinedLabels(in *common.UserDefinedLabels, out *commonv1alpha1.UserDefinedLabels, s conversion.Scope) error {
	return autoConvert_common_UserDefinedLabels_To_v1alpha1_UserDefinedLabels(in, out, s)
}

func autoConvert_common_UserDefinedLabels_To_v1alpha1_UserDefinedLabels(in *common.UserDefinedLabels, out *commonv1alpha1.UserDefinedLabels, _ conversion.Scope) error {
	in.Labels = out.Labels
	return nil
}

func Convert_v1alpha1_UserDefinedLabels_To_common_UserDefinedLabels(in *commonv1alpha1.UserDefinedLabels, out *common.UserDefinedLabels, s conversion.Scope) error {
	return autoConvert_v1alpha1_UserDefinedLabels_To_common_UserDefinedLabels(in, out, s)
}

func autoConvert_v1alpha1_UserDefinedLabels_To_common_UserDefinedLabels(in *commonv1alpha1.UserDefinedLabels, out *common.UserDefinedLabels, _ conversion.Scope) error {
	in.Labels = out.Labels
	return nil
}
