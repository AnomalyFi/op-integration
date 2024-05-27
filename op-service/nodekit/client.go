package nodekit

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	trpc "github.com/AnomalyFi/seq-sdk/client"
	"github.com/AnomalyFi/seq-sdk/types"
	"github.com/ethereum/go-ethereum/log"
)

const ENV_PREFIX = "NODEKIT"

type Client struct {
	baseUrl string
	client  *trpc.JSONRPCClient
	log     log.Logger
	chainID string
	seqAddr string
}

func NewClient(log log.Logger, url string) *Client {
	ss := strings.Split(url, "/")
	chainID := ss[len(ss)-1]

	cli := trpc.NewJSONRPCClient(url, 1337, chainID)

	return &Client{
		//baseUrl: url,
		client:  cli,
		log:     log,
		seqAddr: url,
		chainID: chainID,
	}
}

// TODO really need
func (c *Client) FetchHeadersForWindow(ctx context.Context, start uint64, end uint64) (WindowStart, error) {
	start_time := start
	end_time := end

	cli := c.client

	res, err := cli.GetBlockHeadersByStart(context.Background(), int64(start_time), int64(end_time))

	log.Info("seq info", "chain-id", c.chainID, "sequencer-addr", c.seqAddr)

	//TODO is this causing the error: We skipped an L1 block and the next L1 block is eligible as an origin, advancing by one
	if err != nil {
		c.log.Error("Error in GetBlockHeadersByStart", "error", err)
		return WindowStart{}, err
	}

	blocks := make([]Header, len(res.Blocks))
	for i, blk := range res.Blocks {
		t, err := convertBlockInfoToHeader(blk)
		if err != nil {
			return WindowStart{}, err
		}
		blocks[i] = *t
	}

	if len(res.Prev.BlockId) == 0 {
		err = errors.New("zero length id")
		c.log.Error("Error in FetchHeadersForWindow", "error", err)
	}

	prev, err := convertBlockInfoToHeader(res.Prev)
	if err != nil {
		return WindowStart{}, err
	}

	var next *Header

	if !(res.Next == (types.BlockInfo{})) {
		next, err = convertBlockInfoToHeader(res.Next)
		if err != nil {
			return WindowStart{}, err
		}
	}

	w := WindowStart{
		From:   res.From,
		Window: blocks,
		Prev:   prev,
		Next:   next,
	}

	return w, nil
}

func (c *Client) FetchRemainingHeadersForWindow(ctx context.Context, from uint64, end uint64) (WindowMore, error) {
	var next *Header

	end_time := end

	cli := c.client
	res, err := cli.GetBlockHeadersByHeight(context.Background(), from, int64(end_time))

	if err != nil {
		return WindowMore{}, err
	}

	blocks := make([]Header, len(res.Blocks))
	for i, blk := range res.Blocks {
		if len(blk.BlockId) == 0 {
			err = errors.New("zero length id")
			c.log.Error("Error in FetchRemainingHeadersForWindow", "error", err)
		}
		t, err := convertBlockInfoToHeader(blk)
		if err != nil {
			return WindowMore{}, err
		}
		blocks[i] = *t
	}

	if !(res.Next == (types.BlockInfo{})) {
		next, err = convertBlockInfoToHeader(res.Next)
		if err != nil {
			return WindowMore{}, err
		}
	}

	w := WindowMore{
		Window: blocks,
		Next:   next,
	}

	return w, nil
}

func (c *Client) FetchTransactionsInBlock(ctx context.Context, header *Header, namespace uint64) (TransactionsInBlock, error) {
	cli := c.client

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, namespace)
	hexNamespace := hex.EncodeToString(buf)
	res, err := cli.GetBlockTransactionsByNamespace(context.Background(), header.Height, hexNamespace)

	if err != nil {
		return TransactionsInBlock{}, err
	}

	nRes := convertSEQTransactionResponseToNamespaceResponse(res)
	return nRes.Validate(header, namespace)
}

// Function to convert SEQTransaction to Transaction
func convertSEQTransactionToTransaction(seqTx *types.SEQTransaction) Transaction {
	return Transaction{
		ChainId: seqTx.Namespace, // Use Index as ChainId (you can modify this as needed)
		Data:    seqTx.Transaction,
	}
}

// Function to convert SEQTransaction to Transaction
func convertBlockInfoToHeader(blockInfo types.BlockInfo) (*Header, error) {
	bytes, err := DecodeCB58(blockInfo.BlockId)
	if err != nil {
		return nil, err
	}
	tmp := blockInfo.Timestamp
	return &Header{
		TransactionsRoot: NmtRoot{
			Root: bytes,
		},
		Timestamp:         uint64(tmp),
		TimestampOriginal: uint64(blockInfo.Timestamp),
		L1Head:            blockInfo.L1Head,
		Height:            blockInfo.Height,
	}, nil
}

// Function to convert SEQTransactionResponse to NamespaceResponse
func convertSEQTransactionResponseToNamespaceResponse(seqResponse *types.SEQTransactionResponse) *NamespaceResponse {
	transactions := make([]Transaction, len(seqResponse.Txs))
	for i, seqTx := range seqResponse.Txs {
		transactions[i] = convertSEQTransactionToTransaction(seqTx)
	}

	return &NamespaceResponse{
		Transactions: &transactions,
	}
}

type NamespaceResponse struct {
	Transactions *[]Transaction `json:"transactions"`
}

// Validate a NamespaceResponse and extract the transactions.
// NMT proof validation is currently stubbed out.
func (res *NamespaceResponse) Validate(header *Header, namespace uint64) (TransactionsInBlock, error) {
	if res.Transactions == nil {
		return TransactionsInBlock{}, fmt.Errorf("field transactions of type NamespaceResponse is required")
	}

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, namespace)
	hexNamespace := hex.EncodeToString(buf)
	var txs []Bytes
	for i, tx := range *res.Transactions {
		if tx.ChainId != hexNamespace {
			return TransactionsInBlock{}, fmt.Errorf("transaction %d has wrong namespace (%s, expected %s)", i, tx.ChainId, hexNamespace)
		}
		txs = append(txs, tx.Data)
	}

	return TransactionsInBlock{
		Transactions: txs,
	}, nil
}
