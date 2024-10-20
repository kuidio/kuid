/*
Copyright 2023 The Nephio Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

// UserDefinedLabels define metadata to the resource.
type UserDefinedLabels struct {
	// Labels as user defined labels
	// +optional
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty" protobuf:"bytes,1,rep,name=labels"`
}

// GetUserDefinedLabels returns a map with a copy of the
// user defined labels
func (r *UserDefinedLabels) GetUserDefinedLabels() map[string]string {
	l := map[string]string{}
	if len(r.Labels) != 0 {
		for k, v := range r.Labels {
			l[k] = v
		}
	}
	return l
}

type ClaimLabels struct {
	UserDefinedLabels `json:",inline" yaml:",inline" protobuf:"bytes,1,opt,name=userDefinedLabels"`
	// Selector defines the selector criterias
	// +kubebuilder:validation:Optional
	Selector *metav1.LabelSelector `json:"selector,omitempty" yaml:"selector,omitempty" protobuf:"bytes,2,opt,name=selector"`
}

// GetUserDefinedLabels returns a map with a copy of the
// user defined labels
func (r *ClaimLabels) GetUserDefinedLabels() map[string]string {
	return r.UserDefinedLabels.GetUserDefinedLabels()
}

// GetSelectorLabels returns a map with a copy of the
// selector labels
func (r *ClaimLabels) GetSelectorLabels() map[string]string {
	l := map[string]string{}
	if r.Selector != nil {
		for k, v := range r.Selector.MatchLabels {
			l[k] = v
		}
	}
	return l
}

// GetLabelSelector returns a labels selector based
// on the label selector
func (r *ClaimLabels) GetLabelSelector() (labels.Selector, error) {
	l := r.GetSelectorLabels()
	fullselector := labels.NewSelector()
	for k, v := range l {
		req, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			return nil, err
		}
		fullselector = fullselector.Add(*req)
	}
	return fullselector, nil
}
