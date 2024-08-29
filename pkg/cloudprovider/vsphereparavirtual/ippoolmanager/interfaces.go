package ippoolmanager

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	"k8s.io/cloud-provider-vsphere/pkg/cloudprovider/vsphereparavirtual/ippoolmanager/helper"
)

// IPPoolManager defines an interface that can interact with nsx.vmware.com.ippool
type IPPoolManager interface {
	GetIPPoolListerSynced() cache.InformerSynced
	GetIPPoolInformer() cache.SharedIndexInformer
	StartIPPoolInformers(stopCh <-chan struct{})

	GetIPPool(clusterNS, clusterName string) (helper.NSXIPPool, error)
	GetIPPoolFromIndexer(key string) (helper.NSXIPPool, error)
	CreateIPPool(clusterNS, clusterName string, ownerRef *metav1.OwnerReference) (helper.NSXIPPool, error)

	GetIPPoolSubnets(ippool helper.NSXIPPool) (map[string]string, error)
	AddSubnetToIPPool(node *corev1.Node, ippool helper.NSXIPPool, ownerRef *metav1.OwnerReference) error
	DeleteSubnetFromIPPool(subnetName string, ippool helper.NSXIPPool) error
	DiffIPPoolSubnets(old, cur helper.NSXIPPool) bool
}
