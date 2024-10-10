package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/bitcoin-sv/teranode-operator/internal/utils"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileDeployment is the blockchain service deployment reconciler
func (r *BlockchainReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	blockchain := teranodev1alpha1.Blockchain{}
	if err := r.Get(r.Context, r.NamespacedName, &blockchain); err != nil {
		return false, err
	}
	labels := getAppLabels()
	labels["app"] = "blockchain"
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "blockchain",
			Namespace: r.NamespacedName.Namespace,
			Labels:    labels,
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &blockchain)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BlockchainReconciler) updateDeployment(dep *appsv1.Deployment, blockchain *teranodev1alpha1.Blockchain) error {
	err := controllerutil.SetControllerReference(blockchain, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultBlockchainDeploymentSpec()
	utils.SetDeploymentOverrides(dep, blockchain)

	return nil
}

func defaultBlockchainDeploymentSpec() *appsv1.DeploymentSpec {
	podLabels := map[string]string{
		"app":        "blockchain",
		"deployment": "blockchain",
		"project":    "service",
	}
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "blockchain-service",
		},
	}
	// TODO: set a default
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
											"app": "blockchain",
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
						Args:            []string{"-blockchain=1"},
						Image:           DefaultImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "blockchain",
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("2"),
								corev1.ResourceMemory: resource.MustParse("2Gi"),
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
						TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: BlockchainHTTPPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: BlockchainGRPCPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: ProfilerPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: DebuggerPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: HealthPort,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{},
					},
				},
				Volumes: []corev1.Volume{},
			},
		},
	}
}
