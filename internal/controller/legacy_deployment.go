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

// ReconcileDeployment is the legacy service deployment reconciler
func (r *LegacyReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	legacy := teranodev1alpha1.Legacy{}
	if err := r.Get(r.Context, r.NamespacedName, &legacy); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "legacy",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &legacy)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *LegacyReconciler) updateDeployment(dep *appsv1.Deployment, legacy *teranodev1alpha1.Legacy) error {
	err := controllerutil.SetControllerReference(legacy, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultLegacyDeploymentSpec()
	// If user configures a node selector
	if legacy.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = legacy.Spec.NodeSelector
	}

	// If user configures tolerations
	if legacy.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *legacy.Spec.Tolerations
	}

	// If user configures affinity
	if legacy.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = legacy.Spec.Affinity
	}

	// if user configures resources requests
	if legacy.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *legacy.Spec.Resources
	}

	// if user configures image overrides
	if legacy.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = legacy.Spec.Image
	}
	if legacy.Spec.ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = legacy.Spec.ImagePullPolicy
	}

	// if user configures a service account
	if legacy.Spec.ServiceAccount != "" {
		dep.Spec.Template.Spec.ServiceAccountName = legacy.Spec.ServiceAccount
	}

	// if user configures replicas
	if legacy.Spec.Replicas != nil {
		dep.Spec.Replicas = pointer.Int32(*legacy.Spec.Replicas)
	}

	// if user configures a config map name
	if legacy.Spec.ConfigMapName != "" {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: legacy.Spec.ConfigMapName},
			},
		})
	}

	// if user configures a custom command
	if len(legacy.Spec.Command) > 0 {
		dep.Spec.Template.Spec.Containers[0].Command = legacy.Spec.Command
	}

	// if user configures custom arguments
	if len(legacy.Spec.Args) > 0 {
		dep.Spec.Template.Spec.Containers[0].Args = legacy.Spec.Args
	}
	return nil
}

func defaultLegacyDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "legacy",
		"deployment": "legacy",
		"project":    "service",
	}
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "legacy-service",
		},
	}
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(1),
		Selector: metav1.SetAsLabelSelector(labels),
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
						Args:            []string{"-legacy=1"},
						Image:           DefaultImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "legacy",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("2"),
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
						},
						//ReadinessProbe: &corev1.Probe{ // TODO: verify the lack of /health endpoint
						//	ProbeHandler: corev1.ProbeHandler{
						//		HTTPGet: &corev1.HTTPGetAction{
						//			Path: "/health",
						//			Port: intstr.FromInt32(9091),
						//		},
						//	},
						//	InitialDelaySeconds: 1,
						//	PeriodSeconds:       10,
						//	FailureThreshold:    5,
						//	TimeoutSeconds:      3,
						//},
						//LivenessProbe: &corev1.Probe{
						//	ProbeHandler: corev1.ProbeHandler{
						//		HTTPGet: &corev1.HTTPGetAction{
						//			Path: "/health",
						//			Port: intstr.FromInt32(9091),
						//		},
						//	},
						//	InitialDelaySeconds: 1,
						//	PeriodSeconds:       10,
						//	FailureThreshold:    5,
						//	TimeoutSeconds:      3,
						//},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: DebuggerPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: ProfilerPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: LegacyHttpPort,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/data",
								Name:      SharedPVCName,
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: SharedPVCName,
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: SharedPVCName,
							},
						},
					},
				},
			},
		},
	}
}
