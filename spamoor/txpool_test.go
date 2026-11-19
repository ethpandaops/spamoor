package spamoor

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/txtypes"
)

// OnComplete is documented as always being called once processing ends. The early
// return for an already-cancelled context used to skip it, which left callers that
// release a WaitGroup from the callback (WalletPool.ReclaimFunds) hanging forever.
func TestSendTransactionCallsOnCompleteWhenContextCancelled(t *testing.T) {
	pool := &TxPool{}
	wallet := &Wallet{}
	tx := txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)
	err := pool.SendTransaction(ctx, wallet, tx, &SendTransactionOptions{
		OnComplete: func(_ *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
			if receipt != nil {
				t.Errorf("expected no receipt for a cancelled submission")
			}
			done <- err
		},
	})
	if err == nil {
		t.Fatal("expected SendTransaction to return the context error")
	}

	select {
	case cbErr := <-done:
		if cbErr == nil {
			t.Fatal("expected OnComplete to receive the context error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("OnComplete was not called for an already-cancelled context")
	}
}

// Mirrors the ReclaimFunds wiring: one send per wallet, wg.Done only from OnComplete.
func TestSendTransactionReleasesWaitGroupWhenContextCancelled(t *testing.T) {
	pool := &TxPool{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	const txCount = 20
	var wg sync.WaitGroup
	wg.Add(txCount)

	for i := range txCount {
		wallet := &Wallet{}
		tx := txtypes.NewTx(&txtypes.LegacyTx{Nonce: uint64(i)})

		go func() {
			_ = pool.SendTransaction(ctx, wallet, tx, &SendTransactionOptions{
				OnComplete: func(_ *txtypes.Transaction, _ *txtypes.Receipt, _ error) {
					wg.Done()
				},
			})
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("wait group was never released after the context was cancelled")
	}
}

func TestUpdateSlotDuration(t *testing.T) {
	pool := &TxPool{options: &TxPoolOptions{Logger: logrus.NewEntry(logrus.New()), SlotDuration: 6 * time.Second}}

	pool.updateSlotDuration(&BlockInfo{Timestamp: 1000}, &BlockInfo{Timestamp: 1012})
	if pool.GetSlotDuration() != 6*time.Second {
		t.Fatalf("expected the configured fallback without slot numbers, got %v", pool.GetSlotDuration())
	}

	pool.updateSlotDuration(&BlockInfo{Timestamp: 1000, SlotNumber: 100}, &BlockInfo{Timestamp: 1012, SlotNumber: 101})
	if pool.GetSlotDuration() != 12*time.Second {
		t.Fatalf("expected 12s, got %v", pool.GetSlotDuration())
	}

	pool.updateSlotDuration(&BlockInfo{Timestamp: 2000, SlotNumber: 200}, &BlockInfo{Timestamp: 2027, SlotNumber: 203})
	if pool.GetSlotDuration() != 9*time.Second {
		t.Fatalf("expected 9s across missed slots, got %v", pool.GetSlotDuration())
	}
}

func newTrackingTestPool(ctx context.Context) *TxPool {
	return &TxPool{
		options: &TxPoolOptions{
			Context:    ctx,
			Logger:     logrus.NewEntry(logrus.New()),
			ClientPool: &ClientPool{},
		},
		wallets:      make(map[common.Address]*txPoolWalletRegistration, 4),
		walletPools:  make(map[*WalletPool]struct{}, 2),
		confirmedTxs: make(map[uint64][]*TxInfo, 4),
	}
}

type txResult struct {
	receipt *txtypes.Receipt
	err     error
}

// trackTestTx tracks a transaction without submitting it and returns its result channel.
func trackTestTx(t *testing.T, pool *TxPool, ctx context.Context, wallet *Wallet, tx *txtypes.Transaction) chan txResult {
	t.Helper()

	result := make(chan txResult, 1)

	err := pool.submitTransaction(ctx, wallet, tx, &SendTransactionOptions{
		OnComplete: func(_ *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
			result <- txResult{receipt, err}
		},
	}, false)
	if err != nil {
		t.Fatalf("submitTransaction failed: %v", err)
	}

	return result
}

// includeTestTx processes the inclusion of a transaction the way block processing does.
func includeTestTx(pool *TxPool, from common.Address, tx *txtypes.Transaction) *txtypes.Receipt {
	receipt := &txtypes.Receipt{TxHash: tx.Hash()}

	_, completions := pool.processTransactionInclusion(1, from, tx, receipt, nil)
	completions = pool.processAuthorizations(tx, completions)

	for _, completion := range completions {
		go pool.notifyTx(completion.tx, completion.receipt, nil)
	}

	return receipt
}

func expectTxResult(t *testing.T, name string, result chan txResult, receipt *txtypes.Receipt) {
	t.Helper()

	select {
	case got := <-result:
		if got.receipt != receipt || got.err != nil {
			t.Fatalf("%s: unexpected result %v / %v", name, got.receipt, got.err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("%s: did not complete", name)
	}
}

// expectTxUnresolved expects a transaction whose outcome could not be looked up.
func expectTxUnresolved(t *testing.T, name string, result chan txResult) {
	t.Helper()

	select {
	case got := <-result:
		if got.receipt != nil || !errors.Is(got.err, errReceiptLookupFailed) {
			t.Fatalf("%s: unexpected result %v / %v", name, got.receipt, got.err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("%s: did not complete", name)
	}
}

func expectTxPending(t *testing.T, name string, result chan txResult) {
	t.Helper()

	select {
	case <-result:
		t.Fatalf("%s: completed", name)
	case <-time.After(20 * time.Millisecond):
	}
}

// An included transaction confirms itself and retires everything of its sender at or
// below its nonce; later nonces stay pending.
func TestInclusionConfirmsAndReplaces(t *testing.T) {
	ctx := context.Background()
	pool := newTrackingTestPool(ctx)
	wallet := NewWallet(nil, common.HexToAddress("0xa1"))

	results := make([]chan txResult, 3)
	for nonce := range results {
		results[nonce] = trackTestTx(t, pool, ctx, wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: uint64(nonce)}))
	}

	competitor := txtypes.NewTx(&txtypes.LegacyTx{Nonce: 1, Gas: 21000})
	competitorResult := trackTestTx(t, pool, ctx, wallet, competitor)

	receipt := includeTestTx(pool, wallet.GetAddress(), competitor)

	// the lower nonce was passed unseen; no client can tell here what became of it
	expectTxUnresolved(t, "lower nonce", results[0])
	expectTxResult(t, "replaced", results[1], nil)
	expectTxResult(t, "included", competitorResult, receipt)
	expectTxPending(t, "higher nonce", results[2])

	if txs, senders := pool.tracker.count(); txs != 1 || senders != 1 {
		t.Fatalf("expected 1 tracked transaction, got %d of %d senders", txs, senders)
	}

	// a transaction this process did not send retires the rest
	includeTestTx(pool, wallet.GetAddress(), txtypes.NewTx(&txtypes.LegacyTx{Nonce: 2, Gas: 30000}))
	expectTxResult(t, "replaced from outside", results[2], nil)

	if txs, senders := pool.tracker.count(); txs != 0 || senders != 0 {
		t.Fatalf("tracker not empty: %d transactions of %d senders", txs, senders)
	}
}

// Keyed transactions compete per nonce key: an included one retires every transaction
// of its sender that shares a key at or below its sequence number.
func TestInclusionRetiresKeyedCompetitors(t *testing.T) {
	ctx := context.Background()
	pool := newTrackingTestPool(ctx)
	wallet := NewWallet(nil, common.HexToAddress("0xa1"))

	keyed := func(seq uint64, salt uint64, keys ...uint64) *txtypes.Transaction {
		nonceKeys := make([]*uint256.Int, len(keys))
		for i, key := range keys {
			nonceKeys[i] = uint256.NewInt(key)
		}

		return txtypes.NewTx((&txtypes.FrameTx{Sender: wallet.GetAddress(), ChainID: uint256.NewInt(salt)}).WithNonceKeys(nonceKeys, seq))
	}

	included := keyed(0, 1, 5, 6)
	includedResult := trackTestTx(t, pool, ctx, wallet, included)
	sameKeys := trackTestTx(t, pool, ctx, wallet, keyed(0, 2, 5, 6))
	oneSharedKey := trackTestTx(t, pool, ctx, wallet, keyed(0, 3, 6, 7))
	nextSeq := trackTestTx(t, pool, ctx, wallet, keyed(1, 4, 5))
	otherKey := trackTestTx(t, pool, ctx, wallet, keyed(0, 5, 8))
	singleUse := trackTestTx(t, pool, ctx, wallet, txtypes.NewTx((&txtypes.FrameTx{Sender: wallet.GetAddress()}).WithSingleUseNonceKeys([]*uint256.Int{uint256.NewInt(9)})))
	accountNonce := trackTestTx(t, pool, ctx, wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0}))

	receipt := includeTestTx(pool, wallet.GetAddress(), included)

	expectTxResult(t, "included", includedResult, receipt)
	expectTxResult(t, "same keys", sameKeys, nil)
	expectTxResult(t, "one shared key", oneSharedKey, nil)
	expectTxPending(t, "next sequence", nextSeq)
	expectTxPending(t, "other key", otherKey)
	expectTxPending(t, "single use", singleUse)
	expectTxPending(t, "account nonce", accountNonce)

	// a single-use transaction exhausts its keys, whatever else waits on them
	includeTestTx(pool, wallet.GetAddress(), txtypes.NewTx((&txtypes.FrameTx{Sender: wallet.GetAddress(), ChainID: uint256.NewInt(9)}).WithSingleUseNonceKeys([]*uint256.Int{uint256.NewInt(5), uint256.NewInt(9)})))

	expectTxResult(t, "next sequence", nextSeq, nil)
	expectTxResult(t, "single use", singleUse, nil)
	expectTxPending(t, "other key", otherKey)

	if txs, _ := pool.tracker.count(); txs != 2 {
		t.Fatalf("expected 2 tracked transactions, got %d", txs)
	}
}

// A cancelled or rejected transaction completes once and leaves nothing behind.
func TestCancelledAndRejectedTxsAreDropped(t *testing.T) {
	pool := newTrackingTestPool(context.Background())
	wallet := NewWallet(nil, common.HexToAddress("0xa1"))

	ctx, cancel := context.WithCancel(context.Background())
	cancelled := trackTestTx(t, pool, ctx, wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0}))
	rejectedTx := txtypes.NewTx(&txtypes.LegacyTx{Nonce: 1})
	rejected := trackTestTx(t, pool, ctx, wallet, rejectedTx)

	tracked := pool.tracker.byHash(rejectedTx.Hash())[0]
	rejectErr := context.DeadlineExceeded
	pool.completeTx(tracked, nil, rejectErr)
	pool.completeTx(tracked, nil, rejectErr)

	select {
	case got := <-rejected:
		if got.receipt != nil || got.err != rejectErr {
			t.Fatalf("unexpected result %v / %v", got.receipt, got.err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("rejected transaction did not complete")
	}

	cancel()
	expectTxResult(t, "cancelled", cancelled, nil)
	expectTxPending(t, "rejected twice", rejected)

	// the cancelled one holds its nonce until the stale pass finds its wallet unused
	pool.processStaleSenders(0)

	if txs, senders := pool.tracker.count(); txs != 0 || senders != 0 {
		t.Fatalf("tracker not empty: %d transactions of %d senders", txs, senders)
	}

	if gaps := wallet.GetNonceGaps(); gaps != nil {
		t.Fatalf("unexpected nonce gaps: %v", gaps)
	}
}

// A transaction whose nonce is already used is not left waiting on its lane.
func TestUsedNonceTxCompletes(t *testing.T) {
	ctx := context.Background()
	pool := newTrackingTestPool(ctx)
	wallet := NewWallet(nil, common.HexToAddress("0xa1"))
	wallet.SetNonce(5)

	expectTxUnresolved(t, "used nonce", trackTestTx(t, pool, ctx, wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 4})))

	if txs, _ := pool.tracker.count(); txs != 0 {
		t.Fatalf("tracker not empty: %d transactions", txs)
	}
}

// An applied EIP-7702 authorization uses up the authority's nonce: its wallet moves on
// and what it had pending on that nonce is replaced. One with another nonce is skipped.
func TestAuthorizationAdvancesAuthorityNonce(t *testing.T) {
	ctx := context.Background()
	pool := newTrackingTestPool(ctx)

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	authority := pool.RegisterWallet(NewWallet(key, crypto.PubkeyToAddress(key.PublicKey)), ctx)
	authority.SetNonce(3)

	pending := trackTestTx(t, pool, ctx, authority, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 3}))
	next := trackTestTx(t, pool, ctx, authority, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 4}))

	setCode := func(nonce uint64) *txtypes.Transaction {
		auth, err := txtypes.SignAuthorization(txtypes.SetCodeAuthorization{Nonce: nonce}, key)
		if err != nil {
			t.Fatal(err)
		}

		return txtypes.NewTx(&txtypes.SetCodeTx{ChainID: uint256.NewInt(1), AuthList: []txtypes.SetCodeAuthorization{auth}})
	}

	includeTestTx(pool, common.HexToAddress("0xb1"), setCode(7))
	expectTxPending(t, "skipped authorization", pending)

	includeTestTx(pool, common.HexToAddress("0xb1"), setCode(3))
	expectTxResult(t, "authority nonce used", pending, nil)
	expectTxPending(t, "next nonce", next)

	if authority.GetConfirmedNonce() != 4 || authority.GetNonce() != 4 {
		t.Fatalf("unexpected authority nonces: confirmed %d, pending %d", authority.GetConfirmedNonce(), authority.GetNonce())
	}
}

