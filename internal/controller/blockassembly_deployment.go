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
	BlockAssembly := teranodev1alpha1.BlockAssembly{}
	if err := r.Get(r.Context, r.NamespacedName, &BlockAssembly); err != nil {
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
		return r.updateDeployment(&dep, &BlockAssembly)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BlockAssemblyReconciler) updateDeployment(dep *appsv1.Deployment, BlockAssembly *teranodev1alpha1.BlockAssembly) error {
	err := controllerutil.SetControllerReference(BlockAssembly, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultBlockAssemblyDeploymentSpec()
	// If user configures a node selector
	if BlockAssembly.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = BlockAssembly.Spec.NodeSelector
	}

	// If user configures tolerations
	if BlockAssembly.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *BlockAssembly.Spec.Tolerations
	}

	// If user configures affinity
	if BlockAssembly.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = BlockAssembly.Spec.Affinity
	}

	// if user configures resources requests
	if BlockAssembly.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *BlockAssembly.Spec.Resources
	}

	// if user configures image overrides
	if BlockAssembly.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = BlockAssembly.Spec.Image
	}
	if BlockAssembly.Spec.ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = BlockAssembly.Spec.ImagePullPolicy
	}

	// if user configures a service account
	if BlockAssembly.Spec.ServiceAccount != "" {
		dep.Spec.Template.Spec.ServiceAccountName = BlockAssembly.Spec.ServiceAccount
	}

	return nil
}

func defaultBlockAssemblyDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "block-assembly",
		"deployment": "block-assembly",
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
	}
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
	image := "foo_image"
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(2),
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
						Image:           image,
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
								ContainerPort: 8085,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 9091,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 4040,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/data/",
								Name:      "subtree-storage",
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "subtree-storage",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: "subtree-storage",
							},
						},
					},
				},
			},
		},
	}
}
