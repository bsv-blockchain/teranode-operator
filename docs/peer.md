## Peer API
In addition to the standard configuration values, the Peer API also takes in the following parameters:

| Key                  | Type                              | Description                                             |
|----------------------|-----------------------------------|---------------------------------------------------------|
| `serviceAnnotations` | `map[string]string`               | Annotations to set on the Kubernetes service definition |
| `wsIngress`          | [`IngressDefinition`](ingress.md) | Defined ingress configuration values for ws access      |
| `wssIngress`         | [`IngressDefinition`](ingress.md) | Defined ingress configuration values for wss access     |
| `grpcIngress`        | [`IngressDefinition`](ingress.md) | Defined ingress configuration values for grpc access    |