// A transaction counts for the pool whose scope it is sent in, whichever wallet sends
// it, and without a scope for the pool that prepared the wallet.
func TestWalletPoolAttribution(t *testing.T) {
	baseCtx := context.Background()
	pool := newTrackingTestPool(baseCtx)
	root := &RootWallet{wallet: NewWallet(nil, common.HexToAddress("0x01"))}

	firstCtx := WithWalletPoolScope(baseCtx)
	first := NewWalletPool(firstCtx, logrus.New(), root, nil, pool)
	secondCtx, stopSecond := context.WithCancel(WithWalletPoolScope(baseCtx))
	second := NewWalletPool(secondCtx, logrus.New(), root, nil, pool)

	results := make(chan error, 4)
	second.SetTransactionTracker(func(err error) { results <- err })
	first.SetTransactionTracker(func(error) { t.Error("result reported to the wrong pool") })

	wallet := NewWallet(nil, common.HexToAddress("0xa1"))
	wallet.ownerPool.Store(first)

	scoped := txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0})
	derivedCtx, cancelDerived := context.WithTimeout(secondCtx, time.Minute)
	defer cancelDerived()

	scopedResult := trackTestTx(t, pool, derivedCtx, wallet, scoped)
	rootResult := trackTestTx(t, pool, secondCtx, root.wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0, Gas: 21000}))
	unscoped := txtypes.NewTx(&txtypes.LegacyTx{Nonce: 1})
	unscopedResult := trackTestTx(t, pool, baseCtx, wallet, unscoped)

	if first.GetPendingTxCount() != 1 || second.GetPendingTxCount() != 2 {
		t.Fatalf("unexpected pending counts: %d / %d", first.GetPendingTxCount(), second.GetPendingTxCount())
	}

	if tracked := pool.tracker.byHash(scoped.Hash())[0]; tracked.walletPool != second {
		t.Fatal("scoped transaction not attributed to the sending pool")
	}

	if tracked := pool.tracker.byHash(unscoped.Hash())[0]; tracked.walletPool != first {
		t.Fatal("unscoped transaction not attributed to the wallet's pool")
	}

	// block statistics follow the attribution
	block := &txtypes.Block{}
	confirmed := []*TxInfo{
		{TxHash: scoped.Hash(), From: wallet.GetAddress(), receipt: &txtypes.Receipt{GasUsed: 21000}, walletPool: second},
		{TxHash: unscoped.Hash(), From: wallet.GetAddress(), receipt: &txtypes.Receipt{GasUsed: 30000}, walletPool: first},
		{TxHash: common.HexToHash("0x01"), From: common.HexToAddress("0xff")},
	}

	all := pool.calculateAllWalletPoolStats(confirmed, block, nil)

	for walletPool, gas := range map[*WalletPool]uint64{first: 30000, second: 21000} {
		for _, stats := range []*WalletPoolBlockStats{all.WalletPoolStats[walletPool], pool.calculateWalletPoolStats(walletPool, confirmed, block)} {
			if stats.ConfirmedTxCount != 1 || stats.AffectedWallets != 1 || len(stats.Receipts) != 1 || stats.Receipts[0].GasUsed != gas {
				t.Fatalf("unexpected stats: %+v", stats)
			}
		}
	}

	receipt := includeTestTx(pool, wallet.GetAddress(), scoped)
	expectTxResult(t, "scoped", scopedResult, receipt)

	select {
	case <-results:
	case <-time.After(2 * time.Second):
		t.Fatal("no result reported to the sending pool")
	}

	// a stopped pool's transactions end without a result being reported
	stopSecond()
	expectTxResult(t, "root", rootResult, nil)
	expectTxPending(t, "unscoped", unscopedResult)

	if first.GetPendingTxCount() != 1 || second.GetPendingTxCount() != 0 {
		t.Fatalf("unexpected pending counts: %d / %d", first.GetPendingTxCount(), second.GetPendingTxCount())
	}

	select {
	case <-results:
		t.Fatal("result reported for a stopped pool")
	case <-time.After(20 * time.Millisecond):
	}
}

