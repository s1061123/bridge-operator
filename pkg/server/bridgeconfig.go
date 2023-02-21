/*
Copyright 2022 Red Hat, Inc

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

package server

import (
	"fmt"
	//"os"
	"strings"

	bridgeconfv1alpha1 "github.com/s1061123/bridge-operator/pkg/apis/openshift.io/v1alpha1"

	"k8s.io/klog"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

func generateEgressVlanInterfaceName(egressVlan *bridgeconfv1alpha1.EgressVlanInterface) string {
	if egressVlan.Protocol == "802.1ad" {
		return fmt.Sprintf("%s.%dad", egressVlan.Name, egressVlan.Id)
	}
	return fmt.Sprintf("%s.%d", egressVlan.Name, egressVlan.Id)
}

func getBridgePorts(bridge netlink.Link) ([]netlink.Link, error) {
	bridgeIndex := bridge.Attrs().Index
	bridgePorts := []netlink.Link{}

	allLinks, err := netlink.LinkList()
	if err != nil {
		return nil, nil
	}
	for _, v := range allLinks {
		if v.Attrs().MasterIndex == bridgeIndex {
			bridgePorts = append(bridgePorts, v)
		}
	}
	return bridgePorts, nil
}

func getPortsNameList (links []netlink.Link) string {
	nameList := []string{}

	for _, v := range links {
		nameList = append(nameList, v.Attrs().Name)
	}
	return strings.Join(nameList, ",")
}

func bridgeConfigured(config *bridgeconfv1alpha1.BridgeConfiguration) error {
	boolTrue := true
	bridge := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name: config.Name,
			TxQLen: 1000,
		},
		VlanFiltering: &boolTrue,
	}

	if err := netlink.LinkAdd(bridge); err != nil {
		klog.Errorf("failed to create bridge: %q: %v", config.Name, err)
		return err
	}
	if err := netlink.LinkSetUp(bridge); err != nil {
		klog.Errorf("failed to set bridge link up: %q: %v", config.Name, err)
		return err
	}

	return bridgeUpdate(nil, config)
}

func bridgeUnconfigured(config *bridgeconfv1alpha1.BridgeConfiguration) error {
	bridge, err := netlink.LinkByName(config.Name)
	if err != nil {
		klog.Errorf("failed to get bridge: %q: %v", config.Name, err)
		return err
	}

	//remove all bridge link, managed by bridge-operator
	err = egressInterfaceReconcile(bridge, config.Spec.EgressInterfaces, nil)
	if err != nil {
		klog.Errorf("failed to reconcile egress interfaces: %v", err)
		return err
	}

	err = egressVlanInterfaceReconcile(bridge, config.Spec.EgressVlanInterfaces,  nil)
	if err != nil {
		klog.Errorf("failed to reconcile egress vlan interfaces: %v", err)
		return err
	}

	//XXX: this code could be in validation webhook, or finalizer mechanisms!!
	/*
	// check if some bridge ports exist
	bridgePorts, err := getBridgePorts(bridge) 
	if err != nil {
		klog.Errorf("failed to get bridge ports: %v", err)
		return err
	}
	if len(bridgePorts) != 0 {
		portsName := getPortsNameList(bridgePorts)
		klog.Errorf("failed to delete the bridge because port exist: %s", portsName)
		return fmt.Errorf("failed to delete the bridge because ports exist: %s", portsName)
	}
	*/

	err = netlink.LinkDel(bridge)
	if err != nil {
		klog.Errorf("failed to delete bridge: %q: %v", config.Name, err)
		return err
	}

	return nil
}

