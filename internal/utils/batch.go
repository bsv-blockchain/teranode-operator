/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package utils implements utility functions
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
