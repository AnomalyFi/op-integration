package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/peterbourgon/ff/v3"

	nodekittx "github.com/AnomalyFi/nodekit-sdk/tx"
)

// Environment variables beginning with this prefix can be used to instantiate command line flags
const ENV_PREFIX = "OP_GETH_PROXY"

var (
	fs            = flag.NewFlagSet("proxy", flag.ContinueOnError)
	listenAddr    = fs.String("listen-addr", "127.0.0.1:9090", "proxy's listening address")
	sequencerAddr = fs.String("seq-addr", "http://127.0.0.1:50000", "address of NodeKit SEQ")
	gethAddr      = fs.String("geth-addr", "http://127.0.0.1:8545", "address of the op-geth node")
	vm_id         = fs.Int("vm-id", 1, "VM ID of the OP rollup instance")
	chain_id      = fs.String("chain-id", "node", "Chain ID of SEQ instance")
)

type Transaction struct {
	Vm      int   `json:"vm"`
	Payload []int `json:"payload"`
}

type rpcMessage struct {
	Params []json.RawMessage `json:"params,omitempty"`
	Method string            `json:"method,omitempty"`
}

func ForwardToSequencer(message rpcMessage) {
	var hexString string
	if err := json.Unmarshal(message.Params[0], &hexString); err != nil {
		log.Println("Error decoding params field of the rpc message.")
		return
	}
	var txnBytes, err = hex.DecodeString(hexString[2:])
	if err != nil {
		log.Println("Error decoding hex string. The first raw transaction parameter must be valid hex.")
		return
	}

	log.Println("Transaction received, forwarding to sequencer.")
	err = nodekittx.BuildAndSendTransaction(*sequencerAddr, *chain_id, string(*vm_id), txnBytes)
	if err != nil {
		log.Println("Failed to submit to SEQ: ", err)
		return
	}
}

type baseHandle struct{}

func (h *baseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	toUrl, err := url.Parse(*gethAddr)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(toUrl)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	// Once we've read the body, we need to replace it with another reader because
	// ReadAll can only be called once: https://blog.flexicondev.com/read-go-http-request-body-multiple-times
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var message rpcMessage
	if err := json.Unmarshal(body, &message); err != nil {
		log.Println("Invalid request: expected RPC message")
		return
	}
	// Check for sendRawTransaction
	if message.Method == "eth_sendRawTransaction" {
		ForwardToSequencer(message)
	}
	proxy.ServeHTTP(w, r)
}

func main() {
	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix(ENV_PREFIX)); err != nil {
		panic(err)
	}

	h := &baseHandle{}
	http.Handle("/", h)

	log.Println("Starting proxy server on", *listenAddr)
	server := &http.Server{
		Addr:    *listenAddr,
		Handler: h,
	}
	log.Fatal(server.ListenAndServe())
}