// BenchmarkTrackAndInclude tracks 10k transactions, confirms them and waits for their
// callbacks, spread over a few busy senders and over many.
func BenchmarkTrackAndInclude(b *testing.B) {
	const total = 10000

	for _, senders := range []int{10, 1000} {
		b.Run(fmt.Sprintf("senders=%d", senders), func(b *testing.B) {
			ctx := context.Background()
			pool := newTrackingTestPool(ctx)

			var wg sync.WaitGroup

			options := &SendTransactionOptions{
				OnComplete: func(*txtypes.Transaction, *txtypes.Receipt, error) { wg.Done() },
			}

			wallets := make([]*Wallet, senders)
			for i := range wallets {
				wallets[i] = NewWallet(nil, common.BigToAddress(big.NewInt(int64(i+1))))
			}

			b.ReportAllocs()
			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				b.StopTimer()

				txs := make([]*txtypes.Transaction, total)
				for i := range txs {
					txs[i] = txtypes.NewTx(&txtypes.LegacyTx{Nonce: uint64(n*total + i/senders), Gas: uint64(i%senders + 1)})
					txs[i].Hash()
				}

				wg.Add(total)
				b.StartTimer()

				for i, tx := range txs {
					if err := pool.submitTransaction(ctx, wallets[i%senders], tx, options, false); err != nil {
						b.Fatal(err)
					}
				}

				var completions []txCompletion
				for i, tx := range txs {
					_, completions = pool.processTransactionInclusion(1, wallets[i%senders].GetAddress(), tx, &txtypes.Receipt{}, completions)
				}

				for _, completion := range completions {
					go pool.notifyTx(completion.tx, completion.receipt, nil)
				}

				wg.Wait()
			}

			b.StopTimer()

			if txs, _ := pool.tracker.count(); txs != 0 {
				b.Fatalf("tracker not empty: %d", txs)
			}
		})
	}
}

