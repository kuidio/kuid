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

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/henderiw/logger/log"
	_ "github.com/kuidio/kuid/apis/all"
	asbev1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	"github.com/kuidio/kuid/pkg/backend"
	"github.com/kuidio/kuid/pkg/config"
	"github.com/kuidio/kuid/pkg/generated/openapi"
	"github.com/kuidio/kuid/pkg/reconcilers"
	_ "github.com/kuidio/kuid/pkg/reconcilers/all"
	"github.com/kuidio/kuid/pkg/reconcilers/ctrlconfig"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/component-base/logs"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
// defaultEtcdPathPrefix = "/registry/backend.kuid.dev"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	l := log.NewLogger(&log.HandlerOptions{Name: "kuid-server-logger", AddSource: false})
	slog.SetDefault(l)
	ctx := log.IntoContext(context.Background(), l)
	log := log.FromContext(ctx)

	opts := zap.Options{
		TimeEncoder: zapcore.RFC3339NanoTimeEncoder,
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// if no async we dont have to start any, the reconcilers len will determine this
	recs := map[string]reconcilers.Reconciler{}
	ctrlCfg := &ctrlconfig.ControllerConfig{Backends: map[string]backend.Backend{}}
	kuidConfig, err := config.GetKuidConfig()
	if err != nil {
		log.Error("cannot get kuid config", "err", err)
		os.Exit(1)
	}

	registryOptions, err := config.GetRegistryOptions(ctx, kuidConfig.Storage)
	if err != nil {
		log.Error("cannot get kuid storage registry options", "err", err)
		os.Exit(1)
	}

	// apiserver is only relevant when not using etcd
	var apiserver *builder.Server
	if kuidConfig.Storage != config.StorageType_Etcd {
		apiserver = builder.APIServer.
			WithServerName("kuid-api-server").
			WithOpenAPIDefinitions("Config", "v1alpha1", openapi.GetOpenAPIDefinitions).
			WithoutEtcd()
	}

	backends := map[string]backend.Backend{}
	for _, kuidGroupConfig := range kuidConfig.Groups {
		group := kuidGroupConfig.Group
		log.Info("kuid group configured", "group", group, "enabled", kuidGroupConfig.Enabled, "sync", kuidGroupConfig.Sync)
		if !kuidGroupConfig.Enabled {
			continue
		}
		// check if group is registered
		groupConfig, ok := config.Groups[group]
		if !ok {
			log.Info("group configured in kuidconfig, but not registered", "group", group)
			continue
		}
		// create the storageProvider
		if kuidConfig.Storage != config.StorageType_Etcd {
			if groupConfig.BackendFn != nil {
				be := groupConfig.BackendFn()
				backends[group] = be
				for _, resource := range groupConfig.Resources {
					storageProvider := resource.StorageProviderFn(ctx, be, kuidGroupConfig.Sync, registryOptions)
					for _, resourceVersion := range resource.ResourceVersions {
						apiserver.WithResourceAndHandler(resourceVersion, storageProvider)
					}
				}
			} else {
				for _, resource := range groupConfig.Resources {
					storageProvider := resource.StorageProviderFn(ctx, nil, kuidGroupConfig.Sync, registryOptions)
					for _, resourceVersion := range resource.ResourceVersions {
						apiserver.WithResourceAndHandler(resourceVersion, storageProvider)
					}
				}
			}
		}

		// reconcilers get registered when async operations are configured
		if !kuidGroupConfig.Sync {
			if reconcilers, ok := reconcilers.ReconcilerGroups[group]; ok {
				for reconilerName, reconciler := range reconcilers {
					recs[reconilerName] = reconciler
				}
			}
			// add the backend -> for etcd this needs to be configmaps -> tbd how we handle this
		}
	}

	if kuidConfig.Storage != config.StorageType_Etcd {
		cmd, err := apiserver.Build(ctx)
		if err != nil {
			panic(err)
		}
		for _, kuidGroupConfig := range kuidConfig.Groups {
			group := kuidGroupConfig.Group
			if !kuidGroupConfig.Enabled {
				continue
			}
			// check if group is registered
			groupConfig, ok := config.Groups[group]
			if !ok {
				continue
			}

			if groupConfig.ApplyStorageToBackendFn != nil {
				if err := groupConfig.ApplyStorageToBackendFn(ctx, backends[group], apiserver); err != nil {
					log.Error("cannot apply storage to backend", "error", err.Error())
					os.Exit(1)
				}
			}
		}
		if err := cmd.Execute(); err != nil {
			panic(err)
		}
	}

	if len(recs) != 0 {
		// setup scheme for controllers
		runScheme := runtime.NewScheme()
		// add the core object to the scheme
		for _, api := range (runtime.SchemeBuilder{
			clientgoscheme.AddToScheme,
			//infrabev1alpha1.AddToScheme,
			asbev1alpha1.AddToScheme,
			//ipambev1alpha1.AddToScheme,
			//vlanbev1alpha1.AddToScheme,
			//vxlanbev1alpha1.AddToScheme,
			//extcommbev1alpha1.AddToScheme,
			//genidbev1alpha1.AddToScheme,
		}) {
			if err := api(runScheme); err != nil {
				log.Error("cannot add scheme", "err", err)
				os.Exit(1)
			}
		}

		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), manager.Options{
			Scheme: runScheme,
		})
		if err != nil {
			log.Error("cannot start manager", "err", err)
			os.Exit(1)
		}
		for name, reconciler := range recs {
			log.Info("reconciler", "name", name, "enabled", IsReconcilerEnabled(name))
			if IsReconcilerEnabled(name) {
				_, err := reconciler.SetupWithManager(ctx, mgr, ctrlCfg)
				if err != nil {
					log.Error("cannot add controllers to manager", "err", err.Error())
					os.Exit(1)
				}
			}
		}
		if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
			log.Error("unable to set up health check", "error", err.Error())
			os.Exit(1)
		}
		if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
			log.Error("unable to set up ready check", "error", err.Error())
			os.Exit(1)
		}

		log.Info("starting manager")
		if err := mgr.Start(ctx); err != nil {
			log.Error("problem running manager", "error", err.Error())
			os.Exit(1)
		}
	} else {
		for range ctx.Done() {
			log.Info("context cancelled...")
		}
	}

	/*
		asStorageProviders := as.NewStorageProviders(ctx, true, registryOptions)

		go func() {
			apiserver := builder.APIServer.
				WithServerName("kuid-api-server").
				WithOpenAPIDefinitions("Config", "v1alpha1", openapi.GetOpenAPIDefinitions).
				WithoutEtcd().
				WithResourceAndHandler(&as.ASIndex{}, asStorageProviders.GetIndexStorageProvider()).
				WithResourceAndHandler(&as.ASClaim{}, asStorageProviders.GetClaimStorageProvider()).
				WithResourceAndHandler(&as.ASEntry{}, asStorageProviders.GetEntryStorageProvider()).
				WithResourceAndHandler(&asbev1alpha1.ASIndex{}, asStorageProviders.GetIndexStorageProvider()).
				WithResourceAndHandler(&asbev1alpha1.ASClaim{}, asStorageProviders.GetClaimStorageProvider()).
				WithResourceAndHandler(&asbev1alpha1.ASEntry{}, asStorageProviders.GetEntryStorageProvider())

			cmd, err := apiserver.Build(ctx)
			if err != nil {
				panic(err)
			}
			if err := asStorageProviders.ApplyStorageToBackend(ctx, apiserver); err != nil {
				panic(err)
			}
			if err := cmd.Execute(); err != nil {
				panic(err)
			}
		}()
	*/

}

// IsReconcilerEnabled checks if an environment variable `ENABLE_<reconcilerName>` exists
// return "true" if the var is set and is not equal to "false".
func IsReconcilerEnabled(reconcilerName string) bool {
	if val, found := os.LookupEnv(fmt.Sprintf("ENABLE_%s", strings.ToUpper(reconcilerName))); found {
		if strings.ToLower(val) != "false" {
			return true
		}
	}
	return false
}
