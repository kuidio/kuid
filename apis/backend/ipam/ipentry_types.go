/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "as IS" BasIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipam

import (
	"reflect"

	"github.com/henderiw/iputil"
	"github.com/kform-dev/choreo/apis/condition"
	"github.com/kuidio/kuid/apis/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPEntrySpec defines the desired state of IPEntry
type IPEntrySpec struct {
	// Index defines the index for the IP Entry
	Index string `json:"index" yaml:"index" protobuf:"bytes,1,opt,name=index"`
	// IndexEntry identifies if the entry is originated from an IP Index
	IndexEntry bool `json:"indexEntry" yaml:"indexEntry" protobuf:"bytes,2,opt,name=indexEntry"`
	// PrefixType defines the prefixtype of IPEntry; for address and range claims this is not relevant
	// - network kind is used for physical, virtual nics on a device
	// - pool kind is used for allocating dedicated IP addresses
	// - aggregate kind is used for claiming an aggregate prefix; only used for networkInstance prefixes
	// +kubebuilder:validation:Enum=`network`;`aggregate`;`pool`;
	// +optional
	PrefixType *IPPrefixType `json:"prefixType,omitempty" yaml:"prefixType,omitempty" protobuf:"bytes,3,opt,name=prefixType"`
	// ClaimType defines the claimType of the IP Entry
	// +kubebuilder:validation:Enum=`staticAddress`;`staticPrefix`;`staticRange`;`dynamicPrefix`;`dynamicAddress`;
	ClaimType IPClaimType `json:"claimType,omitempty" yaml:"claimType,omitempty" protobuf:"bytes,4,opt,name=claimType"`
	// Prefix defines the prefix for the IP entry; which can be an expanded prefix from the prefix, range or address
	Prefix string `json:"prefix" yaml:"prefix" protobuf:"bytes,5,opt,name=prefix"`
	// DefaultGateway defines if the address acts as a default gateway
	// +optional
	DefaultGateway *bool `json:"defaultGateway,omitempty" yaml:"defaultGateway,omitempty" protobuf:"varint,6,opt,name=defaultGateway"`
	// AddressFamily defines the address family for the IP claim
	// +kubebuilder:validation:Enum=`ipv4`;`ipv6`
	// +kubebuilder:validation:Optional
	// +optional
	AddressFamily *iputil.AddressFamily `json:"addressFamily,omitempty" yaml:"addressFamily,omitempty" protobuf:"bytes,7,opt,name=addressFamily"`
	// UserDefinedLabels define the user defined labels
	common.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,8,opt,name=userDefinedLabels"`
}

// IPEntryStatus defines the observed state of IPEntry
type IPEntryStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// IPEntry is the Schema for the ipentry API
type IPEntry struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   IPEntrySpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status IPEntryStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// IPEntryList contains a list of IPEntries
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IPEntryList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []IPEntry `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	IPEntryKind     = reflect.TypeOf(IPEntry{}).Name()
	IPEntryListKind = reflect.TypeOf(IPEntryList{}).Name()
)