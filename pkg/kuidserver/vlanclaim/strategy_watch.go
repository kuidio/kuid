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

package vlanclaim

import (
	"context"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/watch"
)

func (r *strategy) BeginWatch(ctx context.Context) error { return nil }

func (r *strategy) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	ctx, cancel := context.WithCancel(ctx)

	w := &watcher{
		cancel:         cancel,
		resultChan:     make(chan watch.Event),
		watcherManager: r.watcherManager,
	}

	go w.listAndWatch(ctx, r, options)

	return w, nil
}
