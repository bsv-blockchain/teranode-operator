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

// ReconcileDeployment is the propagation service deployment reconciler
func (r *PropagationReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	propagation := teranodev1alpha1.Propagation{}
	if err := r.Get(r.Context, r.NamespacedName, &propagation); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "propagation",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &propagation)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PropagationReconciler) updateDeployment(dep *appsv1.Deployment, propagation *teranodev1alpha1.Propagation) error {
	err := controllerutil.SetControllerReference(propagation, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultPropagationDeploymentSpec()
	// If user configures a node selector
	if propagation.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = propagation.Spec.NodeSelector
	}

	// If user configures tolerations
	if propagation.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *propagation.Spec.Tolerations
	}

	// If user configures affinity
	if propagation.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = propagation.Spec.Affinity
	}

	// if user configures resources requests
	if propagation.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *propagation.Spec.Resources
	}

	return nil
}

func defaultPropagationDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "propagation",
		"deployment": "propagation",
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
					Name: "propagation-config-m",
				},
			},
		},*/
	}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "propagation-service",
		},
	}
	image := "foo_image"
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(2), // TODO: verify the replicas number, the spec has 28
		Selector: metav1.SetAsLabelSelector(labels),
		//Strategy: appsv1.DeploymentStrategy{ // TODO: verify if no strategy should be used by default
		//	Type: appsv1.RecreateDeploymentStrategyType,
		//},
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
						Args:            []string{"-propagation=1"},
						Image:           image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "propagation",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("58Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("28"),
								corev1.ResourceMemory: resource.MustParse("40Gi"),
							},
						},
						ReadinessProbe: &corev1.Probe{
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
								ContainerPort: 8084,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 8384,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 8833,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/app/certs",
								Name:      "scaling-tls",
								ReadOnly:  true,
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "scaling-tls",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "scaling-tls",
								Items: []corev1.KeyToPath{
									{
										Key:  "tls.crt",
										Path: "ubsv.crt",
									},
									{
										Key:  "tls.key",
										Path: "ubsv.key",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
