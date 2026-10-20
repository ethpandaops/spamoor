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

// FrameSupport describes a chain's EIP-8141 frame transaction capability, derived from
// the contracts the three EIPs deploy with keyless deployment transactions:
//
//	EIP-8141  EXPIRY_VERIFIER      0x81413f0cF12e9b6a49B1D0439E081c577D57FfFf
//	EIP-8250  NONCE_MANAGER        0x8250968C12e01A19d6F667b9B2F3b3A4d0e51cB7
//	EIP-8272  RECENT_ROOT_ADDRESS  0x8272D9679689Ea2f307140CdF9002D27dC00Ffff
//
// These are ordinary contracts and may be in state before the fork activates (devnet
// genesis files carry them from block 0), so their presence settles the envelope layout
// but not whether the chain accepts frame transactions yet.
//
// The extension set decides the envelope's wire layout, so it has to be known before a
// transaction can be encoded. EIP-8272 leaves the envelope alone, its verifier is an
// ordinary frame, so it is reported separately.
type FrameSupport struct {
	// Active reports whether the chain implements frame transactions.
	Active bool

	// Extensions are the envelope extensions the chain activates.
	Extensions txtypes.FrameExtensions

	// RecentRoots reports whether the EIP-8272 recent root contract is installed, so a
	// transaction may carry a recent root verifier frame.
	RecentRoots bool
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

// GetFrameSupportWithInit returns the chain's frame transaction capability, probing it on
// first use.
//
// A positive result is cached for the process lifetime; a negative one is re-probed,
// because on a chain that deploys the contracts by transaction they appear after genesis
// and a spammer started before then would otherwise be stuck reading an inactive chain.
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
		logrus.Infof("detected frame transaction support: envelope %s, recent roots %v", support.Extensions, support.RecentRoots)
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

	support.RecentRoots = roots

	return support, nil
}

// predeployActive reports whether a predeploy account has been installed. Nonce is
// checked as well as code so that a genesis carrying an older code still counts.
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

// AwaitFrameSupport returns the chain's frame transaction capability, blocking until the
// chain activates it.
//
// On a chain that deploys the contracts by transaction they appear after genesis, so a
// spammer started from genesis simply has to wait for them. Where genesis already carries
// them this returns before the fork; scenarios that read a refusal as an answer have to
// wait for a frame transaction to land as well (see frametx-fuzz's awaitFrameAcceptance).
// That is a normal state rather than a failure, so it is logged once as a warning and
// re-probed until the fork lands or the context is cancelled. Probe errors are treated
// the same way: a node that is not answering yet is not a chain without frames.
func (pool *TxPool) AwaitFrameSupport(ctx context.Context, logger logrus.FieldLogger) (FrameSupport, error) {
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	waiting := false

	for {
		support, err := pool.GetFrameSupportWithInit(ctx)

		switch {
		case err == nil && support.Active:
			if waiting {
				logger.Infof("frame transactions are now active on this chain")
			}

			return support, nil

		case !waiting:
			waiting = true

			if err != nil {
				logger.Warnf("could not probe frame transaction support (%v), waiting for the chain to become available", err)
			} else {
				logger.Warnf("no account at the EIP-8141 expiry verifier predeploy %s, this chain does not implement frame transactions yet, waiting for activation", txtypes.ExpiryVerifier)
			}
		}

		select {
		case <-ctx.Done():
			return FrameSupport{}, ctx.Err()
		case <-time.After(frameSupportRetryInterval):
		}
	}
}
