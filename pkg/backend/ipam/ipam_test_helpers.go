package ipam

import (
	"fmt"

	"github.com/henderiw/iputil"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

type testprefix struct {
	name          string
	claimType     ipambev1alpha1.IPClaimType
	prefixType    *ipambev1alpha1.IPPrefixType
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
	namespace    = "dummy"
	staticPrefix   = ipambev1alpha1.IPClaimType_StaticPrefix
	staticRange    = ipambev1alpha1.IPClaimType_StaticRange
	staticAddress  = ipambev1alpha1.IPClaimType_StaticAddress
	dynamicPrefix  = ipambev1alpha1.IPClaimType_DynamicPrefix
	dynamicAddress = ipambev1alpha1.IPClaimType_DynamicAddress
)

var aggregate = ptr.To[ipambev1alpha1.IPPrefixType](ipambev1alpha1.IPPrefixType_Aggregate)
var network = ptr.To[ipambev1alpha1.IPPrefixType](ipambev1alpha1.IPPrefixType_Network)
var pool = ptr.To[ipambev1alpha1.IPPrefixType](ipambev1alpha1.IPPrefixType_Pool)
var other = ptr.To[ipambev1alpha1.IPPrefixType](ipambev1alpha1.IPPrefixType_Other)

func getIndex(index string) *ipambev1alpha1.IPIndex {
	return ipambev1alpha1.BuildIPIndex(
		metav1.ObjectMeta{Namespace: namespace, Name: index},
		nil,
		nil,
	)
}

func (r testprefix) getIPClaimFromNetworkPrefix(index string) (*ipambev1alpha1.IPClaim, error) {
	idx := ipambev1alpha1.BuildIPIndex(
		metav1.ObjectMeta{Namespace: namespace, Name: index},
		nil,
		nil,
	)
	return idx.GetClaim(ipambev1alpha1.Prefix{Prefix: r.ip, UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels}})
}

func (r testprefix) getStaticPrefixIPClaim(index string) (*ipambev1alpha1.IPClaim, error) {
	pi, err := iputil.New(r.ip)
	if err != nil {
		return nil, err
	}
	ipClaim := ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: pi.GetSubnetName()},
		&ipambev1alpha1.IPClaimSpec{
			Index: index,
			PrefixType:      r.prefixType,
			Prefix:          ptr.To[string](r.ip),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
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

func (r testprefix) getDynamicPrefixIPClaim(index string) (*ipambev1alpha1.IPClaim, error) {
	ipClaim := ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&ipambev1alpha1.IPClaimSpec{
			Index: index,
			PrefixType:      r.prefixType,
			CreatePrefix:    ptr.To[bool](true),
			PrefixLength:    ptr.To[uint32](r.prefixLength),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
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

func (r testprefix) getStaticAddressIPClaim(index string) (*ipambev1alpha1.IPClaim, error) {
	pi, err := iputil.New(r.ip)
	if err != nil {
		return nil, err
	}

	pi = iputil.NewPrefixInfo(pi.GetIPAddressPrefix())

	ipClaim := ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: pi.GetSubnetName()},
		&ipambev1alpha1.IPClaimSpec{
			Index: index,
			Address:         ptr.To[string](r.ip),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
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


func (r testprefix) getDynamicAddressIPClaim(index string) (*ipambev1alpha1.IPClaim, error) {
	ipClaim := ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&ipambev1alpha1.IPClaimSpec{
			Index: index,
			PrefixType:      nil,
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
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

func (r testprefix) getStaticRangeIPClaim(index string) (*ipambev1alpha1.IPClaim, error) {
	ipClaim := ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: namespace, Name: r.name},
		&ipambev1alpha1.IPClaimSpec{
			Index: index,
			Range:           ptr.To[string](r.ip),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
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
