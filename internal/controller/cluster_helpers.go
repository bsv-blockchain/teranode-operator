package controller

import (
	"reflect"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
)

// mergeDeploymentOverrides selectively merges deployment overrides from the cluster spec into
// target. Only fields explicitly set (non-zero) in clusterOverrides override the target; every
// field of DeploymentOverrides is handled generically, so new override fields propagate
// automatically without editing this function.
func mergeDeploymentOverrides(target, clusterOverrides *teranodev1alpha1.DeploymentOverrides) {
	if clusterOverrides == nil {
		return
	}
	t := reflect.ValueOf(target).Elem()
	s := reflect.ValueOf(clusterOverrides).Elem()
	for i := range s.NumField() {
		if f := s.Field(i); !f.IsZero() {
			t.Field(i).Set(f)
		}
	}
}
