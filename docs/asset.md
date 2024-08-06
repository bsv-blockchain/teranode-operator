## Asset API
In addition to the standard configuration values, the Asset API also takes in the following parameters:

| Key                  | Type                              | Description                                             |
|----------------------|-----------------------------------|---------------------------------------------------------|
| `serviceAnnotations` | `map[string]string`               | Annotations to set on the Kubernetes service definition |
| `httpIngress`        | [`IngressDefinition`](ingress.md) | Defined ingress configuration values for http access    |
| `httpsIngress`       | [`IngressDefinition`](ingress.md) | Defined ingress configuration values for https access   |
| `grpcIngress`        | [`IngressDefinition`](ingress.md) | Defined ingress configuration values for grpc access    |
