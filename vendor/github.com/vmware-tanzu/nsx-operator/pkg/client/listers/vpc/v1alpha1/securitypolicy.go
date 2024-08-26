/* Copyright © 2024 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/vmware-tanzu/nsx-operator/pkg/apis/vpc/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// SecurityPolicyLister helps list SecurityPolicies.
// All objects returned here must be treated as read-only.
type SecurityPolicyLister interface {
	// List lists all SecurityPolicies in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.SecurityPolicy, err error)
	// SecurityPolicies returns an object that can list and get SecurityPolicies.
	SecurityPolicies(namespace string) SecurityPolicyNamespaceLister
	SecurityPolicyListerExpansion
}

// securityPolicyLister implements the SecurityPolicyLister interface.
type securityPolicyLister struct {
	indexer cache.Indexer
}

// NewSecurityPolicyLister returns a new SecurityPolicyLister.
func NewSecurityPolicyLister(indexer cache.Indexer) SecurityPolicyLister {
	return &securityPolicyLister{indexer: indexer}
}

// List lists all SecurityPolicies in the indexer.
func (s *securityPolicyLister) List(selector labels.Selector) (ret []*v1alpha1.SecurityPolicy, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.SecurityPolicy))
	})
	return ret, err
}

// SecurityPolicies returns an object that can list and get SecurityPolicies.
func (s *securityPolicyLister) SecurityPolicies(namespace string) SecurityPolicyNamespaceLister {
	return securityPolicyNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// SecurityPolicyNamespaceLister helps list and get SecurityPolicies.
// All objects returned here must be treated as read-only.
type SecurityPolicyNamespaceLister interface {
	// List lists all SecurityPolicies in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.SecurityPolicy, err error)
	// Get retrieves the SecurityPolicy from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.SecurityPolicy, error)
	SecurityPolicyNamespaceListerExpansion
}

// securityPolicyNamespaceLister implements the SecurityPolicyNamespaceLister
// interface.
type securityPolicyNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all SecurityPolicies in the indexer for a given namespace.
func (s securityPolicyNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.SecurityPolicy, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.SecurityPolicy))
	})
	return ret, err
}

// Get retrieves the SecurityPolicy from the indexer for a given namespace and name.
func (s securityPolicyNamespaceLister) Get(name string) (*v1alpha1.SecurityPolicy, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("securitypolicy"), name)
	}
	return obj.(*v1alpha1.SecurityPolicy), nil
}
