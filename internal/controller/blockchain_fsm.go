package controller

import (
	"errors"
	"fmt"

	teranodev1alpha1 "github.com/bsv-blockchain/teranode-operator/api/v1alpha1"
	"github.com/bsv-blockchain/teranode/services/blockchain"
	"github.com/bsv-blockchain/teranode/settings"
	"github.com/bsv-blockchain/teranode/ulogger"
	"github.com/go-logr/logr"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

// ErrFSMDisabled is returned when the finite state machine is disabled
var ErrFSMDisabled = errors.New("finite state machine is disabled")

// GetFSMState gets the current state of the FSM
func (r *BlockchainReconciler) GetFSMState(log logr.Logger) (*blockchain.FSMStateType, error) {
	b := teranodev1alpha1.Blockchain{}
	if err := r.Get(r.Context, r.NamespacedName, &b); err != nil {
		return nil, err
	}
	if b.Spec.FiniteStateMachine != nil && !b.Spec.FiniteStateMachine.Enabled {
		return nil, ErrFSMDisabled
	}
	if r.BlockchainClient == nil {
		uLog := ulogger.New("fsm")
		host := BlockchainServiceName
		if b.Spec.FiniteStateMachine != nil && b.Spec.FiniteStateMachine.Host != "" {
			host = b.Spec.FiniteStateMachine.Host
		}
		blockchainHost := fmt.Sprintf("%s:%d", host, BlockchainGRPCPort)
		tSettings := settings.Settings{
			BlockChain: settings.BlockChainSettings{
				GRPCAddress: blockchainHost,
			},
		}
		bClient, err := blockchain.NewClient(r.Context, uLog, &tSettings, blockchainHost)
		if err != nil {
			return nil, err
		}
		r.BlockchainClient = bClient
		r.Log.Info("Initiating FSM client", "host", blockchainHost)
	}

	return r.BlockchainClient.GetFSMCurrentState(r.Context)
}

func (r *BlockchainReconciler) IsLegacyEnabled() (bool, error) {
	b := teranodev1alpha1.Blockchain{}
	if err := r.Get(r.Context, r.NamespacedName, &b); err != nil {
		return false, err
	}
	legacyEnabled := false
	// Attempt to get the parent Cluster CR to know service configuration
	ownerRefs := b.GetOwnerReferences()
	for _, ownerRef := range ownerRefs {
		if ownerRef.Kind == "Cluster" {
			cluster := teranodev1alpha1.Cluster{}
			if err := r.Get(
				r.Context,
				types.NamespacedName{
					Name:      ownerRef.Name,
					Namespace: r.NamespacedName.Namespace,
				}, &cluster); err != nil && !k8serrors.IsNotFound(err) {
				return false, err
			}
			legacyEnabled = cluster.Spec.Legacy.Enabled
		}
	}
	return legacyEnabled, nil
}

func (r *BlockchainReconciler) ReconcileState(state blockchain.FSMStateType) error {
	switch state {
	case blockchain.FSMStateRUNNING:
		break
	case blockchain.FSMStateIDLE:
		legacyEnabled, err := r.IsLegacyEnabled()
		if err != nil {
			// lets just break for now because we don't know if we are supposed to be doing anything
			return err
		}
		if !legacyEnabled {
			return r.BlockchainClient.Run(r.Context, "")
		}
	}
	return nil
}
