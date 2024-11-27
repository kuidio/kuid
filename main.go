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
	"crypto/tls"
	"log/slog"
	"os"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/henderiw/logger/log"
	_ "github.com/kuidio/kuid/apis/all"
	"github.com/kuidio/kuid/pkg/backend"
	kuidconfig "github.com/kuidio/kuid/pkg/config"
	"github.com/kuidio/kuid/pkg/generated/openapi"
	"github.com/kuidio/kuid/pkg/reconcilers"
	_ "github.com/kuidio/kuid/pkg/reconcilers/all"
	"github.com/kuidio/kuid/pkg/reconcilers/ctrlconfig"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/component-base/logs"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

const (
// defaultEtcdPathPrefix = "/registry/backend.kuid.dev"
)

type ReconcilerGroup struct {
	addToSchema func(*runtime.Scheme) error
	reconcilers []reconcilers.Reconciler
}

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
	groupReconcilers := map[string]*ReconcilerGroup{}
	ctrlCfg := &ctrlconfig.ControllerConfig{Backends: map[string]backend.Backend{}}
	kuidConfig, err := kuidconfig.GetKuidConfig()
	if err != nil {
		log.Error("cannot get kuid config", "err", err)
		os.Exit(1)
	}

	registryOptions, err := kuidconfig.GetRegistryOptions(ctx, kuidConfig.Storage)
	if err != nil {
		log.Error("cannot get kuid storage registry options", "err", err)
		os.Exit(1)
	}

	// apiserver is only relevant when not using etcd
	var apiserver *builder.Server
	if kuidConfig.Storage != kuidconfig.StorageType_Etcd {
		apiserver = builder.NewAPIServer().
			WithServerName("kuid-api-server").
			WithOpenAPIDefinitions("Config", "v1alpha1", openapi.GetOpenAPIDefinitions).
			WithoutEtcd()
	}

	for _, kuidGroupConfig := range kuidConfig.Groups {
		group := kuidGroupConfig.Group
		log.Info("kuid group configured", "group", group, "enabled", kuidGroupConfig.Enabled, "sync", kuidGroupConfig.Sync)
		if !kuidGroupConfig.Enabled {
			continue
		}
		// check if group is registered
		groupConfig, ok := kuidconfig.Groups[group]
		if !ok {
			log.Info("group configured in kuidconfig, but not registered", "group", group)
			continue
		}
		// create the storageProvider
		if kuidConfig.Storage != kuidconfig.StorageType_Etcd {
			if groupConfig.BackendFn != nil {
				be := groupConfig.BackendFn()
				ctrlCfg.Backends[group] = be
				for _, resource := range groupConfig.Resources {
					storageProvider := resource.StorageProviderFn(ctx, resource.Internal, be, kuidGroupConfig.Sync, registryOptions)
					for _, resourceVersion := range resource.ResourceVersions {
						apiserver.WithResourceAndHandler(resourceVersion, storageProvider)
					}
				}
			} else {
				for _, resource := range groupConfig.Resources {
					storageProvider := resource.StorageProviderFn(ctx, resource.Internal, nil, kuidGroupConfig.Sync, registryOptions)
					for _, resourceVersion := range resource.ResourceVersions {
						apiserver.WithResourceAndHandler(resourceVersion, storageProvider)
					}
				}
			}
		}

		// TODO: determine if the reconsilers behave different for sync compared to async
		if recs, ok := reconcilers.ReconcilerGroups[group]; ok {
			log.Info("reconciler group", "group", group)
			groupReconcilers[group] = &ReconcilerGroup{
				addToSchema: groupConfig.AddToScheme,
				reconcilers: make([]reconcilers.Reconciler, 0, len(recs)),
			}
			for name, reconciler := range recs {
				log.Info("reconciler", "group", group, "name", name)
				groupReconcilers[group].reconcilers = append(groupReconcilers[group].reconcilers, reconciler)
			}
		}
		// add the backend -> for etcd this needs to be configmaps -> tbd how we handle this
	}

	if kuidConfig.Storage != kuidconfig.StorageType_Etcd {
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
			groupConfig, ok := kuidconfig.Groups[group]
			if !ok {
				continue
			}

			if groupConfig.ApplyStorageToBackendFn != nil {
				if err := groupConfig.ApplyStorageToBackendFn(ctx, ctrlCfg.Backends[group], apiserver); err != nil {
					log.Error("cannot apply storage to backend", "error", err.Error())
					os.Exit(1)
				}
			}
		}
		go func() {
			if err := cmd.Execute(); err != nil {
				panic(err)
			}
		}()
	}

	log.Info("groupReconcilers", "total", len(groupReconcilers))
	if len(groupReconcilers) != 0 {
		// setup scheme for controllers
		runScheme := runtime.NewScheme()
		if err := clientgoscheme.AddToScheme(runScheme); err != nil {
			log.Error("cannot add scheme", "err", err)
			os.Exit(1)
		}
		// add all schemas for the reconcilers to
		for _, reconcilerGroup := range groupReconcilers {
			if err := reconcilerGroup.addToSchema(runScheme); err != nil {
				log.Error("cannot add scheme", "err", err)
				os.Exit(1)
			}
		}

		var tlsOpts []func(*tls.Config)
		metricsServerOptions := metricsserver.Options{
			BindAddress:   ":8443",
			SecureServing: true,
			// FilterProvider is used to protect the metrics endpoint with authn/authz.
			// These configurations ensure that only authorized users and service accounts
			// can access the metrics endpoint. The RBAC are configured in 'config/rbac/kustomization.yaml'. More info:
			// https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/metrics/filters#WithAuthenticat
			FilterProvider: filters.WithAuthenticationAndAuthorization,
			// If CertDir, CertName, and KeyName are not specified, controller-runtime will automatically
			// generate self-signed certificates for the metrics server. While convenient for development and testing,
			// this setup is not recommended for production.
			TLSOpts: tlsOpts,
		}

		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), manager.Options{
			Scheme:  runScheme,
			Metrics: metricsServerOptions,
			Controller: config.Controller{
				MaxConcurrentReconciles: 16,
			},
		})
		if err != nil {
			log.Error("cannot start manager", "err", err)
			os.Exit(1)
		}
		for _, reconcilerGroup := range groupReconcilers {
			for _, reconciler := range reconcilerGroup.reconcilers {
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
}
