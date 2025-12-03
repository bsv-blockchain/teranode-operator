package controller

import (
	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
)

// mergeDeploymentOverrides selectively merges deployment overrides from cluster spec
// Only fields explicitly set in clusterOverrides will override the target
//
//nolint:gocognit,gocyclo // Function complexity is inherent to handling multiple override fields
func mergeDeploymentOverrides(target *teranodev1alpha1.DeploymentOverrides, clusterOverrides *teranodev1alpha1.DeploymentOverrides) {
	if clusterOverrides.Replicas != nil {
		target.Replicas = clusterOverrides.Replicas
	}
	if clusterOverrides.Image != "" {
		target.Image = clusterOverrides.Image
	}
	if clusterOverrides.ImagePullPolicy != "" {
		target.ImagePullPolicy = clusterOverrides.ImagePullPolicy
	}
	if clusterOverrides.ImagePullSecrets != nil {
		target.ImagePullSecrets = clusterOverrides.ImagePullSecrets
	}
	if clusterOverrides.ServiceAccount != "" {
		target.ServiceAccount = clusterOverrides.ServiceAccount
	}
	if clusterOverrides.ConfigMapName != "" {
		target.ConfigMapName = clusterOverrides.ConfigMapName
	}
	if clusterOverrides.Resources != nil {
		target.Resources = clusterOverrides.Resources
	}
	if clusterOverrides.NodeSelector != nil {
		target.NodeSelector = clusterOverrides.NodeSelector
	}
	if clusterOverrides.Tolerations != nil {
		target.Tolerations = clusterOverrides.Tolerations
	}
	if clusterOverrides.Affinity != nil {
		target.Affinity = clusterOverrides.Affinity
	}
	if clusterOverrides.PodAntiAffinity != nil {
		target.PodAntiAffinity = clusterOverrides.PodAntiAffinity
	}
	if clusterOverrides.Strategy != nil {
		target.Strategy = clusterOverrides.Strategy
	}
	if len(clusterOverrides.Command) > 0 {
		target.Command = clusterOverrides.Command
	}
	if len(clusterOverrides.Args) > 0 {
		target.Args = clusterOverrides.Args
	}
	if len(clusterOverrides.Env) > 0 {
		target.Env = clusterOverrides.Env
	}
	if len(clusterOverrides.EnvFrom) > 0 {
		target.EnvFrom = clusterOverrides.EnvFrom
	}
	if clusterOverrides.ServiceAnnotations != nil {
		target.ServiceAnnotations = clusterOverrides.ServiceAnnotations
	}
}
