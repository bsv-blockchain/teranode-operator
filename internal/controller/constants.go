// Package controller implements all reconciler logic
package controller

// DefaultServiceAccountName defines the name of the service acount created by the bundle
const DefaultServiceAccountName = "teranode-operator-service-runner"

// SharedPVCName is the PVC name for shared storage
const SharedPVCName = "cluster-storage"

// DefaultImage is the default teranode service image
const DefaultImage = "434394763103.dkr.ecr.eu-north-1.amazonaws.com/teranode-public:v0.6.2"

// DefaultCoinbaseImage is the default coinbase service image
const DefaultCoinbaseImage = "434394763103.dkr.ecr.eu-north-1.amazonaws.com/teranode-coinbase:v0.1.0"

// Service Names

const BlockchainServiceName = "blockchain"

// Ports for services

const (
	AlertSystemPort          = 9908
	AlertWebserverPort       = 3000
	AssetGRPCPort            = 8091
	AssetHTTPPort            = 8090
	BlockAssemblyPort        = 8085
	BlockchainGRPCPort       = 8087
	BlockchainHTTPPort       = 8082
	BlockValidationGRPCPort  = 8088
	BlockValidationHTTPPort  = 8188
	BootstrapGRPCPort        = 8089
	BootstrapHTTPPort        = 8099
	CoinbaseGRPCPort         = 8093
	CoinbaseHTTPPort         = 8094
	CoinbaseP2PPort          = 9907
	MinerHTTPPort            = 8092
	PeerPort                 = 9905
	PeerLegacyPort           = 8333
	PeerHTTPPort             = 9906
	PropagationGRPCPort      = 8084
	PropagationHTTPPort      = 8833
	PropagationQuicPort      = 8384
	RPCPort                  = 9292
	SubtreeValidatorGRPCPort = 8086
	LegacyHTTPPort           = 8098
	ProfilerPort             = 9091
	DebuggerPort             = 4040
	HealthPort               = 8000
)

// Deployment Names
const (
	AssetDeploymentName            = "asset"
	BlockAssemblyDeploymentName    = "block-assembly"
	BlockchainDeploymentName       = "blockchain"
	BlockValidationDeploymentName  = "block-validation"
	BootstrapDeploymentName        = "bootstrap"
	CoinbaseDeploymentName         = "coinbase"
	MinerDeploymentName            = "miner"
	PropagationDeploymentName      = "propagation"
	SubtreeValidatorDeploymentName = "subtree-validator"
	AlertSystemDeploymentName      = "alert-system"
)

// Replicas
const (
	DefaultAssetReplicas            = 2
	DefaultBlockAssemblyReplicas    = 2
	DefaultBlockchainReplicas       = 1
	DefaultBlockPersisterReplicas   = 1
	DefaultBlockValidationReplicas  = 1
	DefaultLegacyReplicas           = 1
	DefaultPeerReplicas             = 1
	DefaultRPCReplicas              = 1
	DefaultUtxoPersisterReplicas    = 1
	DefaultPropagationReplicas      = 2
	DefaultSubtreeValidatorReplicas = 2
	DefaultAlertSystemReplicas      = 1
)
