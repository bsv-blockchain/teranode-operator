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

// ReconcileDeployment is the block persister service deployment reconciler
func (r *BlockPersisterReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	blockPersister := teranodev1alpha1.BlockPersister{}
	if err := r.Get(r.Context, r.NamespacedName, &blockPersister); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "block-persister",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &blockPersister)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BlockPersisterReconciler) updateDeployment(dep *appsv1.Deployment, blockPersister *teranodev1alpha1.BlockPersister) error {
	err := controllerutil.SetControllerReference(blockPersister, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultBlockPersisterDeploymentSpec()
	utils.SetDeploymentOverrides(dep, blockPersister)

	return nil
}

func defaultBlockPersisterDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "block-persister",
		"deployment": "block-persister",
		"project":    "service",
	}
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "block-persister-service",
		},
		{
			Name:  "blockPersister_groupLimit", // TODO: this value should be reviewed for what it is, why 100, etc.
			Value: "100",
		},
	}
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(2),
		Selector: metav1.SetAsLabelSelector(labels),
		//Strategy: appsv1.DeploymentStrategy{ // TODO: confirm the lack of defined update strategy
		//	Type: appsv1.RecreateDeploymentStrategyType,
		// },
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
						Args:            []string{"-blockpersister=1"},
						Image:           DefaultImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "block-persister",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("80Gi"), // TODO: these values should be reviewed
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("25"),
								corev1.ResourceMemory: resource.MustParse("80Gi"),
							},
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health/readiness",
									Port: intstr.FromInt32(8000),
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
									Port: intstr.FromInt32(8000),
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
									Port: intstr.FromInt32(8000),
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
