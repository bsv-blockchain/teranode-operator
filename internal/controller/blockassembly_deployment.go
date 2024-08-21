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

// ReconcileDeployment is the block assembly service deployment reconciler
func (r *BlockAssemblyReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	blockAssembly := teranodev1alpha1.BlockAssembly{}
	if err := r.Get(r.Context, r.NamespacedName, &blockAssembly); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "block-assembly",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &blockAssembly)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BlockAssemblyReconciler) updateDeployment(dep *appsv1.Deployment, blockAssembly *teranodev1alpha1.BlockAssembly) error {
	err := controllerutil.SetControllerReference(blockAssembly, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultBlockAssemblyDeploymentSpec()
	// If user configures a node selector
	if blockAssembly.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = blockAssembly.Spec.NodeSelector
	}

	// If user configures tolerations
	if blockAssembly.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *blockAssembly.Spec.Tolerations
	}

	// If user configures affinity
	if blockAssembly.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = blockAssembly.Spec.Affinity
	}

	// if user configures resources requests
	if blockAssembly.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *blockAssembly.Spec.Resources
	}

	// if user configures image overrides
	if blockAssembly.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = blockAssembly.Spec.Image
	}
	if blockAssembly.Spec.ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = blockAssembly.Spec.ImagePullPolicy
	}

	// if user configures a service account
	if blockAssembly.Spec.ServiceAccount != "" {
		dep.Spec.Template.Spec.ServiceAccountName = blockAssembly.Spec.ServiceAccount
	}

	// if user configures replicas
	if blockAssembly.Spec.Replicas != nil {
		dep.Spec.Replicas = pointer.Int32(*blockAssembly.Spec.Replicas)
	}

	// if user configures a config map name
	if blockAssembly.Spec.ConfigMapName != "" {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: blockAssembly.Spec.ConfigMapName},
			},
		})
	}

	return nil
}

func defaultBlockAssemblyDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "block-assembly",
		"deployment": "block-assembly",
		"project":    "service",
	}
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "blockassembly-service",
		},
		{
			Name:  "JAEGER_SERVICE_NAME",
			Value: "blockassembly-service",
		},
	}
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(1),
		Selector: metav1.SetAsLabelSelector(labels),
		//Strategy: appsv1.DeploymentStrategy{ // TODO: confirm the lack of defined update strategy
		//	Type: appsv1.RecreateDeploymentStrategyType,
		//},
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
						Args:            []string{"-blockassembly=1"},
						Image:           DefaultImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "block-assembly",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("1500Gi"), // TODO: these values should be reviewed
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("100"),
								corev1.ResourceMemory: resource.MustParse("1000Gi"),
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
								ContainerPort: BlockAssemblyPort,
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
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/data/",
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
