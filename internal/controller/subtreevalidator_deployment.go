package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileDeployment is the subtreeValidator service deployment reconciler
func (r *SubtreeValidatorReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	subtreeValidator := teranodev1alpha1.SubtreeValidator{}
	if err := r.Get(r.Context, r.NamespacedName, &subtreeValidator); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "subtree-validator",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &subtreeValidator)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *SubtreeValidatorReconciler) updateDeployment(dep *appsv1.Deployment, subtreeValidator *teranodev1alpha1.SubtreeValidator) error {
	err := controllerutil.SetControllerReference(subtreeValidator, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultSubtreeValidatorDeploymentSpec()
	// If user configures a node selector
	if subtreeValidator.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = subtreeValidator.Spec.NodeSelector
	}

	// If user configures tolerations
	if subtreeValidator.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *subtreeValidator.Spec.Tolerations
	}

	// If user configures affinity
	if subtreeValidator.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = subtreeValidator.Spec.Affinity
	}

	// if user configures resources requests
	if subtreeValidator.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *subtreeValidator.Spec.Resources
	}

	// if user configures pod template annotations
	if subtreeValidator.Spec.PodTemplateAnnotations != nil {
		dep.Spec.Template.Annotations = subtreeValidator.Spec.PodTemplateAnnotations
	}

	// if user configures image overrides
	if subtreeValidator.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = subtreeValidator.Spec.Image
	}
	if subtreeValidator.Spec.ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = subtreeValidator.Spec.ImagePullPolicy
	}

	// if user configures a service account
	if subtreeValidator.Spec.ServiceAccount != "" {
		dep.Spec.Template.Spec.ServiceAccountName = subtreeValidator.Spec.ServiceAccount
	}
	return nil
}

func defaultSubtreeValidatorDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "subtree-validator",
		"deployment": "subtree-validator",
		"project":    "service",
	}
	envFrom := []corev1.EnvFromSource{
		{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "shared-config-m",
				},
			},
		},
	}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "subtreevalidation-service",
		},
	}
	image := "foo_image"
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(2),
		Selector: metav1.SetAsLabelSelector(labels),
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.RecreateDeploymentStrategyType,
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
						Args:            []string{"-subtreevalidation=1"},
						Image:           image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "subtree-validator",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("370Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("40"),
								corev1.ResourceMemory: resource.MustParse("250Gi"),
							},
						},
						/*ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.FromInt32(9091),
								},
							},
							InitialDelaySeconds: 1,
							PeriodSeconds:       10,
							FailureThreshold:    5,
							TimeoutSeconds:      3,
						},*/
						/*LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.FromInt32(9091),
								},
							},
							InitialDelaySeconds: 1,
							PeriodSeconds:       10,
							FailureThreshold:    5,
							TimeoutSeconds:      3,
						},*/
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 4040,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 8086,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 9091,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/data/subtreestore",
								Name:      "subtree-storage",
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "subtree-storage",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: "subtree-storage",
							},
						},
					},
				},
			},
		},
	}
}
