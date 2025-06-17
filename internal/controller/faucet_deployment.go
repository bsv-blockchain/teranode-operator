package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileDeployment is the faucet service deployment reconciler
func (r *FaucetReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	faucet := teranodev1alpha1.Faucet{}
	if err := r.Get(r.Context, r.NamespacedName, &faucet); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "faucet",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("faucet"),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &faucet)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *FaucetReconciler) updateDeployment(dep *appsv1.Deployment, faucet *teranodev1alpha1.Faucet) error {
	err := controllerutil.SetControllerReference(faucet, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultFaucetDeploymentSpec()
	// If user configures a node selector
	if faucet.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = faucet.Spec.NodeSelector
	}

	// If user configures tolerations
	if faucet.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *faucet.Spec.Tolerations
	}

	// If user configures affinity
	if faucet.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = faucet.Spec.Affinity
	}

	// if user configures resources requests
	if faucet.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *faucet.Spec.Resources
	}

	// if user configures image overrides
	if faucet.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = faucet.Spec.Image
	}
	if faucet.Spec.ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = faucet.Spec.ImagePullPolicy
	}

	// if user configures a service account
	if faucet.Spec.ServiceAccount != "" {
		dep.Spec.Template.Spec.ServiceAccountName = faucet.Spec.ServiceAccount
	}

	// if user configures a config map name
	if faucet.Spec.ConfigMapName != "" {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: faucet.Spec.ConfigMapName},
			},
		})
	}
	return nil
}

func defaultFaucetDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "faucet",
		"deployment": "faucet",
		"project":    "service",
	}
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "faucet-service",
		},
	}
	image := "foo_image"
	return &appsv1.DeploymentSpec{
		Replicas: ptr.To(int32(2)),
		Selector: metav1.SetAsLabelSelector(labels),
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.RollingUpdateDeploymentStrategyType,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.Time{},
				Labels:            labels,
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: DefaultServiceAccountName,
				Containers: []corev1.Container{
					{
						EnvFrom:         envFrom,
						Env:             env,
						Args:            []string{"-faucet=1"},
						Image:           image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "faucet",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("5Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("2"),
								corev1.ResourceMemory: resource.MustParse("5Gi"),
							},
						},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 4040,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 8097,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 9091,
								Protocol:      corev1.ProtocolTCP,
							},
						},
					},
				},
				Volumes: []corev1.Volume{},
			},
		},
	}
}
