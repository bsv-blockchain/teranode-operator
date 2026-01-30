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

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
)

// deduplicateEnvVars removes duplicate environment variables, keeping the last occurrence
func deduplicateEnvVars(envVars []corev1.EnvVar) []corev1.EnvVar {
	if len(envVars) == 0 {
		return envVars
	}
	seen := make(map[string]int)
	result := make([]corev1.EnvVar, 0, len(envVars))

	// Track last index of each env var name
	for i, env := range envVars {
		seen[env.Name] = i
	}

	// Keep only the last occurrence of each name
	for i, env := range envVars {
		if seen[env.Name] == i {
			result = append(result, env)
		}
	}

	return result
}

// getEnvFromKey generates a unique key for an EnvFromSource
func getEnvFromKey(ef corev1.EnvFromSource) string {
	var (
		kind string
		name string
	)

	if ef.ConfigMapRef != nil {
		kind = "cm"
		name = ef.ConfigMapRef.Name
	} else if ef.SecretRef != nil {
		kind = "secret"
		name = ef.SecretRef.Name
	} else {
		return ""
	}

	// Include Prefix so that entries with the same source but different prefixes are distinct
	return kind + ":" + ef.Prefix + ":" + name
}

// deduplicateEnvFrom removes duplicate EnvFromSource entries based on ConfigMap/Secret names
func deduplicateEnvFrom(envFrom []corev1.EnvFromSource) []corev1.EnvFromSource {
	if len(envFrom) == 0 {
		return envFrom
	}
	seen := make(map[string]int)
	result := make([]corev1.EnvFromSource, 0, len(envFrom))

	// Track last index of each configmap/secret
	for i, ef := range envFrom {
		if key := getEnvFromKey(ef); key != "" {
			seen[key] = i
		}
	}

	// Keep only the last occurrence
	for i, ef := range envFrom {
		if key := getEnvFromKey(ef); key != "" && seen[key] == i {
			result = append(result, ef)
		}
	}

	return result
}

// deduplicateImagePullSecrets removes duplicate ImagePullSecrets
func deduplicateImagePullSecrets(secrets []corev1.LocalObjectReference) []corev1.LocalObjectReference {
	if len(secrets) == 0 {
		return secrets
	}
	seen := make(map[string]int)
	result := make([]corev1.LocalObjectReference, 0, len(secrets))

	// Track last index of each secret name
	for i, secret := range secrets {
		seen[secret.Name] = i
	}

	// Keep only the last occurrence
	for i, secret := range secrets {
		if seen[secret.Name] == i {
			result = append(result, secret)
		}
	}

	return result
}

// deduplicateVolumes removes duplicate volumes by name
func deduplicateVolumes(volumes []corev1.Volume) []corev1.Volume {
	if len(volumes) == 0 {
		return volumes
	}
	seen := make(map[string]int)
	result := make([]corev1.Volume, 0, len(volumes))

	// Track last index of each volume name
	for i, vol := range volumes {
		seen[vol.Name] = i
	}

	// Keep only the last occurrence
	for i, vol := range volumes {
		if seen[vol.Name] == i {
			result = append(result, vol)
		}
	}

	return result
}

// deduplicateVolumeMounts removes duplicate volume mounts by name and mountPath
func deduplicateVolumeMounts(mounts []corev1.VolumeMount) []corev1.VolumeMount {
	if len(mounts) == 0 {
		return mounts
	}
	seen := make(map[string]int)
	result := make([]corev1.VolumeMount, 0, len(mounts))

	// Track last index using name+mountPath as key
	for i, mount := range mounts {
		key := mount.Name + ":" + mount.MountPath
		seen[key] = i
	}

	// Keep only the last occurrence
	for i, mount := range mounts {
		key := mount.Name + ":" + mount.MountPath
		if seen[key] == i {
			result = append(result, mount)
		}
	}

	return result
}

func SetDeploymentOverrides(client client.Client, dep *appsv1.Deployment, cr v1alpha1.TeranodeService) {
	SetDeploymentOverridesWithContext(context.Background(), logr.Logger{}, client, dep, cr, "")
}

//nolint:gocognit,gocyclo // Function complexity is inherent to handling multiple override cases
func SetDeploymentOverridesWithContext(ctx context.Context, log logr.Logger, client client.Client, dep *appsv1.Deployment, cr v1alpha1.TeranodeService, crKind string) {
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
		dep.Spec.Template.Spec.Containers[0].Env = deduplicateEnvVars(
			append(dep.Spec.Template.Spec.Containers[0].Env, cr.DeploymentOverrides().Env...),
		)
	}

	// if user configures envFrom vars
	if len(cr.DeploymentOverrides().EnvFrom) > 0 {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = deduplicateEnvFrom(
			append(dep.Spec.Template.Spec.Containers[0].EnvFrom, cr.DeploymentOverrides().EnvFrom...),
		)
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
		dep.Spec.Template.Spec.ImagePullSecrets = deduplicateImagePullSecrets(
			append(dep.Spec.Template.Spec.ImagePullSecrets, *cr.DeploymentOverrides().ImagePullSecrets...),
		)
	}

	if len(cr.DeploymentOverrides().Volumes) > 0 {
		dep.Spec.Template.Spec.Volumes = deduplicateVolumes(
			append(dep.Spec.Template.Spec.Volumes, cr.DeploymentOverrides().Volumes...),
		)
	}
	if len(cr.DeploymentOverrides().VolumeMounts) > 0 {
		dep.Spec.Template.Spec.Containers[0].VolumeMounts = deduplicateVolumeMounts(
			append(dep.Spec.Template.Spec.Containers[0].VolumeMounts, cr.DeploymentOverrides().VolumeMounts...),
		)
	}
	if cr.DeploymentOverrides().Replicas != nil {
		dep.Spec.Replicas = cr.DeploymentOverrides().Replicas
	}
}