func egressInterfaceReconcile(bridge netlink.Link, oldList, newList []bridgeconfv1alpha1.EgressInterface) error {
	var tobeRemove, tobeAdd []bridgeconfv1alpha1.EgressInterface

	for _, oldItem := range oldList {
		found := false
		for _, newItem := range newList {
			if oldItem == newItem {
				found = true
			}
		}
		if !found {
			tobeRemove = append(tobeRemove, oldItem)
		}
	}

	for _, newItem := range newList {
		found := false
		for _, oldItem := range oldList {
			if oldItem == newItem {
				found = true
			}
		}
		if !found {
			tobeAdd = append(tobeAdd, newItem)
		}
	}

	//XXX: need to gather error and return it as one error
	for _, v := range tobeRemove {
		link, err := netlink.LinkByName(v.Name)
		if err != nil {
			klog.Errorf("failed to get interface: %q: %v", v.Name, err)
			continue
		}

		if link.Attrs().MasterIndex == 0 {
			klog.Errorf("failed to unset interface %q master due to no master: %v", v.Name, err)
			continue
		}
		if err := netlink.LinkSetNoMaster(link); err != nil {
			klog.Errorf("failed to unset interface %q from bridge %q: %v", v.Name, bridge.Attrs().Name, err)
			continue
		}
	}

	for _, v := range tobeAdd {
		link, err := netlink.LinkByName(v.Name)
		if err != nil {
			klog.Errorf("failed to get interface: %q: %v", v.Name, err)
			continue
		}

		if link.Attrs().MasterIndex != 0 {
			klog.Errorf("failed to set interface %q master because this is already in master index %q: %v", v.Name, link.Attrs().MasterIndex, err)
			continue
		}

		if link.Attrs().MasterIndex == bridge.Attrs().Index {
			klog.Errorf("interface %q is already in bridge %q", v.Name, bridge.Attrs().Name)
			continue
		}
		if err := netlink.LinkSetMaster(link, bridge); err != nil {
			klog.Errorf("failed to set interface %q to bridge %q: %v", v.Name, bridge.Attrs().Name, err)
			continue
		}
	}
	return nil
}

func egressVlanInterfaceReconcile(bridge netlink.Link, oldList, newList []bridgeconfv1alpha1.EgressVlanInterface) error {
	var tobeRemove, tobeAdd []bridgeconfv1alpha1.EgressVlanInterface

	for _, oldItem := range oldList {
		found := false
		for _, newItem := range newList {
			if oldItem == newItem {
				found = true
			}
		}
		if !found {
			tobeRemove = append(tobeRemove, oldItem)
		}
	}

	for _, newItem := range newList {
		found := false
		for _, oldItem := range oldList {
			if oldItem == newItem {
				found = true
			}
		}
		if !found {
			tobeAdd = append(tobeAdd, newItem)
		}
	}

	//XXX: need to gather error and return it as one error
	for _, v := range tobeRemove {
		ifName := generateEgressVlanInterfaceName(&v)
		link, err := netlink.LinkByName(ifName)
		if err != nil {
			klog.Errorf("failed to get interface: %q: %v", ifName, err)
			continue
		}

		if link.Attrs().MasterIndex == 0 {
			klog.Errorf("failed to unset interface %q master due to no master: %v", v.Name, err)
			continue
		}
		if err := netlink.LinkSetNoMaster(link); err != nil {
			klog.Errorf("failed to unset interface %q from bridge %q: %v", v.Name, bridge.Attrs().Name, err)
			continue
		}

		if err := netlink.LinkDel(link); err != nil {
			klog.Errorf("failed to delete vlan interface: %q: %v", ifName, err)
			continue
		}
	}
	for _, v := range tobeAdd {
		parentIntf, err := netlink.LinkByName(v.Name)
		if err != nil {
			klog.Errorf("failed to get interface: %q: %v", v.Name, err)
			continue
		}

		ifName := generateEgressVlanInterfaceName(&v)
		link, err := netlink.LinkByName(ifName)

		// create vlan interface if not existed
		if _, ok := err.(netlink.LinkNotFoundError); ok {
			vlanProtocol := netlink.VLAN_PROTOCOL_8021Q
			if v.Protocol == "802.1ad" {
				vlanProtocol = netlink.VLAN_PROTOCOL_8021AD
			}
			vlanIntf := &netlink.Vlan{
				LinkAttrs: netlink.LinkAttrs{
					Name: ifName,
					ParentIndex: parentIntf.Attrs().Index,
					TxQLen: 1000,
				},
				VlanId: v.Id,
				VlanProtocol: vlanProtocol,
			}

			if err := netlink.LinkAdd(vlanIntf); err != nil {
				klog.Errorf("failed to create vlan interface %q: %v", ifName, err)
				continue
			}
			if err := netlink.LinkSetUp(vlanIntf); err != nil {
				klog.Errorf("failed to set vlan interface up: %q: %v", ifName, err)
				continue
			}
			link = netlink.Link(vlanIntf)
		}

		if link.Attrs().MasterIndex == bridge.Attrs().Index {
			klog.Errorf("interface %q is already in bridge %q", v.Name, bridge.Attrs().Name)
			continue
		}
		if err := netlink.LinkSetMaster(link, bridge); err != nil {
			klog.Errorf("failed to set interface %q to bridge %q: %v", v.Name, bridge.Attrs().Name, err)
			continue
		}
	}
	return nil
}

