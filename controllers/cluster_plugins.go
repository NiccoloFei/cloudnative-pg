/*
Copyright The CloudNativePG Contributors

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

// Package controllers contains the controller of the CRD
package controllers

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apiv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"github.com/cloudnative-pg/cloudnative-pg/pkg/management/log"
)

// updatePluginsStatus ensures that we load the plugins that are required to reconcile
// this cluster
func (r *ClusterReconciler) updatePluginsStatus(ctx context.Context, cluster *apiv1.Cluster) error {
	contextLogger := log.FromContext(ctx)

	// Load the plugins
	pluginClient, err := cluster.LoadPluginClient(ctx)
	if err != nil {
		contextLogger.Error(err, "Error loading plugins, retrying")
		return err
	}
	defer func() {
		pluginClient.Close(ctx)
	}()

	// Get the status of the plugins and store it inside the status section
	oldCluster := cluster.DeepCopy()
	metadataList := pluginClient.MetadataList()
	cluster.Status.PluginStatus = make([]apiv1.PluginStatus, len(metadataList))
	for i, entry := range metadataList {
		cluster.Status.PluginStatus[i].Name = entry.Name
		cluster.Status.PluginStatus[i].Version = entry.Version
		cluster.Status.PluginStatus[i].Capabilities = entry.Capabilities
		cluster.Status.PluginStatus[i].OperatorCapabilities = entry.OperatorCapabilities
		cluster.Status.PluginStatus[i].WALCapabilities = entry.WALCapabilities
		cluster.Status.PluginStatus[i].BackupCapabilities = entry.BackupCapabilities
	}

	return r.Client.Status().Patch(ctx, cluster, client.MergeFrom(oldCluster))
}

// preReconcilePluginHooks ensures we call the pre-reconcile plugin hooks
func (r *ClusterReconciler) preReconcilePluginHooks(ctx context.Context, cluster *apiv1.Cluster) (ctrl.Result, error) {
	contextLogger := log.FromContext(ctx)

	// Load the plugins
	pluginClient, err := cluster.LoadPluginClient(ctx)
	if err != nil {
		contextLogger.Error(err, "Error loading plugins, retrying")
		return ctrl.Result{}, err
	}
	defer func() {
		pluginClient.Close(ctx)
	}()

	return pluginClient.PreReconcile(ctx, cluster)
}

// postReconcilePluginHooks ensures we call the post-reconcile plugin hooks
func (r *ClusterReconciler) postReconcilePluginHooks(ctx context.Context, cluster *apiv1.Cluster) (ctrl.Result, error) {
	contextLogger := log.FromContext(ctx)

	// Load the plugins
	pluginClient, err := cluster.LoadPluginClient(ctx)
	if err != nil {
		contextLogger.Error(err, "Error loading plugins, retrying")
		return ctrl.Result{}, err
	}
	defer func() {
		pluginClient.Close(ctx)
	}()

	return pluginClient.PostReconcile(ctx, cluster)
}
