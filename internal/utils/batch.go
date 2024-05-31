package utils

import "github.com/go-logr/logr"

// ReconcileFunc is a reconcile function type
type ReconcileFunc func(logr.Logger) (bool, error)

// ReconcileBatch will reconcile a batch of functions
func ReconcileBatch(l logr.Logger, reconcileFunctions ...ReconcileFunc) (bool, error) {
	for _, f := range reconcileFunctions {
		if cont, err := f(l); !cont || err != nil {
			return cont, err
		}
	}
	return true, nil
}
