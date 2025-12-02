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
			Labels:    getAppLabels("coinbase"),
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
	utils.SetClusterOverrides(r.Client, dep, coinbase)
	utils.SetDeploymentOverrides(r.Client, dep, coinbase)

	return nil
}

func defaultCoinbaseDeploymentSpec() *appsv1.DeploymentSpec {
	labels := getAppLabels("coinbase")
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "coinbase-service",
		},
	}
	return &appsv1.DeploymentSpec{
		Replicas: ptr.To(int32(1)),
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
				ServiceAccountName: DefaultServiceAccountName,
				Containers: []corev1.Container{
					{
						EnvFrom:         envFrom,
						Env:             env,
						Args:            []string{"-coinbase=1"},
						Image:           DefaultCoinbaseImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "coinbase",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("4Gi"), // TODO: verify these numbers
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("1"),
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health/readiness",
									Port: intstr.FromInt32(CoinbaseHTTPPort),
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
									Port: intstr.FromInt32(CoinbaseHTTPPort),
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
									Port: intstr.FromInt32(CoinbaseHTTPPort),
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
								ContainerPort: CoinbaseGRPCPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: CoinbaseHTTPPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: CoinbaseP2PPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: ProfilerPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: HealthPort,
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
