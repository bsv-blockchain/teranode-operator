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

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
)

var _ = Describe("Peer Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		peer := &teranodev1alpha1.Peer{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind Peer")
			err := k8sClient.Get(ctx, typeNamespacedName, peer)
			if err != nil && errors.IsNotFound(err) {
				resource := &teranodev1alpha1.Peer{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: teranodev1alpha1.PeerSpec{},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &teranodev1alpha1.Peer{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Peer")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &PeerReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			fetchedDeployment := &appsv1.Deployment{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      "peer",
				Namespace: "default",
			}, fetchedDeployment)).To(Succeed())
			Expect(len(fetchedDeployment.Spec.Template.Spec.Containers)).To(Equal(1))
			Expect(len(fetchedDeployment.Spec.Template.Spec.Containers[0].VolumeMounts)).To(Equal(1))
			Expect(fetchedDeployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal("/data"))
			Expect(fetchedDeployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name).To(Equal(SharedPVCName))
		})
	})
})
