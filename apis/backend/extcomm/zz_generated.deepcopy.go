//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package extcomm

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMClaim) DeepCopyInto(out *EXTCOMMClaim) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMClaim.
func (in *EXTCOMMClaim) DeepCopy() *EXTCOMMClaim {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EXTCOMMClaim) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMClaimFilter) DeepCopyInto(out *EXTCOMMClaimFilter) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMClaimFilter.
func (in *EXTCOMMClaimFilter) DeepCopy() *EXTCOMMClaimFilter {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMClaimFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMClaimList) DeepCopyInto(out *EXTCOMMClaimList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EXTCOMMClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMClaimList.
func (in *EXTCOMMClaimList) DeepCopy() *EXTCOMMClaimList {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMClaimList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EXTCOMMClaimList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMClaimSpec) DeepCopyInto(out *EXTCOMMClaimSpec) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(uint64)
		**out = **in
	}
	if in.Range != nil {
		in, out := &in.Range, &out.Range
		*out = new(string)
		**out = **in
	}
	in.ClaimLabels.DeepCopyInto(&out.ClaimLabels)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMClaimSpec.
func (in *EXTCOMMClaimSpec) DeepCopy() *EXTCOMMClaimSpec {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMClaimSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMClaimStatus) DeepCopyInto(out *EXTCOMMClaimStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(uint64)
		**out = **in
	}
	if in.Range != nil {
		in, out := &in.Range, &out.Range
		*out = new(string)
		**out = **in
	}
	if in.ExpiryTime != nil {
		in, out := &in.ExpiryTime, &out.ExpiryTime
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMClaimStatus.
func (in *EXTCOMMClaimStatus) DeepCopy() *EXTCOMMClaimStatus {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMClaimStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMDynamicIDSyntaxValidator) DeepCopyInto(out *EXTCOMMDynamicIDSyntaxValidator) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMDynamicIDSyntaxValidator.
func (in *EXTCOMMDynamicIDSyntaxValidator) DeepCopy() *EXTCOMMDynamicIDSyntaxValidator {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMDynamicIDSyntaxValidator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMEntry) DeepCopyInto(out *EXTCOMMEntry) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMEntry.
func (in *EXTCOMMEntry) DeepCopy() *EXTCOMMEntry {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMEntry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EXTCOMMEntry) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMEntryFilter) DeepCopyInto(out *EXTCOMMEntryFilter) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMEntryFilter.
func (in *EXTCOMMEntryFilter) DeepCopy() *EXTCOMMEntryFilter {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMEntryFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMEntryList) DeepCopyInto(out *EXTCOMMEntryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EXTCOMMEntry, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMEntryList.
func (in *EXTCOMMEntryList) DeepCopy() *EXTCOMMEntryList {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMEntryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EXTCOMMEntryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMEntrySpec) DeepCopyInto(out *EXTCOMMEntrySpec) {
	*out = *in
	in.ClaimLabels.DeepCopyInto(&out.ClaimLabels)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMEntrySpec.
func (in *EXTCOMMEntrySpec) DeepCopy() *EXTCOMMEntrySpec {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMEntrySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMEntryStatus) DeepCopyInto(out *EXTCOMMEntryStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMEntryStatus.
func (in *EXTCOMMEntryStatus) DeepCopy() *EXTCOMMEntryStatus {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMEntryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMIndex) DeepCopyInto(out *EXTCOMMIndex) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMIndex.
func (in *EXTCOMMIndex) DeepCopy() *EXTCOMMIndex {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMIndex)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EXTCOMMIndex) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMIndexFilter) DeepCopyInto(out *EXTCOMMIndexFilter) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMIndexFilter.
func (in *EXTCOMMIndexFilter) DeepCopy() *EXTCOMMIndexFilter {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMIndexFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMIndexList) DeepCopyInto(out *EXTCOMMIndexList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EXTCOMMIndex, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMIndexList.
func (in *EXTCOMMIndexList) DeepCopy() *EXTCOMMIndexList {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMIndexList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EXTCOMMIndexList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMIndexSpec) DeepCopyInto(out *EXTCOMMIndexSpec) {
	*out = *in
	if in.MinID != nil {
		in, out := &in.MinID, &out.MinID
		*out = new(uint64)
		**out = **in
	}
	if in.MaxID != nil {
		in, out := &in.MaxID, &out.MaxID
		*out = new(uint64)
		**out = **in
	}
	in.UserDefinedLabels.DeepCopyInto(&out.UserDefinedLabels)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMIndexSpec.
func (in *EXTCOMMIndexSpec) DeepCopy() *EXTCOMMIndexSpec {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMIndexSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMIndexStatus) DeepCopyInto(out *EXTCOMMIndexStatus) {
	*out = *in
	if in.MinID != nil {
		in, out := &in.MinID, &out.MinID
		*out = new(int64)
		**out = **in
	}
	if in.MaxID != nil {
		in, out := &in.MaxID, &out.MaxID
		*out = new(int64)
		**out = **in
	}
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMIndexStatus.
func (in *EXTCOMMIndexStatus) DeepCopy() *EXTCOMMIndexStatus {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMIndexStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMRangeSyntaxValidator) DeepCopyInto(out *EXTCOMMRangeSyntaxValidator) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMRangeSyntaxValidator.
func (in *EXTCOMMRangeSyntaxValidator) DeepCopy() *EXTCOMMRangeSyntaxValidator {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMRangeSyntaxValidator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EXTCOMMStaticIDSyntaxValidator) DeepCopyInto(out *EXTCOMMStaticIDSyntaxValidator) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EXTCOMMStaticIDSyntaxValidator.
func (in *EXTCOMMStaticIDSyntaxValidator) DeepCopy() *EXTCOMMStaticIDSyntaxValidator {
	if in == nil {
		return nil
	}
	out := new(EXTCOMMStaticIDSyntaxValidator)
	in.DeepCopyInto(out)
	return out
}
