package main

import (
	"bytes"
	"context"
	"encoding/binary"
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

	trpc "github.com/AnomalyFi/seq-sdk/client"
)

// Environment variables beginning with this prefix can be used to instantiate command line flags
const ENV_PREFIX = "OP_GETH_PROXY"

// TODO fix this
var (
	fs            = flag.NewFlagSet("proxy", flag.ContinueOnError)
	listenAddr    = fs.String("listen-addr", "127.0.0.1:9090", "proxy's listening address")
	sequencerAddr = fs.String("seq-addr", "http://127.0.0.1:35957/ext/bc/2CAt2HdXkWeTaezB6VCpPfPzDTK2t49Nrp5dtnuWnFY9e8NRaV", "address of NodeKit SEQ")
	gethAddr      = fs.String("geth-addr", "http://127.0.0.1:8545", "address of the op-geth node")
	vm_id         = fs.Int("vm-id", 1, "VM ID of the OP rollup instance")
	chain_id      = fs.String("chain-id", "2CAt2HdXkWeTaezB6VCpPfPzDTK2t49Nrp5dtnuWnFY9e8NRaV", "Chain ID of SEQ instance")
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

	cli := trpc.NewJSONRPCClient(*sequencerAddr, 1337, *chain_id)

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, vm_id)

	if err != nil {
		log.Println("Failed to write bytes before SEQ: ", err)
		return
	}

	_, err = cli.SubmitTx(context.Background(), *chain_id, 1337, buf.Bytes(), txnBytes)

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
