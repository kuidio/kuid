package ipam

import (
	"context"
	"fmt"
	"reflect"

	"github.com/henderiw/apiserver-builder/pkg/builder"
	"github.com/henderiw/apiserver-builder/pkg/builder/resource"
	"github.com/henderiw/apiserver-store/pkg/generic/registry"
	"github.com/henderiw/iputil"
	"github.com/kuidio/kuid/apis/backend/ipam"
	"github.com/kuidio/kuid/apis/backend/ipam/register"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	"github.com/kuidio/kuid/apis/common"
	bebackend "github.com/kuidio/kuid/pkg/backend"
	ipambe "github.com/kuidio/kuid/pkg/backend/ipam"
	"github.com/kuidio/kuid/pkg/config"
	"github.com/kuidio/kuid/pkg/generated/openapi"
	"github.com/kuidio/kuid/pkg/registry/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/utils/ptr"
)

type testprefix struct {
	name          string
	claimType     ipam.IPClaimType
	prefixType    *ipam.IPPrefixType
	ip            string
	prefixLength  uint32
	labels        map[string]string
	selector      *metav1.LabelSelector
	expectedError bool
	expectedDG    string
	expectedIP    string
}

// alias
const (
	namespace      = "dummy"
	staticPrefix   = ipam.IPClaimType_StaticPrefix
	staticRange    = ipam.IPClaimType_StaticRange
	staticAddress  = ipam.IPClaimType_StaticAddress
	dynamicPrefix  = ipam.IPClaimType_DynamicPrefix
	dynamicAddress = ipam.IPClaimType_DynamicAddress
)

var aggregate = ptr.To(ipam.IPPrefixType_Aggregate)
var network = ptr.To(ipam.IPPrefixType_Network)
var pool = ptr.To(ipam.IPPrefixType_Pool)

//var other = ptr.To(ipam.IPPrefixType_Other)

func apiServer() *builder.Server {
	return builder.NewAPIServer().
		WithServerName("kuid-api-server").
		WithOpenAPIDefinitions("Config", "v1alpha1", openapi.GetOpenAPIDefinitions).
		WithoutEtcd()
}

func initBackend(ctx context.Context, apiserver *builder.Server) (bebackend.Backend, error) {
	groupConfig := config.GroupConfig{
		BackendFn:               register.NewBackend,
		ApplyStorageToBackendFn: register.ApplyStorageToBackend,
		Resources: []*config.ResourceConfig{
			{StorageProviderFn: register.NewIndexStorageProvider, Internal: &ipam.IPIndex{}, ResourceVersions: []resource.Object{&ipam.IPIndex{}, &ipambev1alpha1.IPIndex{}}},
			{StorageProviderFn: register.NewClaimStorageProvider, Internal: &ipam.IPClaim{}, ResourceVersions: []resource.Object{&ipam.IPClaim{}, &ipambev1alpha1.IPClaim{}}},
			{StorageProviderFn: register.NewStorageProvider, Internal: &ipam.IPEntry{}, ResourceVersions: []resource.Object{&ipam.IPEntry{}, &ipambev1alpha1.IPEntry{}}},
		},
	}

	be := groupConfig.BackendFn()
	for _, resource := range groupConfig.Resources {
		storageProvider := resource.StorageProviderFn(ctx, resource.Internal, be, true, &options.Options{
			Type: options.StorageType_Memory,
		})
		for _, resourceVersion := range resource.ResourceVersions {
			apiserver.WithResourceAndHandler(resourceVersion, storageProvider)
		}
	}

	if _, err := apiserver.Build(ctx); err != nil {
		return nil, err
	}
	if err := groupConfig.ApplyStorageToBackendFn(ctx, be, apiserver); err != nil {
		return nil, err
	}
	return be, nil
}

func getStorage(ctx context.Context, apiServer *builder.Server, gr schema.GroupResource) (*registry.Store, error) {
	storageProvider := apiServer.StorageProvider[gr]
	storage, err := storageProvider.Get(ctx, apiServer.Schemes[0], &Getter{})
	if err != nil {
		return nil, err
	}
	registryStore, ok := storage.(*registry.Store)
	if !ok {
		return nil, fmt.Errorf("index store is not a *registry.Store, got: %v", reflect.TypeOf(storage).Name())
	}
	return registryStore, nil
}