// The work queue runs everything it is given, in order per worker, on at most the
// given number of workers, and keeps none of them around afterwards.
func TestWorkQueue(t *testing.T) {
	const limit, total = 4, 2000

	var (
		queue   workQueue
		wg      sync.WaitGroup
		mu      sync.Mutex
		running int
		peak    int
		done    int
	)

	wg.Add(total)

	for i := 0; i < total; i++ {
		queue.push(limit, func() {
			defer wg.Done()

			mu.Lock()
			running++
			peak = max(peak, running)
			mu.Unlock()

			time.Sleep(10 * time.Microsecond)

			mu.Lock()
			running--
			done++
			mu.Unlock()
		})
	}

	wg.Wait()

	if done != total || peak > limit {
		t.Fatalf("ran %d of %d with up to %d at once (limit %d)", done, total, peak, limit)
	}

	deadline := time.Now().Add(2 * time.Second)
	for {
		queue.mu.Lock()
		workers, items := queue.workers, len(queue.items)
		queue.mu.Unlock()

		if workers == 0 && items == 0 {
			break
		}

		if time.Now().After(deadline) {
			t.Fatalf("queue not idle: %d workers, %d items", workers, items)
		}

		time.Sleep(time.Millisecond)
	}
}

// A transaction at a lower nonce than the included one had its turn in a block that was
// not processed. It is not written off as replaced with the block: its receipt decides.
func TestPassedNonceIsResolvedByReceipt(t *testing.T) {
	ctx := context.Background()
	pool := newTrackingTestPool(ctx)
	wallet := NewWallet(nil, common.HexToAddress("0xa1"))

	passed := trackTestTx(t, pool, ctx, wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0}))
	replaced := trackTestTx(t, pool, ctx, wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 1}))
	included := txtypes.NewTx(&txtypes.LegacyTx{Nonce: 1, Gas: 21000})
	includedResult := trackTestTx(t, pool, ctx, wallet, included)

	receipt := &txtypes.Receipt{TxHash: included.Hash()}
	_, completions := pool.processTransactionInclusion(1, wallet.GetAddress(), included, receipt, nil)

	if len(completions) != 2 {
		t.Fatalf("expected the included and the replaced transaction to complete with the block, got %d", len(completions))
	}

	for _, completion := range completions {
		if completion.tx.seq != 1 {
			t.Fatalf("transaction at nonce %d completed with the block", completion.tx.seq)
		}

		go pool.notifyTx(completion.tx, completion.receipt, nil)
	}

	expectTxResult(t, "included", includedResult, receipt)
	expectTxResult(t, "replaced", replaced, nil)

	// there is no client to ask here, so it ends with its outcome unknown
	expectTxUnresolved(t, "passed", passed)

	if txs, senders := pool.tracker.count(); txs != 0 || senders != 0 {
		t.Fatalf("tracker not empty: %d transactions of %d senders", txs, senders)
	}
}

