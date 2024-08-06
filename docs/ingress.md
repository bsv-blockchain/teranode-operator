## Ingress Configuration

There are 3 services that require ingress to work properly in the Teranode environment: `Asset`, `Peer`, and `Propagation`. All 3 of these services expose their respective ingress definitions with a custom type `IngressDef` which has the following key values:
## Configuration
| Key               | Type                          | Description                                                                |
|-------------------|-------------------------------|----------------------------------------------------------------------------|
| `className`       | `string`                      | Ingress class to be used for this ingress, if left blank it is the default |
| `annotations`     | `map[string]string`           | Custom annotations to be applied to the ingress resource                   |
| `host`            | `string`                      | Host value to be used on the ingress                                       |

This provides the user with flexibility in deferring to their preferred ingress provider while using the native Kubernetes `Ingress` resource.