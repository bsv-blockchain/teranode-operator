package controller

import (
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	"github.com/bsv-blockchain/teranode-operator/internal/utils"
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
			Labels:    getAppLabels("legacy"),
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
	utils.SetDeploymentOverrides(r.Client, dep, legacy)
	utils.SetClusterOverrides(r.Client, dep, legacy)

	return nil
}

func defaultLegacyDeploymentSpec() *appsv1.DeploymentSpec {
	labels := getAppLabels("legacy")
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "legacy-service",
		},
	}
	return &appsv1.DeploymentSpec{
		Replicas: ptr.To(int32(1)),
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
						Args:            []string{"-legacy=1"},
						Image:           DefaultImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "legacy",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("8Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("2"),
								corev1.ResourceMemory: resource.MustParse("8Gi"),
							},
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health/readiness",
									Port: intstr.FromInt32(HealthPort),
								},
							},
							InitialDelaySeconds: 1,
							PeriodSeconds:       5,
							FailureThreshold:    2,
							TimeoutSeconds:      3,
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health/liveness",
									Port: intstr.FromInt32(HealthPort),
								},
							},
							InitialDelaySeconds: 1,
							PeriodSeconds:       5,
							FailureThreshold:    2,
							TimeoutSeconds:      3,
						},
						StartupProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health/readiness",
									Port: intstr.FromInt32(HealthPort),
								},
							},
							FailureThreshold: 30,
							PeriodSeconds:    10,
						},
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
								ContainerPort: LegacyHTTPPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: HealthPort,
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
