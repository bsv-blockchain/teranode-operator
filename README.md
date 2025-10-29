# teranode-operator
An operator to manage the Teranode services for Kubernetes

## Description
This operator controls the management of each microservice associated with a Teranode cluster. It currently supports deployment via bundle.

# Installation

## TODO: Dylan to update with helm documentation

## Running a node
Once you have the operator installed, modify `config/samples/teranode_v1alpha1_node.yaml` with your needed configuration values, then create the instance in the cluster:
```bash
$ kubectl create config/samples/teranode_v1alpha1_cluster.yaml
```
This step assumes you have created a prerequisite `configmap` and specified it on the above CR.

This will create the associated services, and you should see something like:
```bash
$ kubectl get pods
NAME                                                              READY   STATUS      RESTARTS   AGE
asset-5cc5745c75-6m5gf                                            1/1     Running     0          3d11h
asset-5cc5745c75-84p58                                            1/1     Running     0          3d11h
block-assembly-649dfd8596-k8q29                                   1/1     Running     0          3d11h
block-assembly-649dfd8596-njdgn                                   1/1     Running     0          3d11h
block-persister-57784567d6-tdln7                                  1/1     Running     0          3d11h
block-persister-57784567d6-wdx84                                  1/1     Running     0          3d11h
block-validator-6c4bf46f8b-bvxmm                                  1/1     Running     0          3d11h
blockchain-ccbbd894c-k95z9                                        1/1     Running     0          3d11h
coinbase-6d769f5f4d-zkb4s                                         1/1     Running     0          3d11h
dkr-ecr-eu-north-1-amazonaws-com-teranode-operator-bundle-v0-1    1/1     Running     0          3d11h
ede69fe8f248328195a7b76b2fc4c65a4ae7b7185126cdfd54f61c7eadffnzv   0/1     Completed   0          3d11h
miner-6b454ff67c-jsrgv                                            1/1     Running     0          3d11h
peer-6845bc4749-24ms4                                             1/1     Running     0          3d11h
propagation-648cd4cc56-cw5bp                                      1/1     Running     0          3d11h
propagation-648cd4cc56-sllxb                                      1/1     Running     0          3d11h
subtree-validator-7879f559d5-9gg9c                                1/1     Running     0          3d11h
subtree-validator-7879f559d5-x2dd4                                1/1     Running     0          3d11h
teranode-operator-controller-manager-768f498c4d-mk49k             2/2     Running     0          3d11h
```

## Getting Started With Development

### Prerequisites
- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/teranode-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/teranode-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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