//nolint:gocognit,gocyclo // Function complexity is inherent to handling multiple override cases
func SetClusterOverrides(client client.Client, dep *appsv1.Deployment, cr v1alpha1.TeranodeService) {
	// if parent cluster CR has a configmap or env vars set, append it first
	clusterOwner := GetClusterOwner(client, context.Background(), cr.Metadata())
	if clusterOwner == nil {
		return
	}
	if clusterOwner.Spec.ConfigMapName != "" {
		exists := false
		for _, envFrom := range dep.Spec.Template.Spec.Containers[0].EnvFrom {
			if envFrom.ConfigMapRef != nil && envFrom.ConfigMapRef.Name == clusterOwner.Spec.ConfigMapName {
				// ConfigMap already present, skip adding
				exists = true
			}
		}
		if !exists {
			dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, corev1.EnvFromSource{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: clusterOwner.Spec.ConfigMapName},
				},
			})
		}
	}
	if len(clusterOwner.Spec.Env) > 0 {
		dep.Spec.Template.Spec.Containers[0].Env = deduplicateEnvVars(
			append(dep.Spec.Template.Spec.Containers[0].Env, clusterOwner.Spec.Env...),
		)
	}
	if clusterOwner.Spec.ImagePullSecrets != nil {
		if dep.Spec.Template.Spec.ImagePullSecrets == nil {
			dep.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{}
		}
		dep.Spec.Template.Spec.ImagePullSecrets = deduplicateImagePullSecrets(
			append(dep.Spec.Template.Spec.ImagePullSecrets, *clusterOwner.Spec.ImagePullSecrets...),
		)
	}
	if clusterOwner.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = clusterOwner.Spec.Image
	}
	if len(clusterOwner.Spec.EnvFrom) == 0 {
		return
	}

	// Append cluster envFrom vars, avoiding duplicates
	for _, clusterEnvFrom := range clusterOwner.Spec.EnvFrom {
		found := false
		for _, envFrom := range dep.Spec.Template.Spec.Containers[0].EnvFrom {
			if envFrom.ConfigMapRef != nil && clusterEnvFrom.ConfigMapRef != nil &&
				envFrom.ConfigMapRef.Name == clusterEnvFrom.ConfigMapRef.Name {
				// ConfigMap already present, skip adding
				found = true
			}
			if envFrom.SecretRef != nil && clusterEnvFrom.SecretRef != nil &&
				envFrom.SecretRef.Name == clusterEnvFrom.SecretRef.Name {
				// Secret already present, skip adding
				found = true
			}
		}
		if !found {
			dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, clusterEnvFrom)
		}
	}
}

// ScaleStatus defines the interface for CRs that support scale subresource
type ScaleStatus interface {
	SetReplicas(replicas int32)
	SetSelector(selector string)
}

// UpdateScaleStatus updates the scale-related status fields (replicas and selector) for a CR
// based on the actual deployment state. This is required for the /scale subresource to work properly.
func UpdateScaleStatus(ctx context.Context, c client.Client, namespace, deploymentName string, status ScaleStatus) error {
	deployment := &appsv1.Deployment{}
	err := c.Get(ctx, types.NamespacedName{
		Name:      deploymentName,
		Namespace: namespace,
	}, deployment)
	if err != nil {
		return err
	}

	// Update replicas from deployment status
	if deployment.Status.Replicas != 0 {
		status.SetReplicas(deployment.Status.Replicas)
	} else if deployment.Spec.Replicas != nil {
		status.SetReplicas(*deployment.Spec.Replicas)
	}

	// Update selector from deployment spec
	if deployment.Spec.Selector != nil {
		selector := metav1.FormatLabelSelector(deployment.Spec.Selector)
		status.SetSelector(selector)
	}

	return nil
}

// GetScaleStatusFromDeployment extracts replicas and selector from a deployment
func GetScaleStatusFromDeployment(deployment *appsv1.Deployment) (replicas int32, selector string) {
	if deployment.Status.ReadyReplicas != 0 {
		replicas = deployment.Status.ReadyReplicas
	} else if deployment.Spec.Replicas != nil {
		replicas = *deployment.Spec.Replicas
	}

	if deployment.Spec.Selector != nil {
		// Use FormatLabelSelector for proper formatting
		selector = metav1.FormatLabelSelector(deployment.Spec.Selector)
	}

	return replicas, selector
}
