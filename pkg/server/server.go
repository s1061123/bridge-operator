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
	"context"
	"fmt"
	"time"
	"github.com/s1061123/bridge-operator/pkg/controllers"
	bridgeconfv1alpha1 "github.com/s1061123/bridge-operator/pkg/apis/openshift.io/v1alpha1"
	bridgeconfclient "github.com/s1061123/bridge-operator/pkg/client/clientset/versioned"
	bridgeconfinformer "github.com/s1061123/bridge-operator/pkg/client/informers/externalversions"
	bridgeconflisterv1alpha1 "github.com/s1061123/bridge-operator/pkg/client/listers/openshift.io/v1alpha1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	/*
	discovery "k8s.io/api/discovery/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/selection"
	*/
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	//"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	/*
	corelisters "k8s.io/client-go/listers/core/v1"
	discoverylisters "k8s.io/client-go/listers/discovery/v1"
	*/	
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	//api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/util/async"
	utilnode "k8s.io/kubernetes/pkg/util/node"
	//"k8s.io/utils/exec"
)

type Server struct {
	bridgeConfChanges *controllers.BridgeConfChangeTracker
	bridgeConfMap controllers.BridgeConfMap
	bridgeConfLister bridgeconflisterv1alpha1.BridgeConfigurationLister
	bridgeConfClient bridgeconfclient.Interface

	Client clientset.Interface
	Hostname string
	hostPrefix string
	Broadcaster record.EventBroadcaster
	Recorder record.EventRecorder
	Options *Options
	ConfigSyncPeriod time.Duration
	NodeRef             *v1.ObjectReference

	bridgeConfSynced bool

	syncRunner *async.BoundedFrequencyRunner
}

func NewServer(o *Options) (*Server, error) {
	var kubeConfig *rest.Config
	var err error

	if len(o.Kubeconfig) == 0 {
		klog.Info("Neither kubeconfig file nor master URL was specified. Falling back to in-cluster config.")
		kubeConfig, err = rest.InClusterConfig()
	} else {
		kubeConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: o.Kubeconfig},
			&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: o.master}},
		).ClientConfig()
	}
	if err != nil {
		return nil, err
	}

	client, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	bridgeConfClient, err := bridgeconfclient.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}


	hostname, err := utilnode.GetHostname(o.hostnameOverride)
	if err != nil {
		return nil, err
	}

	eventBroadcaster := record.NewBroadcaster()
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "multus-proxy", Host: hostname})

	nodeRef := &v1.ObjectReference{
		Kind:      "Node",
		Name:      hostname,
		UID:       types.UID(hostname),
		Namespace: "",
	}

	syncPeriod := 3 * time.Second
	minSyncPeriod := 500 * time.Millisecond
	burstSyncs := 2

	bridgeConfChanges := controllers.NewBridgeConfChangeTracker()
	if bridgeConfChanges == nil {
		return nil, fmt.Errorf("cannot create bridgeconf change tracker")
	}

	server := &Server{
		Options: o,
		Client: client,
		bridgeConfClient: bridgeConfClient,
		Hostname: hostname,
		Broadcaster: eventBroadcaster,
		Recorder: recorder,
		NodeRef: nodeRef,
		bridgeConfChanges: bridgeConfChanges,
		bridgeConfMap: make(controllers.BridgeConfMap),
	}
	server.syncRunner = async.NewBoundedFrequencyRunner("sync-runner", server.syncBridges, minSyncPeriod, syncPeriod, burstSyncs)

	return server, nil
}

func (s *Server) OnBridgeConfigSynced() {
	klog.Info("Bridge Synced")
	if s.syncRunner != nil {
		s.syncRunner.Loop(wait.NeverStop)
		s.syncRunner.Run()
	}
}

//XXX may change from syncRunner to monitor bridges!
func (s *Server) syncBridges() {
	bridgeConfigs, err := s.bridgeConfLister.BridgeConfigurations(metav1.NamespaceAll).List(labels.Everything())
	if err != nil {
		klog.Errorf("failed to get bridgeConfig: %v", err)
		return
	}

	for _, bridgeConfig := range bridgeConfigs {
		bridgeReconcile(bridgeConfig)
	}

}

