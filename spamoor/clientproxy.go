package spamoor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

// ErrNoReadyClient is returned by a proxy client when no client matches its selection options.
var ErrNoReadyClient = errors.New("no ready client available")

// clientProxy forwards every call of a proxy Client to a client that is at the
// pool's highest known block at the time of the call.
type clientProxy struct {
	pool    *ClientPool
	options *clientSelectionOptions
}

// newProxyClient builds a Client whose RPC calls, including those made through
// GetEthClient, are routed per call via selectReadyClient.
func newProxyClient(pool *ClientPool, options *clientSelectionOptions) (*Client, error) {
	client := &Client{
		rpchost:      "proxy",
		clientGroups: []string{"default"},
		clientType:   ClientTypeClient,
		enabled:      true,
		nameOverride: "proxy",
		proxy:        &clientProxy{pool: pool, options: options},
	}
	client.logger = pool.logger.WithField("rpc", "proxy")

	if options.group != "" {
		client.clientGroups = []string{options.group}
	}

	// The ethclient path still needs an rpc.Client; its HTTP transport hands each
	// request back to the proxy instead of a socket.
	rpcClient, err := rpc.DialOptions(context.Background(), "http://proxy", rpc.WithHTTPClient(&http.Client{
		Transport: &proxyTransport{client: client},
	}))
	if err != nil {
		return nil, err
	}

	client.rpcClient = rpcClient

	return client, nil
}

// backend picks the client that serves the next call of this proxy.
func (proxy *clientProxy) backend() (*Client, error) {
	client := proxy.pool.selectReadyClient(proxy.options)
	if client == nil {
		return nil, ErrNoReadyClient
	}

	return client, nil
}

// proxyTransport serves the proxy's rpc.Client: it decodes each JSON-RPC request,
// forwards it through a ready backend and encodes the backend's answer.
type proxyTransport struct {
	client *Client
}

type proxyMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *proxyError     `json:"error,omitempty"`
}

type proxyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (t *proxyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return nil, err
	}

	isBatch := bytes.HasPrefix(bytes.TrimSpace(body), []byte("["))
	var requests []*proxyMessage
	if isBatch {
		err = json.Unmarshal(body, &requests)
	} else {
		var single proxyMessage
		err = json.Unmarshal(body, &single)
		requests = []*proxyMessage{&single}
	}
	if err != nil {
		return nil, fmt.Errorf("proxy: invalid json-rpc request: %w", err)
	}

	// One backend per HTTP request so a batch is answered from one state view.
	backend, err := t.client.proxy.backend()
	if err != nil {
		return nil, err
	}

	ctx, cancel := backend.getContext(req.Context())
	defer cancel()

	batch := make([]rpc.BatchElem, len(requests))
	results := make([]json.RawMessage, len(requests))
	for i, request := range requests {
		args, err := proxyArgs(request.Params)
		if err != nil {
			return nil, fmt.Errorf("proxy: invalid params for %s: %w", request.Method, err)
		}

		batch[i] = rpc.BatchElem{Method: request.Method, Args: args, Result: &results[i]}
	}

	if len(batch) == 1 {
		batch[0].Error = backend.rpcClient.CallContext(ctx, batch[0].Result, batch[0].Method, batch[0].Args...)
		backend.incrementRequestStats(batch[0].Method == "eth_sendRawTransaction", batch[0].Error != nil)
	} else if err := backend.rpcClient.BatchCallContext(ctx, batch); err != nil {
		backend.incrementRequestStats(false, true)
		return nil, err
	} else {
		backend.incrementRequestStats(false, false)
	}

	responses := make([]*proxyMessage, len(requests))
	for i, request := range requests {
		responses[i] = &proxyMessage{Version: "2.0", ID: request.ID}
		switch {
		case errors.Is(batch[i].Error, rpc.ErrNoResult):
			responses[i].Result = json.RawMessage("null")
		case batch[i].Error != nil:
			responses[i].Error = toProxyError(batch[i].Error)
		default:
			responses[i].Result = results[i]
		}
	}

	var out []byte
	if isBatch {
		out, err = json.Marshal(responses)
	} else {
		out, err = json.Marshal(responses[0])
	}
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        http.Header{"Content-Type": []string{"application/json"}},
		Body:          io.NopCloser(bytes.NewReader(out)),
		ContentLength: int64(len(out)),
		Request:       req,
	}, nil
}

// proxyArgs splits an encoded params array into raw arguments that re-encode verbatim.
func proxyArgs(params json.RawMessage) ([]any, error) {
	if len(params) == 0 || bytes.Equal(bytes.TrimSpace(params), []byte("null")) {
		return nil, nil
	}

	var raw []json.RawMessage
	if err := json.Unmarshal(params, &raw); err != nil {
		return nil, err
	}

	args := make([]any, len(raw))
	for i, arg := range raw {
		args[i] = arg
	}

	return args, nil
}

// toProxyError keeps the backend's JSON-RPC error code and data when it has them.
func toProxyError(err error) *proxyError {
	proxyErr := &proxyError{Code: -32603, Message: err.Error()}

	var rpcErr rpc.Error
	if errors.As(err, &rpcErr) {
		proxyErr.Code = rpcErr.ErrorCode()
	}

	var dataErr rpc.DataError
	if errors.As(err, &dataErr) {
		proxyErr.Data = dataErr.ErrorData()
	}

	return proxyErr
}
