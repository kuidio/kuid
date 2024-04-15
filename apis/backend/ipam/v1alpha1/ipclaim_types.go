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
	"reflect"

	"github.com/henderiw/iputil"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPClaimSpec defines the desired state of IPClaim
type IPClaimSpec struct {
	// Kind defines the kind of prefix for the IP Claim
	// - network kind is used for physical, virtual nics on a device
	// - loopback kind is used for loopback interfaces within a device
	// - pool kind is used for pools for dhcp/radius/bng/upf/etc
	// - aggregate kind is used for claiming an aggregate prefix
	// +kubebuilder:validation:Enum=`network`;`loopback`;`pool`;`aggregate`
	// +kubebuilder:default=network
	Kind PrefixKind `json:"kind" yaml:"kind" protobuf:"bytes,1,opt,name=kind"`
	// NetworkInstance defines the networkInstance context for the IP claim
	// The NetworkInstance must exist within the IPClaim namespace to succeed
	// in claiming the ip
	NetworkInstance string `json:"networkInstance" yaml:"networkInstance" protobuf:"bytes,2,opt,name=networkInstance"`
	// AddressFamily defines the address family for the IP claim
	// +kubebuilder:validation:Enum=`ipv4`;`ipv6`
	// +kubebuilder:validation:Optional
	// +optional
	AddressFamily *iputil.AddressFamily `json:"addressFamily,omitempty" yaml:"addressFamily,omitempty" protobuf:"bytes,3,opt,name=addressFamily"`
	// Prefix defines the prefix for the IP claim
	// Used for specific prefix claim or used as a hint for a dynamic prefix claim in case of restart
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])/(([0-9])|([1-2][0-9])|(3[0-2]))|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))`
	// +kubebuilder:validation:Optional
	// +optional
	Prefix *string `json:"prefix,omitempty" yaml:"prefix,omitempty" protobuf:"bytes,4,opt,name=prefix"`
	// Gateway defines if the prefix/address is a gateway
	// +kubebuilder:validation:Optional
	// +optional
	Gateway *bool `json:"gateway,omitempty" yaml:"gateway,omitempty" protobuf:"varint,5,opt,name=gateway"`
	// PrefixLength defines the prefix length for the IP Claim
	// If not present we use assume /32 for ipv4 and /128 for ipv6
	// +kubebuilder:validation:Optional
	// +optional
	PrefixLength *uint32 `json:"prefixLength,omitempty" yaml:"prefixLength,omitempty" protobuf:"varint,6,opt,name=prefixLength"`
	// Index defines the index of the IP Claim, used to get a deterministic IP from a prefix
	// If not present we claim a random prefix from a prefix
	// +kubebuilder:validation:Optional
	// +optional
	Index *uint32 `json:"index,omitempty" yaml:"index,omitempty" protobuf:"varint,7,opt,name=index"`
	// CreatePrefix defines if this prefix must be created. Only used for non address prefixes
	// e.g. non /32 ipv4 and non /128 ipv6 prefixes
	// +kubebuilder:validation:Optional
	// +optional
	CreatePrefix *bool `json:"createPrefix,omitempty" yaml:"createPrefix,omitempty" protobuf:"varint,8,opt,name=createPrefix"`
	// ClaimLabels define the user defined labels and selector labels used
	// in resource claim
	commonv1alpha1.ClaimLabels `json:",inline" yaml:",inline" protobuf:"bytes,9,opt,name=claimLabels"`

	Owner *commonv1alpha1.OwnerReference `json:"owner,omitempty" yaml:"owner,omitempty" protobuf:"bytes,10,opt,name=owner"`
}

// IPClaimStatus defines the observed state of IPClaim
type IPClaimStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// Prefix defines the prefix, claimed through the IPAM backend
	// +kubebuilder:validation:Optional
	// +optional
	Prefix *string `json:"prefix,omitempty" yaml:"prefix,omitempty" protobuf:"bytes,2,opt,name=prefix"`
	// Gateway defines the gateway IP for the claimed prefix
	// Gateway is only relevant for prefix kind = network
	// +kubebuilder:validation:Optional
	// +optional
	Gateway *string `json:"gateway,omitempty" yaml:"gateway,omitempty" protobuf:"bytes,3,opt,name=gateway"`
	// ExpiryTime defines when the claim expires
	// +kubebuilder:validation:Optional
	// +optional
	ExpiryTime *string `json:"expiryTime,omitempty" yaml:"expiryTime,omitempty" protobuf:"bytes,4,opt,name=expiryTime"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IPClaim is the Schema for the ipclaim API
//
// +k8s:openapi-gen=true
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
	IPClaimKind = reflect.TypeOf(IPClaim{}).Name()
)
