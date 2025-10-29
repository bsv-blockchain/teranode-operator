## Cluster Custom Resource
This resource is what a user is expected to create in order to get an instance of Teranode. While each service is managed by separate APIs, this API allows the user to declare all of their services in one place with shared artifacts.

The following is a sample Cluster CR:
```yaml
apiVersion: teranode.bsvblockchain.org/v1alpha1
kind: Cluster
metadata:
  name: cluster-sample
spec:
  configMapName: "my-config"
  bootstrap:
    enabled: false
    spec: {}
  asset:
    enabled: true
    spec:
      serviceAnnotations:
        traefik.ingress.kubernetes.io/service.serversscheme: h2c
      resources:
        requests:
          cpu: 1
          memory: 1Gi
        limits:
          memory: 2Gi
      grpcIngress:
        className: &my-class traefik-internal
        annotations:
          traefik.ingress.kubernetes.io/router.entrypoints: ast-grpc
        host: &my-host t3.testing.ubsv.dev
      httpIngress:
        className: *my-class
        annotations:
          traefik.ingress.kubernetes.io/router.entrypoints: web
        host: *my-host
      httpsIngress:
        className: *my-class
        annotations:
          traefik.ingress.kubernetes.io/router.entrypoints: websecure
          cert-manager.io/cluster-issuer: letsencrypt-prod
        host: *my-host
  blockValidator:
    enabled: true
    spec:
      resources:
        requests:
          cpu: 1
          memory: 1Gi
        limits:
          memory: 2Gi
  blockPersister:
    enabled: true
    spec:
      resources:
        requests:
          cpu: 1
          memory: 1Gi
        limits:
          memory: 2Gi
      storageResources:
        requests:
          storage: 5Gi
      storageClass: "fsx-sc"
  blockchain:
    enabled: true
    spec:
      resources:
        requests:
          cpu: 1
          memory: 1Gi
        limits:
          memory: 2Gi
  blockAssembly:
    enabled: true
    spec:
      resources:
        requests:
          cpu: 1
          memory: 1Gi
        limits:
          memory: 2Gi
  miner:
    enabled: true
    spec:
      resources:
        requests:
          cpu: 1
          memory: 1Gi
        limits:
          memory: 2Gi
  peer:
    enabled: true
    spec:
      resources:
        requests:
          cpu: 1
          memory: 1Gi
        limits:
          memory: 2Gi
      grpcIngress:
        className: *my-class
        annotations:
          traefik.ingress.kubernetes.io/router.entrypoints: peer-grpc
        host: *my-host
      wsIngress:
        className: *my-class
        annotations:
          traefik.ingress.kubernetes.io/router.entrypoints: p2p-ws
        host: *my-host
      wssIngress:
        className: *my-class
        annotations:
          traefik.ingress.kubernetes.io/router.entrypoints: p2p-wss
        host: *my-host
  propagation:
    enabled: true
    spec:
      serviceAnnotations:
        traefik.ingress.kubernetes.io/service.serversscheme: h2c
      grpcIngress:
        className: *my-class
        annotations:
          traefik.ingress.kubernetes.io/router.entrypoints: prop-grpc
        host: *my-host
      resources:
        requests:
          cpu: 1
          memory: 1Gi
        limits:
          memory: 2Gi
  subtreeValidator:
    enabled: true
    spec:
      resources:
        requests:
          cpu: 1
          memory: 2Gi
        limits:
          memory: 6Gi
      storageResources:
        requests:
          storage: 5Gi
      storageClass: "fsx-sc"
      storageVolume: "shared-storage-1"
  validator:
    enabled: false
    spec: {}
  coinbase:
    enabled: true
    spec:
      resources:
        requests:
          cpu: 1
          memory: 1Gi
        limits:
          memory: 2Gi
```

At the root level, `configMapName` allows the user to set a configmap that will be mounted as environment variables for each service.


