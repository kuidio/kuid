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

package options

import (
	"context"

	"github.com/dgraph-io/badger/v4"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	//"sigs.k8s.io/controller-runtime/pkg/client"
)

type StorageType int

const (
	StorageType_Memory StorageType = iota
	StorageType_File
	StorageType_KV
)

type Options struct {
	// Storage
	Prefix string
	Type   StorageType
	DB     *badger.DB
	// Target
	//Client client.Client
	// specific functions
	DryRunner      DryRunner
	BackendInvoker BackendInvoker
}

type DryRunner interface {
	DryRunCreate(ctx context.Context, key types.NamespacedName, obj runtime.Object, dryrun bool) (runtime.Object, error)
	DryRunUpdate(ctx context.Context, key types.NamespacedName, obj, old runtime.Object, dryrun bool) (runtime.Object, error)
	DryRunDelete(ctx context.Context, key types.NamespacedName, obj runtime.Object, dryrun bool) (runtime.Object, error)
}

type BackendInvoker interface {
	InvokeCreate(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error)
	InvokeUpdate(ctx context.Context, obj, old runtime.Object, recursion bool) (runtime.Object, runtime.Object, error)
	InvokeDelete(ctx context.Context, obj runtime.Object, recursion bool) (runtime.Object, error)
}
