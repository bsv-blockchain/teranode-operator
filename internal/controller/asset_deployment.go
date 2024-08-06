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

// ReconcileDeployment is the asset service deployment reconciler
func (r *AssetReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	asset := teranodev1alpha1.Asset{}
	if err := r.Get(r.Context, r.NamespacedName, &asset); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "asset",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &asset)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *AssetReconciler) updateDeployment(dep *appsv1.Deployment, asset *teranodev1alpha1.Asset) error {
	err := controllerutil.SetControllerReference(asset, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultAssetDeploymentSpec()
	// If user configures a node selector
	if asset.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = asset.Spec.NodeSelector
	}

	// If user configures tolerations
	if asset.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *asset.Spec.Tolerations
	}

	// If user configures affinity
	if asset.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = asset.Spec.Affinity
	}

	// if user configures resources requests
	if asset.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *asset.Spec.Resources
	}

	// if user configures image or image pull policy
	if asset.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = asset.Spec.Image
	}
	if asset.Spec.ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = asset.Spec.ImagePullPolicy
	}

	// if user configures a service account
	if asset.Spec.ServiceAccount != "" {
		dep.Spec.Template.Spec.ServiceAccountName = asset.Spec.ServiceAccount
	}

	// if user configures a config map name
	if asset.Spec.ConfigMapName != "" {
		dep.Spec.Template.Spec.Containers[0].EnvFrom = append(dep.Spec.Template.Spec.Containers[0].EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: asset.Spec.ConfigMapName},
			},
		})
	}

	return nil
}

func defaultAssetDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "asset",
		"deployment": "asset",
		"project":    "service",
	}
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "asset-service",
		},
	}
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(2),
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
						Args:            []string{"-asset=1"},
						Image:           DefaultImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "asset",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("28Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("8"),
								corev1.ResourceMemory: resource.MustParse("20Gi"),
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
								ContainerPort: DebuggerPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: AssetHTTPPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: AssetGRPCPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: ProfilerPort,
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
