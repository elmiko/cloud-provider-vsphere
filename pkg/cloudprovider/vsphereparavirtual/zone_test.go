package vsphereparavirtual

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/rest"
	cloudprovider "k8s.io/cloud-provider"

	vmopclient "k8s.io/cloud-provider-vsphere/pkg/cloudprovider/vsphereparavirtual/vmoperator/client"
)

var (
	vmName     = types.NodeName("test-vm")
	fakeVMName = types.NodeName("fake-vm")
	vmuuid     = "421960e7-3041-f44a-4b3f-ed99748c12d0"
	providerid = "vsphere://" + vmuuid
)

func TestNewZones(t *testing.T) {
	testCases := []struct {
		name        string
		config      *rest.Config
		expectedErr error
		testVM      *vmopv1.VirtualMachine
	}{
		{
			name:        "NewZone: when everything is ok",
			config:      &rest.Config{},
			testVM:      createTestVMWithZone(string(vmName), testClusterNameSpace),
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//initVMopClient(testCase.testVM)
			_, err := NewZones(testClusterNameSpace, testCase.config)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestZonesByProviderID(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult string
		expectedErr    error
		testVM         *vmopv1.VirtualMachine
	}{
		{
			name:           "TestZonesByProviderID should return true",
			testVM:         createTestVMWithZoneID(string(vmName), testClusterNameSpace, vmuuid),
			expectedResult: "zone-a",
			expectedErr:    nil,
		},
		{
			name:           "TestZonesByProviderID should return error",
			testVM:         createTestVMWithZoneID(string(vmName), testClusterNameSpace, "fakeuuid"),
			expectedResult: "",
			expectedErr:    cloudprovider.InstanceNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()

			zone, _, err := initVMopClient(testCase.testVM)
			assert.NoError(t, err)
			z, err := zone.GetZoneByProviderID(ctx, providerid)

			if testCase.expectedErr != nil {
				assert.Equal(t, cloudprovider.InstanceNotFound, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.expectedResult, z.FailureDomain)
		})
	}
}

func TestZonesByNodeName(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult string
		expectedErr    error
		testVM         *vmopv1.VirtualMachine
		vmName         types.NodeName
	}{
		{
			name:           "TestZonesByNodeName should return true",
			testVM:         createTestVMWithZoneID(string(vmName), testClusterNameSpace, vmuuid),
			vmName:         vmName,
			expectedResult: "zone-a",
			expectedErr:    nil,
		},
		{
			name:           "TestZonesByNodeName should return error",
			testVM:         createTestVMWithZoneID(string(vmName), testClusterNameSpace, "fakeuuid"),
			vmName:         fakeVMName,
			expectedResult: "",
			expectedErr:    cloudprovider.InstanceNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()

			zone, _, err := initVMopClient(testCase.testVM)
			assert.NoError(t, err)
			z, err := zone.GetZoneByNodeName(ctx, testCase.vmName)

			if testCase.expectedErr != nil {
				assert.Equal(t, cloudprovider.InstanceNotFound, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.expectedResult, z.FailureDomain)
		})
	}
}

func initVMopClient(testVM *vmopv1.VirtualMachine) (zones, *dynamicfake.FakeDynamicClient, error) {
	scheme := runtime.NewScheme()
	_ = vmopv1.AddToScheme(scheme)
	fc := dynamicfake.NewSimpleDynamicClient(scheme)
	fcw := vmopclient.NewFakeClientSet(fc)
	zone := zones{
		vmClient:  fcw,
		namespace: testClusterNameSpace,
	}
	_, err := fcw.V1alpha2().VirtualMachines(testVM.Namespace).Create(context.TODO(), testVM, metav1.CreateOptions{})
	return zone, fc, err
}

func createTestVMWithZone(name, namespace string) *vmopv1.VirtualMachine {
	labels := make(map[string]string)
	labels["topology.kubernetes.io/zone"] = "zone-a"
	return &vmopv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
	}
}

func createTestVMWithZoneID(name, namespace, biosUUID string) *vmopv1.VirtualMachine {
	labels := make(map[string]string)
	labels["topology.kubernetes.io/zone"] = "zone-a"
	return &vmopv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Status: vmopv1.VirtualMachineStatus{
			BiosUUID: biosUUID,
		},
	}
}
