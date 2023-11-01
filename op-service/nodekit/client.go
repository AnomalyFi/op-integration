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

type Client struct {
	baseUrl string
	client  *trpc.JSONRPCClient
	log     log.Logger
}

func NewClient(log log.Logger, url string) *Client {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	id := "g4BiH3gadeH88eaGV2Gzr4QVot416jnkxcsqziSUkGnmGQavv"

	urlNew := "http://127.0.0.1:42723/ext/bc/g4BiH3gadeH88eaGV2Gzr4QVot416jnkxcsqziSUkGnmGQavv"

	cli := trpc.NewJSONRPCClient(urlNew, 1337, id)

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
	var next *Header

	start_time := start * 1000
	end_time := end * 1000

	id := "g4BiH3gadeH88eaGV2Gzr4QVot416jnkxcsqziSUkGnmGQavv"

	urlNew := "http://127.0.0.1:42723/ext/bc/g4BiH3gadeH88eaGV2Gzr4QVot416jnkxcsqziSUkGnmGQavv"

	cli := trpc.NewJSONRPCClient(urlNew, 1337, id)

	res, err := cli.GetBlockHeadersByStart(context.Background(), int64(start_time), int64(end_time))

	//res, err := c.client.GetBlockHeadersByStart(context.Background(), int64(start_time), int64(end_time))

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
		err = errors.New("Zero Length Id")
		c.log.Error("Error in FetchHeadersForWindow", "error", err)
	}

	prev, err := convertBlockInfoToHeader(res.Prev)
	if err != nil {
		return WindowStart{}, err
	}

	if !(res.Next == (types.BlockInfo{})) {
		//! TODO this is where the error is. It's on the c.log.Error line
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
	//var res WindowMore
	var next *Header
	//getBlockHeadersByHeight

	id := "g4BiH3gadeH88eaGV2Gzr4QVot416jnkxcsqziSUkGnmGQavv"

	urlNew := "http://127.0.0.1:42723/ext/bc/g4BiH3gadeH88eaGV2Gzr4QVot416jnkxcsqziSUkGnmGQavv"

	cli := trpc.NewJSONRPCClient(urlNew, 1337, id)

	end_time := end * 1000

	res, err := cli.GetBlockHeadersByHeight(context.Background(), from, int64(end_time))
	//c.client.GetBlockHeadersByHeight(context.Background(), from, int64(end))

	if err != nil {
		return WindowMore{}, err
	}

	blocks := make([]Header, len(res.Blocks))
	for i, blk := range res.Blocks {
		if len(blk.BlockId) == 0 {
			err = errors.New("Zero Length Id")
			c.log.Error("Error in FetchRemainingHeadersForWindow", "error", err)
		}
		t, err := convertBlockInfoToHeader(blk)
		if err != nil {
			return WindowMore{}, err
		}
		blocks[i] = *t
	}

	if !(res.Next == (types.BlockInfo{})) {
		//! TODO this is where the error is. It's on the c.log.Error line
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

func (c *Client) FetchTransactionsInBlock(ctx context.Context, block uint64, header *Header, namespace uint64) (TransactionsInBlock, error) {
	//var res NamespaceResponse
	id := "g4BiH3gadeH88eaGV2Gzr4QVot416jnkxcsqziSUkGnmGQavv"

	urlNew := "http://127.0.0.1:42723/ext/bc/g4BiH3gadeH88eaGV2Gzr4QVot416jnkxcsqziSUkGnmGQavv"

	cli := trpc.NewJSONRPCClient(urlNew, 1337, id)
	// TODO First I encode the integer to bytes form. Then I use hex.EncodeToString on it.
	//hex.EncodeToString(action.ChainId)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, namespace)
	hexNamespace := hex.EncodeToString(buf)
	res, err := cli.GetBlockTransactionsByNamespace(context.Background(), block, hexNamespace)

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
func convertBlockInfoToHeader(blockInfo types.BlockInfo) (*Header, error) {
	//TODO check this later
	bytes, err := DecodeCB58(blockInfo.BlockId)
	if err != nil {
		return nil, err
	}
	return &Header{
		TransactionsRoot: NmtRoot{
			Root: bytes,
		}, // Use Index as ChainId (you can modify this as needed)
		Metadata: Metadata{
			Timestamp: uint64(blockInfo.Timestamp),
			L1Head:    blockInfo.L1Head,
		},
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
