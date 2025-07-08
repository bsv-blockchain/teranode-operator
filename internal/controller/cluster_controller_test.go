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
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	teranodev1alpha1 "github.com/bitcoin-sv/teranode-operator/api/v1alpha1"
)

var _ = Describe("Cluster Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		cluster := &teranodev1alpha1.Cluster{
			Spec: teranodev1alpha1.ClusterSpec{
				Legacy: teranodev1alpha1.LegacyConfig{
					Spec: &teranodev1alpha1.LegacySpec{},
				},
				Asset: teranodev1alpha1.AssetConfig{
					Spec: &teranodev1alpha1.AssetSpec{},
				},
				BlockAssembly: teranodev1alpha1.BlockAssemblyConfig{
					Spec: &teranodev1alpha1.BlockAssemblySpec{},
				},
				Blockchain: teranodev1alpha1.BlockchainConfig{
					Spec: &teranodev1alpha1.BlockchainSpec{},
				},
				BlockPersister: teranodev1alpha1.BlockPersisterConfig{
					Spec: &teranodev1alpha1.BlockPersisterSpec{},
				},
				BlockValidator: teranodev1alpha1.BlockValidatorConfig{
					Spec: &teranodev1alpha1.BlockValidatorSpec{},
				},
			},
		}

		BeforeEach(func() {
			By("creating the custom resource for the Kind Cluster")
			err := k8sClient.Get(ctx, typeNamespacedName, cluster)
			if err != nil && errors.IsNotFound(err) {
				resource := &teranodev1alpha1.Cluster{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: teranodev1alpha1.ClusterSpec{
						Asset: teranodev1alpha1.AssetConfig{
							Spec: &teranodev1alpha1.AssetSpec{},
						},
						AlertSystem: teranodev1alpha1.AlertSystemConfig{
							Spec: &teranodev1alpha1.AlertSystemSpec{},
						},
						BlockAssembly: teranodev1alpha1.BlockAssemblyConfig{
							Spec: &teranodev1alpha1.BlockAssemblySpec{},
						},
						Blockchain: teranodev1alpha1.BlockchainConfig{
							Spec: &teranodev1alpha1.BlockchainSpec{},
						},
						BlockPersister: teranodev1alpha1.BlockPersisterConfig{
							Spec: &teranodev1alpha1.BlockPersisterSpec{},
						},
						BlockValidator: teranodev1alpha1.BlockValidatorConfig{
							Spec: &teranodev1alpha1.BlockValidatorSpec{},
						},
						Bootstrap: teranodev1alpha1.BootstrapConfig{
							Spec: &teranodev1alpha1.BootstrapSpec{},
						},
						Coinbase: teranodev1alpha1.CoinbaseConfig{
							Spec: &teranodev1alpha1.CoinbaseSpec{},
						},
						Legacy: teranodev1alpha1.LegacyConfig{
							Spec: &teranodev1alpha1.LegacySpec{},
						},
						Peer: teranodev1alpha1.PeerConfig{
							Spec: &teranodev1alpha1.PeerSpec{},
						},
						Propagation: teranodev1alpha1.PropagationConfig{
							Spec: &teranodev1alpha1.PropagationSpec{},
						},
						RPC: teranodev1alpha1.RPCConfig{
							Spec: &teranodev1alpha1.RPCSpec{},
						},
						SubtreeValidator: teranodev1alpha1.SubtreeValidatorConfig{
							Spec: &teranodev1alpha1.SubtreeValidatorSpec{},
						},
						UtxoPersister: teranodev1alpha1.UtxoPersisterConfig{
							Spec: &teranodev1alpha1.UtxoPersisterSpec{},
						},
						Validator: teranodev1alpha1.ValidatorConfig{
							Spec: &teranodev1alpha1.ValidatorSpec{},
						},
					},
					// TODO(user): Specify other spec details if needed.
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &teranodev1alpha1.Cluster{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Cluster")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ClusterReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
		It("should create only enabled components", func() {
			cluster := &teranodev1alpha1.Cluster{}
			err := k8sClient.Get(ctx, typeNamespacedName, cluster)
			Expect(err).NotTo(HaveOccurred())
			// Create a cluster with only some components enabled
			cluster.Spec.Asset.Enabled = true
			cluster.Spec.BlockAssembly.Enabled = true

			Expect(k8sClient.Update(ctx, cluster)).To(Succeed())

			// Reconcile
			controllerReconciler := &ClusterReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify enabled components exist
			asset := &teranodev1alpha1.Asset{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-asset", cluster.Name),
				Namespace: "default",
			}, asset)).To(Succeed())

			blockAssembly := &teranodev1alpha1.BlockAssembly{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-blockassembly", cluster.Name),
				Namespace: "default",
			}, blockAssembly)).To(Succeed())

			// Verify disabled component doesn't exist
			alertSystem := &teranodev1alpha1.AlertSystem{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-alert-system", cluster.Name),
				Namespace: "default",
			}, alertSystem)
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})

		It("should create a PVC when cluster is created regardless of spec", func() {
			pvc := &v1.PersistentVolumeClaim{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      SharedPVCName,
				Namespace: "default",
			}, pvc)).To(Succeed(), "PVC should be created when cluster is created")
		})

		It("should delete components when disabled", func() {
			// Update the cluster with AlertSystem enabled
			cluster := &teranodev1alpha1.Cluster{}
			err := k8sClient.Get(ctx, typeNamespacedName, cluster)
			Expect(err).NotTo(HaveOccurred())

			// Enable alert system
			cluster.Spec.AlertSystem.Enabled = true

			// Update the cluster
			Expect(k8sClient.Update(ctx, cluster)).To(Succeed())

			// Reconcile to create the resources
			controllerReconciler := &ClusterReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify AlertSystem exists
			alertSystem := &teranodev1alpha1.AlertSystem{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-alert-system", cluster.Name),
				Namespace: "default",
			}, alertSystem)).To(Succeed())

			// Now disable the component
			fetchedCluster := &teranodev1alpha1.Cluster{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, fetchedCluster)).To(Succeed())
			fetchedCluster.Spec.AlertSystem.Enabled = false
			Expect(k8sClient.Update(ctx, fetchedCluster)).To(Succeed())

			// Reconcile again
			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify AlertSystem was deleted
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-alert-system", cluster.Name),
				Namespace: "default",
			}, alertSystem)
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})

		It("should apply image from cluster to components", func() {
			// Update the cluster with custom image
			cluster := &teranodev1alpha1.Cluster{}
			err := k8sClient.Get(ctx, typeNamespacedName, cluster)
			Expect(err).NotTo(HaveOccurred())

			// Set a custom image for the cluster
			testImage := "custom-image:v1"
			cluster.Spec.Image = testImage
			cluster.Spec.Asset.Enabled = true // Asset should have this image

			// Set a second custom image for BlockPersister
			// This should take precedence over the cluster image
			testImage2 := "custom-image2:v2"
			cluster.Spec.BlockPersister.Spec = &teranodev1alpha1.BlockPersisterSpec{
				DeploymentOverrides: &teranodev1alpha1.DeploymentOverrides{
					Image: testImage2, // Custom image for BlockPersister
				},
			}
			cluster.Spec.BlockPersister.Enabled = true

			// Update the cluster
			Expect(k8sClient.Update(ctx, cluster)).To(Succeed())

			// Reconcile
			controllerReconciler := &ClusterReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify Asset got the cluster image
			asset := &teranodev1alpha1.Asset{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-asset", cluster.Name),
				Namespace: "default",
			}, asset)).To(Succeed())
			Expect(asset.Spec.DeploymentOverrides.Image).To(Equal(testImage))

			// Verify BlockPersister kept its custom image
			blockPersister := &teranodev1alpha1.BlockPersister{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-blockpersister", cluster.Name),
				Namespace: "default",
			}, blockPersister)).To(Succeed())
			Expect(blockPersister.Spec.DeploymentOverrides.Image).To(Equal(testImage2))
		})

		It("should create all components when all are enabled", func() {
			cluster := &teranodev1alpha1.Cluster{}
			err := k8sClient.Get(ctx, typeNamespacedName, cluster)
			Expect(err).NotTo(HaveOccurred())

			cluster.Spec.Asset.Enabled = true
			cluster.Spec.AlertSystem.Enabled = true
			cluster.Spec.Blockchain.Enabled = true
			cluster.Spec.Peer.Enabled = true
			cluster.Spec.RPC.Enabled = true
			cluster.Spec.BlockPersister.Enabled = true
			cluster.Spec.BlockValidator.Enabled = true
			cluster.Spec.Coinbase.Enabled = true
			cluster.Spec.Legacy.Enabled = true
			cluster.Spec.Propagation.Enabled = true
			cluster.Spec.SubtreeValidator.Enabled = true
			cluster.Spec.UtxoPersister.Enabled = true
			cluster.Spec.Validator.Enabled = true

			// Update the cluster
			Expect(k8sClient.Update(ctx, cluster)).To(Succeed())

			// Reconcile
			controllerReconciler := &ClusterReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify components were created (checking a subset)
			components := []struct {
				name   string
				object client.Object
			}{
				{fmt.Sprintf("%s-alert-system", cluster.Name), &teranodev1alpha1.AlertSystem{}},
				{fmt.Sprintf("%s-asset", cluster.Name), &teranodev1alpha1.Asset{}},
				{fmt.Sprintf("%s-blockchain", cluster.Name), &teranodev1alpha1.Blockchain{}},
				{fmt.Sprintf("%s-blockpersister", cluster.Name), &teranodev1alpha1.BlockPersister{}},
				{fmt.Sprintf("%s-blockvalidator", cluster.Name), &teranodev1alpha1.BlockValidator{}},
				{fmt.Sprintf("%s-coinbase", cluster.Name), &teranodev1alpha1.Coinbase{}},
				{fmt.Sprintf("%s-legacy", cluster.Name), &teranodev1alpha1.Legacy{}},
				{fmt.Sprintf("%s-propagation", cluster.Name), &teranodev1alpha1.Propagation{}},
				{fmt.Sprintf("%s-subtreevalidator", cluster.Name), &teranodev1alpha1.SubtreeValidator{}},
				{fmt.Sprintf("%s-utxo-persister", cluster.Name), &teranodev1alpha1.UtxoPersister{}},
				{fmt.Sprintf("%s-validator", cluster.Name), &teranodev1alpha1.Validator{}},
				{fmt.Sprintf("%s-peer", cluster.Name), &teranodev1alpha1.Peer{}},
				{fmt.Sprintf("%s-rpc", cluster.Name), &teranodev1alpha1.RPC{}},
			}

			for _, component := range components {
				Expect(k8sClient.Get(ctx, types.NamespacedName{
					Name:      component.name,
					Namespace: "default",
				}, component.object)).To(Succeed(), "Component %s should exist", component.name)
			}
		})
		It("should disable all services when cluster is disabled", func() {
			cluster := &teranodev1alpha1.Cluster{}
			err := k8sClient.Get(ctx, typeNamespacedName, cluster)
			Expect(err).NotTo(HaveOccurred())

			// Enable some components
			cluster.Spec.Asset.Enabled = true
			cluster.Spec.BlockAssembly.Enabled = true
			cluster.Spec.Blockchain.Enabled = true

			// Update the cluster
			Expect(k8sClient.Update(ctx, cluster)).To(Succeed())

			// Reconcile to create the resources
			controllerReconciler := &ClusterReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}
			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify components were created
			asset := &teranodev1alpha1.Asset{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-asset", cluster.Name),
				Namespace: "default",
			}, asset)).To(Succeed())

			blockAssembly := &teranodev1alpha1.BlockAssembly{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-blockassembly", cluster.Name),
				Namespace: "default",
			}, blockAssembly)).To(Succeed())

			blockchain := &teranodev1alpha1.Blockchain{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-blockchain", cluster.Name),
				Namespace: "default",
			}, blockchain)).To(Succeed())

			// Now disable the cluster
			fetchedCluster := &teranodev1alpha1.Cluster{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, fetchedCluster)).To(Succeed())
			fetchedCluster.Spec.Enabled = ptr.To(false)
			Expect(k8sClient.Update(ctx, fetchedCluster)).To(Succeed())

			// Reconcile again
			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify all components were deleted
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-asset", cluster.Name),
				Namespace: "default",
			}, asset)
			Expect(errors.IsNotFound(err)).To(BeTrue(), "Asset should be deleted when cluster is disabled")

			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-blockassembly", cluster.Name),
				Namespace: "default",
			}, blockAssembly)
			Expect(errors.IsNotFound(err)).To(BeTrue(), "BlockAssembly should be deleted when cluster is disabled")

			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      fmt.Sprintf("%s-blockchain", cluster.Name),
				Namespace: "default",
			}, blockchain)
			Expect(errors.IsNotFound(err)).To(BeTrue(), "Blockchain should be deleted when cluster is disabled")
		})
	})
})