// A cancelled transaction may still be in a mempool, so its nonce stays taken on the
// lane until it is passed, for as long as its wallet is registered.
func TestCancelledTxKeepsItsNonce(t *testing.T) {
	pool := newTrackingTestPool(context.Background())

	walletCtx, unregister := context.WithCancel(context.Background())
	wallet := pool.RegisterWallet(NewWallet(nil, common.HexToAddress("0xa1")), walletCtx)

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancelledTx := txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0})
	cancelled := trackTestTx(t, pool, cancelCtx, wallet, cancelledTx)
	keyedTx := txtypes.NewTx((&txtypes.FrameTx{Sender: wallet.GetAddress()}).WithNonceKeys([]*uint256.Int{uint256.NewInt(5)}, 0))
	cancelledKeyed := trackTestTx(t, pool, cancelCtx, wallet, keyedTx)
	onlyTx := txtypes.NewTx(&txtypes.LegacyTx{Nonce: 1})
	only := trackTestTx(t, pool, cancelCtx, wallet, onlyTx)

	cancel()
	expectTxResult(t, "cancelled", cancelled, nil)
	expectTxResult(t, "cancelled keyed", cancelledKeyed, nil)
	expectTxResult(t, "cancelled last", only, nil)

	// nothing waits for an outcome anymore, but both nonces are still taken
	if wallet.GetPendingTx(cancelledTx) == nil || wallet.GetPendingTx(onlyTx) == nil {
		t.Fatal("cancelled transaction is not pending anymore")
	}

	if wallet.GetPendingTx(keyedTx) != nil {
		t.Fatal("cancelled keyed transaction is still on its lane")
	}

	if lowest, ok := wallet.GetLowestPendingNonce(); !ok || lowest != 0 {
		t.Fatalf("unexpected lowest pending nonce: %d %v", lowest, ok)
	}

	if txs, senders := pool.tracker.count(); txs != 0 || senders != 1 {
		t.Fatalf("expected a sender holding only placeholders, got %d transactions of %d senders", txs, senders)
	}

	pool.processStaleSenders(0)

	if wallet.GetPendingTx(onlyTx) == nil {
		t.Fatal("placeholder of a registered wallet was dropped")
	}

	// its nonce being included releases a placeholder, and nothing completes twice
	includeTestTx(pool, wallet.GetAddress(), cancelledTx)
	expectTxPending(t, "cancelled", cancelled)

	if wallet.GetPendingTx(cancelledTx) != nil || wallet.GetPendingTx(onlyTx) == nil {
		t.Fatal("wrong placeholder released")
	}

	// the rest goes once the wallet is not in use anymore
	unregister()
	pool.processStaleSenders(0)

	if wallet.GetPendingTx(onlyTx) != nil {
		t.Fatal("placeholder outlived its wallet")
	}

	if txs, senders := pool.tracker.count(); txs != 0 || senders != 0 {
		t.Fatalf("tracker not empty: %d transactions of %d senders", txs, senders)
	}
}

