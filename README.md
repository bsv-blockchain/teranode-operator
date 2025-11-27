<div align="center">

# ☸️&nbsp;&nbsp;teranode-operator

**Kubernetes operator for orchestrating Teranode blockchain infrastructure.**

<br/>

<a href="https://github.com/bsv-blockchain/teranode-operator/releases"><img src="https://img.shields.io/github/release-pre/bsv-blockchain/teranode-operator?include_prereleases&style=flat-square&logo=github&color=black" alt="Release"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/bsv-blockchain/teranode-operator?style=flat-square&logo=go&color=00ADD8" alt="Go Version"></a>
<a href="https://github.com/bsv-blockchain/teranode-operator/blob/master/LICENSE"><img src="https://img.shields.io/github/license/bsv-blockchain/teranode-operator?style=flat-square&color=blue" alt="License"></a>

<br/>

<table align="center" border="0">
  <tr>
    <td align="right">
       <code>CI / CD</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/bsv-blockchain/teranode-operator/actions"><img src="https://img.shields.io/github/actions/workflow/status/bsv-blockchain/teranode-operator/fortress.yml?branch=master&label=build&logo=github&style=flat-square" alt="Build"></a>
       <a href="https://github.com/bsv-blockchain/teranode-operator/actions"><img src="https://img.shields.io/github/last-commit/bsv-blockchain/teranode-operator?style=flat-square&logo=git&logoColor=white&label=last%20update" alt="Last Commit"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Quality</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://goreportcard.com/report/github.com/bsv-blockchain/teranode-operator"><img src="https://goreportcard.com/badge/github.com/bsv-blockchain/teranode-operator?style=flat-square" alt="Go Report"></a>
       <!-- <a href="https://codecov.io/gh/bsv-blockchain/teranode-operator"><img src="https://codecov.io/gh/bsv-blockchain/teranode-operator/branch/master/graph/badge.svg?style=flat-square" alt="Coverage"></a> -->
    </td>
  </tr>

  <tr>
    <td align="right">
       <code>Security</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://scorecard.dev/viewer/?uri=github.com/bsv-blockchain/teranode-operator"><img src="https://api.scorecard.dev/projects/github.com/bsv-blockchain/teranode-operator/badge?style=flat-square" alt="Scorecard"></a>
       <a href=".github/SECURITY.md"><img src="https://img.shields.io/badge/policy-active-success?style=flat-square&logo=security&logoColor=white" alt="Security"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Community</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/bsv-blockchain/teranode-operator/graphs/contributors"><img src="https://img.shields.io/github/contributors/bsv-blockchain/teranode-operator?style=flat-square&color=orange" alt="Contributors"></a>
       <a href="https://github.com/sponsors/bsv-blockchain"><img src="https://img.shields.io/badge/sponsor-BSV-181717.svg?logo=github&style=flat-square" alt="Sponsor"></a>
    </td>
  </tr>
</table>

</div>

<br/>
<br/>

<div align="center">

### <code>Project Navigation</code>

</div>

<table align="center">
  <tr>
    <td align="center" width="33%">
       📦&nbsp;<a href="#-installation"><code>Installation</code></a>
    </td>
    <td align="center" width="33%">
       📚&nbsp;<a href="#-documentation"><code>Documentation</code></a>
    </td>
    <td align="center" width="33%">
       🧪&nbsp;<a href="#-examples--tests"><code>Examples&nbsp;&&nbsp;Tests</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       🤝&nbsp;<a href="#-contributing"><code>Contributing</code></a>
    </td>
    <td align="center">
       🛠️&nbsp;<a href="#-code-standards"><code>Code&nbsp;Standards</code></a>
    </td>
    <td align="center">
       ⚡&nbsp;<a href="#-benchmarks"><code>Benchmarks</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       🤖&nbsp;<a href="#-ai-compliance"><code>AI&nbsp;Compliance</code></a>
    </td>
    <td align="center">
       ⚖️&nbsp;<a href="#-license"><code>License</code></a>
    </td>
    <td align="center">
       👥&nbsp;<a href="#-maintainers"><code>Maintainers</code></a>
    </td>
  </tr>
</table>
<br/>

## 📦 Installation

