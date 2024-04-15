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

// IPEntrySpec defines the desired state of IPEntry
type IPEntrySpec struct {
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
	// AddressFamily defines the address family for the IP Entry
	// +kubebuilder:validation:Enum=`ipv4`;`ipv6`
	AddressFamily iputil.AddressFamily `json:"addressFamily" yaml:"addressFamily" protobuf:"bytes,3,opt,name=addressFamily"`
	// Prefix defines the prefix for the IP Entry; can be address or prefix
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])/(([0-9])|([1-2][0-9])|(3[0-2]))|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))`
	Prefix string `json:"prefix" yaml:"prefix" protobuf:"bytes,4,opt,name=prefix"`
	// ParentPrefix defines the parent prefix for the IP Entry
	// Used for specific prefix claim or used as a hint for a dynamic prefix claim in case of restart
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])/(([0-9])|([1-2][0-9])|(3[0-2]))|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))`
	// +kubebuilder:validation:Optional
	Subnet string `json:"subnet" yaml:"subnet" protobuf:"bytes,5,opt,name=subnet"`
	// Index defines the index of the IP Entry, used to get a deterministic IP from a prefix
	// If not present we claim a random prefix from a prefix
	// +kubebuilder:validation:Optional
	Index *uint32 `json:"index,omitempty" yaml:"index,omitempty" protobuf:"varint,6,opt,name=index"`
	// Gateway defines if the prefix/address is a gateway
	// +kubebuilder:validation:Optional
	Gateway *bool `json:"gateway,omitempty" yaml:"gateway,omitempty" protobuf:"varint,7,opt,name=gateway"`
	// IPClaim defines the name of the ip claim that is the origin of this ip entry
	IPClaim string `json:"ipClaim" yaml:"ipClaim" protobuf:"bytes,8,opt,name=ipClaim"`
	// UserDefinedLabels define the user defined labels
	commonv1alpha1.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,9,opt,name=userDefinedLabels"`

	Owner *commonv1alpha1.OwnerReference `json:"owner,omitempty" yaml:"owner,omitempty" protobuf:"bytes,10,opt,name=owner"`
}

// IPEntryStatus defines the observed state of IPEntry
type IPEntryStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IPEntry is the Schema for the ipentry API
//
// +k8s:openapi-gen=true
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
	IPEntryKind = reflect.TypeOf(IPEntry{}).Name()
)
