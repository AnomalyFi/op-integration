package arcadia

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/flashbots/go-boost-utils/bls"
)

const HEADER_ROLLUP_SIG = "X-ROLLUP-SEQ-SIG"

const (
	ROLLUP_REGISTERED = iota
	ROLLUP_EXITED
	ROLLUP_NOT_REGISTERED
)

type ClientConfig struct {
	ArcadiaUrl         string
	Logger             log.Logger
	SequencerPubkey    []byte
	SequencerSecretKey []byte

	ChainID string
}

type RPCInterface interface {
	GetPayload(height uint64) ([]hexutil.Bytes, error)
	RollupStatus() (int, error)
}

var _ RPCInterface = (*Client)(nil)

type Client struct {
	cfg ClientConfig
	log log.Logger

	sk         *bls.SecretKey
	httpClient *http.Client
}

func NewArcadiaClient(cfg *ClientConfig) (*Client, error) {
	sk, err := bls.SecretKeyFromBytes(cfg.SequencerSecretKey)
	if err != nil {
		return nil, err
	}

	return &Client{
		cfg:        *cfg,
		log:        cfg.Logger,
		sk:         sk,
		httpClient: &http.Client{},
	}, nil
}

type GetPayloadRequest struct {
	ChainID     string `json:"chainId"`
	BlockNumber uint64 `json:"blockNumber"`
}

// the response we send back to the sidecar
type GetPayloadResponse struct {
	// First TOB txs then ROB txs
	Transactions []hexutil.Bytes `json:"transactions"`
}

func (c *Client) GetPayload(height uint64) ([]hexutil.Bytes, error) {
	payloadReq := GetPayloadRequest{
		ChainID:     c.cfg.ChainID,
		BlockNumber: height,
	}

	reqBytes, err := json.Marshal(payloadReq)
	if err != nil {
		return nil, err
	}
	reqHash, err := sha256HashPayload(reqBytes)
	if err != nil {
		return nil, err
	}
	sig := bls.Sign(c.sk, reqHash)
	sigBytes := sig.Bytes()

	req, err := http.NewRequest("POST", c.cfg.ArcadiaUrl, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set(HEADER_ROLLUP_SIG, hexutil.Encode(sigBytes[:]))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not 200: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payloadResp GetPayloadResponse
	if err := json.Unmarshal(body, &payloadResp); err != nil {
		return nil, err
	}

	return payloadResp.Transactions, nil
}

func (c *Client) RollupStatus() (int, error) {
	return ROLLUP_NOT_REGISTERED, nil
}

func sha256HashPayload(payload []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write(payload)
	if err != nil {
		return nil, err
	}
	bs := h.Sum(nil)
	return bs, nil
}
