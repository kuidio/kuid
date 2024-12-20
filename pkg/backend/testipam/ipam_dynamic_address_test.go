package ipam

import (
	"testing"

	"github.com/kuidio/kuid/apis/backend"
	"github.com/kuidio/kuid/apis/backend/ipam"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIPAMDynamicAddress(t *testing.T) {
	tests := map[string]prefixTest{
		"FromOther": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/8"},
			},
			prefixes: []testprefix{
				{claimType: dynamicAddress, name: "addrClaim1", expectedError: false, expectedIP: "10.0.0.0/32"},
			},
		},
		"FromNetwork": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/8"},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, name: "network1", ip: "10.0.0.0/24", prefixType: network, expectedError: false},
				{claimType: dynamicAddress, name: "addrClaim1", expectedError: false, expectedIP: "10.0.0.1/24", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "network1"},
				}},
				//{claimType: dynamicAddress, name: "addrClaim1", expectedError: false, expectedIP: "10.0.0.1/24", selector: &metav1.LabelSelector{
				//	MatchLabels: map[string]string{backend.KuidClaimNameKey: "network1"},
				//}}, // we explicitly reclaim the same ip
				{claimType: dynamicAddress, name: "addrClaim2", expectedError: false, expectedIP: "10.0.0.254/24", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "network1"},
				}},
				{claimType: dynamicAddress, name: "addrClaim3", expectedError: false, expectedIP: "10.0.0.2/24", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "network1"},
				}},
				{claimType: dynamicAddress, name: "addrClaim4", expectedError: false, expectedIP: "10.0.0.3/24", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "network1"},
				}},
				{claimType: dynamicAddress, name: "addrClaim5", expectedError: false, expectedIP: "10.0.0.252/24", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "network1"},
				}},
			},
		},
		"FromRange": {
			index: "a",
			indexPrefixes: []ipam.Prefix{
				{Prefix: "10.0.0.0/8"},
			},
			prefixes: []testprefix{
				{claimType: staticPrefix, name: "network1", ip: "10.0.0.0/24", expectedError: false},
				{claimType: staticRange, name: "range1", ip: "10.0.0.10-10.0.0.100", expectedError: false},
				{claimType: dynamicAddress, name: "addrClaim1", expectedError: false, expectedIP: "10.0.0.10/32", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "range1"},
				}},
				{claimType: dynamicAddress, name: "addrClaim1", expectedError: false, expectedIP: "10.0.0.10/32", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "range1"}, // we explicitly reclaim the same ip
				}},
				{claimType: dynamicAddress, name: "addrClaim2", expectedError: false, expectedIP: "10.0.0.11/32", selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{backend.KuidClaimNameKey: "range1"},
				}},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := prefixTestRun(name, tc); err != nil {
				t.Errorf("test %s failed err: %v", name, err)
			}
		})
	}
}
