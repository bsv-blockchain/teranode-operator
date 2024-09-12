package utils

import (
	"github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

func SetDeploymentOverrides(dep *appsv1.Deployment, cr v1alpha1.TeranodeService) {
	// If user configures a node selector
	if cr.NodeSelector() != nil {
		dep.Spec.Template.Spec.NodeSelector = cr.NodeSelector()
	}

	// If user configures tolerations
	if cr.Tolerations() != nil {
		dep.Spec.Template.Spec.Tolerations = *cr.Tolerations()
	}

	// If user configures affinity
	if cr.Affinity() != nil {
		dep.Spec.Template.Spec.Affinity = cr.Affinity()
	}

	// if user configures resources requests
	if cr.Resources() != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *cr.Resources()
	}

	// if user configures replicas
	if cr.Replicas() != nil {
		dep.Spec.Replicas = pointer.Int32(*cr.Replicas())
	}

	// if user configures image or image pull policy
	if cr.Image() != "" {
		dep.Spec.Template.Spec.Containers[0].Image = cr.Image()
	}
	if cr.ImagePullPolicy() != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = cr.ImagePullPolicy()
	}

	// if user configures a service account
	if cr.ServiceAccountName() != "" {
		dep.Spec.Template.Spec.ServiceAccountName = cr.ServiceAccountName()
	}

	// if user configures a config map name
	if cr.ConfigMapName() != "" {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: cr.ConfigMapName()},
			},
		})
	}

	// if user configures a custom command
	if len(cr.Command()) > 0 {
		dep.Spec.Template.Spec.Containers[0].Command = cr.Command()
	}

	// if user configures custom arguments
	if len(cr.Args()) > 0 {
		dep.Spec.Template.Spec.Containers[0].Args = cr.Args()
	}
	return
}
