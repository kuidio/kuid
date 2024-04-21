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
// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	iputil "github.com/henderiw/iputil"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPClaim) DeepCopyInto(out *IPClaim) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPClaim.
func (in *IPClaim) DeepCopy() *IPClaim {
	if in == nil {
		return nil
	}
	out := new(IPClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IPClaim) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPClaimList) DeepCopyInto(out *IPClaimList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IPClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPClaimList.
func (in *IPClaimList) DeepCopy() *IPClaimList {
	if in == nil {
		return nil
	}
	out := new(IPClaimList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IPClaimList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPClaimSpec) DeepCopyInto(out *IPClaimSpec) {
	*out = *in
	if in.PrefixType != nil {
		in, out := &in.PrefixType, &out.PrefixType
		*out = new(IPPrefixType)
		**out = **in
	}
	if in.Prefix != nil {
		in, out := &in.Prefix, &out.Prefix
		*out = new(string)
		**out = **in
	}
	if in.Address != nil {
		in, out := &in.Address, &out.Address
		*out = new(string)
		**out = **in
	}
	if in.Range != nil {
		in, out := &in.Range, &out.Range
		*out = new(string)
		**out = **in
	}
	if in.DefaultGateway != nil {
		in, out := &in.DefaultGateway, &out.DefaultGateway
		*out = new(bool)
		**out = **in
	}
	if in.CreatePrefix != nil {
		in, out := &in.CreatePrefix, &out.CreatePrefix
		*out = new(bool)
		**out = **in
	}
	if in.PrefixLength != nil {
		in, out := &in.PrefixLength, &out.PrefixLength
		*out = new(uint32)
		**out = **in
	}
	if in.AddressFamily != nil {
		in, out := &in.AddressFamily, &out.AddressFamily
		*out = new(iputil.AddressFamily)
		**out = **in
	}
	if in.Index != nil {
		in, out := &in.Index, &out.Index
		*out = new(uint32)
		**out = **in
	}
	in.ClaimLabels.DeepCopyInto(&out.ClaimLabels)
	if in.Owner != nil {
		in, out := &in.Owner, &out.Owner
		*out = new(commonv1alpha1.OwnerReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPClaimSpec.
func (in *IPClaimSpec) DeepCopy() *IPClaimSpec {
	if in == nil {
		return nil
	}
	out := new(IPClaimSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPClaimStatus) DeepCopyInto(out *IPClaimStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	if in.Range != nil {
		in, out := &in.Range, &out.Range
		*out = new(string)
		**out = **in
	}
	if in.Address != nil {
		in, out := &in.Address, &out.Address
		*out = new(string)
		**out = **in
	}
	if in.Prefix != nil {
		in, out := &in.Prefix, &out.Prefix
		*out = new(string)
		**out = **in
	}
	if in.DefaultGateway != nil {
		in, out := &in.DefaultGateway, &out.DefaultGateway
		*out = new(string)
		**out = **in
	}
	if in.ExpiryTime != nil {
		in, out := &in.ExpiryTime, &out.ExpiryTime
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPClaimStatus.
func (in *IPClaimStatus) DeepCopy() *IPClaimStatus {
	if in == nil {
		return nil
	}
	out := new(IPClaimStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPEntry) DeepCopyInto(out *IPEntry) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPEntry.
func (in *IPEntry) DeepCopy() *IPEntry {
	if in == nil {
		return nil
	}
	out := new(IPEntry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IPEntry) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPEntryList) DeepCopyInto(out *IPEntryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IPEntry, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPEntryList.
func (in *IPEntryList) DeepCopy() *IPEntryList {
	if in == nil {
		return nil
	}
	out := new(IPEntryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IPEntryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPEntrySpec) DeepCopyInto(out *IPEntrySpec) {
	*out = *in
	if in.PrefixType != nil {
		in, out := &in.PrefixType, &out.PrefixType
		*out = new(IPPrefixType)
		**out = **in
	}
	if in.DefaultGateway != nil {
		in, out := &in.DefaultGateway, &out.DefaultGateway
		*out = new(bool)
		**out = **in
	}
	if in.AddressFamily != nil {
		in, out := &in.AddressFamily, &out.AddressFamily
		*out = new(iputil.AddressFamily)
		**out = **in
	}
	in.UserDefinedLabels.DeepCopyInto(&out.UserDefinedLabels)
	if in.Owner != nil {
		in, out := &in.Owner, &out.Owner
		*out = new(commonv1alpha1.OwnerReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPEntrySpec.
func (in *IPEntrySpec) DeepCopy() *IPEntrySpec {
	if in == nil {
		return nil
	}
	out := new(IPEntrySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPEntryStatus) DeepCopyInto(out *IPEntryStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPEntryStatus.
func (in *IPEntryStatus) DeepCopy() *IPEntryStatus {
	if in == nil {
		return nil
	}
	out := new(IPEntryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPIndex) DeepCopyInto(out *IPIndex) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPIndex.
func (in *IPIndex) DeepCopy() *IPIndex {
	if in == nil {
		return nil
	}
	out := new(IPIndex)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IPIndex) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPIndexList) DeepCopyInto(out *IPIndexList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IPIndex, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPIndexList.
func (in *IPIndexList) DeepCopy() *IPIndexList {
	if in == nil {
		return nil
	}
	out := new(IPIndexList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IPIndexList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPIndexSpec) DeepCopyInto(out *IPIndexSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPIndexSpec.
func (in *IPIndexSpec) DeepCopy() *IPIndexSpec {
	if in == nil {
		return nil
	}
	out := new(IPIndexSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPIndexStatus) DeepCopyInto(out *IPIndexStatus) {
	*out = *in
	in.ConditionedStatus.DeepCopyInto(&out.ConditionedStatus)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPIndexStatus.
func (in *IPIndexStatus) DeepCopy() *IPIndexStatus {
	if in == nil {
		return nil
	}
	out := new(IPIndexStatus)
	in.DeepCopyInto(out)
	return out
}
