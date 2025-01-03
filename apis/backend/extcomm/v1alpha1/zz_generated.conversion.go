//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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
// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	backend "github.com/kuidio/kuid/apis/backend"
	asv1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	extcomm "github.com/kuidio/kuid/apis/backend/extcomm"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*EXTCOMMClaim)(nil), (*extcomm.EXTCOMMClaim)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMClaim_To_extcomm_EXTCOMMClaim(a.(*EXTCOMMClaim), b.(*extcomm.EXTCOMMClaim), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMClaim)(nil), (*EXTCOMMClaim)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMClaim_To_v1alpha1_EXTCOMMClaim(a.(*extcomm.EXTCOMMClaim), b.(*EXTCOMMClaim), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMClaimList)(nil), (*extcomm.EXTCOMMClaimList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMClaimList_To_extcomm_EXTCOMMClaimList(a.(*EXTCOMMClaimList), b.(*extcomm.EXTCOMMClaimList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMClaimList)(nil), (*EXTCOMMClaimList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMClaimList_To_v1alpha1_EXTCOMMClaimList(a.(*extcomm.EXTCOMMClaimList), b.(*EXTCOMMClaimList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMClaimSpec)(nil), (*extcomm.EXTCOMMClaimSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMClaimSpec_To_extcomm_EXTCOMMClaimSpec(a.(*EXTCOMMClaimSpec), b.(*extcomm.EXTCOMMClaimSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMClaimSpec)(nil), (*EXTCOMMClaimSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMClaimSpec_To_v1alpha1_EXTCOMMClaimSpec(a.(*extcomm.EXTCOMMClaimSpec), b.(*EXTCOMMClaimSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMClaimStatus)(nil), (*extcomm.EXTCOMMClaimStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMClaimStatus_To_extcomm_EXTCOMMClaimStatus(a.(*EXTCOMMClaimStatus), b.(*extcomm.EXTCOMMClaimStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMClaimStatus)(nil), (*EXTCOMMClaimStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMClaimStatus_To_v1alpha1_EXTCOMMClaimStatus(a.(*extcomm.EXTCOMMClaimStatus), b.(*EXTCOMMClaimStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMEntry)(nil), (*extcomm.EXTCOMMEntry)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMEntry_To_extcomm_EXTCOMMEntry(a.(*EXTCOMMEntry), b.(*extcomm.EXTCOMMEntry), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMEntry)(nil), (*EXTCOMMEntry)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMEntry_To_v1alpha1_EXTCOMMEntry(a.(*extcomm.EXTCOMMEntry), b.(*EXTCOMMEntry), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMEntryList)(nil), (*extcomm.EXTCOMMEntryList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMEntryList_To_extcomm_EXTCOMMEntryList(a.(*EXTCOMMEntryList), b.(*extcomm.EXTCOMMEntryList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMEntryList)(nil), (*EXTCOMMEntryList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMEntryList_To_v1alpha1_EXTCOMMEntryList(a.(*extcomm.EXTCOMMEntryList), b.(*EXTCOMMEntryList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMEntrySpec)(nil), (*extcomm.EXTCOMMEntrySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMEntrySpec_To_extcomm_EXTCOMMEntrySpec(a.(*EXTCOMMEntrySpec), b.(*extcomm.EXTCOMMEntrySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMEntrySpec)(nil), (*EXTCOMMEntrySpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMEntrySpec_To_v1alpha1_EXTCOMMEntrySpec(a.(*extcomm.EXTCOMMEntrySpec), b.(*EXTCOMMEntrySpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMEntryStatus)(nil), (*extcomm.EXTCOMMEntryStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMEntryStatus_To_extcomm_EXTCOMMEntryStatus(a.(*EXTCOMMEntryStatus), b.(*extcomm.EXTCOMMEntryStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMEntryStatus)(nil), (*EXTCOMMEntryStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMEntryStatus_To_v1alpha1_EXTCOMMEntryStatus(a.(*extcomm.EXTCOMMEntryStatus), b.(*EXTCOMMEntryStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMIndex)(nil), (*extcomm.EXTCOMMIndex)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMIndex_To_extcomm_EXTCOMMIndex(a.(*EXTCOMMIndex), b.(*extcomm.EXTCOMMIndex), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMIndex)(nil), (*EXTCOMMIndex)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMIndex_To_v1alpha1_EXTCOMMIndex(a.(*extcomm.EXTCOMMIndex), b.(*EXTCOMMIndex), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMIndexClaim)(nil), (*extcomm.EXTCOMMIndexClaim)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMIndexClaim_To_extcomm_EXTCOMMIndexClaim(a.(*EXTCOMMIndexClaim), b.(*extcomm.EXTCOMMIndexClaim), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMIndexClaim)(nil), (*EXTCOMMIndexClaim)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMIndexClaim_To_v1alpha1_EXTCOMMIndexClaim(a.(*extcomm.EXTCOMMIndexClaim), b.(*EXTCOMMIndexClaim), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMIndexList)(nil), (*extcomm.EXTCOMMIndexList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMIndexList_To_extcomm_EXTCOMMIndexList(a.(*EXTCOMMIndexList), b.(*extcomm.EXTCOMMIndexList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMIndexList)(nil), (*EXTCOMMIndexList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMIndexList_To_v1alpha1_EXTCOMMIndexList(a.(*extcomm.EXTCOMMIndexList), b.(*EXTCOMMIndexList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMIndexSpec)(nil), (*extcomm.EXTCOMMIndexSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMIndexSpec_To_extcomm_EXTCOMMIndexSpec(a.(*EXTCOMMIndexSpec), b.(*extcomm.EXTCOMMIndexSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMIndexSpec)(nil), (*EXTCOMMIndexSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMIndexSpec_To_v1alpha1_EXTCOMMIndexSpec(a.(*extcomm.EXTCOMMIndexSpec), b.(*EXTCOMMIndexSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EXTCOMMIndexStatus)(nil), (*extcomm.EXTCOMMIndexStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EXTCOMMIndexStatus_To_extcomm_EXTCOMMIndexStatus(a.(*EXTCOMMIndexStatus), b.(*extcomm.EXTCOMMIndexStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*extcomm.EXTCOMMIndexStatus)(nil), (*EXTCOMMIndexStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_extcomm_EXTCOMMIndexStatus_To_v1alpha1_EXTCOMMIndexStatus(a.(*extcomm.EXTCOMMIndexStatus), b.(*EXTCOMMIndexStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_EXTCOMMClaim_To_extcomm_EXTCOMMClaim(in *EXTCOMMClaim, out *extcomm.EXTCOMMClaim, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_EXTCOMMClaimSpec_To_extcomm_EXTCOMMClaimSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_EXTCOMMClaimStatus_To_extcomm_EXTCOMMClaimStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMClaim_To_extcomm_EXTCOMMClaim is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMClaim_To_extcomm_EXTCOMMClaim(in *EXTCOMMClaim, out *extcomm.EXTCOMMClaim, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMClaim_To_extcomm_EXTCOMMClaim(in, out, s)
}

func autoConvert_extcomm_EXTCOMMClaim_To_v1alpha1_EXTCOMMClaim(in *extcomm.EXTCOMMClaim, out *EXTCOMMClaim, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_extcomm_EXTCOMMClaimSpec_To_v1alpha1_EXTCOMMClaimSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_extcomm_EXTCOMMClaimStatus_To_v1alpha1_EXTCOMMClaimStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_extcomm_EXTCOMMClaim_To_v1alpha1_EXTCOMMClaim is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMClaim_To_v1alpha1_EXTCOMMClaim(in *extcomm.EXTCOMMClaim, out *EXTCOMMClaim, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMClaim_To_v1alpha1_EXTCOMMClaim(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMClaimList_To_extcomm_EXTCOMMClaimList(in *EXTCOMMClaimList, out *extcomm.EXTCOMMClaimList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]extcomm.EXTCOMMClaim, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_EXTCOMMClaim_To_extcomm_EXTCOMMClaim(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMClaimList_To_extcomm_EXTCOMMClaimList is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMClaimList_To_extcomm_EXTCOMMClaimList(in *EXTCOMMClaimList, out *extcomm.EXTCOMMClaimList, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMClaimList_To_extcomm_EXTCOMMClaimList(in, out, s)
}

func autoConvert_extcomm_EXTCOMMClaimList_To_v1alpha1_EXTCOMMClaimList(in *extcomm.EXTCOMMClaimList, out *EXTCOMMClaimList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EXTCOMMClaim, len(*in))
		for i := range *in {
			if err := Convert_extcomm_EXTCOMMClaim_To_v1alpha1_EXTCOMMClaim(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_extcomm_EXTCOMMClaimList_To_v1alpha1_EXTCOMMClaimList is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMClaimList_To_v1alpha1_EXTCOMMClaimList(in *extcomm.EXTCOMMClaimList, out *EXTCOMMClaimList, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMClaimList_To_v1alpha1_EXTCOMMClaimList(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMClaimSpec_To_extcomm_EXTCOMMClaimSpec(in *EXTCOMMClaimSpec, out *extcomm.EXTCOMMClaimSpec, s conversion.Scope) error {
	out.Index = in.Index
	out.ID = (*uint64)(unsafe.Pointer(in.ID))
	out.Range = (*string)(unsafe.Pointer(in.Range))
	if err := asv1alpha1.Convert_v1alpha1_ClaimLabels_To_common_ClaimLabels(&in.ClaimLabels, &out.ClaimLabels, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMClaimSpec_To_extcomm_EXTCOMMClaimSpec is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMClaimSpec_To_extcomm_EXTCOMMClaimSpec(in *EXTCOMMClaimSpec, out *extcomm.EXTCOMMClaimSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMClaimSpec_To_extcomm_EXTCOMMClaimSpec(in, out, s)
}

func autoConvert_extcomm_EXTCOMMClaimSpec_To_v1alpha1_EXTCOMMClaimSpec(in *extcomm.EXTCOMMClaimSpec, out *EXTCOMMClaimSpec, s conversion.Scope) error {
	out.Index = in.Index
	out.ID = (*uint64)(unsafe.Pointer(in.ID))
	out.Range = (*string)(unsafe.Pointer(in.Range))
	if err := asv1alpha1.Convert_common_ClaimLabels_To_v1alpha1_ClaimLabels(&in.ClaimLabels, &out.ClaimLabels, s); err != nil {
		return err
	}
	return nil
}

// Convert_extcomm_EXTCOMMClaimSpec_To_v1alpha1_EXTCOMMClaimSpec is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMClaimSpec_To_v1alpha1_EXTCOMMClaimSpec(in *extcomm.EXTCOMMClaimSpec, out *EXTCOMMClaimSpec, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMClaimSpec_To_v1alpha1_EXTCOMMClaimSpec(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMClaimStatus_To_extcomm_EXTCOMMClaimStatus(in *EXTCOMMClaimStatus, out *extcomm.EXTCOMMClaimStatus, s conversion.Scope) error {
	if err := asv1alpha1.Convert_v1alpha1_ConditionedStatus_To_condition_ConditionedStatus(&in.ConditionedStatus, &out.ConditionedStatus, s); err != nil {
		return err
	}
	out.ID = (*uint64)(unsafe.Pointer(in.ID))
	out.Range = (*string)(unsafe.Pointer(in.Range))
	out.ExpiryTime = (*string)(unsafe.Pointer(in.ExpiryTime))
	return nil
}

// Convert_v1alpha1_EXTCOMMClaimStatus_To_extcomm_EXTCOMMClaimStatus is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMClaimStatus_To_extcomm_EXTCOMMClaimStatus(in *EXTCOMMClaimStatus, out *extcomm.EXTCOMMClaimStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMClaimStatus_To_extcomm_EXTCOMMClaimStatus(in, out, s)
}

func autoConvert_extcomm_EXTCOMMClaimStatus_To_v1alpha1_EXTCOMMClaimStatus(in *extcomm.EXTCOMMClaimStatus, out *EXTCOMMClaimStatus, s conversion.Scope) error {
	if err := asv1alpha1.Convert_condition_ConditionedStatus_To_v1alpha1_ConditionedStatus(&in.ConditionedStatus, &out.ConditionedStatus, s); err != nil {
		return err
	}
	out.ID = (*uint64)(unsafe.Pointer(in.ID))
	out.Range = (*string)(unsafe.Pointer(in.Range))
	out.ExpiryTime = (*string)(unsafe.Pointer(in.ExpiryTime))
	return nil
}

// Convert_extcomm_EXTCOMMClaimStatus_To_v1alpha1_EXTCOMMClaimStatus is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMClaimStatus_To_v1alpha1_EXTCOMMClaimStatus(in *extcomm.EXTCOMMClaimStatus, out *EXTCOMMClaimStatus, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMClaimStatus_To_v1alpha1_EXTCOMMClaimStatus(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMEntry_To_extcomm_EXTCOMMEntry(in *EXTCOMMEntry, out *extcomm.EXTCOMMEntry, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_EXTCOMMEntrySpec_To_extcomm_EXTCOMMEntrySpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_EXTCOMMEntryStatus_To_extcomm_EXTCOMMEntryStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMEntry_To_extcomm_EXTCOMMEntry is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMEntry_To_extcomm_EXTCOMMEntry(in *EXTCOMMEntry, out *extcomm.EXTCOMMEntry, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMEntry_To_extcomm_EXTCOMMEntry(in, out, s)
}

func autoConvert_extcomm_EXTCOMMEntry_To_v1alpha1_EXTCOMMEntry(in *extcomm.EXTCOMMEntry, out *EXTCOMMEntry, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_extcomm_EXTCOMMEntrySpec_To_v1alpha1_EXTCOMMEntrySpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_extcomm_EXTCOMMEntryStatus_To_v1alpha1_EXTCOMMEntryStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_extcomm_EXTCOMMEntry_To_v1alpha1_EXTCOMMEntry is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMEntry_To_v1alpha1_EXTCOMMEntry(in *extcomm.EXTCOMMEntry, out *EXTCOMMEntry, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMEntry_To_v1alpha1_EXTCOMMEntry(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMEntryList_To_extcomm_EXTCOMMEntryList(in *EXTCOMMEntryList, out *extcomm.EXTCOMMEntryList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]extcomm.EXTCOMMEntry, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_EXTCOMMEntry_To_extcomm_EXTCOMMEntry(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMEntryList_To_extcomm_EXTCOMMEntryList is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMEntryList_To_extcomm_EXTCOMMEntryList(in *EXTCOMMEntryList, out *extcomm.EXTCOMMEntryList, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMEntryList_To_extcomm_EXTCOMMEntryList(in, out, s)
}

func autoConvert_extcomm_EXTCOMMEntryList_To_v1alpha1_EXTCOMMEntryList(in *extcomm.EXTCOMMEntryList, out *EXTCOMMEntryList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EXTCOMMEntry, len(*in))
		for i := range *in {
			if err := Convert_extcomm_EXTCOMMEntry_To_v1alpha1_EXTCOMMEntry(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_extcomm_EXTCOMMEntryList_To_v1alpha1_EXTCOMMEntryList is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMEntryList_To_v1alpha1_EXTCOMMEntryList(in *extcomm.EXTCOMMEntryList, out *EXTCOMMEntryList, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMEntryList_To_v1alpha1_EXTCOMMEntryList(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMEntrySpec_To_extcomm_EXTCOMMEntrySpec(in *EXTCOMMEntrySpec, out *extcomm.EXTCOMMEntrySpec, s conversion.Scope) error {
	out.Index = in.Index
	out.IndexEntry = in.IndexEntry
	out.ClaimType = backend.ClaimType(in.ClaimType)
	out.ID = in.ID
	if err := asv1alpha1.Convert_v1alpha1_ClaimLabels_To_common_ClaimLabels(&in.ClaimLabels, &out.ClaimLabels, s); err != nil {
		return err
	}
	out.Claim = in.Claim
	return nil
}

// Convert_v1alpha1_EXTCOMMEntrySpec_To_extcomm_EXTCOMMEntrySpec is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMEntrySpec_To_extcomm_EXTCOMMEntrySpec(in *EXTCOMMEntrySpec, out *extcomm.EXTCOMMEntrySpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMEntrySpec_To_extcomm_EXTCOMMEntrySpec(in, out, s)
}

func autoConvert_extcomm_EXTCOMMEntrySpec_To_v1alpha1_EXTCOMMEntrySpec(in *extcomm.EXTCOMMEntrySpec, out *EXTCOMMEntrySpec, s conversion.Scope) error {
	out.Index = in.Index
	out.IndexEntry = in.IndexEntry
	out.ClaimType = backend.ClaimType(in.ClaimType)
	out.ID = in.ID
	if err := asv1alpha1.Convert_common_ClaimLabels_To_v1alpha1_ClaimLabels(&in.ClaimLabels, &out.ClaimLabels, s); err != nil {
		return err
	}
	out.Claim = in.Claim
	return nil
}

// Convert_extcomm_EXTCOMMEntrySpec_To_v1alpha1_EXTCOMMEntrySpec is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMEntrySpec_To_v1alpha1_EXTCOMMEntrySpec(in *extcomm.EXTCOMMEntrySpec, out *EXTCOMMEntrySpec, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMEntrySpec_To_v1alpha1_EXTCOMMEntrySpec(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMEntryStatus_To_extcomm_EXTCOMMEntryStatus(in *EXTCOMMEntryStatus, out *extcomm.EXTCOMMEntryStatus, s conversion.Scope) error {
	if err := asv1alpha1.Convert_v1alpha1_ConditionedStatus_To_condition_ConditionedStatus(&in.ConditionedStatus, &out.ConditionedStatus, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMEntryStatus_To_extcomm_EXTCOMMEntryStatus is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMEntryStatus_To_extcomm_EXTCOMMEntryStatus(in *EXTCOMMEntryStatus, out *extcomm.EXTCOMMEntryStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMEntryStatus_To_extcomm_EXTCOMMEntryStatus(in, out, s)
}

func autoConvert_extcomm_EXTCOMMEntryStatus_To_v1alpha1_EXTCOMMEntryStatus(in *extcomm.EXTCOMMEntryStatus, out *EXTCOMMEntryStatus, s conversion.Scope) error {
	if err := asv1alpha1.Convert_condition_ConditionedStatus_To_v1alpha1_ConditionedStatus(&in.ConditionedStatus, &out.ConditionedStatus, s); err != nil {
		return err
	}
	return nil
}

// Convert_extcomm_EXTCOMMEntryStatus_To_v1alpha1_EXTCOMMEntryStatus is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMEntryStatus_To_v1alpha1_EXTCOMMEntryStatus(in *extcomm.EXTCOMMEntryStatus, out *EXTCOMMEntryStatus, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMEntryStatus_To_v1alpha1_EXTCOMMEntryStatus(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMIndex_To_extcomm_EXTCOMMIndex(in *EXTCOMMIndex, out *extcomm.EXTCOMMIndex, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_EXTCOMMIndexSpec_To_extcomm_EXTCOMMIndexSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_EXTCOMMIndexStatus_To_extcomm_EXTCOMMIndexStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMIndex_To_extcomm_EXTCOMMIndex is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMIndex_To_extcomm_EXTCOMMIndex(in *EXTCOMMIndex, out *extcomm.EXTCOMMIndex, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMIndex_To_extcomm_EXTCOMMIndex(in, out, s)
}

func autoConvert_extcomm_EXTCOMMIndex_To_v1alpha1_EXTCOMMIndex(in *extcomm.EXTCOMMIndex, out *EXTCOMMIndex, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_extcomm_EXTCOMMIndexSpec_To_v1alpha1_EXTCOMMIndexSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_extcomm_EXTCOMMIndexStatus_To_v1alpha1_EXTCOMMIndexStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_extcomm_EXTCOMMIndex_To_v1alpha1_EXTCOMMIndex is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMIndex_To_v1alpha1_EXTCOMMIndex(in *extcomm.EXTCOMMIndex, out *EXTCOMMIndex, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMIndex_To_v1alpha1_EXTCOMMIndex(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMIndexClaim_To_extcomm_EXTCOMMIndexClaim(in *EXTCOMMIndexClaim, out *extcomm.EXTCOMMIndexClaim, s conversion.Scope) error {
	out.Name = in.Name
	out.ID = (*uint64)(unsafe.Pointer(in.ID))
	out.Range = (*string)(unsafe.Pointer(in.Range))
	if err := asv1alpha1.Convert_v1alpha1_UserDefinedLabels_To_common_UserDefinedLabels(&in.UserDefinedLabels, &out.UserDefinedLabels, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMIndexClaim_To_extcomm_EXTCOMMIndexClaim is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMIndexClaim_To_extcomm_EXTCOMMIndexClaim(in *EXTCOMMIndexClaim, out *extcomm.EXTCOMMIndexClaim, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMIndexClaim_To_extcomm_EXTCOMMIndexClaim(in, out, s)
}

func autoConvert_extcomm_EXTCOMMIndexClaim_To_v1alpha1_EXTCOMMIndexClaim(in *extcomm.EXTCOMMIndexClaim, out *EXTCOMMIndexClaim, s conversion.Scope) error {
	out.Name = in.Name
	out.ID = (*uint64)(unsafe.Pointer(in.ID))
	out.Range = (*string)(unsafe.Pointer(in.Range))
	if err := asv1alpha1.Convert_common_UserDefinedLabels_To_v1alpha1_UserDefinedLabels(&in.UserDefinedLabels, &out.UserDefinedLabels, s); err != nil {
		return err
	}
	return nil
}

// Convert_extcomm_EXTCOMMIndexClaim_To_v1alpha1_EXTCOMMIndexClaim is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMIndexClaim_To_v1alpha1_EXTCOMMIndexClaim(in *extcomm.EXTCOMMIndexClaim, out *EXTCOMMIndexClaim, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMIndexClaim_To_v1alpha1_EXTCOMMIndexClaim(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMIndexList_To_extcomm_EXTCOMMIndexList(in *EXTCOMMIndexList, out *extcomm.EXTCOMMIndexList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]extcomm.EXTCOMMIndex, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_EXTCOMMIndex_To_extcomm_EXTCOMMIndex(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMIndexList_To_extcomm_EXTCOMMIndexList is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMIndexList_To_extcomm_EXTCOMMIndexList(in *EXTCOMMIndexList, out *extcomm.EXTCOMMIndexList, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMIndexList_To_extcomm_EXTCOMMIndexList(in, out, s)
}

func autoConvert_extcomm_EXTCOMMIndexList_To_v1alpha1_EXTCOMMIndexList(in *extcomm.EXTCOMMIndexList, out *EXTCOMMIndexList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EXTCOMMIndex, len(*in))
		for i := range *in {
			if err := Convert_extcomm_EXTCOMMIndex_To_v1alpha1_EXTCOMMIndex(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_extcomm_EXTCOMMIndexList_To_v1alpha1_EXTCOMMIndexList is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMIndexList_To_v1alpha1_EXTCOMMIndexList(in *extcomm.EXTCOMMIndexList, out *EXTCOMMIndexList, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMIndexList_To_v1alpha1_EXTCOMMIndexList(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMIndexSpec_To_extcomm_EXTCOMMIndexSpec(in *EXTCOMMIndexSpec, out *extcomm.EXTCOMMIndexSpec, s conversion.Scope) error {
	out.MinID = (*uint64)(unsafe.Pointer(in.MinID))
	out.MaxID = (*uint64)(unsafe.Pointer(in.MaxID))
	if err := asv1alpha1.Convert_v1alpha1_UserDefinedLabels_To_common_UserDefinedLabels(&in.UserDefinedLabels, &out.UserDefinedLabels, s); err != nil {
		return err
	}
	out.Transitive = in.Transitive
	out.Type = in.Type
	out.SubType = in.SubType
	out.GlobalID = in.GlobalID
	if in.Claims != nil {
		in, out := &in.Claims, &out.Claims
		*out = make([]extcomm.EXTCOMMIndexClaim, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_EXTCOMMIndexClaim_To_extcomm_EXTCOMMIndexClaim(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Claims = nil
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMIndexSpec_To_extcomm_EXTCOMMIndexSpec is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMIndexSpec_To_extcomm_EXTCOMMIndexSpec(in *EXTCOMMIndexSpec, out *extcomm.EXTCOMMIndexSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMIndexSpec_To_extcomm_EXTCOMMIndexSpec(in, out, s)
}

func autoConvert_extcomm_EXTCOMMIndexSpec_To_v1alpha1_EXTCOMMIndexSpec(in *extcomm.EXTCOMMIndexSpec, out *EXTCOMMIndexSpec, s conversion.Scope) error {
	out.MinID = (*uint64)(unsafe.Pointer(in.MinID))
	out.MaxID = (*uint64)(unsafe.Pointer(in.MaxID))
	if err := asv1alpha1.Convert_common_UserDefinedLabels_To_v1alpha1_UserDefinedLabels(&in.UserDefinedLabels, &out.UserDefinedLabels, s); err != nil {
		return err
	}
	out.Transitive = in.Transitive
	out.Type = in.Type
	out.SubType = in.SubType
	out.GlobalID = in.GlobalID
	if in.Claims != nil {
		in, out := &in.Claims, &out.Claims
		*out = make([]EXTCOMMIndexClaim, len(*in))
		for i := range *in {
			if err := Convert_extcomm_EXTCOMMIndexClaim_To_v1alpha1_EXTCOMMIndexClaim(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Claims = nil
	}
	return nil
}

// Convert_extcomm_EXTCOMMIndexSpec_To_v1alpha1_EXTCOMMIndexSpec is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMIndexSpec_To_v1alpha1_EXTCOMMIndexSpec(in *extcomm.EXTCOMMIndexSpec, out *EXTCOMMIndexSpec, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMIndexSpec_To_v1alpha1_EXTCOMMIndexSpec(in, out, s)
}

func autoConvert_v1alpha1_EXTCOMMIndexStatus_To_extcomm_EXTCOMMIndexStatus(in *EXTCOMMIndexStatus, out *extcomm.EXTCOMMIndexStatus, s conversion.Scope) error {
	out.MinID = (*int64)(unsafe.Pointer(in.MinID))
	out.MaxID = (*int64)(unsafe.Pointer(in.MaxID))
	if err := asv1alpha1.Convert_v1alpha1_ConditionedStatus_To_condition_ConditionedStatus(&in.ConditionedStatus, &out.ConditionedStatus, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_EXTCOMMIndexStatus_To_extcomm_EXTCOMMIndexStatus is an autogenerated conversion function.
func Convert_v1alpha1_EXTCOMMIndexStatus_To_extcomm_EXTCOMMIndexStatus(in *EXTCOMMIndexStatus, out *extcomm.EXTCOMMIndexStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_EXTCOMMIndexStatus_To_extcomm_EXTCOMMIndexStatus(in, out, s)
}

func autoConvert_extcomm_EXTCOMMIndexStatus_To_v1alpha1_EXTCOMMIndexStatus(in *extcomm.EXTCOMMIndexStatus, out *EXTCOMMIndexStatus, s conversion.Scope) error {
	out.MinID = (*int64)(unsafe.Pointer(in.MinID))
	out.MaxID = (*int64)(unsafe.Pointer(in.MaxID))
	if err := asv1alpha1.Convert_condition_ConditionedStatus_To_v1alpha1_ConditionedStatus(&in.ConditionedStatus, &out.ConditionedStatus, s); err != nil {
		return err
	}
	return nil
}

// Convert_extcomm_EXTCOMMIndexStatus_To_v1alpha1_EXTCOMMIndexStatus is an autogenerated conversion function.
func Convert_extcomm_EXTCOMMIndexStatus_To_v1alpha1_EXTCOMMIndexStatus(in *extcomm.EXTCOMMIndexStatus, out *EXTCOMMIndexStatus, s conversion.Scope) error {
	return autoConvert_extcomm_EXTCOMMIndexStatus_To_v1alpha1_EXTCOMMIndexStatus(in, out, s)
}
