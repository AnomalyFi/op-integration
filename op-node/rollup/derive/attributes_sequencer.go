package derive

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

type L1OriginSelectorIface interface {
	FindL1Origin(ctx context.Context, l2Head eth.L2BlockRef) (eth.L1BlockRef, error)
}

type AttributesSequencer struct {
	log     log.Logger
	metrics Metrics

	l1OriginSelector      L1OriginSelectorIface
	attrBuilder           AttributesBuilder
	broadcastPayloadAttrs func(id string, data []byte)
}

func NewAttributesSequencer(log log.Logger, l1OriginSelector L1OriginSelectorIface, attrBuilder AttributesBuilder, broadcastPayloadAttrs func(id string, data []byte), metrics Metrics) *AttributesSequencer {
	return &AttributesSequencer{
		log:     log,
		metrics: metrics,

		l1OriginSelector:      l1OriginSelector,
		attrBuilder:           attrBuilder,
		broadcastPayloadAttrs: broadcastPayloadAttrs,
	}
}

func (as *AttributesSequencer) PreparePayloadAttributes(ctx context.Context, l2Head eth.L2BlockRef) (*eth.PayloadAttributes, error) {
	// Figure out which L1 origin block we're going to be building on top of.
	l1Origin, err := as.l1OriginSelector.FindL1Origin(ctx, l2Head)
	if err != nil {
		as.log.Error("Error finding next L1 Origin", "err", err)
		return nil, err
	}

	if !(l2Head.L1Origin.Hash == l1Origin.ParentHash || l2Head.L1Origin.Hash == l1Origin.Hash) {
		// TODO: use different metrics for this
		// as.metrics.RecordSequencerInconsistentL1Origin(l2Head.L1Origin, l1Origin.ID())
		return nil, NewResetError(fmt.Errorf("cannot build new L2 block with L1 origin %s (parent L1 %s) on current L2 head %s with L1 origin %s", l1Origin, l1Origin.ParentHash, l2Head, l2Head.L1Origin))
	}

	as.log.Info("creating new block", "parent", l2Head, "l1Origin", l1Origin)

	fetchCtx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	attrs, err := as.attrBuilder.PreparePayloadAttributes(fetchCtx, l2Head, l1Origin.ID(), nil)
	if err != nil {
		return nil, err
	}

	as.log.Info("unmarshaling txns", "len(attrs.Transactions)", len(attrs.Transactions))
	txs := make(types.Transactions, 0, len(attrs.Transactions))
	for _, otx := range attrs.Transactions {
		tx := new(types.Transaction)
		tx.UnmarshalBinary(otx)
		txs = append(txs, tx)
	}

	builderAttrs := &eth.BuilderPayloadAttributes{
		Timestamp:             attrs.Timestamp,
		Random:                common.Hash(attrs.PrevRandao),
		SuggestedFeeRecipient: attrs.SuggestedFeeRecipient,
		Slot:                  l2Head.Number + 1,
		HeadHash:              l2Head.Hash,
		Withdrawals:           *attrs.Withdrawals,
		ParentBeaconBlockRoot: nil,
		Transactions:          txs,
		GasLimit:              uint64(*attrs.GasLimit),
	}

	// attrsEvent := &eth.BuilderPayloadAttributesEvent{
	// 	Version: "",
	// 	Data: eth.BuilderPayloadAttributesEventData{
	// 		ProposalSlot:    l2Head.Number + 1,
	// 		ParentBlockHash: l2Head.Hash,
	// 		PayloadAttributes: eth.BuilderPayloadAttributes{
	// 			Timestamp:             uint64(attrs.Timestamp),
	// 			PrevRandao:            common.Hash(attrs.PrevRandao),
	// 			SuggestedFeeRecipient: attrs.SuggestedFeeRecipient,
	// 			GasLimit:              uint64(*attrs.GasLimit),
	// 			Transactions:          txs,
	// 		},
	// 	},
	// }

	attrsData, err := json.Marshal(builderAttrs)
	if err != nil {
		return nil, err
	}

	as.log.Info("broadcasting new payload attributes", "json", attrsData)
	as.broadcastPayloadAttrs("payload_attributes", attrsData)
	return attrs, nil
}
