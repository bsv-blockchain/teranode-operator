package controller

import (
	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// ReconcileGrpcIngress is the ingress for the peer grpc server
func (r *PeerReconciler) ReconcileGrpcIngress(log logr.Logger) (bool, error) {

	peer := teranodev1alpha1.Peer{}
	if err := r.Get(r.Context, r.NamespacedName, &peer); err != nil {
		return false, err
	}
	// Skip if GrpcIngress isn't set
	if peer.Spec.GrpcIngress == nil {
		return false, nil
	}
	labels := getAppLabels()
	prefix := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "peer-grpc",
			Namespace:   r.NamespacedName.Namespace,
			Annotations: peer.Spec.GrpcIngress.Annotations,
			Labels:      labels,
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: peer.Spec.GrpcIngress.Host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: "peer",
											Port: v1.ServiceBackendPort{
												Name: "p2p",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateGrpcIngress(ingress, &peer)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PeerReconciler) updateGrpcIngress(ingress *v1.Ingress, peer *teranodev1alpha1.Peer) error {
	err := controllerutil.SetControllerReference(peer, ingress, r.Scheme)
	if err != nil {
		return err
	}
	if peer.Spec.GrpcIngress == nil {
		return nil
	}
	if peer.Spec.GrpcIngress.Annotations == nil {
		ingress.Annotations = peer.Spec.GrpcIngress.Annotations
	}
	if peer.Spec.GrpcIngress.ClassName != nil {
		ingress.Spec.IngressClassName = peer.Spec.GrpcIngress.ClassName
	}

	return nil
}

// ReconcileWsIngress is the ingress for the peer ws server
func (r *PeerReconciler) ReconcileWsIngress(log logr.Logger) (bool, error) {

	peer := teranodev1alpha1.Peer{}
	if err := r.Get(r.Context, r.NamespacedName, &peer); err != nil {
		return false, err
	}
	// Skip if WsIngress isn't set
	if peer.Spec.WsIngress == nil {
		return false, nil
	}
	labels := getAppLabels()
	prefix := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "peer-ws",
			Namespace:   r.NamespacedName.Namespace,
			Annotations: peer.Spec.WsIngress.Annotations,
			Labels:      labels,
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: peer.Spec.GrpcIngress.Host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: "asset",
											Port: v1.ServiceBackendPort{
												Name: "asset-http",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateWsIngress(ingress, &peer)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PeerReconciler) updateWsIngress(ingress *v1.Ingress, peer *teranodev1alpha1.Peer) error {
	err := controllerutil.SetControllerReference(peer, ingress, r.Scheme)
	if err != nil {
		return err
	}
	if peer.Spec.WsIngress == nil {
		return nil
	}
	if peer.Spec.WsIngress.Annotations == nil {
		ingress.Annotations = peer.Spec.WsIngress.Annotations
	}
	if peer.Spec.WsIngress.ClassName != nil {
		ingress.Spec.IngressClassName = peer.Spec.WsIngress.ClassName
	}

	return nil
}

// ReconcileWssIngress is the ingress for the peer wss server
func (r *PeerReconciler) ReconcileWssIngress(log logr.Logger) (bool, error) {

	peer := teranodev1alpha1.Peer{}
	if err := r.Get(r.Context, r.NamespacedName, &peer); err != nil {
		return false, err
	}
	// Skip if domain isn't set
	if peer.Spec.WssIngress == nil {
		return false, nil
	}

	labels := getAppLabels()
	prefix := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "peer-wss",
			Namespace:   r.NamespacedName.Namespace,
			Annotations: peer.Spec.WssIngress.Annotations,
			Labels:      labels,
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: peer.Spec.GrpcIngress.Host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: "asset",
											Port: v1.ServiceBackendPort{
												Name: "asset-http",
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: peer.Spec.GrpcIngress.Host,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: "peer",
											Port: v1.ServiceBackendPort{
												Name: "p2p-http",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(r.Context, r.Client, ingress, func() error {
		return r.updateWssIngress(ingress, &peer)
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *PeerReconciler) updateWssIngress(ingress *v1.Ingress, peer *teranodev1alpha1.Peer) error {
	err := controllerutil.SetControllerReference(peer, ingress, r.Scheme)
	if err != nil {
		return err
	}
	if peer.Spec.WssIngress == nil {
		return nil
	}
	if peer.Spec.WssIngress.Annotations == nil {
		ingress.Annotations = peer.Spec.WssIngress.Annotations
	}
	if peer.Spec.WssIngress.ClassName != nil {
		ingress.Spec.IngressClassName = peer.Spec.WssIngress.ClassName
	}

	return nil
}
