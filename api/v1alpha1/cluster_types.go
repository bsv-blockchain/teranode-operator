/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AssetConfig defines the asset configuration
type AssetConfig struct {
	Enabled bool       `json:"enabled"`
	Spec    *AssetSpec `json:"spec"`
}

// AlertSystemConfig defines the alert system configuration
type AlertSystemConfig struct {
	Enabled bool             `json:"enabled"`
	Spec    *AlertSystemSpec `json:"spec"`
}

// BlockAssemblyConfig defines the blockassembly configuration
type BlockAssemblyConfig struct {
	Enabled bool               `json:"enabled"`
	Spec    *BlockAssemblySpec `json:"spec"`
}

// BlockchainConfig defines the blockchain configuration
type BlockchainConfig struct {
	Enabled bool            `json:"enabled"`
	Spec    *BlockchainSpec `json:"spec"`
}

// BlockPersisterConfig defines the blockpersister configuration
type BlockPersisterConfig struct {
	Enabled bool                `json:"enabled"`
	Spec    *BlockPersisterSpec `json:"spec"`
}

// BlockValidatorConfig defines the blockvalidator configuration
type BlockValidatorConfig struct {
	Enabled bool                `json:"enabled"`
	Spec    *BlockValidatorSpec `json:"spec"`
}

// BootstrapConfig defines the bootstrap configuration
type BootstrapConfig struct {
	Enabled bool           `json:"enabled"`
	Spec    *BootstrapSpec `json:"spec"`
}

// CoinbaseConfig defines the coinbase configuration
type CoinbaseConfig struct {
	Enabled bool          `json:"enabled"`
	Spec    *CoinbaseSpec `json:"spec"`
}

// LegacyConfig defines the legacy configuration
type LegacyConfig struct {
	Enabled bool        `json:"enabled"`
	Spec    *LegacySpec `json:"spec"`
}

// PeerConfig defines the miner configuration
type PeerConfig struct {
	Enabled bool      `json:"enabled"`
	Spec    *PeerSpec `json:"spec"`
}

// PropagationConfig defines the propagation configuration
type PropagationConfig struct {
	Enabled bool             `json:"enabled"`
	Spec    *PropagationSpec `json:"spec"`
}

// RPCConfig defines the rpc configuration
type RPCConfig struct {
	Enabled bool     `json:"enabled"`
	Spec    *RPCSpec `json:"spec"`
}

// SubtreeValidatorConfig defines the subtreevalidator configuration
type SubtreeValidatorConfig struct {
	Enabled bool                  `json:"enabled"`
	Spec    *SubtreeValidatorSpec `json:"spec"`
}

// UtxoPersisterConfig defines the utxo persister configuration
type UtxoPersisterConfig struct {
	Enabled bool               `json:"enabled"`
	Spec    *UtxoPersisterSpec `json:"spec"`
}

// ValidatorConfig defines the validator configuration
type ValidatorConfig struct {
	Enabled bool           `json:"enabled"`
	Spec    *ValidatorSpec `json:"spec"`
}

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	AlertSystem      AlertSystemConfig      `json:"alertSystem"`
	Asset            AssetConfig            `json:"asset"`
	BlockAssembly    BlockAssemblyConfig    `json:"blockAssembly"`
	Blockchain       BlockchainConfig       `json:"blockchain"`
	BlockPersister   BlockPersisterConfig   `json:"blockPersister"`
	BlockValidator   BlockValidatorConfig   `json:"blockValidator"`
	Bootstrap        BootstrapConfig        `json:"bootstrap"`
	Coinbase         CoinbaseConfig         `json:"coinbase"`
	Legacy           LegacyConfig           `json:"legacy"`
	Peer             PeerConfig             `json:"peer"`
	Propagation      PropagationConfig      `json:"propagation"`
	RPC              RPCConfig              `json:"rpc"`
	SubtreeValidator SubtreeValidatorConfig `json:"subtreeValidator"`
	UtxoPersister    UtxoPersisterConfig    `json:"utxoPersister"`
	Validator        ValidatorConfig        `json:"validator"`

	ConfigMapName string                 `json:"configMapName"`
	Env           []corev1.EnvVar        `json:"env,omitempty"`
	EnvFrom       []corev1.EnvFromSource `json:"envFrom,omitempty"`
	Image         string                 `json:"image,omitempty"`

	SharedStorage StorageConfig `json:"sharedStorage"`
}

type StorageConfig struct {
	StorageResources *corev1.VolumeResourceRequirements `json:"storageResources,omitempty"`
	StorageClass     string                             `json:"storageClass,omitempty"`
	StorageVolume    string                             `json:"storageVolume,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Cluster is the Schema for the nodes API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
