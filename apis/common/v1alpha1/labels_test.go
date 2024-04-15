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
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetUserDefinedLabels(t *testing.T) {
	cases := map[string]struct {
		labels map[string]string
		want   map[string]string
	}{
		"Labels": {
			labels: map[string]string{"a": "b", "c": "d"},
			want:   map[string]string{"a": "b", "c": "d"},
		},
		"Nil": {
			labels: nil,
			want:   map[string]string{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			o := UserDefinedLabels{
				Labels: tc.labels,
			}

			got := o.GetUserDefinedLabels()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("-want, +got:\n%s", diff)
			}
		})
	}
}

func TestGetSelectorLabels(t *testing.T) {
	cases := map[string]struct {
		selector *metav1.LabelSelector
		want     map[string]string
	}{
		"Labels": {
			selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"a": "b", "c": "d"},
			},
			want: map[string]string{"a": "b", "c": "d"},
		},
		"Nil": {
			selector: nil,
			want:     map[string]string{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			o := ClaimLabels{
				Selector: tc.selector,
			}

			got := o.GetSelectorLabels()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("-want, +got:\n%s", diff)
			}
		})
	}
}

/*
func TestGetFullLabels(t *testing.T) {
	cases := map[string]struct {
		selector *metav1.LabelSelector
		labels   map[string]string
		want     map[string]string
	}{
		"Labels": {
			selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"a": "b", "c": "d"},
			},
			labels: map[string]string{"e": "f", "g": "h"},
			want:   map[string]string{"a": "b", "c": "d", "e": "f", "g": "h"},
		},
		"Nil": {
			selector: nil,
			want:     map[string]string{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			o := ClaimLabels{
				Selector:          tc.selector,
				UserDefinedLabels: UserDefinedLabels{tc.labels},
			}

			got := o.GetFullLabels()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("-want, +got:\n%s", diff)
			}
		})
	}
}
*/

func TestGetLabelSelector(t *testing.T) {
	cases := map[string]struct {
		selector *metav1.LabelSelector
		want     string
	}{
		"Labels": {
			selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"a": "b", "c": "d"},
			},
			want: "a=b,c=d",
		},
		"Nil": {
			selector: nil,
			want:     "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			o := ClaimLabels{
				Selector: tc.selector,
			}

			got, err := o.GetLabelSelector()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got.String()); diff != "" {
				t.Errorf("-want, +got:\n%s", diff)
			}
		})
	}
}
