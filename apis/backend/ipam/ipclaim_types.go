/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "VLAN IS" BVLANIS,
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

// IPClaimSpec defines the desired state of IPClaim
type IPClaimSpec struct {
	// Index defines the index for the IP Entry
	Index string `json:"index" yaml:"index" protobuf:"bytes,1,opt,name=index"`
	// PrefixType network indicates a special type of prefix for which network and broadcast addresses
	// are claimed in the ipam, used for physical, virtual nics devices
	// If no prefixes type is defined the internally this is defaulted to other
	// +kubebuilder:validation:Enum=`network`;`other`;
	// +optional
	PrefixType *IPPrefixType `json:"prefixType,omitempty" yaml:"prefixType,omitempty" protobuf:"bytes,2,opt,name=prefixType"`
	// Prefix defines the prefix for the IP claim
	// +optional
	Prefix *string `json:"prefix,omitempty" yaml:"prefix,omitempty" protobuf:"bytes,3,opt,name=prefix"`
	// Address defines the address for the IP claim
	// +optional
	Address *string `json:"address,omitempty" yaml:"address,omitempty" protobuf:"bytes,4,opt,name=address"`
	// Range defines the range for the IP claim
	// +optional
	Range *string `json:"range,omitempty" yaml:"range,omitempty" protobuf:"bytes,5,opt,name=range"`
	// DefaultGateway defines if the address acts as a default gateway
	// +optional
	DefaultGateway *bool `json:"defaultGateway,omitempty" yaml:"defaultGateway,omitempty" protobuf:"varint,6,opt,name=defaultGateway"`
	// CreatePrefix defines if this prefix must be created. Only used for dynamic prefixes
	// e.g. non /32 ipv4 and non /128 ipv6 prefixes
	// +optional
	CreatePrefix *bool `json:"createPrefix,omitempty" yaml:"createPrefix,omitempty" protobuf:"varint,7,opt,name=createPrefix"`
	// PrefixLength defines the prefix length for the IP Claim, Must be set when CreatePrefic is set
	// If not present we use assume /32 for ipv4 and /128 for ipv6
	// +optional
	PrefixLength *uint32 `json:"prefixLength,omitempty" yaml:"prefixLength,omitempty" protobuf:"varint,8,opt,name=prefixLength"`
	// AddressFamily defines the address family for the IP claim
	// +kubebuilder:validation:Enum=`ipv4`;`ipv6`
	// +kubebuilder:validation:Optional
	// +optional
	AddressFamily *iputil.AddressFamily `json:"addressFamily,omitempty" yaml:"addressFamily,omitempty" protobuf:"bytes,9,opt,name=addressFamily"`
	// Index defines the index of the IP Claim, used to get a deterministic IP from a prefix
	// If not present we claim a random prefix from a prefix
	// +kubebuilder:validation:Optional
	// +optional
	Idx *uint32 `json:"idx,omitempty" yaml:"idx,omitempty" protobuf:"varint,10,opt,name=idx"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	common.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,11,opt,name=claimLabels"`
}

// IPClaimStatus defines the observed state of IPClaim
type IPClaimStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	condition.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// Range defines the range, claimed through the IPAM backend
	// +optional
	Range *string `json:"range,omitempty" yaml:"range,omitempty" protobuf:"bytes,2,opt,name=range"`
	// Address defines the address, claimed through the IPAM backend
	// +optional
	Address *string `json:"address,omitempty" yaml:"address,omitempty" protobuf:"bytes,3,opt,name=address"`
	// Prefix defines the prefix, claimed through the IPAM backend
	// +optional
	Prefix *string `json:"prefix,omitempty" yaml:"prefix,omitempty" protobuf:"bytes,4,opt,name=prefix"`
	// DefaultGateway defines the default gateway IP for the claimed prefix
	// DefaultGateway is only relevant for prefix kind = network
	// +optional
	DefaultGateway *string `json:"defaultGateway,omitempty" yaml:"defaultGateway,omitempty" protobuf:"bytes,5,opt,name=defaultGateway"`
	// ExpiryTime defines when the claim expires
	// +kubebuilder:validation:Optional
	// +optional
	ExpiryTime *string `json:"expiryTime,omitempty" yaml:"expiryTime,omitempty" protobuf:"bytes,6,opt,name=expiryTime"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// IPClaim is the Schema for the ipclaim API
type IPClaim struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   IPClaimSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status IPClaimStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// IPClaimList contains a list of IPClaims
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IPClaimList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []IPClaim `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	IPClaimKind     = reflect.TypeOf(IPClaim{}).Name()
	IPClaimListKind = reflect.TypeOf(IPClaimList{}).Name()
)
