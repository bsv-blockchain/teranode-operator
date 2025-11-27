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

// ReconcileDeployment is the rpc service deployment reconciler
func (r *RPCReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	rpc := teranodev1alpha1.RPC{}
	if err := r.Get(r.Context, r.NamespacedName, &rpc); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rpc",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels("rpc"),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &rpc)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RPCReconciler) updateDeployment(dep *appsv1.Deployment, rpc *teranodev1alpha1.RPC) error {
	err := controllerutil.SetControllerReference(rpc, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultRPCDeploymentSpec()
	utils.SetDeploymentOverrides(r.Client, dep, rpc)
	utils.SetClusterOverrides(r.Client, dep, rpc)

	return nil
}

func defaultRPCDeploymentSpec() *appsv1.DeploymentSpec {
	labels := getAppLabels("rpc")
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "rpc-service",
		},
	}
	return &appsv1.DeploymentSpec{
		Replicas: ptr.To(int32(1)),
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
						Args:            []string{"-rpc=1"},
						Image:           DefaultImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "rpc",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("1"),
								corev1.ResourceMemory: resource.MustParse("1Gi"),
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
								ContainerPort: RPCPort,
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
