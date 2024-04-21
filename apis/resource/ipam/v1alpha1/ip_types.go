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

	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	resourcev1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	conditionv1alpha1 "github.com/kuidio/kuid/apis/condition/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPSpec defines the desired state of the IP
type IPSpec struct {
	// NetworkInstance defines the networkInstance context of the IP.
	// A NetworkInstance is a dedicated routing table instance
	NetworkInstance string `json:"networkInstance" yaml:"networkInstance" protobuf:"bytes,2,opt,name=networkInstance"`
	// PrefixType defines the type of prefix.
	// - network kind is used for physical, virtual nics on a device (allocates .net and .broadcast address, return address/prefixlength)
	// - pool kind is used for pool claims
	// +kubebuilder:validation:Enum=`network`;`pool`;
	// +optional
	PrefixType *ipambev1alpha1.IPPrefixType `json:"prefixType,omitempty" yaml:"prefixType,omitempty" protobuf:"bytes,2,opt,name=prefixType"`
	// Prefix defines the ip prefix in cidr notation.
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])/(([0-9])|([1-2][0-9])|(3[0-2]))|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))`
	Prefix *string `json:"prefix,omitempty" yaml:"prefix,omitempty" protobuf:"bytes,3,opt,name=prefix"`
	// Address defines the ip address in cidr notation
	// +kubebuilder:validation:Pattern=`(([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])/(([0-9])|([1-2][0-9])|(3[0-2]))|((:|[0-9a-fA-F]{0,4}):)([0-9a-fA-F]{0,4}:){0,5}((([0-9a-fA-F]{0,4}:)?(:|[0-9a-fA-F]{0,4}))|(((25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])))(/(([0-9])|([0-9]{2})|(1[0-1][0-9])|(12[0-8])))`
	Address *string `json:"address,omitempty" yaml:"address,omitempty" protobuf:"bytes,3,opt,name=address"`
	// Range defines the ip range. e.g. 10.0.0.10-10.0.0.100 or 2000::100-2000::1000
	Range *string `json:"range,omitempty" yaml:"range,omitempty" protobuf:"bytes,2,opt,name=range"`
	// DefaultGateway defines if the address acts as a default gateway
	// +kubebuilder:validation:Optional
	// +optional
	DefaultGateway *bool `json:"defaultGateway,omitempty" yaml:"defaultGateway,omitempty" protobuf:"varint,4,opt,name=defaultGateway"`
	// +optional
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	resourcev1alpha1.UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,5,opt,name=userDefinedLabels"`
}

// IPStatus defines the observed state of IP
type IPStatus struct {
	// ConditionedStatus provides the status of the IP using conditions
	// - a ready condition indicates the overall status of the resource
	conditionv1alpha1.ConditionedStatus `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// Prefix defines the prefix, claimed through the IPAM backend
	// +kubebuilder:validation:Optional
	Prefix *string `json:"prefix,omitempty" yaml:"prefix,omitempty" protobuf:"bytes,2,opt,name=prefix"`
	// Address defines the address, claimed through the IPAM backend
	// +kubebuilder:validation:Optional
	Address *string `json:"address,omitempty" yaml:"address,omitempty" protobuf:"bytes,3,opt,name=address"`
	// Range defines the range, claimed through the IPAM backend
	// +kubebuilder:validation:Optional
	Range *string `json:"range,omitempty" yaml:"range,omitempty" protobuf:"bytes,4,opt,name=range"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="NETWORK-INSTANCE",type="string",JSONPath=".spec.networkInstance"
// +kubebuilder:printcolumn:name="PREFIX-TYPE",type="string",JSONPath=".spec.prefixType"
// +kubebuilder:printcolumn:name="PREFIX-REQ",type="string",JSONPath=".spec.prefix"
// +kubebuilder:printcolumn:name="PREFIX-CLAIM",type="string",JSONPath=".status.prefix"
// +kubebuilder:printcolumn:name="RANGE-REQ",type="string",JSONPath=".spec.range"
// +kubebuilder:printcolumn:name="RANGE-CLAIM",type="string",JSONPath=".status.range"
// +kubebuilder:printcolumn:name="ADDRESS-REQ",type="string",JSONPath=".spec.address"
// +kubebuilder:printcolumn:name="ADDRESS-CLAIM",type="string",JSONPath=".status.address"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:categories={kuid,ipam}

// IP is the Schema for the ip API
type IP struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   IPSpec   `json:"spec,omitempty" yaml:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status IPStatus `json:"status,omitempty" yaml:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

//+kubebuilder:object:root=true

// IPList contains a list of IP(s)
type IPList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []IP `json:"items" yaml:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	IPKind = reflect.TypeOf(IP{}).Name()
)
