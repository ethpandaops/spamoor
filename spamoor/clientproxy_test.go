package spamoor

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
)

func testClient(name string, head uint64, clientType ClientType) *Client {
	return &Client{
		rpchost:      name,
		nameOverride: name,
		clientGroups: []string{"default"},
		clientType:   clientType,
		enabled:      true,
		blockHeight:  head,
	}
}

func TestSelectReadyClient(t *testing.T) {
	ahead := testClient("ahead", 100, ClientTypeClient)
	ahead2 := testClient("ahead2", 100, ClientTypeClient)
	behind := testClient("behind", 98, ClientTypeClient)
	builder := testClient("builder", 0, ClientTypeBuilder)
	pool := &ClientPool{logger: logrus.New(), goodClients: []*Client{behind, ahead, builder, ahead2}}

	opts := defaultClientSelectionOptions()
	seen := map[string]int{}
	for i := 0; i < 4; i++ {
		seen[pool.selectReadyClient(opts).GetName()]++
	}
	if seen["ahead"] != 2 || seen["ahead2"] != 2 {
		t.Fatalf("expected round-robin over the clients at the pool head, got %v", seen)
	}

	// A group that only has lagging clients uses its least stale member.
	behind.clientGroups = []string{"slow"}
	slower := testClient("slower", 90, ClientTypeClient)
	slower.clientGroups = []string{"slow"}
	pool.goodClients = append(pool.goodClients, slower)
	slowOpts := defaultClientSelectionOptions()
	slowOpts.group = "slow"
	if got := pool.selectReadyClient(slowOpts); got != behind {
		t.Fatalf("expected the least stale slow client, got %v", got.GetName())
	}

	// Builders are only used when nothing else matches.
	builderOpts := defaultClientSelectionOptions()
	builderOpts.excludeTypes[ClientTypeClient] = true
	if got := pool.selectReadyClient(builderOpts); got != builder {
		t.Fatalf("expected the builder, got %v", got.GetName())
	}
}

type proxyTestService struct {
	head uint64
}

func (s *proxyTestService) BlockNumber() hexutil.Uint64 {
	return hexutil.Uint64(s.head)
}

func (s *proxyTestService) ChainId() *hexutil.Big {
	return (*hexutil.Big)(big.NewInt(1337))
}

func (s *proxyTestService) GetBalance(address string, block string) (*hexutil.Big, error) {
	return nil, &proxyTestError{}
}

type proxyTestError struct{}

func (e *proxyTestError) Error() string          { return "nope" }
func (e *proxyTestError) ErrorCode() int         { return -32042 }
func (e *proxyTestError) ErrorData() interface{} { return "detail" }

func TestProxyClientRoutesCalls(t *testing.T) {
	newBackend := func(name string, head uint64) *Client {
		server := rpc.NewServer()
		if err := server.RegisterName("eth", &proxyTestService{head: head}); err != nil {
			t.Fatal(err)
		}

		client := testClient(name, head, ClientTypeClient)
		client.rpcClient = rpc.DialInProc(server)

		return client
	}

	behind := newBackend("behind", 5)
	ahead := newBackend("ahead", 7)
	pool := &ClientPool{logger: logrus.New(), goodClients: []*Client{behind, ahead}}

	proxy := pool.GetClient(WithClientSelectionMode(SelectClientProxy))
	if proxy == nil || proxy.proxy == nil {
		t.Fatal("expected a proxy client")
	}
	if again := pool.GetClient(WithClientSelectionMode(SelectClientProxy)); again != proxy {
		t.Fatal("expected the proxy to be reused for the same options")
	}

	ctx := context.Background()

	// Direct RPC methods go through callRPC to the client at the pool head.
	for i := 0; i < 3; i++ {
		head, err := proxy.GetBlockNumber(ctx)
		if err != nil || head != 7 {
			t.Fatalf("expected head 7 from the ahead client, got %v (%v)", head, err)
		}
	}

	// The ethclient path goes through the HTTP transport.
	chainID, err := proxy.GetEthClient().ChainID(ctx)
	if err != nil || chainID.Int64() != 1337 {
		t.Fatalf("expected chain id 1337 via ethclient, got %v (%v)", chainID, err)
	}

	// JSON-RPC error codes and data survive the round trip.
	_, err = proxy.GetEthClient().BalanceAt(ctx, [20]byte{}, nil)
	var rpcErr rpc.DataError
	if err == nil || !errors.As(err, &rpcErr) || rpcErr.ErrorData() != "detail" {
		t.Fatalf("expected the backend error data, got %v", err)
	}

	// Batches are forwarded as a batch.
	var blockNumber, chainIDRaw json.RawMessage
	batch := []rpc.BatchElem{
		{Method: "eth_blockNumber", Result: &blockNumber},
		{Method: "eth_chainId", Result: &chainIDRaw},
	}
	if err := proxy.rpcClient.BatchCallContext(ctx, batch); err != nil {
		t.Fatal(err)
	}
	if string(blockNumber) != `"0x7"` || string(chainIDRaw) != `"0x539"` || batch[0].Error != nil || batch[1].Error != nil {
		t.Fatalf("unexpected batch results: %s %s %v %v", blockNumber, chainIDRaw, batch[0].Error, batch[1].Error)
	}

	// The ahead client, not the proxy, accounted for the requests.
	if total, _, _ := ahead.GetRequestStats(); total == 0 {
		t.Fatal("expected the backend to record the request stats")
	}
	if total, _, _ := behind.GetRequestStats(); total != 0 {
		t.Fatalf("expected the lagging client to stay idle, got %v requests", total)
	}
}
