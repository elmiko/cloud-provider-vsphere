/* Copyright © 2024 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: Apache-2.0 */

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/vmware-tanzu/nsx-operator/pkg/apis/vpc/v1alpha1"
	scheme "github.com/vmware-tanzu/nsx-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VPCNetworkConfigurationsGetter has a method to return a VPCNetworkConfigurationInterface.
// A group's client should implement this interface.
type VPCNetworkConfigurationsGetter interface {
	VPCNetworkConfigurations() VPCNetworkConfigurationInterface
}

// VPCNetworkConfigurationInterface has methods to work with VPCNetworkConfiguration resources.
type VPCNetworkConfigurationInterface interface {
	Create(ctx context.Context, vPCNetworkConfiguration *v1alpha1.VPCNetworkConfiguration, opts v1.CreateOptions) (*v1alpha1.VPCNetworkConfiguration, error)
	Update(ctx context.Context, vPCNetworkConfiguration *v1alpha1.VPCNetworkConfiguration, opts v1.UpdateOptions) (*v1alpha1.VPCNetworkConfiguration, error)
	UpdateStatus(ctx context.Context, vPCNetworkConfiguration *v1alpha1.VPCNetworkConfiguration, opts v1.UpdateOptions) (*v1alpha1.VPCNetworkConfiguration, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.VPCNetworkConfiguration, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.VPCNetworkConfigurationList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.VPCNetworkConfiguration, err error)
	VPCNetworkConfigurationExpansion
}

// vPCNetworkConfigurations implements VPCNetworkConfigurationInterface
type vPCNetworkConfigurations struct {
	client rest.Interface
}

// newVPCNetworkConfigurations returns a VPCNetworkConfigurations
func newVPCNetworkConfigurations(c *CrdV1alpha1Client) *vPCNetworkConfigurations {
	return &vPCNetworkConfigurations{
		client: c.RESTClient(),
	}
}

// Get takes name of the vPCNetworkConfiguration, and returns the corresponding vPCNetworkConfiguration object, and an error if there is any.
func (c *vPCNetworkConfigurations) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.VPCNetworkConfiguration, err error) {
	result = &v1alpha1.VPCNetworkConfiguration{}
	err = c.client.Get().
		Resource("vpcnetworkconfigurations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VPCNetworkConfigurations that match those selectors.
func (c *vPCNetworkConfigurations) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.VPCNetworkConfigurationList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.VPCNetworkConfigurationList{}
	err = c.client.Get().
		Resource("vpcnetworkconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested vPCNetworkConfigurations.
func (c *vPCNetworkConfigurations) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("vpcnetworkconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a vPCNetworkConfiguration and creates it.  Returns the server's representation of the vPCNetworkConfiguration, and an error, if there is any.
func (c *vPCNetworkConfigurations) Create(ctx context.Context, vPCNetworkConfiguration *v1alpha1.VPCNetworkConfiguration, opts v1.CreateOptions) (result *v1alpha1.VPCNetworkConfiguration, err error) {
	result = &v1alpha1.VPCNetworkConfiguration{}
	err = c.client.Post().
		Resource("vpcnetworkconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(vPCNetworkConfiguration).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a vPCNetworkConfiguration and updates it. Returns the server's representation of the vPCNetworkConfiguration, and an error, if there is any.
func (c *vPCNetworkConfigurations) Update(ctx context.Context, vPCNetworkConfiguration *v1alpha1.VPCNetworkConfiguration, opts v1.UpdateOptions) (result *v1alpha1.VPCNetworkConfiguration, err error) {
	result = &v1alpha1.VPCNetworkConfiguration{}
	err = c.client.Put().
		Resource("vpcnetworkconfigurations").
		Name(vPCNetworkConfiguration.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(vPCNetworkConfiguration).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *vPCNetworkConfigurations) UpdateStatus(ctx context.Context, vPCNetworkConfiguration *v1alpha1.VPCNetworkConfiguration, opts v1.UpdateOptions) (result *v1alpha1.VPCNetworkConfiguration, err error) {
	result = &v1alpha1.VPCNetworkConfiguration{}
	err = c.client.Put().
		Resource("vpcnetworkconfigurations").
		Name(vPCNetworkConfiguration.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(vPCNetworkConfiguration).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the vPCNetworkConfiguration and deletes it. Returns an error if one occurs.
func (c *vPCNetworkConfigurations) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("vpcnetworkconfigurations").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *vPCNetworkConfigurations) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("vpcnetworkconfigurations").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched vPCNetworkConfiguration.
func (c *vPCNetworkConfigurations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.VPCNetworkConfiguration, err error) {
	result = &v1alpha1.VPCNetworkConfiguration{}
	err = c.client.Patch(pt).
		Resource("vpcnetworkconfigurations").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
