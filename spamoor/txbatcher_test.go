package spamoor

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// txBatcherWithNoClients builds a real TxBatcher whose deployment fails
// deterministically with no network: GetContractDeployment needs a client
// from the pool, and an empty ClientPool always returns nil.
func txBatcherWithNoClients() *TxBatcher {
	pool := &TxPool{options: &TxPoolOptions{ClientPool: &ClientPool{}}}
	factory := &DeploymentFactory{txpool: pool}
	return &TxBatcher{txpool: pool, factory: factory}
}

// A failed deployment must not leave the batcher marked as deployed with a
// zero address - that would make every subsequent Deploy() call silently
// report success while every funding transaction built against GetAddress()
// becomes a real-value transfer to the zero address instead of reaching the
// batcher contract.
func TestDeploy_FailedDeploymentDoesNotMarkAsDeployed(t *testing.T) {
	batcher := txBatcherWithNoClients()

	err := batcher.Deploy(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected Deploy to fail with no client available")
	}

	if batcher.isDeployed {
		t.Fatal("expected isDeployed to stay false after a failed deploy")
	}
	if batcher.GetAddress() != (common.Address{}) {
		t.Fatalf("expected address to stay zero after a failed deploy, got %s", batcher.GetAddress())
	}
}

// A failed deployment must be retriable: a later Deploy() call has to
// actually attempt deployment again rather than short-circuiting to a stale
// "success" left over from the earlier failure.
func TestDeploy_FailedDeploymentAllowsRetry(t *testing.T) {
	batcher := txBatcherWithNoClients()

	_ = batcher.Deploy(context.Background(), nil, nil)

	err := batcher.Deploy(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected the retry to also fail (still no client available), not silently report success")
	}
	if batcher.GetAddress() != (common.Address{}) {
		t.Fatalf("expected address to stay zero, got %s", batcher.GetAddress())
	}
}
