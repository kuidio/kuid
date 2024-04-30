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

package all

import (
	_ "github.com/kuidio/kuid/pkg/reconcilers/ipamipclaim"
	_ "github.com/kuidio/kuid/pkg/reconcilers/ipamnetworkinstance"
	_ "github.com/kuidio/kuid/pkg/reconcilers/vlanclaim"
	_ "github.com/kuidio/kuid/pkg/reconcilers/vlanindex"
	_ "github.com/kuidio/kuid/pkg/reconcilers/vxlanclaim"
	_ "github.com/kuidio/kuid/pkg/reconcilers/vxlanindex"
	_ "github.com/kuidio/kuid/pkg/reconcilers/asclaim"
	_ "github.com/kuidio/kuid/pkg/reconcilers/asindex"
)
