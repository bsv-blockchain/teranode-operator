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

// ReconcileDeployment is the bootstrap service deployment reconciler
func (r *BootstrapReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	bs := teranodev1alpha1.Bootstrap{}
	if err := r.Get(r.Context, r.NamespacedName, &bs); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bootstrap",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &bs)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BootstrapReconciler) updateDeployment(dep *appsv1.Deployment, bs *teranodev1alpha1.Bootstrap) error {
	err := controllerutil.SetControllerReference(bs, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultBootstrapDeploymentSpec()
	if bs.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = bs.Spec.Image
	}
	if bs.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *bs.Spec.Resources
	}
	if bs.Spec.ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = bs.Spec.ImagePullPolicy
	}

	// if user configures replicas
	if bs.Spec.Replicas != nil {
		dep.Spec.Replicas = pointer.Int32(*bs.Spec.Replicas)
	}

	// if user configures a service account
	if bs.Spec.ServiceAccount != "" {
		dep.Spec.Template.Spec.ServiceAccountName = bs.Spec.ServiceAccount
	}

	// if user configures a config map name
	if bs.Spec.ConfigMapName != "" {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: bs.Spec.ConfigMapName},
			},
		})
	}
	return nil
}

func defaultBootstrapDeploymentSpec() *appsv1.DeploymentSpec {
	podLabels := map[string]string{
		"app":        "bootstrap",
		"deployment": "bootstrap",
		"project":    "service",
	}
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "bootstrap-service",
		},
		{
			Name:  "startBootstrap",
			Value: "true",
		},
	}
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(1),
		Selector: metav1.SetAsLabelSelector(podLabels),
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxUnavailable: &intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 0,
				},
				MaxSurge: &intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 1,
				},
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.Time{},
				Labels:            podLabels,
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: DefaultServiceAccountName,
				Affinity: &corev1.Affinity{
					PodAntiAffinity: &corev1.PodAntiAffinity{
						PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
							{
								PodAffinityTerm: corev1.PodAffinityTerm{
									LabelSelector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											"app": "bootstrap",
										},
									},
									TopologyKey: "kubernetes.io/hostname",
								},
								Weight: 100,
							},
						},
					},
				},
				Containers: []corev1.Container{
					{
						EnvFrom:         envFrom,
						Env:             env,
						Args:            []string{},
						Image:           DefaultImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "bootstrap",
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("4Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("2"),
								corev1.ResourceMemory: resource.MustParse("4Gi"),
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
						TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: DebuggerPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: BootstrapGRPCPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: BootstrapHTTPPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: ProfilerPort,
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
