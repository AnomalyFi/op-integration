package derive

import (
	"context"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/nodekit"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

type NodeKitProvider struct {
	SequencerAddress common.Address
	L1Fetcher        L1Fetcher
	log              log.Logger
}

func NewNodeKitProvider(log log.Logger, SequencerAddress common.Address, l1Fetcher L1Fetcher) *NodeKitProvider {
	return &NodeKitProvider{
		SequencerAddress: SequencerAddress,
		L1Fetcher:        l1Fetcher,
		log:              log,
	}

}

func (provider *NodeKitProvider) VerifyCommitments(firstHeight uint64, comms []nodekit.Commitment) (bool, error) {
	fetchedComms, err := provider.L1Fetcher.L1SequencerCommitmentsFromHeight(firstHeight, uint64(len(comms)), provider.SequencerAddress)
	if err != nil {
		return false, err
	}

	if len(fetchedComms) != len(comms) {
		return false, fmt.Errorf("fetched commitments has a different length than provided commitments (%d vs %d)", len(fetchedComms), len(comms))
	}

	for i, comm := range comms {
		if !comm.Equals(fetchedComms[i]) {
			provider.log.Warn("commitment does not match expected", "first", firstHeight, "i", i, "comm", comm, "expected", fetchedComms[i])
			return false, nil
		}
	}

	return true, nil
}

func (provider *NodeKitProvider) L1BlockRefByNumber(ctx context.Context, num uint64) (eth.L1BlockRef, error) {
	return provider.L1Fetcher.L1BlockRefByNumber(ctx, num)
}

func (provider *NodeKitProvider) FetchReceipts(ctx context.Context, blockHash common.Hash) (eth.BlockInfo, types.Receipts, error) {
	return provider.L1Fetcher.FetchReceipts(ctx, blockHash)
}
