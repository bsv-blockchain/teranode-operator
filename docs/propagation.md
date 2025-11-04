## Propagation API
In addition to the standard configuration values, the Propagation API also takes in the following parameters:

| Key                  | Type                              | Description                                             |
|----------------------|-----------------------------------|---------------------------------------------------------|
| `serviceAnnotations` | `map[string]string`               | Annotations to set on the Kubernetes service definition |
| `grpcIngress`        | [`IngressDefinition`](ingress.md) | Defined ingress configuration values for grpc access    |