var _ generic.RESTOptionsGetter = &Getter{}

type Getter struct{}

func (r *Getter) GetRESTOptions(resource schema.GroupResource, example runtime.Object) (generic.RESTOptions, error) {
	return generic.RESTOptions{}, nil
}

func getIndex(index string, prefixes []ipam.Prefix) *ipam.IPIndex {
	return ipam.BuildIPIndex(
		metav1.ObjectMeta{Namespace: namespace, Name: index},
		&ipam.IPIndexSpec{
			Prefixes: prefixes,
		},
		nil,
	)
}

func (r testprefix) getStaticPrefixIPClaim(index string) (*ipam.IPClaim, error) {
	pi, err := iputil.New(r.ip)
	if err != nil {
		return nil, err
	}
	ipClaim := ipam.BuildIPClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: pi.GetSubnetName()},
		&ipam.IPClaimSpec{
			Index:      index,
			PrefixType: r.prefixType,
			Prefix:     ptr.To(r.ip),
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := ipClaim.ValidateSyntax("") // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return ipClaim, nil
}

func (r testprefix) getDynamicPrefixIPClaim(index string) (*ipam.IPClaim, error) {
	ipClaim := ipam.BuildIPClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&ipam.IPClaimSpec{
			Index:        index,
			PrefixType:   r.prefixType,
			CreatePrefix: ptr.To[bool](true),
			PrefixLength: ptr.To[uint32](r.prefixLength),
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
				Selector:          r.selector,
			},
		},
		nil,
	)
	fielErrList := ipClaim.ValidateSyntax("") // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return ipClaim, nil
}

func (r testprefix) getStaticAddressIPClaim(index string) (*ipam.IPClaim, error) {
	pi, err := iputil.New(r.ip)
	if err != nil {
		return nil, err
	}

	pi = iputil.NewPrefixInfo(pi.GetIPAddressPrefix())

	ipClaim := ipam.BuildIPClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: pi.GetSubnetName()},
		&ipam.IPClaimSpec{
			Index:   index,
			Address: ptr.To[string](r.ip),
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := ipClaim.ValidateSyntax("") // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return ipClaim, nil
}

func (r testprefix) getDynamicAddressIPClaim(index string) (*ipam.IPClaim, error) {
	ipClaim := ipam.BuildIPClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&ipam.IPClaimSpec{
			Index:      index,
			PrefixType: nil,
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
				Selector:          r.selector,
			},
		},
		nil,
	)
	fielErrList := ipClaim.ValidateSyntax("") // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return ipClaim, nil
}

