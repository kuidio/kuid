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

package resource

import (
	"k8s.io/apimachinery/pkg/api/errors"
)

const (
	errUpdateObject = "cannot update k8s resource"
)

type ErrorIs func(err error) bool

func Ignore(is ErrorIs, err error) error {
	if is(err) {
		return nil
	}
	return err
}

// IgnoreAny ignores errors that satisfy any of the supplied ErrorIs functions
// by returning nil. Errors that do not satisfy any of the supplied functions
// are returned unmodified.
func IgnoreAny(err error, is ...ErrorIs) error {
	for _, f := range is {
		if f(err) {
			return nil
		}
	}
	return err
}

// IgnoreNotFound returns the supplied error, or nil if the error indicates a
// Kubernetes resource was not found.
func IgnoreNotFound(err error) error {
	return Ignore(errors.IsNotFound, err)
}

// IgnoreAlreadyExist returns the supplied error, or nil if the error indicates a
// Kubernetes resource that  already exists.
func IgnoreAlreadyExist(err error) error {
	return Ignore(errors.IsAlreadyExists, err)
}