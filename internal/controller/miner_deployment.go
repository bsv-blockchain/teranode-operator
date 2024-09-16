package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/bitcoin-sv/teranode-operator/internal/utils"
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
	utils.SetDeploymentOverrides(dep, miner)

	return nil
}

func defaultMinerDeploymentSpec() *appsv1.DeploymentSpec {
	labels := map[string]string{
		"app":        "miner",
		"deployment": "miner",
		"project":    "service",
	}
	envFrom := []corev1.EnvFromSource{}
	env := []corev1.EnvVar{
		{
			Name:  "SERVICE_NAME",
			Value: "miner-service",
		},
	}
	return &appsv1.DeploymentSpec{
		Replicas: pointer.Int32(1),
		Selector: metav1.SetAsLabelSelector(labels),
		//Strategy: appsv1.DeploymentStrategy{ // TODO: verify no default deployment strategy
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
						Args:            []string{"-miner=1"},
						Image:           DefaultImage,
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
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: DebuggerPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: BootstrapGRPCPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: MinerHTTPPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: BootstrapHTTPPort,
								Protocol:      corev1.ProtocolTCP,
							},
							{
								ContainerPort: ProfilerPort,
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
