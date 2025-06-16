package testingutils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      interface{}   `json:"id"`
}

// JSONRPCError represents a JSON-RPC error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MockRPCServer implements a basic Ethereum JSON-RPC server for testing
type MockRPCServer struct {
	server *httptest.Server
	mu     sync.RWMutex

	// Mock data
	chainId         *big.Int
	nonce           uint64
	balance         *big.Int
	blockHeight     uint64
	blockHeightTime time.Time
	gasCap          *big.Int
	tipCap          *big.Int
	clientVersion   string
	receipt         *types.Receipt
	block           *types.Block
	blockReceipts   []*types.Receipt
	gasLimit        uint64

	// Mock errors
	err error
}

// NewMockRPCServer creates a new test RPC server
func NewMockRPCServer() *MockRPCServer {
	srv := &MockRPCServer{
		chainId:         big.NewInt(1337),
		nonce:           0,
		balance:         big.NewInt(0),
		blockHeight:     0,
		blockHeightTime: time.Now(),
		gasCap:          big.NewInt(1000000000),
		tipCap:          big.NewInt(1000000000),
		clientVersion:   "test-client/v1.0.0",
		gasLimit:        30000000,
	}

	// Create HTTP server with JSON-RPC handler
	httpServer := httptest.NewServer(http.HandlerFunc(srv.handleJSONRPC))
	srv.server = httpServer

	return srv
}

// handleJSONRPC handles JSON-RPC requests
func (s *MockRPCServer) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON-RPC request", http.StatusBadRequest)
		return
	}

	response := s.processJSONRPCRequest(&req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// processJSONRPCRequest processes a JSON-RPC request and returns a response
func (s *MockRPCServer) processJSONRPCRequest(req *JSONRPCRequest) *JSONRPCResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	response := &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	if s.err != nil {
		response.Error = &JSONRPCError{
			Code:    -32000,
			Message: s.err.Error(),
		}
		return response
	}

	switch req.Method {
	case "eth_chainId":
		response.Result = "0x" + s.chainId.Text(16)
	case "eth_getTransactionCount":
		response.Result = "0x" + strconv.FormatUint(s.nonce, 16)
	case "eth_getBalance":
		response.Result = "0x" + s.balance.Text(16)
	case "eth_gasPrice":
		response.Result = "0x" + s.gasCap.Text(16)
	case "eth_maxPriorityFeePerGas":
		response.Result = "0x" + s.tipCap.Text(16)
	case "eth_sendRawTransaction", "eth_sendTransaction":
		response.Result = "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	case "eth_getTransactionReceipt":
		if s.receipt == nil {
			response.Result = nil
		} else {
			// Ensure logs field is initialized
			if s.receipt.Logs == nil {
				s.receipt.Logs = []*types.Log{}
			}
			response.Result = s.receipt
		}
	case "eth_blockNumber":
		response.Result = "0x" + strconv.FormatUint(s.blockHeight, 16)
	case "web3_clientVersion":
		response.Result = s.clientVersion
	case "eth_getBlockByNumber":
		if s.block == nil {
			// Return a minimal block with just the header
			header := &types.Header{
				ParentHash:  common.Hash{},
				UncleHash:   types.EmptyUncleHash,
				Coinbase:    common.Address{},
				Root:        types.EmptyRootHash,
				TxHash:      types.EmptyTxsHash,
				ReceiptHash: types.EmptyReceiptsHash,
				Bloom:       types.Bloom{},
				Difficulty:  big.NewInt(0),
				Number:      big.NewInt(int64(s.blockHeight)),
				GasLimit:    s.gasLimit,
				GasUsed:     0,
				Time:        uint64(s.blockHeightTime.Unix()),
				Extra:       []byte{},
				MixDigest:   common.Hash{},
				Nonce:       types.BlockNonce{},
			}
			response.Result = types.NewBlock(header, nil, nil, nil)
		} else {
			response.Result = s.block
		}
	case "eth_getBlockReceipts":
		if s.blockReceipts == nil {
			response.Result = []*types.Receipt{}
		} else {
			// Ensure logs field is initialized for all receipts
			for _, receipt := range s.blockReceipts {
				if receipt.Logs == nil {
					receipt.Logs = []*types.Log{}
				}
			}
			response.Result = s.blockReceipts
		}
	default:
		response.Error = &JSONRPCError{
			Code:    -32601,
			Message: fmt.Sprintf("Method not found: %s", req.Method),
		}
	}

	return response
}

// Close closes the test server
func (s *MockRPCServer) Close() {
	s.server.Close()
}

// URL returns the server's URL
func (s *MockRPCServer) URL() string {
	return s.server.URL
}

// SetMockError sets an error to be returned by mock methods
func (s *MockRPCServer) SetMockError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.err = err
}

// SetMockChainId sets the mock chain ID
func (s *MockRPCServer) SetMockChainId(chainId *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.chainId = chainId
}

// SetMockNonce sets the mock nonce
func (s *MockRPCServer) SetMockNonce(nonce uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nonce = nonce
}

// SetMockBalance sets the mock balance
func (s *MockRPCServer) SetMockBalance(balance *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.balance = balance
}

// SetMockBlockHeight sets the mock block height
func (s *MockRPCServer) SetMockBlockHeight(height uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockHeight = height
	s.blockHeightTime = time.Now()
	// Update block number in block if it exists
	if s.block != nil {
		header := s.block.Header()
		header.Number = big.NewInt(int64(height))
		s.block = types.NewBlockWithHeader(header)
	}
}

// SetMockGasFees sets the mock gas fees
func (s *MockRPCServer) SetMockGasFees(gasCap, tipCap *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gasCap = gasCap
	s.tipCap = tipCap
}

// SetMockReceipt sets the mock receipt
func (s *MockRPCServer) SetMockReceipt(receipt *types.Receipt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.receipt = receipt
}

// SetMockBlock sets the mock block
func (s *MockRPCServer) SetMockBlock(block *types.Block) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.block = block
}

// SetMockBlockReceipts sets the mock block receipts
func (s *MockRPCServer) SetMockBlockReceipts(receipts []*types.Receipt) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockReceipts = receipts
}

// SetMockGasLimit sets the mock gas limit
func (s *MockRPCServer) SetMockGasLimit(gasLimit uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gasLimit = gasLimit
}
