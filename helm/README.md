# teranode-operator-helm

[Teranode](https://www.bsvblockchain.org/teranode) [Helm](https://helm.sh/docs/) chart to deploy a [teranode-operator](https://github.com/bitcoin-sv/teranode-operator) controller on [Kubernetes](https://kubernetes.io/docs/home/).

## Versions

| Version |    Date    | Comments                                                                                    |
| :-----: | :--------: | :------------------------------------------------------------------------------------------ |
|  0.1.0  | 04/11/2025 | Initial version created                                                                     |

## Example Usage

Minimal installation example.

```shell
helm registry login --username <user> --password-stdin ghcr.io
helm install teranode-operator oci://ghcr.io/bsv-blockchain/teranode-operator -n teranode-operator --set deployment.env.watch_namespace="some-namespace"
```

Minimal local installation example.

```shell
git clone git@github.com:bsv-blockchain/teranode-operator.git
cd teranode-operator/helm
kubectl create ns teranode-operator
helm install teranode-operator ./ -n teranode-operator --set deployment.env.watch_namespace="some-namespace"
```

Extended local installation example.

```shell
git clone git@github.com:bsv-blockchain/teranode-operator.git
cd teranode-operator/helm
helm install teranode-operator ./ -n some-namespace --create-namespace some-namespace --set deployment.env.watch_namespace="some-namespace, some-other-namespace, some-completely-different-namespace" --set deployment.image.tag=v2.0.1
```

Local uninstallation example.

```shell
cd teranode-operator/helm
helm uninstall teranode-operator -n teranode-operator
```

## Caveats

### Helm & Kubernetes CRDs

Due to how Helm works, currently there is no way to completely clear the Kubernetes cluster Helm installed the teranode-operator on because the [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) required by the teranode-operator are not being deleted with the `helm uninstall` command.

If you need to ensure a complete uninstallation of teranode-operator (ie. for major version upgrades) you need to execute the following command after the `helm uninstall` step:

```shell
kubectl get crd -o name | grep '\.teranode\.bsvblockchain\.org' | xargs kubectl delete
```

Please, always check for any upgrade instructions in the README.md or release notes when updating/upgrading teranode-operator.

### Sample cluster

This Helm chart deploys the Teranode Operator that requires further actions from the user, but to make it easier for new users to explore Teranode it is possible to make it deploy a sample `Cluster` definition along with sample `ConfigMap` it needs.

**Be aware that this sample cluster is very opinionated and depends on the whole suite of predefined applications deployed in BSVA style, that can be deployed using BSVA Terraform/OpenTofu modules on AWS only.**

If you won't use these modules to set up the dependant services exactly like the BSVA reference architecture looks like you won't get a working Teranode cluster this way.

In order to deploy such sample cluster you need to set at least two additional parameters, `sampleCluster.enable`, `sampleCluster.hostname` and third optional for network type `sampleCluster.network`:

```shell
helm install teranode-operator ./ -n teranode-operator --set deployment.env.watch_namespace="some-namespace" --set sampleCluster.enable=true --set sampleCluster.hostname=example.com --set sampleCluster.network=mainnet
```

### Multiple clusters

The operator requires a string with a name of the Kubernetes namespace to watch for Cluster resource being defined via the `deployment.env.watch_namespace` variable. If needed, that value can be set to multiple namespaces with a comma separated string in order for the teranode-operator to watch all given Kubernetes namespaces for multiple instances being deployed and managed by a single teranode-operator on a single Kuberneted cluster.

## Requirements

| Name       | Version |
| ---------- | ------- |
| Helm       | ~> 3.15 |
| Kubernetes | ~> 1.30 |
