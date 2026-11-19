package controller

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
)

func newRPCTestReconciler(t *testing.T) *ClusterReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := teranodev1alpha1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to build scheme: %v", err)
	}
	return &ClusterReconciler{Scheme: scheme}
}

func newRPCTestCluster(image string, rpcSpec *teranodev1alpha1.RPCSpec) *teranodev1alpha1.Cluster {
	return &teranodev1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "default", UID: "cluster-uid"},
		Spec: teranodev1alpha1.ClusterSpec{
			Image: image,
			RPC:   teranodev1alpha1.RPCConfig{Enabled: true, Spec: rpcSpec},
		},
	}
}

// An existing RPC CR already carries the image from its first reconcile. A later
// cluster image bump must still reach it, as it does for every other service.
func TestUpdateRPCPropagatesClusterImageBump(t *testing.T) {
	r := newRPCTestReconciler(t)
	rpc := &teranodev1alpha1.RPC{ObjectMeta: metav1.ObjectMeta{Name: "test-rpc", Namespace: "default"}}

	if err := r.updateRPC(rpc, newRPCTestCluster("teranode:old", nil)); err != nil {
		t.Fatalf("first reconcile failed: %v", err)
	}
	if err := r.updateRPC(rpc, newRPCTestCluster("teranode:new", nil)); err != nil {
		t.Fatalf("second reconcile failed: %v", err)
	}

	if got := rpc.Spec.DeploymentOverrides.Image; got != "teranode:new" {
		t.Errorf("cluster image bump not propagated to RPC, got %q", got)
	}
}

// A per-service image override in the cluster spec takes precedence over the
// cluster-wide image.
func TestUpdateRPCServiceImageOverrideWins(t *testing.T) {
	r := newRPCTestReconciler(t)
	rpc := &teranodev1alpha1.RPC{ObjectMeta: metav1.ObjectMeta{Name: "test-rpc", Namespace: "default"}}
	rpcSpec := &teranodev1alpha1.RPCSpec{
		DeploymentOverrides: &teranodev1alpha1.DeploymentOverrides{Image: "teranode:pinned"},
	}

	if err := r.updateRPC(rpc, newRPCTestCluster("teranode:new", rpcSpec)); err != nil {
		t.Fatalf("reconcile failed: %v", err)
	}

	if got := rpc.Spec.DeploymentOverrides.Image; got != "teranode:pinned" {
		t.Errorf("per-service RPC image override not honored, got %q", got)
	}
}
