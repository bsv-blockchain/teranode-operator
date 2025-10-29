/*
Copyright 2024.

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

package utils

import (
	"context"

	"github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SetDeploymentOverrides(client client.Client, dep *appsv1.Deployment, cr v1alpha1.TeranodeService) {
	if cr.DeploymentOverrides() == nil {
		return
	}
	// If user configures a node selector
	if cr.DeploymentOverrides().NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = cr.DeploymentOverrides().NodeSelector
	}

	// If user configures tolerations
	if cr.DeploymentOverrides().Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *cr.DeploymentOverrides().Tolerations
	}

	// If user configures affinity
	if cr.DeploymentOverrides().Affinity != nil {
		dep.Spec.Template.Spec.Affinity = cr.DeploymentOverrides().Affinity
	}

	// if user configures resources requests
	if cr.DeploymentOverrides().Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *cr.DeploymentOverrides().Resources
	}

	// if user configures replicas
	if cr.DeploymentOverrides().Replicas != nil {
		dep.Spec.Replicas = cr.DeploymentOverrides().Replicas
	}

	// if user configures image or image pull policy
	if cr.DeploymentOverrides().Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = cr.DeploymentOverrides().Image
	}
	if cr.DeploymentOverrides().ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = cr.DeploymentOverrides().ImagePullPolicy
	}

	// if user configures a service account
	if cr.DeploymentOverrides().ServiceAccount != "" {
		dep.Spec.Template.Spec.ServiceAccountName = cr.DeploymentOverrides().ServiceAccount
	}

	// if user configures env vars
	if len(cr.DeploymentOverrides().Env) > 0 {
		dep.Spec.Template.Spec.Containers[0].Env = append(dep.Spec.Template.Spec.Containers[0].Env, cr.DeploymentOverrides().Env...)
	}

	// if user configures envFrom vars
	if len(cr.DeploymentOverrides().EnvFrom) > 0 {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, cr.DeploymentOverrides().EnvFrom...)
	}

	// if user configures a config map for the service, append it next
	if cr.DeploymentOverrides().ConfigMapName != "" {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: cr.DeploymentOverrides().ConfigMapName},
			},
		})
	}

	// if user configures a custom command
	if len(cr.DeploymentOverrides().Command) > 0 {
		dep.Spec.Template.Spec.Containers[0].Command = cr.DeploymentOverrides().Command
	}

	// if user configures custom arguments
	if len(cr.DeploymentOverrides().Args) > 0 {
		dep.Spec.Template.Spec.Containers[0].Args = cr.DeploymentOverrides().Args
	}

	if cr.DeploymentOverrides().Strategy != nil {
		dep.Spec.Strategy = *cr.DeploymentOverrides().Strategy
	}

	if cr.DeploymentOverrides().ImagePullSecrets != nil {
		if dep.Spec.Template.Spec.ImagePullSecrets == nil {
			dep.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{}
		}
		dep.Spec.Template.Spec.ImagePullSecrets = append(dep.Spec.Template.Spec.ImagePullSecrets, *cr.DeploymentOverrides().ImagePullSecrets...)
	}
}

func SetClusterOverrides(client client.Client, dep *appsv1.Deployment, cr v1alpha1.TeranodeService) {
	// if parent cluster CR has a configmap or env vars set, append it first
	clusterOwner := GetClusterOwner(client, context.Background(), cr.Metadata())
	if clusterOwner == nil {
		return
	}
	if clusterOwner.Spec.ConfigMapName != "" {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: clusterOwner.Spec.ConfigMapName},
			},
		})
	}
	if len(clusterOwner.Spec.Env) > 0 {
		dep.Spec.Template.Spec.Containers[0].Env = append(dep.Spec.Template.Spec.Containers[0].Env, clusterOwner.Spec.Env...)
	}
	if len(clusterOwner.Spec.EnvFrom) > 0 {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, clusterOwner.Spec.EnvFrom...)
	}
	if clusterOwner.Spec.ImagePullSecrets != nil {
		if dep.Spec.Template.Spec.ImagePullSecrets == nil {
			dep.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{}
		}
		dep.Spec.Template.Spec.ImagePullSecrets = append(dep.Spec.Template.Spec.ImagePullSecrets, *clusterOwner.Spec.ImagePullSecrets...)
	}
}
