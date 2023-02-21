/*
Copyright 2022 The Kubernetes Authors.

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

package v1alpha1

import (
    //"k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    //"k8s.io/apimachinery/pkg/util/intstr"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:noStatus
// +resourceName=bridge-configurations
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BridgeConfiguration ...
type BridgeConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec BridgeConfigurationSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BridgeConfigurationList ...
type BridgeConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	
	Items []BridgeConfiguration `json:"items"`
}

// BridgeConfigurationSpec ...
type BridgeConfigurationSpec struct {
	// +optional
	NodeSelector metav1.LabelSelector `json:"nodeSelector,omitempty"`

	Name string `json:"name"`

	// +optional
	EgressVlanInterfaces []EgressVlanInterface `json:"egressVlanInterfaces,omitempty"`

	// +optional
	EgressInterfaces []EgressInterface `json:"egressInterfaces,omitempty"`
}

// EgressVlanInterface ...
type EgressVlanInterface struct {
	Name string `json:"name"`

	// Protocol should be '802.1q'(.1Q vlan) or '802.1ad' (Q-in-Q, dual tag)
	// default is '802.1q'
	Protocol string `json:"protocol"`

	Id int `json:"id"`
}

// EgressInterface
type EgressInterface struct {
	Name string `json:"name"`
}

// +genclient
// +genclient:nonNamespaced
// +resourceName=bridge-informations
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=TypeMeta

// BridgeInformation contains bridge information for each hosts.
// A BridgeInformation contains all bridges which is used by bridge-operator
// in one node.
type BridgeInformation struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status BridgeInformationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BridgeInformationList
type BridgeInformationList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []BridgeInformation `json:"items"`
}

// BridgeInformationStatus
type BridgeInformationStatus struct {
	Name string `json:"name"`
	Node string `json:"node"`
	Managed bool `json:"managed"`
	Ports []BridgeInformationPortStatus `json:"ports"`
}

// BridgeInformationPortStatus
type BridgeInformationPortStatus struct {
	Name string `json:"name"`
	Managed bool `json:"managed"`
}
