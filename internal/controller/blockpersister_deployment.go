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
	// If user configures a node selector
	if blockPersister.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = blockPersister.Spec.NodeSelector
	}

	// If user configures tolerations
	if blockPersister.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *blockPersister.Spec.Tolerations
	}

	// If user configures affinity
	if blockPersister.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = blockPersister.Spec.Affinity
	}

	// if user configures resources requests
	if blockPersister.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *blockPersister.Spec.Resources
	}

	// if user configures image overrides
	if blockPersister.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = blockPersister.Spec.Image
	}
	if blockPersister.Spec.ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = blockPersister.Spec.ImagePullPolicy
	}

	// if user configures a service account
	if blockPersister.Spec.ServiceAccount != "" {
		dep.Spec.Template.Spec.ServiceAccountName = blockPersister.Spec.ServiceAccount
	}

	return nil
}

func defaultBlockPersisterDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "block-persister",
		"deployment": "block-persister",
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
					Name: "block-persister-config-m",
				},
			},
		},*/
	}
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
						Args:            []string{"-blockpersister=1"},
						Image:           image,
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
								ContainerPort: 4040,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/data/subtreestore",
								Name:      "subtree-storage",
							},
							{
								MountPath: "/data/blockstore",
								Name:      "block-persister-storage",
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
					{
						Name: "block-persister-storage",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: "block-persister-storage",
							},
						},
					},
				},
			},
		},
	}
}
