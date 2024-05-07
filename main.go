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

//go:generate apiserver-runtime-gen
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
	asbev1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	vlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	vxlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
	"github.com/kuidio/kuid/apis/generated/clientset/versioned/scheme"
	kuidopenapi "github.com/kuidio/kuid/apis/generated/openapi"
	bebackend "github.com/kuidio/kuid/pkg/backend/backend"
	"github.com/kuidio/kuid/pkg/backend/ipam"
	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/pkg/kuidserver/asclaim"
	"github.com/kuidio/kuid/pkg/kuidserver/asentry"
	"github.com/kuidio/kuid/pkg/kuidserver/asindex"
	"github.com/kuidio/kuid/pkg/kuidserver/ipindex"
	"github.com/kuidio/kuid/pkg/kuidserver/ipclaim"
	"github.com/kuidio/kuid/pkg/kuidserver/ipentry"
	serverstore "github.com/kuidio/kuid/pkg/kuidserver/store"
	"github.com/kuidio/kuid/pkg/kuidserver/vlanclaim"
	"github.com/kuidio/kuid/pkg/kuidserver/vlanentry"
	"github.com/kuidio/kuid/pkg/kuidserver/vlanindex"
	"github.com/kuidio/kuid/pkg/kuidserver/vxlanclaim"
	"github.com/kuidio/kuid/pkg/kuidserver/vxlanentry"
	"github.com/kuidio/kuid/pkg/kuidserver/vxlanindex"
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
	if err := scheme.AddToScheme(runScheme); err != nil {
		log.Error("cannot initialize schema", "error", err)
		os.Exit(1)
	}
	// add the core object to the scheme
	for _, api := range (runtime.SchemeBuilder{
		clientgoscheme.AddToScheme,
		ipambev1alpha1.AddToScheme,
		vlanbev1alpha1.AddToScheme,
		vxlanbev1alpha1.AddToScheme,
		asbev1alpha1.AddToScheme,
	}) {
		if err := api(runScheme); err != nil {
			log.Error("cannot add scheme", "err", err)
			os.Exit(1)
		}
	}

	runScheme.AddFieldLabelConversionFunc(
		ipambev1alpha1.SchemeGroupVersion.WithKind(ipambev1alpha1.IPClaimKind),
		ipambev1alpha1.ConvertIPClaimFieldSelector,
	)
	runScheme.AddFieldLabelConversionFunc(
		ipambev1alpha1.SchemeGroupVersion.WithKind(ipambev1alpha1.IPEntryKind),
		ipambev1alpha1.ConvertIPEntryFieldSelector,
	)
	runScheme.AddFieldLabelConversionFunc(
		vlanbev1alpha1.SchemeGroupVersion.WithKind(vlanbev1alpha1.VLANClaimKind),
		vlanbev1alpha1.ConvertVLANClaimFieldSelector,
	)
	runScheme.AddFieldLabelConversionFunc(
		vlanbev1alpha1.SchemeGroupVersion.WithKind(vlanbev1alpha1.VLANEntryKind),
		vlanbev1alpha1.ConvertVLANEntryFieldSelector,
	)
	runScheme.AddFieldLabelConversionFunc(
		vxlanbev1alpha1.SchemeGroupVersion.WithKind(vxlanbev1alpha1.VXLANClaimKind),
		vxlanbev1alpha1.ConvertVXLANClaimFieldSelector,
	)
	runScheme.AddFieldLabelConversionFunc(
		vxlanbev1alpha1.SchemeGroupVersion.WithKind(vxlanbev1alpha1.VXLANEntryKind),
		vxlanbev1alpha1.ConvertVXLANEntryFieldSelector,
	)
	runScheme.AddFieldLabelConversionFunc(
		asbev1alpha1.SchemeGroupVersion.WithKind(asbev1alpha1.ASClaimKind),
		asbev1alpha1.ConvertASClaimFieldSelector,
	)
	runScheme.AddFieldLabelConversionFunc(
		asbev1alpha1.SchemeGroupVersion.WithKind(asbev1alpha1.ASEntryKind),
		asbev1alpha1.ConvertASEntryFieldSelector,
	)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), manager.Options{
		Scheme: runScheme,
	})
	if err != nil {
		log.Error("cannot start manager", "err", err)
		os.Exit(1)
	}


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

	ctrlCfg := &ctrlconfig.ControllerConfig{
		IPAMBackend:  ipbe,
		VLANBackend:  vlanbe,
		VXLANBackend: vxlanbe,
		ASBackend:    asbe,
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

	go func() {
		if err := builder.APIServer.
			WithServerName("kuid-server").
			WithEtcdPath(defaultEtcdPathPrefix).
			WithOpenAPIDefinitions("Kuid", "v1alpha1", kuidopenapi.GetOpenAPIDefinitions).
			WithResourceAndHandler(ctx, &ipambev1alpha1.IPClaim{}, ipclaim.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, ipbe)).
			WithResourceAndHandler(ctx, &ipambev1alpha1.IPEntry{}, ipentry.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, ipbe)).
			WithResourceAndHandler(ctx, &ipambev1alpha1.IPIndex{}, ipindex.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, ipbe)).
			WithResourceAndHandler(ctx, &vlanbev1alpha1.VLANClaim{}, vlanclaim.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, vlanbe)).
			WithResourceAndHandler(ctx, &vlanbev1alpha1.VLANEntry{}, vlanentry.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, vlanbe)).
			WithResourceAndHandler(ctx, &vlanbev1alpha1.VLANIndex{}, vlanindex.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, vlanbe)).
			WithResourceAndHandler(ctx, &vxlanbev1alpha1.VXLANClaim{}, vxlanclaim.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, vxlanbe)).
			WithResourceAndHandler(ctx, &vxlanbev1alpha1.VXLANEntry{}, vxlanentry.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, vxlanbe)).
			WithResourceAndHandler(ctx, &vxlanbev1alpha1.VXLANIndex{}, vxlanindex.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, vxlanbe)).
			WithResourceAndHandler(ctx, &asbev1alpha1.ASClaim{}, asclaim.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, asbe)).
			WithResourceAndHandler(ctx, &asbev1alpha1.ASEntry{}, asentry.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, asbe)).
			WithResourceAndHandler(ctx, &asbev1alpha1.ASIndex{}, asindex.NewProvider(ctx, mgr.GetClient(), &serverstore.Config{
				Prefix: configDir,
				Type:   serverstore.StorageType_KV,
				DB:     db,
			}, asbe)).
			WithoutEtcd().
			Execute(ctx); err != nil {
			log.Info("cannot start config-server")
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
