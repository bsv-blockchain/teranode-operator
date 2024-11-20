// Package controller implements all reconciler logic
package controller

// DefaultServiceAccountName defines the name of the service acount created by the bundle
const DefaultServiceAccountName = "teranode-operator-service-runner"

// SharedPVCName is the PVC name for shared storage
const SharedPVCName = "cluster-storage"

// DefaultImage is the default teranode service image
const DefaultImage = "434394763103.dkr.ecr.eu-north-1.amazonaws.com/teranode-public:v0.4.2"

// Service Names

const BlockchainServiceName = "blockchain"

// Ports for services

const AlertSystemPort = 9908
const AssetGRPCPort = 8091
const AssetHTTPPort = 8090
const BlockAssemblyPort = 8085
const BlockchainGRPCPort = 8087
const BlockchainHTTPPort = 8082
const BlockValidationGRPCPort = 8088
const BlockValidationHTTPPort = 8188
const BootstrapGRPCPort = 8089
const BootstrapHTTPPort = 8099
const CoinbaseGRPCPort = 8093
const MinerHTTPPort = 8092
const PeerPort = 9095
const PeerHTTPPort = 9096
const PropagationGRPCPort = 8084
const PropagationHTTPPort = 8833
const PropagationQuicPort = 8384
const RPCPort = 9292
const SubtreeValidatorGRPCPort = 8086
const LegacyHTTPPort = 8098
const ProfilerPort = 9091
const DebuggerPort = 4040
const HealthPort = 8000