// A transaction that was retired from its lane with its outcome still open is not on
// the lane anymore. Cancelling it then must not leave a placeholder behind that nothing
// can release.
func TestCancelOfRetiredTxLeavesNoPlaceholder(t *testing.T) {
	pool := newTrackingTestPool(context.Background())
	wallet := pool.RegisterWallet(NewWallet(nil, common.HexToAddress("0xa1")), context.Background())

	cancelCtx, cancel := context.WithCancel(context.Background())
	passed := trackTestTx(t, pool, cancelCtx, wallet, txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0}))

	// what an inclusion of a higher nonce does to it, up to the receipt lookup
	sender := pool.tracker.getSender(wallet.GetAddress())
	sender.mu.Lock()
	retired := sender.retire(nil, 1)
	sender.mu.Unlock()

	if len(retired) != 1 {
		t.Fatalf("expected the transaction to be retired, got %d", len(retired))
	}

	cancel()
	expectTxResult(t, "cancelled", passed, nil)

	sender.mu.Lock()
	placeholders := sender.placeholders
	sender.mu.Unlock()

	if placeholders != 0 {
		t.Fatalf("%d placeholders counted for a transaction that is not on its lane", placeholders)
	}

	if txs, senders := pool.tracker.count(); txs != 0 || senders != 0 {
		t.Fatalf("tracker not empty: %d transactions of %d senders", txs, senders)
	}
}

