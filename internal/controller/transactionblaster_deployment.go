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

// ReconcileDeployment is the transaction blaster service deployment reconciler
func (r *TransactionBlasterReconciler) ReconcileDeployment(log logr.Logger) (bool, error) {
	transactionBlaster := teranodev1alpha1.TransactionBlaster{}
	if err := r.Get(r.Context, r.NamespacedName, &transactionBlaster); err != nil {
		return false, err
	}
	dep := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "transaction-blaster",
			Namespace: r.NamespacedName.Namespace,
			Labels:    getAppLabels(),
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, &dep, func() error {
		return r.updateDeployment(&dep, &transactionBlaster)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *TransactionBlasterReconciler) updateDeployment(dep *appsv1.Deployment, transactionBlaster *teranodev1alpha1.TransactionBlaster) error {
	err := controllerutil.SetControllerReference(transactionBlaster, dep, r.Scheme)
	if err != nil {
		return err
	}
	dep.Spec = *defaultTransactionBlasterDeploymentSpec()
	// If user configures a node selector
	if transactionBlaster.Spec.NodeSelector != nil {
		dep.Spec.Template.Spec.NodeSelector = transactionBlaster.Spec.NodeSelector
	}

	// If user configures tolerations
	if transactionBlaster.Spec.Tolerations != nil {
		dep.Spec.Template.Spec.Tolerations = *transactionBlaster.Spec.Tolerations
	}

	// If user configures affinity
	if transactionBlaster.Spec.Affinity != nil {
		dep.Spec.Template.Spec.Affinity = transactionBlaster.Spec.Affinity
	}

	// if user configures resources requests
	if transactionBlaster.Spec.Resources != nil {
		dep.Spec.Template.Spec.Containers[0].Resources = *transactionBlaster.Spec.Resources
	}

	return nil
}

func defaultTransactionBlasterDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "transaction-blaster",
		"deployment": "transaction-blaster",
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
					Name: "transaction-blaster-config-m",
				},
			},
		},*/
	}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "transaction-blaster-service",
		},
	}
	image := "foo_image"
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32Ptr(2),
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
				ServiceAccountName: "sa-m",
				Containers: []corev1.Container{
					{
						EnvFrom:         envFrom,
						Env:             env,
						Command:         []string{"./blaster.run"},
						Args:            []string{"-workers=10000"}, // TODO: this should be a configurable
						Image:           image,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "transaction-blaster",
						// Make sane defaults, and this should be configurable
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("10Gi"),
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("12"),
								corev1.ResourceMemory: resource.MustParse("10Gi"),
							},
						},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.FromInt(9091),
								},
							},
							InitialDelaySeconds: 1,
							PeriodSeconds:       10,
							FailureThreshold:    5,
							TimeoutSeconds:      3,
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.FromInt(9091),
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
							{
								ContainerPort: 9091,
								Protocol:      corev1.ProtocolTCP,
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/app/certs",
								Name:      "scaling-tls",
								ReadOnly:  true,
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "scaling-tls",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "scaling-tls",
								Items: []corev1.KeyToPath{
									{
										Key:  "tls.crt",
										Path: "ubsv.crt",
									},
									{
										Key:  "tls.key",
										Path: "ubsv.key",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
