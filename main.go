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
	"github.com/henderiw/apiserver-store/pkg/db/badgerdb"
	"github.com/henderiw/logger/log"
	"github.com/kuidio/kuid/apis/backend/as"
	asbev1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	"github.com/kuidio/kuid/pkg/generated/openapi"
	"github.com/kuidio/kuid/pkg/reconcilers"
	_ "github.com/kuidio/kuid/pkg/reconcilers/all"
	"github.com/kuidio/kuid/pkg/reconcilers/ctrlconfig"
	"github.com/kuidio/kuid/pkg/registry/options"
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
	defaultEtcdPathPrefix = "/registry/backend.kuid.dev"
)

var (
	configDir = "/config"
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

	// setup controllers
	runScheme := runtime.NewScheme()
	// add the core object to the scheme
	for _, api := range (runtime.SchemeBuilder{
		clientgoscheme.AddToScheme,
		//infrabev1alpha1.AddToScheme,
		//ipambev1alpha1.AddToScheme,
		//vlanbev1alpha1.AddToScheme,
		//vxlanbev1alpha1.AddToScheme,
		asbev1alpha1.AddToScheme,
		//extcommbev1alpha1.AddToScheme,
		//genidbev1alpha1.AddToScheme,
	}) {
		if err := api(runScheme); err != nil {
			log.Error("cannot add scheme", "err", err)
			os.Exit(1)
		}
	}

	/*
		runScheme.AddFieldLabelConversionFunc(
			ipambev1alpha1.SchemeGroupVersion.WithKind(ipambev1alpha1.IPClaimKind),
			ipambev1alpha1.IPClaimConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			ipambev1alpha1.SchemeGroupVersion.WithKind(ipambev1alpha1.IPEntryKind),
			ipambev1alpha1.IPEntryConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			vlanbev1alpha1.SchemeGroupVersion.WithKind(vlanbev1alpha1.VLANClaimKind),
			vlanbev1alpha1.VLANClaimConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			vlanbev1alpha1.SchemeGroupVersion.WithKind(vlanbev1alpha1.VLANEntryKind),
			vlanbev1alpha1.VLANEntryConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			vxlanbev1alpha1.SchemeGroupVersion.WithKind(vxlanbev1alpha1.VXLANClaimKind),
			vxlanbev1alpha1.VXLANClaimConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			vxlanbev1alpha1.SchemeGroupVersion.WithKind(vxlanbev1alpha1.VXLANEntryKind),
			vxlanbev1alpha1.VXLANEntryConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			asbev1alpha1.SchemeGroupVersion.WithKind(asbev1alpha1.ASClaimKind),
			asbev1alpha1.ASClaimConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			asbev1alpha1.SchemeGroupVersion.WithKind(asbev1alpha1.ASEntryKind),
			asbev1alpha1.ASEntryConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			extcommbev1alpha1.SchemeGroupVersion.WithKind(extcommbev1alpha1.EXTCOMMClaimKind),
			extcommbev1alpha1.EXTCOMMClaimConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			extcommbev1alpha1.SchemeGroupVersion.WithKind(extcommbev1alpha1.EXTCOMMEntryKind),
			extcommbev1alpha1.EXTCOMMEntryConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			genidbev1alpha1.SchemeGroupVersion.WithKind(genidbev1alpha1.GENIDClaimKind),
			genidbev1alpha1.GENIDClaimConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			genidbev1alpha1.SchemeGroupVersion.WithKind(genidbev1alpha1.GENIDEntryKind),
			genidbev1alpha1.GENIDEntryConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			infrabev1alpha1.SchemeGroupVersion.WithKind(infrabev1alpha1.NodeKind),
			infrabev1alpha1.NodeConvertFieldSelector,
		)
		runScheme.AddFieldLabelConversionFunc(
			infrabev1alpha1.SchemeGroupVersion.WithKind(infrabev1alpha1.LinkKind),
			infrabev1alpha1.LinkConvertFieldSelector,
		)
	*/

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), manager.Options{
		Scheme: runScheme,
	})
	if err != nil {
		log.Error("cannot start manager", "err", err)
		os.Exit(1)
	}

	/*
		ipbe := ipam.New(mgr.GetClient())
		vlanbe := bebackend.New(
			mgr.GetClient(),
			func() backend.IndexObject { return &vlanbev1alpha1.VLANIndex{} },
			func() backend.ObjectList { return &vlanbev1alpha1.VLANEntryList{} },
			func() backend.ObjectList { return &vlanbev1alpha1.VLANClaimList{} },
			vlanbev1alpha1.GetVLANEntry,
			vlanbev1alpha1.SchemeGroupVersion.WithKind(vlanbev1alpha1.VLANIndexKind),
			vlanbev1alpha1.SchemeGroupVersion.WithKind(vlanbev1alpha1.VLANClaimKind),
		)
		vxlanbe := bebackend.New(
			mgr.GetClient(),
			func() backend.IndexObject { return &vxlanbev1alpha1.VXLANIndex{} },
			func() backend.ObjectList { return &vxlanbev1alpha1.VXLANEntryList{} },
			func() backend.ObjectList { return &vxlanbev1alpha1.VXLANClaimList{} },
			vxlanbev1alpha1.GetVXLANEntry,
			vxlanbev1alpha1.SchemeGroupVersion.WithKind(vxlanbev1alpha1.VXLANIndexKind),
			vxlanbev1alpha1.SchemeGroupVersion.WithKind(vxlanbev1alpha1.VXLANClaimKind),
		)
		asbe := bebackend.New(
			mgr.GetClient(),
			func() backend.IndexObject { return &asbev1alpha1.ASIndex{} },
			func() backend.ObjectList { return &asbev1alpha1.ASEntryList{} },
			func() backend.ObjectList { return &asbev1alpha1.ASClaimList{} },
			asbev1alpha1.GetASEntry,
			asbev1alpha1.SchemeGroupVersion.WithKind(asbev1alpha1.ASIndexKind),
			asbev1alpha1.SchemeGroupVersion.WithKind(asbev1alpha1.ASClaimKind),
		)
		extcommbe := bebackend.New(
			mgr.GetClient(),
			func() backend.IndexObject { return &extcommbev1alpha1.EXTCOMMIndex{} },
			func() backend.ObjectList { return &extcommbev1alpha1.EXTCOMMEntryList{} },
			func() backend.ObjectList { return &extcommbev1alpha1.EXTCOMMClaimList{} },
			extcommbev1alpha1.GetEXTCOMMEntry,
			extcommbev1alpha1.SchemeGroupVersion.WithKind(extcommbev1alpha1.EXTCOMMIndexKind),
			extcommbev1alpha1.SchemeGroupVersion.WithKind(extcommbev1alpha1.EXTCOMMClaimKind),
		)
		genidbe := bebackend.New(
			mgr.GetClient(),
			func() backend.IndexObject { return &genidbev1alpha1.GENIDIndex{} },
			func() backend.ObjectList { return &genidbev1alpha1.GENIDEntryList{} },
			func() backend.ObjectList { return &genidbev1alpha1.GENIDClaimList{} },
			genidbev1alpha1.GetGENIDEntry,
			genidbev1alpha1.SchemeGroupVersion.WithKind(genidbev1alpha1.GENIDIndexKind),
			genidbev1alpha1.SchemeGroupVersion.WithKind(genidbev1alpha1.GENIDClaimKind),
		)
	*/

	ctrlCfg := &ctrlconfig.ControllerConfig{
		//IPAMBackend:    ipbe,
		//VLANBackend:    vlanbe,
		//VXLANBackend:   vxlanbe,
		//ASBackend:      asbe,
		//EXTCOMMBackend: extcommbe,
		//GENIDBackend:   genidbe,
	}
	for name, reconciler := range reconcilers.Reconcilers {
		log.Info("reconciler", "name", name, "enabled", IsReconcilerEnabled(name))
		if IsReconcilerEnabled(name) {
			_, err := reconciler.SetupWithManager(ctx, mgr, ctrlCfg)
			if err != nil {
				log.Error("cannot add controllers to manager", "err", err.Error())
				os.Exit(1)
			}
		}
	}

	db, err := badgerdb.OpenDB(ctx, configDir)
	if err != nil {
		log.Error("cannot open db", "err", err.Error())
		os.Exit(1)
	}

	registryOptions := &options.Options{
		Prefix: configDir,
		Type:   options.StorageType_KV,
		DB:     db,
	}

	asStorageProviders := as.NewStorageProviders(ctx, true, registryOptions)

	go func() {
		apiserver := builder.APIServer.
			WithServerName("kuid-api-server").
			WithOpenAPIDefinitions("Config", "v1alpha1", openapi.GetOpenAPIDefinitions).
			WithResourceAndHandler(&as.ASIndex{}, asStorageProviders.GetIndexStorageProvider()).
			WithResourceAndHandler(&as.ASClaim{}, asStorageProviders.GetClaimStorageProvider()).
			WithResourceAndHandler(&as.ASEntry{}, asStorageProviders.GetEntryStorageProvider()).
			WithResourceAndHandler(&asbev1alpha1.ASIndex{}, asStorageProviders.GetIndexStorageProvider()).
			WithResourceAndHandler(&asbev1alpha1.ASClaim{}, asStorageProviders.GetClaimStorageProvider()).
			WithResourceAndHandler(&asbev1alpha1.ASEntry{}, asStorageProviders.GetEntryStorageProvider()).
			WithoutEtcd()

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