func bridgeVlanIDReconcile(link netlink.Link, vlanInfo map[int32][]*nl.BridgeVlanInfo, newVlanIDs map[uint16]struct{}) {
	for _, vlanInfo := range vlanInfo[int32(link.Attrs().Index)]{
		_, ok := newVlanIDs[vlanInfo.Vid]
		// delete vlan if not in newVlanIDs
		if !ok {
			err := netlink.BridgeVlanDel(link, vlanInfo.Vid, vlanInfo.PortVID(), vlanInfo.EngressUntag(), false, true)
			if err !=  nil {
				klog.Errorf("failed to remove bridge vlan on %s: %v\n", link.Attrs().Name, err)
			}
		} else {
			// if exists, skip it from newVlanIDs
			delete(newVlanIDs, vlanInfo.Vid)
		}
	}

	for vid, _ := range newVlanIDs {
		err := netlink.BridgeVlanAdd(link, vid, false, false, false, true)
		if err !=  nil {
			klog.Errorf("failed to add bridge vlan on %s: %v\n", link.Attrs().Name, err)
		}
	}
}

func bridgeVlanFilteringReconcile(config *bridgeconfv1alpha1.BridgeConfiguration, bridge netlink.Link) error {
	// configure vlan_filtering
	err := netlink.BridgeSetVlanFiltering(bridge, true)
	if err != nil {
		return fmt.Errorf("failed to set bridge vlan_filtering: %v", err)
	}

	// retrieve vlan information for all ports
	vlanInfo, err := netlink.BridgeVlanList()
	if err != nil {
		return fmt.Errorf("failed to get bridge vlan info: %v", err)
	}

	egressInterfaceNameSet := make(map[string]struct{})

	for _, v := range config.Spec.EgressInterfaces {
		egressInterfaceNameSet[v.Name] = struct{}{}
	}
	for _, v := range config.Spec.EgressVlanInterfaces {
		egressInterfaceNameSet[generateEgressVlanInterfaceName(&v)] = struct{}{}
	}

	// get veth's vlan information as list
	bridgePorts, err := getBridgePorts(bridge)

	vlanIDSet := make(map[uint16]struct{})
	// get vlan IDs from other interfaces (e.g. veth)
	for _, port := range bridgePorts {
		if _, ok := egressInterfaceNameSet[port.Attrs().Name]; !ok {
			for _, vlanInfo := range vlanInfo[int32(port.Attrs().Index)] {
				if vlanInfo.PortVID() {
					vlanIDSet[vlanInfo.Vid] = struct{}{}
				}
			}
		}
	}

	for _, port := range bridgePorts {
		portName := port.Attrs().Name
		if _, ok := egressInterfaceNameSet[portName]; ok {
			bridgeVlanIDReconcile(port, vlanInfo, vlanIDSet)
		}
	}

	return nil
}

func bridgeReconcile(config *bridgeconfv1alpha1.BridgeConfiguration) {
	bridge, err := netlink.LinkByName(config.Name)
	if err != nil {
		klog.Errorf("failed to get bridge %s: %v", config.Name, err)
	}

	bridgeVlanFilteringReconcile(config, bridge)
}

func bridgeUpdate(oldConfig, config *bridgeconfv1alpha1.BridgeConfiguration) error {
	bridge, err := netlink.LinkByName(config.Name)
	if err != nil {
		klog.Errorf("failed to get bridge: %q: %v", config.Name, err)
		return err
	}

	// egress interface
	oldList := []bridgeconfv1alpha1.EgressInterface(nil)
	if oldConfig != nil {
		oldList = oldConfig.Spec.EgressInterfaces
	}
	err = egressInterfaceReconcile(bridge, oldList, config.Spec.EgressInterfaces)
	if err != nil {
		klog.Errorf("failed to reconcile egress interfaces: %v", err)
	}

	// egressvlan interface
	oldVlanList := []bridgeconfv1alpha1.EgressVlanInterface(nil)
	if oldConfig != nil {
		oldVlanList = oldConfig.Spec.EgressVlanInterfaces
	}
	err = egressVlanInterfaceReconcile(bridge, oldVlanList, config.Spec.EgressVlanInterfaces)
	if err != nil {
		klog.Errorf("failed to reconcile egress vlan interfaces: %v", err)
	}

	// vlan_filtering
	bridgeVlanFilteringReconcile(config, bridge)

	return nil
}
