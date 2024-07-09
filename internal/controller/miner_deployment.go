package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileDeployment is the miner service deployment reconciler
func (r *MinerReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	miner := teranodev1alpha1.Miner{}
	if err := r.Get(r.Context, r.NamespacedName, &miner); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "miner",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &miner)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *MinerReconciler) updateDeployment(dep *appsv1.Deployment, miner *teranodev1alpha1.Miner) error {
	err := controllerutil.SetControllerReference(miner, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultMinerDeploymentSpec()
	// If user configures a node selector
	if miner.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = miner.Spec.NodeSelector
	}

	// If user configures tolerations
	if miner.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *miner.Spec.Tolerations
	}

	// If user configures affinity
	if miner.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = miner.Spec.Affinity
	}

	// if user configures resources requests
	if miner.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *miner.Spec.Resources
	}

	// if user configures image overrides
	if miner.Spec.Image != "" {
		dep.Spec.Template.Spec.Containers[0].Image = miner.Spec.Image
	}
	if miner.Spec.ImagePullPolicy != "" {
		dep.Spec.Template.Spec.Containers[0].ImagePullPolicy = miner.Spec.ImagePullPolicy
	}

	// if user configures a service account
	if miner.Spec.ServiceAccount != "" {
		dep.Spec.Template.Spec.ServiceAccountName = miner.Spec.ServiceAccount
	}
	return nil
}

func defaultMinerDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "miner",
		"deployment": "miner",
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
					Name: "miner-config-m",
				},
			},
		},*/
	}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "miner-service",
		},
	}
	image := "foo_image"
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(1),
		Selector: metav1.SetAsLabelSelector(labels),
		//Strategy: appsv1.DeploymentStrategy{ // TODO: verify no default deployment strategy
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
						Args:            []string{"-miner=1"},
						Image:           image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "miner",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("2"),
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
						},
						//ReadinessProbe: &corev1.Probe{ // TODO: verify the lack of /health endpoint
						//	ProbeHandler: corev1.ProbeHandler{
						//		HTTPGet: &corev1.HTTPGetAction{
						//			Path: "/health",
						//			Port: intstr.FromInt32(9091),
						//		},
						//	},
						//	InitialDelaySeconds: 1,
						//	PeriodSeconds:       10,
						//	FailureThreshold:    5,
						//	TimeoutSeconds:      3,
						//},
						//LivenessProbe: &corev1.Probe{
						//	ProbeHandler: corev1.ProbeHandler{
						//		HTTPGet: &corev1.HTTPGetAction{
						//			Path: "/health",
						//			Port: intstr.FromInt32(9091),
						//		},
						//	},
						//	InitialDelaySeconds: 1,
						//	PeriodSeconds:       10,
						//	FailureThreshold:    5,
						//	TimeoutSeconds:      3,
						//},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 4040,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 8089,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 8092,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 8099,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: 9091,
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
