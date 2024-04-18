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
	prefix        string
	prefixLength  uint32
	claimType     string
	claimInfo     string
	labels        map[string]string
	selector      *metav1.LabelSelector
	expectedError bool
}

func getIPIndex(niName string) *ipambev1alpha1.IPIndex {
	return ipambev1alpha1.BuildIPIndex(
		metav1.ObjectMeta{Namespace: "dummy", Name: niName},
		nil,
		nil,
	)
}

func getIPClaimFromNetworkPrefix(niName, prefix string, labels map[string]string) (*ipambev1alpha1.IPClaim, error) {
	ni := ipamresv1alpha1.BuildNetworkInstance(
		metav1.ObjectMeta{Namespace: "dummy", Name: niName},
		nil,
		nil,
	)

	return ipambev1alpha1.GetIPClaimFromPrefix(
		ipambev1alpha1.GetIPClaimTypeFromString("aggregate"),
		niName,
		prefix,
		commonv1alpha1.UserDefinedLabels{Labels: labels},
		ni,
	)
}

func getIPClaimFromPrefix(niName, prefix string, claimType *ipambev1alpha1.IPClaimType, labels map[string]string) (*ipambev1alpha1.IPClaim, error) {
	pi, err := iputil.New(prefix)
	if err != nil {
		return nil, err
	}

	pfx := ipamresv1alpha1.BuildIPPrefix(
		metav1.ObjectMeta{Namespace: "dummy", Name: pi.GetSubnetName()},
		&ipamresv1alpha1.IPPrefixSpec{
			NetworkInstance:   niName,
			Type:              claimType,
			Prefix:            prefix,
			UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: labels},
		},
		nil,
	)

	return ipambev1alpha1.GetIPClaimFromPrefix(
		claimType,
		niName,
		prefix,
		commonv1alpha1.UserDefinedLabels{Labels: labels},
		pfx,
	)
}

func getIPClaimFromDynamicPrefix(name, niName string, claimType *ipambev1alpha1.IPClaimType, pl uint32, labels map[string]string, selector *metav1.LabelSelector) (*ipambev1alpha1.IPClaim, error) {
	return ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: name},
		&ipambev1alpha1.IPClaimSpec{
			NetworkInstance: niName,
			Type:            claimType,
			CreatePrefix:    ptr.To[bool](true),
			PrefixLength:    ptr.To[uint32](pl),
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: labels},
				Selector:          selector,
			},
			Owner: &commonv1alpha1.OwnerReference{
				Group:     ipambev1alpha1.SchemeGroupVersion.Group,
				Version:   ipambev1alpha1.SchemeGroupVersion.Version,
				Kind:      ipambev1alpha1.IPClaimKind,
				Namespace: "dummy",
				Name:      name,
			},
		},
		nil,
	), nil
}

func getIPClaimFromAddress(niName, address string, labels map[string]string) (*ipambev1alpha1.IPClaim, error) {
	pi, err := iputil.New(address)
	if err != nil {
		return nil, err
	}

	addr := ipamresv1alpha1.BuildIPAddress(
		metav1.ObjectMeta{Namespace: "dummy", Name: pi.GetSubnetName()},
		&ipamresv1alpha1.IPAddressSpec{
			NetworkInstance:   niName,
			Address:           address,
			UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: labels},
		},
		nil,
	)

	return ipambev1alpha1.GetIPClaimFromAddress(
		niName,
		address,
		commonv1alpha1.UserDefinedLabels{Labels: labels},
		addr,
	)
}

func getIPClaimFromDynamicAddress(name, niName string, labels map[string]string, selector *metav1.LabelSelector) (*ipambev1alpha1.IPClaim, error) {
	fmt.Println("selector", selector)
	
	return ipambev1alpha1.BuildIPClaim(
		metav1.ObjectMeta{Namespace: "dummy", Name: name},
		&ipambev1alpha1.IPClaimSpec{
			NetworkInstance: niName,
			Type:            nil,
			ClaimLabels: commonv1alpha1.ClaimLabels{
				UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: labels},
				Selector:          selector,
			},
			Owner: &commonv1alpha1.OwnerReference{
				Group:     ipambev1alpha1.SchemeGroupVersion.Group,
				Version:   ipambev1alpha1.SchemeGroupVersion.Version,
				Kind:      ipambev1alpha1.IPClaimKind,
				Namespace: "dummy",
				Name:      name,
			},
		},
		nil,
	), nil
}

func getIPClaimFromRange(name, niName, ipRange string, labels map[string]string) (*ipambev1alpha1.IPClaim, error) {
	r := ipamresv1alpha1.BuildIPRange(
		metav1.ObjectMeta{Namespace: "dummy", Name: name},
		&ipamresv1alpha1.IPRangeSpec{
			NetworkInstance:   niName,
			Range:             ipRange,
			UserDefinedLabels: commonv1alpha1.UserDefinedLabels{Labels: labels},
		},
		nil,
	)

	return ipambev1alpha1.GetIPClaimFromRange(
		name,
		niName,
		ipRange,
		commonv1alpha1.UserDefinedLabels{Labels: labels},
		r,
	)
}
