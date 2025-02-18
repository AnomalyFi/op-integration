package sidecar

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/flashbots/go-boost-utils/bls"
)

const HEADER_ROLLUP_SIG = "X-ROLLUP-SEQ-SIG"

const (
	RollupStatusUnknown = iota
	RollupManagedByNodeKit
	RollupNotManagedByNodeKit
)

const (
	pathGetPayload   = "/rollup/getpayload"
	pathRollupStatus = "/rollup/status"
)

var (
	ErrPayloadNotManagedByNodekit = errors.New("payload at given height is not produced by NodeKit")
	ErrArcadiaDown                = errors.New("arcadia is down, need reorg")
)

type ClientConfig struct {
	SidecarUrl         string
	Logger             log.Logger
	SequencerPubkey    *bls.PublicKey
	SequencerSecretKey *bls.SecretKey

	ChainID string
}

type RPCInterface interface {
	GetPayload(height uint64) ([]hexutil.Bytes, error)
	GetPayloadFromDA(height uint64) ([]hexutil.Bytes, error)
	RollupStatus(height uint64) (int, error)
}

var _ RPCInterface = (*Client)(nil)

type Client struct {
	cfg ClientConfig
	log log.Logger

	chainID string

	sk         *bls.SecretKey
	pk         *bls.PublicKey
	httpClient *http.Client
}

func NewSidecarClient(cfg *ClientConfig) (*Client, error) {
	return &Client{
		cfg: *cfg,
		log: cfg.Logger,

		sk:         cfg.SequencerSecretKey,
		pk:         cfg.SequencerPubkey,
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
	endpoint := c.cfg.SidecarUrl + pathGetPayload

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

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(reqBytes))
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
		if resp.StatusCode != http.StatusBadRequest {
			return nil, ErrArcadiaDown
		}
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.log.Warn("unable to read response body", "err", err)
			return nil, fmt.Errorf("status code not 200: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("status code: %d, err: %s", resp.StatusCode, string(respBody))
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

// TODO: to be implemented
func (c *Client) GetPayloadFromDA(height uint64) ([]hexutil.Bytes, error) {
	return nil, nil
}

type RollupStatusRequest struct {
	ChainID string `json:"chainID"`
	Height  uint64 `json:"height"`
}

type RollupStatusResponse struct {
	Status int `json:"status"`
}

func (c *Client) RollupStatus(height uint64) (int, error) {
	endpoint := c.cfg.SidecarUrl + pathRollupStatus

	statusReq := RollupStatusRequest{
		ChainID: c.chainID,
		Height:  height,
	}

	reqBytes, err := json.Marshal(statusReq)
	if err != nil {
		return RollupStatusUnknown, err
	}
	reqHash, err := sha256HashPayload(reqBytes)
	if err != nil {
		return RollupStatusUnknown, err
	}
	sig := bls.Sign(c.sk, reqHash)
	sigBytes := sig.Bytes()

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return RollupStatusUnknown, err
	}
	req.Header.Set(HEADER_ROLLUP_SIG, hexutil.Encode(sigBytes[:]))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RollupStatusUnknown, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return RollupStatusUnknown, nil
	}

	if resp.StatusCode != http.StatusOK {
		c.log.Error("unable to query rollup status", "err", string(respBody))
		// either signature not correct or sidecar is down
		return RollupStatusUnknown, nil
	}

	statusResp := new(RollupStatusResponse)
	if err := json.Unmarshal(respBody, statusResp); err != nil {
		return RollupStatusUnknown, nil
	}
	return statusResp.Status, nil
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
