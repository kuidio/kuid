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

package vlan

import (
	"github.com/henderiw/idxtable/pkg/vlantable"
)

type CacheContext struct {
	table vlantable.VLANTable
}

func NewCacheContext() (*CacheContext, error) {
	table, err := vlantable.New()
	if err != nil {
		return nil, err
	}
	return &CacheContext{
		table: table,
	}, nil

}

func (r *CacheContext) Size() int {
	return r.table.Count()
}
