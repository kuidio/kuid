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
	"github.com/kuidio/kuid/apis/backend"
	asbev1alpha1 "github.com/kuidio/kuid/apis/backend/as/v1alpha1"
	extcommbev1alpha1 "github.com/kuidio/kuid/apis/backend/extcomm/v1alpha1"
	genidbev1alpha1 "github.com/kuidio/kuid/apis/backend/genid/v1alpha1"
	infrabev1alpha1 "github.com/kuidio/kuid/apis/backend/infra/v1alpha1"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	vlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vlan/v1alpha1"
	vxlanbev1alpha1 "github.com/kuidio/kuid/apis/backend/vxlan/v1alpha1"
	"github.com/kuidio/kuid/apis/generated/clientset/versioned/scheme"
	kuidopenapi "github.com/kuidio/kuid/apis/generated/openapi"
	bebackend "github.com/kuidio/kuid/pkg/backend/backend"
	"github.com/kuidio/kuid/pkg/backend/ipam"
	"github.com/kuidio/kuid/pkg/kuidserver/claimserver"
	"github.com/kuidio/kuid/pkg/kuidserver/entryserver"
	"github.com/kuidio/kuid/pkg/kuidserver/genericserver"
	"github.com/kuidio/kuid/pkg/kuidserver/indexserver"
	serverstore "github.com/kuidio/kuid/pkg/kuidserver/store"
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
		infrabev1alpha1.AddToScheme,
		ipambev1alpha1.AddToScheme,
		vlanbev1alpha1.AddToScheme,
		vxlanbev1alpha1.AddToScheme,
		asbev1alpha1.AddToScheme,
		extcommbev1alpha1.AddToScheme,
		genidbev1alpha1.AddToScheme,
	}) {
		if err := api(runScheme); err != nil {
			log.Error("cannot add scheme", "err", err)
			os.Exit(1)
		}
	}

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

	ctrlCfg := &ctrlconfig.ControllerConfig{
		IPAMBackend:    ipbe,
		VLANBackend:    vlanbe,
		VXLANBackend:   vxlanbe,
		ASBackend:      asbe,
		EXTCOMMBackend: extcommbe,
		GENIDBackend:   genidbe,
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
			WithResourceAndHandler(ctx, &ipambev1alpha1.IPClaim{}, claimserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&claimserver.ServerObjContext{
					TracerString:   "ipclaim-server",
					Obj:            &ipambev1alpha1.IPClaim{},
					ConversionFunc: ipambev1alpha1.IPClaimConvertFieldSelector,
					TableConverter: ipambev1alpha1.IPClaimTableConvertor,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				ipbe)).
			WithResourceAndHandler(ctx, &ipambev1alpha1.IPEntry{}, entryserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&entryserver.ServerObjContext{
					TracerString:   "ipentry-server",
					Obj:            &ipambev1alpha1.IPEntry{},
					ConversionFunc: ipambev1alpha1.IPEntryConvertFieldSelector,
					TableConverter: ipambev1alpha1.IPEntryTableConvertor,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				ipbe)).
			WithResourceAndHandler(ctx, &ipambev1alpha1.IPIndex{}, indexserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&indexserver.ServerObjContext{
					TracerString:   "ipindex-server",
					Obj:            &ipambev1alpha1.IPIndex{},
					ConversionFunc: ipambev1alpha1.IPIndexConvertFieldSelector,
					TableConverter: ipambev1alpha1.IPIndexTableConvertor,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				ipbe)).
			WithResourceAndHandler(ctx, &vlanbev1alpha1.VLANClaim{}, claimserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&claimserver.ServerObjContext{
					TracerString:   "vlanclaim-server",
					Obj:            &vlanbev1alpha1.VLANClaim{},
					ConversionFunc: vlanbev1alpha1.VLANClaimConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				vlanbe)).
			WithResourceAndHandler(ctx, &vlanbev1alpha1.VLANEntry{}, entryserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&entryserver.ServerObjContext{
					TracerString:   "vlanentry-server",
					Obj:            &vlanbev1alpha1.VLANEntry{},
					ConversionFunc: vlanbev1alpha1.VLANEntryConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				vlanbe)).
			WithResourceAndHandler(ctx, &vlanbev1alpha1.VLANIndex{}, indexserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&indexserver.ServerObjContext{
					TracerString: "vlanindex-server",
					Obj:          &vlanbev1alpha1.VLANIndex{},
					//NewIndexFn:   func() backend.IndexObject { return &vlanbev1alpha1.VLANIndex{} },
					ConversionFunc: vlanbev1alpha1.VLANIndexConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				vlanbe)).
			WithResourceAndHandler(ctx, &vxlanbev1alpha1.VXLANClaim{}, claimserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&claimserver.ServerObjContext{
					TracerString: "vxlanclaim-server",
					Obj:          &vxlanbev1alpha1.VXLANClaim{},
					//NewIndexFn:   func() backend.IndexObject { return &vlanbev1alpha1.VLANIndex{} },
					ConversionFunc: vxlanbev1alpha1.VXLANClaimConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				vxlanbe)).
			WithResourceAndHandler(ctx, &vxlanbev1alpha1.VXLANEntry{}, entryserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&entryserver.ServerObjContext{
					TracerString:   "vxlanentry-server",
					Obj:            &vxlanbev1alpha1.VXLANEntry{},
					ConversionFunc: vxlanbev1alpha1.VXLANEntryConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				vxlanbe)).
			WithResourceAndHandler(ctx, &vxlanbev1alpha1.VXLANIndex{}, indexserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&indexserver.ServerObjContext{
					TracerString:   "vxlanindex-server",
					Obj:            &vxlanbev1alpha1.VXLANIndex{},
					ConversionFunc: vxlanbev1alpha1.VXLANIndexConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				vxlanbe)).
			WithResourceAndHandler(ctx, &asbev1alpha1.ASClaim{}, claimserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&claimserver.ServerObjContext{
					TracerString: "asclaim-server",
					Obj:          &asbev1alpha1.ASClaim{},
					//NewIndexFn:   func() backend.IndexObject { return &vlanbev1alpha1.VLANIndex{} },
					ConversionFunc: asbev1alpha1.ASClaimConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				asbe)).
			WithResourceAndHandler(ctx, &asbev1alpha1.ASEntry{}, entryserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&entryserver.ServerObjContext{
					TracerString: "asentry-server",
					Obj:          &asbev1alpha1.ASEntry{},
					//NewIndexFn:   func() backend.IndexObject { return &vlanbev1alpha1.VLANIndex{} },
					ConversionFunc: asbev1alpha1.ASEntryConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				asbe)).
			WithResourceAndHandler(ctx, &asbev1alpha1.ASIndex{}, indexserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&indexserver.ServerObjContext{
					TracerString:   "asindex-server",
					Obj:            &asbev1alpha1.ASIndex{},
					ConversionFunc: asbev1alpha1.ASIndexConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				asbe)).
			WithResourceAndHandler(ctx, &extcommbev1alpha1.EXTCOMMClaim{}, claimserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&claimserver.ServerObjContext{
					TracerString:   "extcommclaim-server",
					Obj:            &extcommbev1alpha1.EXTCOMMClaim{},
					NewIndexFn:     func() backend.IndexObject { return &extcommbev1alpha1.EXTCOMMIndex{} },
					ConversionFunc: extcommbev1alpha1.EXTCOMMClaimConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				extcommbe)).
			WithResourceAndHandler(ctx, &extcommbev1alpha1.EXTCOMMEntry{}, entryserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&entryserver.ServerObjContext{
					TracerString:   "extcommentry-server",
					Obj:            &extcommbev1alpha1.EXTCOMMEntry{},
					ConversionFunc: extcommbev1alpha1.EXTCOMMEntryConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				extcommbe)).
			WithResourceAndHandler(ctx, &extcommbev1alpha1.EXTCOMMIndex{}, indexserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&indexserver.ServerObjContext{
					TracerString:   "extcommindex-server",
					Obj:            &extcommbev1alpha1.EXTCOMMIndex{},
					ConversionFunc: extcommbev1alpha1.EXTCOMMIndexConvertFieldSelector,
					TableConverter: extcommbev1alpha1.EXTCOMMIndexTableConvertor,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				extcommbe)).
			WithResourceAndHandler(ctx, &genidbev1alpha1.GENIDClaim{}, claimserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&claimserver.ServerObjContext{
					TracerString:   "genidclaim-server",
					Obj:            &genidbev1alpha1.GENIDClaim{},
					NewIndexFn:     func() backend.IndexObject { return &genidbev1alpha1.GENIDIndex{} },
					ConversionFunc: genidbev1alpha1.GENIDClaimConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &genidbev1alpha1.GENIDEntry{}, entryserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&entryserver.ServerObjContext{
					TracerString:   "genidentry-server",
					Obj:            &genidbev1alpha1.GENIDEntry{},
					ConversionFunc: genidbev1alpha1.GENIDEntryConvertFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &genidbev1alpha1.GENIDIndex{}, indexserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&indexserver.ServerObjContext{
					TracerString:   "genidindex-server",
					Obj:            &genidbev1alpha1.GENIDIndex{},
					ConversionFunc: genidbev1alpha1.GENIDIndexConvertFieldSelector,
					TableConverter: genidbev1alpha1.GENIDIndexTableConvertor,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.Cluster{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "cluster-server",
					Obj:            &infrabev1alpha1.Cluster{},
					ConversionFunc: infrabev1alpha1.ClusterConvertFieldSelector,
					TableConverter: infrabev1alpha1.ClusterTableConvertor,
					FielSelector:   infrabev1alpha1.ClusterParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.Endpoint{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "endpoint-server",
					Obj:            &infrabev1alpha1.Endpoint{},
					ConversionFunc: infrabev1alpha1.EndpointConvertFieldSelector,
					TableConverter: infrabev1alpha1.EndpointTableConvertor,
					FielSelector:   infrabev1alpha1.EndpointParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.EndpointSet{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "endpointset-server",
					Obj:            &infrabev1alpha1.EndpointSet{},
					ConversionFunc: infrabev1alpha1.EndpointSetConvertFieldSelector,
					TableConverter: infrabev1alpha1.EndpointSetTableConvertor,
					FielSelector:   infrabev1alpha1.EndpointSetParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.Link{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "link-server",
					Obj:            &infrabev1alpha1.Link{},
					ConversionFunc: infrabev1alpha1.LinkConvertFieldSelector,
					TableConverter: infrabev1alpha1.LinkTableConvertor,
					FielSelector:   infrabev1alpha1.LinkParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.LinkSet{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "linkset-server",
					Obj:            &infrabev1alpha1.LinkSet{},
					ConversionFunc: infrabev1alpha1.LinkSetConvertFieldSelector,
					TableConverter: infrabev1alpha1.LinkSetTableConvertor,
					FielSelector:   infrabev1alpha1.LinkSetParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.Module{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "module-server",
					Obj:            &infrabev1alpha1.Module{},
					ConversionFunc: infrabev1alpha1.ModuleConvertFieldSelector,
					TableConverter: infrabev1alpha1.ModuleTableConvertor,
					FielSelector:   infrabev1alpha1.ModuleParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.ModuleBay{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "modulebay-server",
					Obj:            &infrabev1alpha1.ModuleBay{},
					ConversionFunc: infrabev1alpha1.ModuleBayConvertFieldSelector,
					TableConverter: infrabev1alpha1.ModuleBayTableConvertor,
					FielSelector:   infrabev1alpha1.ModuleBayParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.Node{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "node-server",
					Obj:            &infrabev1alpha1.Node{},
					ConversionFunc: infrabev1alpha1.NodeConvertFieldSelector,
					TableConverter: infrabev1alpha1.NodeTableConvertor,
					FielSelector:   infrabev1alpha1.NodeParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.NodeGroup{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "nodegroup-server",
					Obj:            &infrabev1alpha1.NodeGroup{},
					ConversionFunc: infrabev1alpha1.NodeGroupConvertFieldSelector,
					TableConverter: infrabev1alpha1.NodeGroupTableConvertor,
					FielSelector:   infrabev1alpha1.NodeGroupParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.NodeItem{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "nodeitem-server",
					Obj:            &infrabev1alpha1.NodeItem{},
					ConversionFunc: infrabev1alpha1.NodeItemConvertFieldSelector,
					TableConverter: infrabev1alpha1.NodeItemTableConvertor,
					FielSelector:   infrabev1alpha1.NodeItemParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.NodeSet{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "nodeset-server",
					Obj:            &infrabev1alpha1.NodeSet{},
					ConversionFunc: infrabev1alpha1.NodeSetConvertFieldSelector,
					TableConverter: infrabev1alpha1.NodeSetTableConvertor,
					FielSelector:   infrabev1alpha1.NodeSetParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.Rack{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "rack-server",
					Obj:            &infrabev1alpha1.Rack{},
					ConversionFunc: infrabev1alpha1.RackConvertFieldSelector,
					TableConverter: infrabev1alpha1.RackTableConvertor,
					FielSelector:   infrabev1alpha1.RackParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.Region{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "region-server",
					Obj:            &infrabev1alpha1.Region{},
					ConversionFunc: infrabev1alpha1.RegionConvertFieldSelector,
					TableConverter: infrabev1alpha1.RegionTableConvertor,
					FielSelector:   infrabev1alpha1.RegionParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
			WithResourceAndHandler(ctx, &infrabev1alpha1.Site{}, genericserver.NewProvider(
				ctx,
				mgr.GetClient(),
				&genericserver.ServerObjContext{
					TracerString:   "site-server",
					Obj:            &infrabev1alpha1.Site{},
					ConversionFunc: infrabev1alpha1.SiteConvertFieldSelector,
					TableConverter: infrabev1alpha1.SiteTableConvertor,
					FielSelector:   infrabev1alpha1.SiteParseFieldSelector,
				},
				&serverstore.Config{
					Prefix: configDir,
					Type:   serverstore.StorageType_KV,
					DB:     db,
				},
				genidbe)).
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
