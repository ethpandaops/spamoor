package spamoor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/txtypes"
)

// frameSupportRetryInterval bounds how often an inactive chain is re-probed. Frame
// support only ever turns on, so a positive result is final and only negatives are
// retried.
const frameSupportRetryInterval = 30 * time.Second

// FrameSupport describes a chain's EIP-8141 frame transaction capability.
//
// It is derived from the predeploys the three EIPs install at activation, not from what
// a client says about a rejected transaction. Each EIP requires its address to be empty
// beforehand and installs code and nonce 1 at the fork, so the presence of that account
// is an exact and client-independent signal:
//
//	EIP-8141  EXPIRY_VERIFIER      0x…8141
//	EIP-8250  NONCE_MANAGER        0x…8250
//	EIP-8272  RECENT_ROOT_ADDRESS  0x…8272
//
// The extension set decides the envelope's wire layout, which is why it has to be known
// before a transaction can be encoded at all.
type FrameSupport struct {
	// Active reports whether the chain implements frame transactions.
	Active bool

	// Extensions are the envelope extensions the chain activates.
	Extensions txtypes.FrameExtensions
}

// frameSupportState is the pool's cached capability and the probe that fills it.
type frameSupportState struct {
	mutex     sync.Mutex
	support   FrameSupport
	lastProbe time.Time
}

// GetFrameSupport returns the cached frame transaction capability. It reports an
// inactive chain until a probe has run, so callers that need an answer should use
// GetFrameSupportWithInit.
func (pool *TxPool) GetFrameSupport() FrameSupport {
	pool.frameSupport.mutex.Lock()
	defer pool.frameSupport.mutex.Unlock()

	return pool.frameSupport.support
}

// GetFrameSupportWithInit returns the chain's frame transaction capability, probing it
// on first use.
//
// A positive result is cached for the process lifetime. A negative one is re-probed at
// most every frameSupportRetryInterval, because the predeploys appear at the fork rather
// than at genesis: a spammer started before activation would otherwise be stuck reading
// an inactive chain forever, and a transaction encoded on that reading is rejected on
// every send rather than once at startup.
func (pool *TxPool) GetFrameSupportWithInit(ctx context.Context) (FrameSupport, error) {
	pool.frameSupport.mutex.Lock()
	defer pool.frameSupport.mutex.Unlock()

	if pool.frameSupport.support.Active {
		return pool.frameSupport.support, nil
	}

	if !pool.frameSupport.lastProbe.IsZero() && time.Since(pool.frameSupport.lastProbe) < frameSupportRetryInterval {
		return pool.frameSupport.support, nil
	}

	client := pool.options.ClientPool.GetClient(WithoutBuilder())
	if client == nil {
		return FrameSupport{}, fmt.Errorf("no client available to probe frame transaction support")
	}

	probeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	support, err := probeFrameSupport(probeCtx, client)
	if err != nil {
		return FrameSupport{}, err
	}

	pool.frameSupport.support = support
	pool.frameSupport.lastProbe = time.Now()

	if support.Active {
		logrus.Infof("detected frame transaction support: envelope %s", support.Extensions)
	}

	return support, nil
}

// probeFrameSupport reads the capability from the predeploys.
func probeFrameSupport(ctx context.Context, client *Client) (FrameSupport, error) {
	frames, err := predeployActive(ctx, client, txtypes.ExpiryVerifier)
	if err != nil {
		return FrameSupport{}, fmt.Errorf("failed reading the EIP-8141 expiry verifier predeploy: %w", err)
	}

	if !frames {
		return FrameSupport{}, nil
	}

	support := FrameSupport{Active: true}

	keyed, err := predeployActive(ctx, client, txtypes.NonceManager)
	if err != nil {
		return FrameSupport{}, fmt.Errorf("failed reading the EIP-8250 nonce manager predeploy: %w", err)
	}

	if keyed {
		support.Extensions |= txtypes.FrameExtKeyedNonces
	}

	roots, err := predeployActive(ctx, client, txtypes.RecentRootAddress)
	if err != nil {
		return FrameSupport{}, fmt.Errorf("failed reading the EIP-8272 recent root predeploy: %w", err)
	}

	if roots {
		support.Extensions |= txtypes.FrameExtRecentRoots
	}

	return support, nil
}

// predeployActive reports whether a predeploy account has been installed.
//
// Activation sets both code and nonce 1. Nonce is checked as well as code because one of
// these codes is still TBD in its EIP, and an account with a non-zero nonce at an address
// the fork configuration required to be empty is unambiguous either way.
func predeployActive(ctx context.Context, client *Client, address common.Address) (bool, error) {
	code, err := client.GetCodeAt(ctx, address)
	if err != nil {
		return false, err
	}

	if len(code) > 0 {
		return true, nil
	}

	nonce, err := client.GetNonceAt(ctx, address, nil)
	if err != nil {
		return false, err
	}

	return nonce > 0, nil
}