// Waiting for a transaction does not make it a transaction of the waiting pool.
func TestAwaitTransactionIsNotAttributed(t *testing.T) {
	baseCtx := context.Background()
	pool := newTrackingTestPool(baseCtx)

	scopeCtx, cancel := context.WithCancel(WithWalletPoolScope(baseCtx))
	walletPool := NewWalletPool(scopeCtx, logrus.New(), &RootWallet{wallet: NewWallet(nil, common.HexToAddress("0x01"))}, nil, pool)
	walletPool.SetTransactionTracker(func(error) { t.Error("result reported for a transaction that was only awaited") })

	wallet := NewWallet(nil, common.HexToAddress("0xa1"))
	tx := txtypes.NewTx(&txtypes.LegacyTx{Nonce: 0})

	done := make(chan error, 1)

	go func() {
		_, err := pool.AwaitTransaction(scopeCtx, wallet, tx)
		done <- err
	}()

	deadline := time.Now().Add(2 * time.Second)
	for len(pool.tracker.byHash(tx.Hash())) == 0 {
		if time.Now().After(deadline) {
			t.Fatal("awaited transaction is not tracked")
		}

		time.Sleep(time.Millisecond)
	}

	if walletPool.GetPendingTxCount() != 0 {
		t.Fatal("awaited transaction counts as pending for the pool")
	}

	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("unexpected error %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("AwaitTransaction did not return")
	}
}