func (r testprefix) getStaticRangeIPClaim(index string) (*ipam.IPClaim, error) {
	ipClaim := ipam.BuildIPClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&ipam.IPClaimSpec{
			Index: index,
			Range: ptr.To[string](r.ip),
			ClaimLabels: common.ClaimLabels{
				UserDefinedLabels: common.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := ipClaim.ValidateSyntax("") // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return ipClaim, nil
}

type prefixTest struct {
	index         string
	indexPrefixes []ipam.Prefix
	prefixes      []testprefix
}

func prefixTestRun(name string, tc prefixTest) error {
	ctx := context.Background()
	apiserver := apiServer()
	be, err := initBackend(ctx, apiserver)
	if err != nil {
		return fmt.Errorf("cannot get backend, err: %v", err)
	}

	indexStorage, err := getStorage(ctx, apiserver, schema.GroupResource{
		Group:    ipam.SchemeGroupVersion.Group,
		Resource: ipam.IPIndexPlural,
	})
	if err != nil {
		return fmt.Errorf("cannot get index storage, err: %v", err)
	}

	claimStorage, err := getStorage(ctx, apiserver, schema.GroupResource{
		Group:    ipam.SchemeGroupVersion.Group,
		Resource: ipam.IPClaimPlural,
	})
	if err != nil {
		return fmt.Errorf("cannot get claim storage, err: %v", err)
	}
	index := getIndex(tc.index, tc.indexPrefixes)
	ctx = genericapirequest.WithNamespace(ctx, index.GetNamespace())
	newIndex, err := indexStorage.Create(ctx, index, nil, &metav1.CreateOptions{
		FieldManager: "backend",
	})
	if err != nil {
		return fmt.Errorf("cannot create index, err: %v", err)
	}
	index, ok := newIndex.(*ipam.IPIndex)
	if !ok {
		return fmt.Errorf("not an ip index, got: %s", reflect.TypeOf(index).Name())
	}

	for _, p := range tc.prefixes {
		var claim *ipam.IPClaim
		var err error

		switch p.claimType {
		case staticPrefix:
			claim, err = p.getStaticPrefixIPClaim(tc.index)
		case staticRange:
			claim, err = p.getStaticRangeIPClaim(tc.index)
		case staticAddress:
			claim, err = p.getStaticAddressIPClaim(tc.index)
		case dynamicPrefix:
			claim, err = p.getDynamicPrefixIPClaim(tc.index)
		case dynamicAddress:
			claim, err = p.getDynamicAddressIPClaim(tc.index)
		}
		if err != nil {
			return fmt.Errorf("wrong prefix type, err: %v", err)
		}

		ctx = genericapirequest.WithNamespace(ctx, claim.GetNamespace())

		exists := true
		oldClaim, err := claimStorage.Get(ctx, claim.GetName(), &metav1.GetOptions{})
		if err != nil {
			exists = false
		}
		var newClaim runtime.Object
		if !exists {
			newClaim, err = claimStorage.Create(ctx, claim, nil, &metav1.CreateOptions{FieldManager: "test"})
		} else {
			defaultObjInfo := rest.DefaultUpdatedObjectInfo(oldClaim, ipambe.ClaimTransformer)
			newClaim, _, err = claimStorage.Update(ctx, claim.GetName(), defaultObjInfo, nil, nil, false, &metav1.UpdateOptions{
				FieldManager: "backend",
			})
		}
		if p.expectedError {
			if err == nil {
				return fmt.Errorf("expected error. got nil")
			}
			continue
		}
		if err != nil {
			return fmt.Errorf("unexpected error, got: %v", err)
		}
		ipClaim, ok := newClaim.(*ipam.IPClaim)
		if !ok {
			return fmt.Errorf("expecting ipClaim, got: %v", reflect.TypeOf(newClaim).Name())
		}

		switch p.claimType {
		case staticPrefix, dynamicPrefix:
			if ipClaim.Status.Prefix == nil {
				return fmt.Errorf("expecting prefix status got nil")
			} else {
				expectedIP := p.ip
				if p.expectedIP != "" {
					expectedIP = p.expectedIP
				}
				if *ipClaim.Status.Prefix != expectedIP {
					return fmt.Errorf("expecting prefix got %s, want %s", *ipClaim.Status.Prefix, expectedIP)
				}
			}
		case staticAddress, dynamicAddress:
			if ipClaim.Status.Address == nil {
				return fmt.Errorf("expecting address status got nil")
			} else {
				expectedIP := p.ip
				if p.expectedIP != "" {
					expectedIP = p.expectedIP
				}
				if *ipClaim.Status.Address != expectedIP {
					return fmt.Errorf("expecting address got %s, want %s", *ipClaim.Status.Address, expectedIP)
				}
			}
			if ipClaim.Status.DefaultGateway == nil {
				if p.expectedDG != "" {
					return fmt.Errorf("expecting defaultGateway %s got nil", p.expectedDG)
				}
			} else {
				if p.expectedDG == "" {
					return fmt.Errorf("unexpected defaultGateway got %s", *ipClaim.Status.DefaultGateway)
				}
				if *ipClaim.Status.DefaultGateway != p.expectedDG {
					return fmt.Errorf("expecting defaultGateway got %s, want %s", *ipClaim.Status.DefaultGateway, p.expectedDG)
				}
			}
		case staticRange:
			if ipClaim.Status.Range == nil {
				return fmt.Errorf("expecting range status got nil")
			} else {
				if *ipClaim.Status.Range != p.ip {
					return fmt.Errorf("expecting prefix got %s, want %s", *ipClaim.Status.Range, p.ip)
				}
			}
		}
	}
	if name == "" {
		fmt.Println("###############")
		be.PrintEntries(ctx, tc.index)
		fmt.Println("###############")
	}
	return nil
}
