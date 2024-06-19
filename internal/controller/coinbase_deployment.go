package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileDeployment is the coinbase service deployment reconciler
func (r *CoinbaseReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	coinbase := teranodev1alpha1.Coinbase{}
	if err := r.Get(r.Context, r.NamespacedName, &coinbase); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "coinbase",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &coinbase)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *CoinbaseReconciler) updateDeployment(dep *appsv1.Deployment, coinbase *teranodev1alpha1.Coinbase) error {
	err := controllerutil.SetControllerReference(coinbase, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultCoinbaseDeploymentSpec()
	// If user configures a node selector
	if coinbase.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = coinbase.Spec.NodeSelector
	}

	// If user configures tolerations
	if coinbase.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *coinbase.Spec.Tolerations
	}

	// If user configures affinity
	if coinbase.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = coinbase.Spec.Affinity
	}

	// if user configures resources requests
	if coinbase.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *coinbase.Spec.Resources
	}

	// if user configures image overrides
	if coinbase.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = coinbase.Spec.Image
	}
	if coinbase.Spec.ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = coinbase.Spec.ImagePullPolicy
	}

	return nil
}

func defaultCoinbaseDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "coinbase",
		"deployment": "coinbase",
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
		// For now, don't override the default config
		/*{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "coinbase-config-m",
				},
			},
		},*/
	}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "coinbase-service",
		},
	}
	image := "foo_image"
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(1), // TODO: validate this number
		Selector: metav1.SetAsLabelSelector(labels),
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxUnavailable: &intstr.IntOrString{
					IntVal: int32(0),
				},
				MaxSurge: &intstr.IntOrString{
					IntVal: int32(1),
				},
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.Time{},
				Labels:            labels,
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: "sa-m",
				Containers: []corev1.Container{
					{
						EnvFrom:         envFrom,
						Env:             env,
						Args:            []string{"-coinbase=1"},
						Image:           image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "coinbase",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("28Gi"), // TODO: verify these numbers
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("8"),
								corev1.ResourceMemory: resource.MustParse("20Gi"),
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
						LivenessProbe: &corev1.Probe{
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
						},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 4040,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 8093,
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
