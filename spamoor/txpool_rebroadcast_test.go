package spamoor

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethpandaops/spamoor/txtypes"
	"github.com/ethpandaops/spamoor/utils"
)

// rpcNode is a JSON-RPC endpoint that records the transactions sent to it. A node with
// knowsAll answers every transaction with "already known", like one still holding it.
type rpcNode struct {
	knowsAll bool
	nonceOf  map[string]uint64

	mu           sync.Mutex
	onChainNonce uint64
	received     []uint64
}

func (n *rpcNode) serve(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     json.RawMessage `json:"id"`
		Method string          `json:"method"`
		Params []string        `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	resp := map[string]any{"jsonrpc": "2.0", "id": req.ID}
	switch req.Method {
	case "eth_getTransactionCount":
		resp["result"] = hexutil.Uint64(n.onChainNonce)
	case "eth_sendRawTransaction":
		nonce := n.nonceOf[req.Params[0]]
		known := n.knowsAll || slices.Contains(n.received, nonce)
		n.received = append(n.received, nonce)
		if known {
			resp["error"] = map[string]any{"code": -32000, "message": "already known"}
		} else {
			resp["result"] = "0x"
		}
	default:
		resp["error"] = map[string]any{"code": -32601, "message": "method not found"}
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (n *rpcNode) receivedNonces() []uint64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	return slices.Clone(n.received)
}

func (n *rpcNode) setOnChainNonce(nonce uint64) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.onChainNonce = nonce
}

type rebroadcastFixture struct {
	pool   *TxPool
	wallet *Wallet
	txs    []*txtypes.Transaction
	nodes  []*rpcNode
}

// newRebroadcastFixture creates a wallet with txCount pending transactions (nonces from
// 0) and a pool connected to one client per entry of knowsAll.
func newRebroadcastFixture(t *testing.T, txCount int, knowsAll ...bool) *rebroadcastFixture {
	t.Helper()

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	wallet := NewWallet(key, crypto.PubkeyToAddress(key.PublicKey))
	wallet.SetChainId(big.NewInt(1337))
	wallet.SetBalance(big.NewInt(1e18))

	nonceOf := map[string]uint64{}
	txs := make([]*txtypes.Transaction, txCount)
	for i := range txs {
		tx, err := wallet.BuildDynamicFeeTx(&txtypes.DynamicFeeTx{
			GasTipCap: big.NewInt(1),
			GasFeeCap: big.NewInt(1),
			Gas:       21000,
			To:        &wallet.address,
			Value:     big.NewInt(0),
		})
		if err != nil {
			t.Fatal(err)
		}

		raw, err := tx.MarshalNetwork()
		if err != nil {
			t.Fatal(err)
		}

		nonceOf[hexutil.Encode(raw)] = tx.Nonce()
		txs[i] = tx
		wallet.getTxNonceChan(tx, &SendTransactionOptions{Rebroadcast: true})
	}

	clientPool := &ClientPool{}
	nodes := make([]*rpcNode, len(knowsAll))
	for i, known := range knowsAll {
		nodes[i] = &rpcNode{knowsAll: known, nonceOf: nonceOf}

		server := httptest.NewServer(http.HandlerFunc(nodes[i].serve))
		t.Cleanup(server.Close)

		client, err := NewClient(&ClientOptions{RpcHost: server.URL})
		if err != nil {
			t.Fatal(err)
		}

		clientPool.goodClients = append(clientPool.goodClients, client)
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	pool := &TxPool{options: &TxPoolOptions{Context: ctx, ClientPool: clientPool}}

	return &rebroadcastFixture{pool: pool, wallet: wallet, txs: txs, nodes: nodes}
}

func waitForNonces(t *testing.T, node *rpcNode, want []uint64) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for !slices.Equal(node.receivedNonces(), want) {
		if time.Now().After(deadline) {
			t.Fatalf("node received nonces %v, want %v", node.receivedNonces(), want)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func nonceRange(from, to uint64) []uint64 {
	nonces := make([]uint64, 0, to-from)
	for nonce := from; nonce < to; nonce++ {
		nonces = append(nonces, nonce)
	}

	return nonces
}

// A client that already knows a transaction will not gossip it again, so stopping there
// left the peers that dropped it without it for good. Only a fresh acceptance ends the
// rebroadcast.
func TestRebroadcastContinuesPastClientsThatAlreadyKnowTheTx(t *testing.T) {
	f := newRebroadcastFixture(t, 1, true, false, false)

	f.pool.rebroadcastTransaction(context.Background(), f.txs[0], &SendTransactionOptions{}, 1, 0)

	for i, want := range [][]uint64{{0}, {0}, nil} {
		if got := f.nodes[i].receivedNonces(); !slices.Equal(got, want) {
			t.Errorf("node %d received %v, want %v", i, got, want)
		}
	}
}

// The start client used to be derived from the tx hash alone, so every retry hit the same
// client first. Successive retries must start at successive clients.
func TestRebroadcastStartsAtGivenClientOffset(t *testing.T) {
	f := newRebroadcastFixture(t, 1, false, false, false)

	for offset := range 3 {
		f.pool.rebroadcastTransaction(context.Background(), f.txs[0], &SendTransactionOptions{}, 1, offset)

		if got := f.nodes[offset].receivedNonces(); !slices.Equal(got, []uint64{0}) {
			t.Errorf("offset %d: first client received %v, want [0]", offset, got)
		}
	}
}

// Overloaded pools trim the tails of long accounts, leaving a run of nonces that no pool
// holds. A stuck wallet must resend a window of its lowest pending transactions in nonce
// order, and when the chain stops right at the window's end, resend the next window
// without waiting for another stale period.
func TestStaleWalletRebroadcastsWindowAndResumesAtItsEnd(t *testing.T) {
	const txCount = 40

	f := newRebroadcastFixture(t, txCount, false)
	node := f.nodes[0]

	f.pool.processStaleConfirmations(staleConfirmationBlocks+1, f.wallet)
	waitForNonces(t, node, nonceRange(0, rebroadcastWindow))

	// The chain includes the window and stops at its end.
	node.setOnChainNonce(rebroadcastWindow)
	lastTx := f.txs[rebroadcastWindow-1]
	f.pool.processTransactionInclusion(staleConfirmationBlocks+2, f.wallet, lastTx, &txtypes.Receipt{}, &utils.TxFees{})

	f.pool.processStaleConfirmations(staleConfirmationBlocks+3, f.wallet)
	waitForNonces(t, node, nonceRange(0, 2*rebroadcastWindow))

	// Without further progress the next round waits for the stale period again.
	f.pool.processStaleConfirmations(staleConfirmationBlocks+4, f.wallet)
	time.Sleep(200 * time.Millisecond)

	if got := node.receivedNonces(); len(got) != 2*rebroadcastWindow {
		t.Fatalf("expected no rebroadcast before the stale period, node received %v", got[2*rebroadcastWindow:])
	}

	if end := f.wallet.rebroadcastEnd; end != 2*rebroadcastWindow {
		t.Fatalf("rebroadcastEnd = %d, want %d", end, 2*rebroadcastWindow)
	}
}

// The window is paced by its lowest transaction: while that one is backing off, none of
// the window is resent.
func TestStaleWindowWaitsForLowestTxBackoff(t *testing.T) {
	f := newRebroadcastFixture(t, 4, false)
	node := f.nodes[0]

	f.pool.processStaleConfirmations(staleConfirmationBlocks+1, f.wallet)
	waitForNonces(t, node, nonceRange(0, 4))

	f.pool.processStaleConfirmations(2*staleConfirmationBlocks+2, f.wallet)
	time.Sleep(200 * time.Millisecond)

	if got := node.receivedNonces(); len(got) != 4 {
		t.Fatalf("expected the window to wait for its backoff, node received %v", got[4:])
	}
}
