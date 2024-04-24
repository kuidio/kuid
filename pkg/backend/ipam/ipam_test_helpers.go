package ipam

import (
	"fmt"

	"github.com/henderiw/iputil"
	ipambev1alpha1 "github.com/kuidio/kuid/apis/backend/ipam/v1alpha1"
	commonv1alpha1 "github.com/kuidio/kuid/apis/common/v1alpha1"
	ipamresv1alpha1 "github.com/kuidio/kuid/apis/resource/ipam/v1alpha1"
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

func getNI(niName string) *ipamresv1alpha1.NetworkInstance {
	return ipamresv1alpha1.BuildNetworkInstance(
		metav1.ObjectMeta{Namespace: "dummy", Name: niName},
		nil,
		nil,
	)
}

func (r testprefix) getIPClaimFromNetworkPrefix(niName string) (*ipambev1alpha1.IPClaim, error) {
	ni := ipamresv1alpha1.BuildNetworkInstance(
		metav1.ObjectMeta{Namespace: "dummy", Name: niName},
		nil,
		nil,
	)
	return ni.GetIPClaim(ipamresv1alpha1.Prefix{Prefix: r.ip, UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels}})
}

func (r testprefix) getStaticPrefixIPClaim(niName string) (*ipambev1alpha1.IPClaim, error) {
	pi, err := iputil.New(r.ip)
	if err != nil {
		return nil, err
	}
	ipClaim := ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: pi.GetSubnetName()},
		&ipambev1alpha1.IPClaimSpec{
			NetworkInstance: niName,
			PrefixType:      r.prefixType,
			Prefix:          ptr.To[string](r.ip),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := ipClaim.ValidateSyntax() // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return ipClaim, nil
}

func (r testprefix) getDynamicPrefixIPClaim(niName string) (*ipambev1alpha1.IPClaim, error) {
	ipClaim := ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: r.name},
		&ipambev1alpha1.IPClaimSpec{
			NetworkInstance: niName,
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
	fielErrList := ipClaim.ValidateSyntax() // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return ipClaim, nil
}

func (r testprefix) getStaticAddressIPClaim(niName string) (*ipambev1alpha1.IPClaim, error) {
	pi, err := iputil.New(r.ip)
	if err != nil {
		return nil, err
	}

	pi = iputil.NewPrefixInfo(pi.GetIPAddressPrefix())

	ipClaim := ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: pi.GetSubnetName()},
		&ipambev1alpha1.IPClaimSpec{
			NetworkInstance: niName,
			Address:         ptr.To[string](r.ip),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := ipClaim.ValidateSyntax() // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return ipClaim, nil
}


func (r testprefix) getDynamicAddressIPClaim(niName string) (*ipambev1alpha1.IPClaim, error) {
	ipClaim := ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: r.name},
		&ipambev1alpha1.IPClaimSpec{
			NetworkInstance: niName,
			PrefixType:      nil,
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
				Selector:          r.selector,
			},
		},
		nil,
	)
	fielErrList := ipClaim.ValidateSyntax() // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return ipClaim, nil
}

func (r testprefix) getStaticRangeIPClaim(niName string) (*ipambev1alpha1.IPClaim, error) {
	ipClaim := ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: r.name},
		&ipambev1alpha1.IPClaimSpec{
			NetworkInstance: niName,
			Range:           ptr.To[string](r.ip),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: r.labels},
			},
		},
		nil,
	)
	fielErrList := ipClaim.ValidateSyntax() // this expands the ownerRef in the spec
	if len(fielErrList) != 0 {
		return nil, fmt.Errorf("invalid syntax %v", fielErrList)
	}
	return ipClaim, nil
}