func (s *Server) Run(hostname string, stopCh chan struct{}) error {
	if s.Broadcaster != nil {
		s.Broadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: s.Client.CoreV1().Events("")})
	}

	bridgeConfInformerFactory := bridgeconfinformer.NewSharedInformerFactoryWithOptions(
		s.bridgeConfClient, s.ConfigSyncPeriod)
	s.bridgeConfLister = bridgeConfInformerFactory.OpenShiftIO().V1alpha1().BridgeConfigurations().Lister()
	bridgeConfig := controllers.NewBridgeConfConfig(
		bridgeConfInformerFactory.OpenShiftIO().V1alpha1().BridgeConfigurations(), s.ConfigSyncPeriod)
	bridgeConfig.RegisterEventHandler(s)

	go bridgeConfig.Run(wait.NeverStop)
	bridgeConfInformerFactory.Start(wait.NeverStop)

	<-stopCh

	return nil
}

func bridgeConfNamespacedName(config *bridgeconfv1alpha1.BridgeConfiguration) string {
	if config == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s/%s", config.Namespace, config.Name)
}

func bridgeConfMatchesNodeLabel(node *v1.Node, config *bridgeconfv1alpha1.BridgeConfiguration) (bool, error) {
	if config.Spec.NodeSelector.Size() != 0 {
		nodeSelectorMap, err := metav1.LabelSelectorAsMap(&config.Spec.NodeSelector)
		if err != nil {
			klog.Errorf("bad label selector for node %q: %v", bridgeConfNamespacedName(config), err)
			return false, err
		}
		bridgeNodeSelector := labels.Set(nodeSelectorMap).AsSelectorPreValidated()
		if !bridgeNodeSelector.Matches(labels.Set(node.Labels)) {
			return false, nil
		}
		return true, nil
	}

	// if no nodeSelector, it must be matched.
	return true, nil
}

func (s *Server) OnBridgeConfigAdd(config *bridgeconfv1alpha1.BridgeConfiguration) {
	node, err := s.Client.CoreV1().Nodes().Get(context.TODO(), s.Hostname, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("cannot get node: %q: %v", s.Hostname, err)
	}

	matches, err := bridgeConfMatchesNodeLabel(node, config)
	if err != nil {
		klog.Errorf("failed to get label selector: %v", err)
	}
	if matches {
		klog.Infof("node matches!")
		bridgeConfigured(config)
	}

}

func (s *Server) OnBridgeConfigUpdate(oldConfig, config *bridgeconfv1alpha1.BridgeConfiguration) {
	klog.Info("Bridge Update")
	node, err := s.Client.CoreV1().Nodes().Get(context.TODO(), s.Hostname, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("cannot get node: %q: %v", s.Hostname, err)
	}

	oldMatches, err := bridgeConfMatchesNodeLabel(node, oldConfig)
	if err != nil {
		klog.Errorf("failed to get label selector: %v", err)
	}

	matches, err := bridgeConfMatchesNodeLabel(node, oldConfig)
	if err != nil {
		klog.Errorf("failed to get label selector: %v", err)
	}

	// label selector is changed
	if matches != oldMatches {
		// false -> true: create new bridge
		if matches == true {
			bridgeConfigured(config)
		} else { // true -> false: delete new bridge
			bridgeConfigured(config)
		}
	}
	err = bridgeUpdate(oldConfig, config)
	if err != nil {
		klog.Errorf("failed to sync bridge config: %v", err)
	}
}

func (s *Server) OnBridgeConfigDelete(config *bridgeconfv1alpha1.BridgeConfiguration) {
	klog.Info("Bridge Delete")
	node, err := s.Client.CoreV1().Nodes().Get(context.TODO(), s.Hostname, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("cannot get node: %q: %v", s.Hostname, err)
	}

	matches, err := bridgeConfMatchesNodeLabel(node, config)
	if err != nil {
		klog.Errorf("failed to get label selector: %v", err)
	}
	if matches {
		klog.Infof("node matches!")
		bridgeUnconfigured(config)
	}
}

