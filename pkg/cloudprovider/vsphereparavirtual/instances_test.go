/*
Copyright 2021 The Kubernetes Authors.

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

package vsphereparavirtual

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	clientgotesting "k8s.io/client-go/testing"
	cloudprovider "k8s.io/cloud-provider"
	vmopclient "k8s.io/cloud-provider-vsphere/pkg/cloudprovider/vsphereparavirtual/vmoperator/client"

	dynamicfake "k8s.io/client-go/dynamic/fake"
)

var (
	testVMName     = types.NodeName("test-vm")
	testVMUUID     = "1bbf49a7-fbce-4502-bb4c-4c3544cacc9e"
	testProviderID = providerPrefix + testVMUUID
)

func createTestVM(name, namespace, biosUUID string) *vmopv1.VirtualMachine {
	return &vmopv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Status: vmopv1.VirtualMachineStatus{
			BiosUUID: biosUUID,
		},
	}
}

func createTestVMWithVMIPAndHost(name, namespace, biosUUID string) *vmopv1.VirtualMachine {
	// TODO: Currently, dual-stack (IPv4 and IPv6) is not supported.
	// Cluster will be assumed as IPv4 Primary by default.
	// In the future, when dual-stack support is implemented, this code should be updated to
	// dynamically determine the IP format based on the cluster's IP family.
	// https://github.com/kubernetes/cloud-provider-vsphere/issues/1129
	return &vmopv1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Status: vmopv1.VirtualMachineStatus{
			BiosUUID: biosUUID,
			Host:     "test-host",
			Network: &vmopv1.VirtualMachineNetworkStatus{
				PrimaryIP4: "1.2.3.4",
			},
		},
	}
}

func TestNewInstances(t *testing.T) {
	testCases := []struct {
		name        string
		config      *rest.Config
		expectedErr error
	}{
		{
			name:        "NewInstance: when everything is ok",
			config:      &rest.Config{},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := NewInstances(testClusterNameSpace, testCase.config)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func initTest(testVM *vmopv1.VirtualMachine) (*instances, *dynamicfake.FakeDynamicClient, error) {
	scheme := runtime.NewScheme()
	_ = vmopv1.AddToScheme(scheme)
	fc := dynamicfake.NewSimpleDynamicClient(scheme)
	fcw := vmopclient.NewFakeClientSet(fc)
	instance := &instances{
		vmClient:  fcw,
		namespace: testClusterNameSpace,
	}
	_, err := fcw.V1alpha2().VirtualMachines(testVM.Namespace).Create(context.TODO(), testVM, metav1.CreateOptions{})
	return instance, fc, err
}

func TestInstanceID(t *testing.T) {
	testCases := []struct {
		name                string
		testVM              *vmopv1.VirtualMachine
		expectInternalError bool
		expectedInstanceID  string
		expectedErr         error
	}{
		{
			name:               "test Instance interface: should not return error",
			testVM:             createTestVM(string(testVMName), testClusterNameSpace, testVMUUID),
			expectedInstanceID: testVMUUID,
			expectedErr:        nil,
		},
		{
			name:               "cannot find virtualmachine with node name",
			testVM:             createTestVM("bogus", testClusterNameSpace, testVMUUID),
			expectedInstanceID: "",
			expectedErr:        cloudprovider.InstanceNotFound,
		},
		{
			name:               "cannot find virtualmachine with namespace",
			testVM:             createTestVM(string(testVMName), "bogus", testVMUUID),
			expectedInstanceID: "",
			expectedErr:        cloudprovider.InstanceNotFound,
		},
		{
			name:               "cannot find virtualmachine with empty bios uuid",
			testVM:             createTestVM(string(testVMName), testClusterNameSpace, ""),
			expectedInstanceID: "",
			expectedErr:        errBiosUUIDEmpty,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance, _, err := initTest(testCase.testVM)
			assert.NoError(t, err)
			instanceID, err := instance.InstanceID(context.Background(), testVMName)
			assert.Equal(t, testCase.expectedErr, err)
			assert.Equal(t, testCase.expectedInstanceID, instanceID)
		})
	}
}

func TestInstanceIDThrowsErr(t *testing.T) {
	testCases := []struct {
		name               string
		testVM             *vmopv1.VirtualMachine
		expectedInstanceID string
	}{
		{
			name:               "test Instance interface: throws an error in client.Get()",
			testVM:             createTestVM(string(testVMName), testClusterNameSpace, testVMUUID),
			expectedInstanceID: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance, fc, err := initTest(testCase.testVM)
			assert.NoError(t, err)
			fc.PrependReactor("get", "virtualmachines", func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, nil, fmt.Errorf("Internal error getting VMs")
			})
			instanceID, err := instance.InstanceID(context.Background(), testVMName)
			assert.NotEqual(t, nil, err)
			assert.NotEqual(t, cloudprovider.InstanceNotFound, err)
			assert.Equal(t, testCase.expectedInstanceID, instanceID)
		})
	}
}

func TestInstanceExistsByProviderID(t *testing.T) {
	testCases := []struct {
		name           string
		testVM         *vmopv1.VirtualMachine
		expectedResult bool
		expectedErr    error
	}{
		{
			name:           "InstanceExistsByProviderID should return true",
			testVM:         createTestVM(string(testVMName), testClusterNameSpace, testVMUUID),
			expectedResult: true,
			expectedErr:    nil,
		},
		{
			name:           "InstanceExistsByProviderID should return false",
			testVM:         createTestVM(string(testVMName), testClusterNameSpace, "bogus"),
			expectedResult: false,
			expectedErr:    nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance, _, err := initTest(testCase.testVM)
			assert.NoError(t, err)
			providerID, err := instance.InstanceExistsByProviderID(context.Background(), testProviderID)
			assert.Equal(t, testCase.expectedErr, err)
			assert.Equal(t, testCase.expectedResult, providerID)
		})
	}
}

func TestInstanceShutdownByProviderID(t *testing.T) {
	testCases := []struct {
		name             string
		testVM           *vmopv1.VirtualMachine
		testVMPowerState string
		expectedResult   bool
		expectedErr      error
	}{
		{
			name:             "InstanceShutdownByProviderID should return true for powered-off VM",
			testVM:           createTestVM(string(testVMName), testClusterNameSpace, testVMUUID),
			testVMPowerState: "PoweredOff",
			expectedResult:   true,
			expectedErr:      nil,
		},
		{
			name:             "InstanceShutdownByProviderID should return false for powered-on VM",
			testVM:           createTestVM(string(testVMName), testClusterNameSpace, testVMUUID),
			testVMPowerState: "PoweredOn",
			expectedResult:   false,
			expectedErr:      nil,
		},
		{
			name:             "InstanceShutdownByProviderID node not found for powered-on VM",
			testVM:           createTestVM(string(testVMName), testClusterNameSpace, "bogus"),
			testVMPowerState: "PoweredOn",
			expectedResult:   false,
			expectedErr:      cloudprovider.InstanceNotFound,
		},
		{
			name:             "InstanceShutdownByProviderID node not found for powered-off VM",
			testVM:           createTestVM(string(testVMName), testClusterNameSpace, "bogus"),
			testVMPowerState: "PoweredOff",
			expectedResult:   false,
			expectedErr:      cloudprovider.InstanceNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.testVMPowerState == "PoweredOn" {
				testCase.testVM.Status.PowerState = vmopv1.VirtualMachinePowerStateOn
			} else {
				testCase.testVM.Status.PowerState = vmopv1.VirtualMachinePowerStateOff
			}

			instance, _, err := initTest(testCase.testVM)
			assert.NoError(t, err)
			ret, err := instance.InstanceShutdownByProviderID(context.Background(), testProviderID)
			assert.Equal(t, testCase.expectedErr, err)
			assert.Equal(t, testCase.expectedResult, ret)
		})
	}
}

func TestNodeAddressesByProviderID(t *testing.T) {
	testCases := []struct {
		name                string
		testVM              *vmopv1.VirtualMachine
		expectedNodeAddress []v1.NodeAddress
		expectedErr         error
	}{
		{
			name:                "NodeAddressesByProviderID returns an empty address for found node with no IP",
			testVM:              createTestVM(string(testVMName), testClusterNameSpace, testVMUUID),
			expectedNodeAddress: []v1.NodeAddress{},
			expectedErr:         nil,
		},
		{
			name:                "NodeAddressesByProviderID returns a NotFound error for a not found node",
			testVM:              createTestVM(string(testVMName), testClusterNameSpace, "bogus"),
			expectedNodeAddress: nil,
			expectedErr:         cloudprovider.InstanceNotFound,
		},
		{
			name:   "NodeAddressesByProviderID returns a valid address for a found node",
			testVM: createTestVMWithVMIPAndHost(string(testVMName), testClusterNameSpace, testVMUUID),
			expectedNodeAddress: []v1.NodeAddress{
				{
					Type:    v1.NodeInternalIP,
					Address: "1.2.3.4",
				},
				{
					Type:    v1.NodeHostName,
					Address: "",
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance, _, err := initTest(testCase.testVM)
			assert.NoError(t, err)
			ret, err := instance.NodeAddressesByProviderID(context.Background(), testProviderID)
			assert.Equal(t, testCase.expectedErr, err)
			assert.Equal(t, testCase.expectedNodeAddress, ret)
		})
	}
}

func TestNodeAddressesByProviderIDInternalErr(t *testing.T) {
	testCases := []struct {
		name                string
		testVM              *vmopv1.VirtualMachine
		expectedNodeAddress []v1.NodeAddress
	}{
		{
			name:                "NodeAddressesByProviderID returns a general error for an internal error",
			testVM:              createTestVM(string(testVMName), testClusterNameSpace, testVMUUID),
			expectedNodeAddress: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance, fc, err := initTest(testCase.testVM)
			assert.NoError(t, err)
			fc.PrependReactor("list", "virtualmachines", func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, nil, fmt.Errorf("Internal error listing VMs")
			})
			ret, err := instance.NodeAddressesByProviderID(context.Background(), testProviderID)
			assert.NotEqual(t, nil, err)
			assert.NotEqual(t, cloudprovider.InstanceNotFound, err)
			assert.Equal(t, testCase.expectedNodeAddress, ret)
		})
	}
}

func TestNodeAddresses(t *testing.T) {
	testCases := []struct {
		name                string
		testVM              *vmopv1.VirtualMachine
		expectedNodeAddress []v1.NodeAddress
		expectedErr         error
	}{
		{
			name:                "NodeAddresses returns an empty address for found node with no IP",
			testVM:              createTestVM(string(testVMName), testClusterNameSpace, testVMUUID),
			expectedNodeAddress: []v1.NodeAddress{},
			expectedErr:         nil,
		},
		{
			name:                "NodeAddresses returns a NotFound error for a not found node",
			testVM:              createTestVM("bogus", testClusterNameSpace, testVMUUID),
			expectedNodeAddress: nil,
			expectedErr:         cloudprovider.InstanceNotFound,
		},
		{
			name:   "NodeAddresses returns a valid address for a found node",
			testVM: createTestVMWithVMIPAndHost(string(testVMName), testClusterNameSpace, testVMUUID),
			expectedNodeAddress: []v1.NodeAddress{
				{
					Type:    v1.NodeInternalIP,
					Address: "1.2.3.4",
				},
				{
					Type:    v1.NodeHostName,
					Address: "",
				},
			},
			expectedErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance, _, err := initTest(testCase.testVM)
			assert.NoError(t, err)
			ret, err := instance.NodeAddresses(context.Background(), testVMName)
			assert.Equal(t, testCase.expectedErr, err)
			assert.Equal(t, testCase.expectedNodeAddress, ret)
		})
	}
}

func TestNodeAddressesInternalErr(t *testing.T) {
	testCases := []struct {
		name                string
		testVM              *vmopv1.VirtualMachine
		expectedNodeAddress []v1.NodeAddress
	}{
		{
			name:                "NodeAddresses returns a general error for an internal error",
			testVM:              createTestVM(string(testVMName), testClusterNameSpace, testVMUUID),
			expectedNodeAddress: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance, fc, err := initTest(testCase.testVM)
			assert.NoError(t, err)
			fc.PrependReactor("get", "virtualmachines", func(action clientgotesting.Action) (handled bool, ret runtime.Object, err error) {
				return true, nil, fmt.Errorf("Internal error getting VMs")
			})
			ret, err := instance.NodeAddresses(context.Background(), testVMName)
			assert.NotEqual(t, nil, err)
			assert.NotEqual(t, cloudprovider.InstanceNotFound, err)
			assert.Equal(t, testCase.expectedNodeAddress, ret)
		})
	}
}
