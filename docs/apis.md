## API Summary

The Teranode Operator exposes the following APIs in the group `teranode.bsvblockchain.org`:
* [`Asset`](./asset.md)
* `BlockAssembly`
* `Blockchain`
* `BlockPersister`
* `BlockValidator`
* `Bootstrap`
* `Coinbase`
* `Miner`
* [`Node`](./node.md)
* `Peer`
* [`Propagation`](./propagation.md)
* `SubtreeValidator`

Each of these APIs share a common set of configuration values that can be set directly on their spec (with the exception of [Node](./node.md)). These configuration values allow the user to customize the deployment of each service to suit their needs.

## Configuration
| Key               | Type                          | Description                                                          |
|-------------------|-------------------------------|----------------------------------------------------------------------|
| `nodeSelector`    | `map[string]string`           | Node Selector field for deployment                                   |
| `tolerations`     | `[]corev1.Toleration`         | Tolerations for deployment                                           |
| `affinity`        | `corev1.Affinity`             | Affinity for deployment                                              |
| `resources`       | `corev1.ResourceRequirements` | Set resource requests/limits on deployment                           |
| `image`           | `string`                      | Image to use for Teranode service                                    |
| `imagePullPolicy` | `corev1.PullPolicy`           | Image Pull Policy to use for Teranode image                          |
| `serviceAccount`  | `string`                      | Service account to run the deployment with                           |
| `configMapName`   | `string`                      | Name of configmap to inject as environment variables for the service |

