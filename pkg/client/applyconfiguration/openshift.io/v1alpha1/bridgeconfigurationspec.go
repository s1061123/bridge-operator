/*
Copyright 2023 Red Hat, Inc

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BridgeConfigurationSpecApplyConfiguration represents an declarative configuration of the BridgeConfigurationSpec type for use
// with apply.
type BridgeConfigurationSpecApplyConfiguration struct {
	NodeSelector         *v1.LabelSelector                       `json:"nodeSelector,omitempty"`
	Name                 *string                                 `json:"name,omitempty"`
	EgressVlanInterfaces []EgressVlanInterfaceApplyConfiguration `json:"egressVlanInterfaces,omitempty"`
	EgressInterfaces     []EgressInterfaceApplyConfiguration     `json:"egressInterfaces,omitempty"`
}

// BridgeConfigurationSpecApplyConfiguration constructs an declarative configuration of the BridgeConfigurationSpec type for use with
// apply.
func BridgeConfigurationSpec() *BridgeConfigurationSpecApplyConfiguration {
	return &BridgeConfigurationSpecApplyConfiguration{}
}

// WithNodeSelector sets the NodeSelector field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NodeSelector field is set to the value of the last call.
func (b *BridgeConfigurationSpecApplyConfiguration) WithNodeSelector(value v1.LabelSelector) *BridgeConfigurationSpecApplyConfiguration {
	b.NodeSelector = &value
	return b
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *BridgeConfigurationSpecApplyConfiguration) WithName(value string) *BridgeConfigurationSpecApplyConfiguration {
	b.Name = &value
	return b
}

// WithEgressVlanInterfaces adds the given value to the EgressVlanInterfaces field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the EgressVlanInterfaces field.
func (b *BridgeConfigurationSpecApplyConfiguration) WithEgressVlanInterfaces(values ...*EgressVlanInterfaceApplyConfiguration) *BridgeConfigurationSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithEgressVlanInterfaces")
		}
		b.EgressVlanInterfaces = append(b.EgressVlanInterfaces, *values[i])
	}
	return b
}

// WithEgressInterfaces adds the given value to the EgressInterfaces field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the EgressInterfaces field.
func (b *BridgeConfigurationSpecApplyConfiguration) WithEgressInterfaces(values ...*EgressInterfaceApplyConfiguration) *BridgeConfigurationSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithEgressInterfaces")
		}
		b.EgressInterfaces = append(b.EgressInterfaces, *values[i])
	}
	return b
}