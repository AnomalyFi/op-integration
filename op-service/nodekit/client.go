package nodekit

import (
	"context"
	"fmt"
	"strings"

	trpc "github.com/AnomalyFi/nodekit-seq/rpc"
	"github.com/AnomalyFi/nodekit-seq/types"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ethereum/go-ethereum/log"
)

type Client struct {
	baseUrl string
	client  *trpc.JSONRPCClient
	log     log.Logger
}

func NewClient(log log.Logger, url string) *Client {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	id := "22tNRzLv1PKasNMx9E6PmPy9QDXDkrwAqNnBRa3ZVBHSMBbBq6"

	urlNew := "http://127.0.0.1:37657/ext/bc/22tNRzLv1PKasNMx9E6PmPy9QDXDkrwAqNnBRa3ZVBHSMBbBq6"

	idReal, _ := ids.FromString(id)
	cli := trpc.NewJSONRPCClient(urlNew, 1337, idReal)

	return &Client{
		//baseUrl: url,
		client: cli,
		log:    log,
	}
}

// TODO really need
func (c *Client) FetchHeadersForWindow(ctx context.Context, start uint64, end uint64) (WindowStart, error) {
	//var res WindowStart
	//getBlockHeadersByStart

	res, err := c.client.GetBlockHeadersByStart(context.Background(), int64(start), int64(end))

	if err != nil {
		return WindowStart{}, err
	}

	blocks := make([]Header, len(res.Blocks))
	for i, blk := range res.Blocks {
		blocks[i] = convertBlockInfoToHeader(blk)
	}

	prev := convertBlockInfoToHeader(res.Prev)
	next := convertBlockInfoToHeader(res.Next)

	w := WindowStart{
		From:   res.From,
		Window: blocks,
		Prev:   &prev,
		Next:   &next,
	}

	return w, nil
}

func (c *Client) FetchRemainingHeadersForWindow(ctx context.Context, from uint64, end uint64) (WindowMore, error) {
	//var res WindowMore
	//getBlockHeadersByHeight
	res, err := c.client.GetBlockHeadersByHeight(context.Background(), from, int64(end))

	if err != nil {
		return WindowMore{}, err
	}

	blocks := make([]Header, len(res.Blocks))
	for i, blk := range res.Blocks {
		blocks[i] = convertBlockInfoToHeader(blk)
	}

	next := convertBlockInfoToHeader(res.Next)

	w := WindowMore{
		Window: blocks,
		Next:   &next,
	}

	return w, nil
}

func (c *Client) FetchTransactionsInBlock(ctx context.Context, block uint64, header *Header, namespace uint64) (TransactionsInBlock, error) {
	//var res NamespaceResponse
	res, err := c.client.GetBlockTransactionsByNamespace(context.Background(), block, string(namespace))

	if err != nil {
		return TransactionsInBlock{}, err
	}

	nRes := convertSEQTransactionResponseToNamespaceResponse(res)
	return nRes.Validate(header, namespace)
	//res.Validate(header, namespace)
}

// Function to convert SEQTransaction to Transaction
func convertSEQTransactionToTransaction(seqTx *types.SEQTransaction) Transaction {
	return Transaction{
		ChainId: seqTx.Namespace, // Use Index as ChainId (you can modify this as needed)
		Data:    seqTx.Transaction,
	}
}

// Function to convert SEQTransaction to Transaction
func convertBlockInfoToHeader(blockInfo trpc.BlockInfo) Header {
	return Header{
		TransactionsRoot: NmtRoot{
			Root: blockInfo.BlockId[:],
		}, // Use Index as ChainId (you can modify this as needed)
		Metadata: Metadata{
			Timestamp: uint64(blockInfo.Timestamp),
			L1Head:    blockInfo.L1Head,
		},
	}
}

// Function to convert SEQTransactionResponse to NamespaceResponse
func convertSEQTransactionResponseToNamespaceResponse(seqResponse *trpc.SEQTransactionResponse) *NamespaceResponse {
	transactions := make([]Transaction, len(seqResponse.Txs))
	for i, seqTx := range seqResponse.Txs {
		transactions[i] = convertSEQTransactionToTransaction(seqTx)
	}

	return &NamespaceResponse{
		Transactions: &transactions,
	}
}

type NamespaceResponse struct {
	//Proof        *json.RawMessage `json:"proof"`
	Transactions *[]Transaction `json:"transactions"`
}

// Validate a NamespaceResponse and extract the transactions.
// NMT proof validation is currently stubbed out.
func (res *NamespaceResponse) Validate(header *Header, namespace uint64) (TransactionsInBlock, error) {
	// if res.Proof == nil {
	// 	return TransactionsInBlock{}, fmt.Errorf("field proof of type NamespaceResponse is required")
	// }
	if res.Transactions == nil {
		return TransactionsInBlock{}, fmt.Errorf("field transactions of type NamespaceResponse is required")
	}

	// Check that these transactions are only and all of the transactions from `namespace` in the
	// block with `header`.
	//proof := NmtProof(*res.Proof)
	// if err := proof.Validate(header.TransactionsRoot, *res.Transactions); err != nil {
	// 	return TransactionsInBlock{}, err
	// }

	// Extract the transactions.
	var txs []Bytes
	for i, tx := range *res.Transactions {
		if tx.ChainId != string(namespace) {
			return TransactionsInBlock{}, fmt.Errorf("transaction %d has wrong namespace (%d, expected %d)", i, tx.ChainId, namespace)
		}
		txs = append(txs, tx.Data)
	}

	return TransactionsInBlock{
		Transactions: txs,
		//Proof:        proof,
	}, nil
}

// func (c *Client) get(ctx context.Context, out any, format string, args ...any) error {
// 	url := c.baseUrl + fmt.Sprintf(format, args...)

// 	c.log.Debug("get", "url", url)
// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
// 	if err != nil {
// 		c.log.Error("failed to build request", "err", err, "url", url)
// 		return err
// 	}
// 	res, err := c.client.Do(req)
// 	if err != nil {
// 		c.log.Error("error in request", "err", err, "url", url)
// 		return err
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != 200 {
// 		// Try to get the response body to include in the error message, as it may have useful
// 		// information about why the request failed. If this call fails, the response will be `nil`,
// 		// which is fine to include in the log, so we can ignore errors.
// 		body, _ := io.ReadAll(res.Body)
// 		c.log.Error("request failed", "err", err, "url", url, "status", res.StatusCode, "response", string(body))
// 		return fmt.Errorf("request failed with status %d", res.StatusCode)
// 	}

// 	// Read the response body into memory before we unmarshal it, rather than passing the io.Reader
// 	// to the json decoder, so that we still have the body and can inspect it if unmarshalling
// 	// failed.
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		c.log.Error("failed to read response body", "err", err, "url", url)
// 		return err
// 	}
// 	if err := json.Unmarshal(body, out); err != nil {
// 		c.log.Error("failed to parse body as json", "err", err, "url", url, "response", string(body))
// 		return err
// 	}
// 	c.log.Debug("request completed successfully", "url", url, "res", res, "body", string(body), "out", out)
// 	return nil
// }