See [teranode-operator-helm](https://github.com/bsv-blockchain/teranode-operator-helm) guide for Helm install details.

<br/>

### Running a node
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
miner-6b454ff67c-jsrgv                                            1/1     Running     0          3d11h
peer-6845bc4749-24ms4                                             1/1     Running     0          3d11h
propagation-648cd4cc56-cw5bp                                      1/1     Running     0          3d11h
propagation-648cd4cc56-sllxb                                      1/1     Running     0          3d11h
subtree-validator-7879f559d5-9gg9c                                1/1     Running     0          3d11h
subtree-validator-7879f559d5-x2dd4                                1/1     Running     0          3d11h
teranode-operator-controller-manager-768f498c4d-mk49k             2/2     Running     0          3d11h
```

<br/>

## 📚 Documentation
This operator controls the management of each microservice associated with a Teranode cluster. It currently supports deployment via bundle.

<br>

### Getting Started with Development

#### Prerequisites
- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

<details>
<summary><strong><code>Deploy on the cluster</code></strong></summary>

**Runs controller manager locally on your machine**
```shell
make run
```

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


</details>

<details>
<summary><strong><code>Uninstall</code></strong></summary>

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

</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

This project uses [goreleaser](https://github.com/goreleaser/goreleaser) for streamlined binary and library deployment to GitHub. To get started, install it via:

```bash
brew install goreleaser
```

The release process is defined in the [.goreleaser.yml](.goreleaser.yml) configuration file.


Then create and push a new Git tag using:

```bash
magex version:bump push=true bump=patch branch=master
```

This process ensures consistent, repeatable releases with properly versioned artifacts and citation metadata.

</details>

<details>
<summary><strong><code>Pre-commit Hooks</code></strong></summary>
<br/>

Set up the Go-Pre-commit System to run the same formatting, linting, and tests defined in [AGENTS.md](.github/AGENTS.md) before every commit:

```bash
go install github.com/mrz1836/go-pre-commit/cmd/go-pre-commit@latest
go-pre-commit install
```

The system is configured via [.env.base](.github/.env.base) and can be customized using also using [.env.custom](.github/.env.custom) and provides 17x faster execution than traditional Python-based pre-commit hooks. See the [complete documentation](http://github.com/mrz1836/go-pre-commit) for details.

</details>

<br>

## 🧪 Examples & Tests

All unit tests and examples run via [GitHub Actions](https://github.com/bsv-blockchain/teranode-operator/actions) and use [Go version 1.25.x](https://go.dev/doc/go1.25).

Run all tests (fast):

```bash script
make test
```

<br/>

## ⚡ Benchmarks

(Coming Soon!)

<br/>

## 🛠️ Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## 🤖 AI Compliance
This project documents expectations for AI assistants using a few dedicated files:

- [AGENTS.md](.github/AGENTS.md) — canonical rules for coding style, workflows, and pull requests used by [Codex](https://chatgpt.com/codex).
- [CLAUDE.md](.github/CLAUDE.md) — quick checklist for the [Claude](https://www.anthropic.com/product) agent.
- [.cursorrules](.cursorrules) — machine-readable subset of the policies for [Cursor](https://www.cursor.so/) and similar tools.
- [sweep.yaml](.github/sweep.yaml) — rules for [Sweep](https://github.com/sweepai/sweep), a tool for code review and pull request management.

Edit `AGENTS.md` first when adjusting these policies, and keep the other files in sync within the same pull request.

<br/>

## 👥 Maintainers
| [<img src="https://github.com/icellan.png" height="50" alt="Siggi" />](https://github.com/icellan) | [<img src="https://github.com/galt-tr.png" height="50" alt="Dylan" />](https://github.com/galt-tr) | [<img src="https://github.com/oskarszoon.png" height="50" alt="Oli" />](https://github.com/oskarszoon) | [<img src="https://github.com/mrz1836.png" height="50" width="50" alt="MrZ" />](https://github.com/mrz1836) |
|:--------------------------------------------------------------------------------------------------:|:--------------------------------------------------------------------------------------------------:|:------------------------------------------------------------------------------------------------------:|:-----------------------------------------------------------------------------------------------------------:|
|                                [Siggi](https://github.com/icellan)                                 |                                [Dylan](https://github.com/galt-tr)                                 |                                  [Oli](https://github.com/oskarszoon)                                  |                                      [MrZ](https://github.com/mrz1836)                                      |

<br/>

## 🤝 Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.

[![Stars](https://img.shields.io/github/stars/bsv-blockchain/teranode-operator?label=Please%20like%20us&style=social&v=1)](https://github.com/bsv-blockchain/teranode-operator/stargazers)

<br/>

## 📝 License

[![License](https://img.shields.io/github/license/bsv-blockchain/teranode-operator.svg?style=flat&v=1)](LICENSE)
