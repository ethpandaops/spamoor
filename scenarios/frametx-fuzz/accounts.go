package frametxfuzz

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/txtypes"
)

// Account contracts: generated code that plays the sender and paymaster roles.
//
// They are deployed by ordinary body frames and kept in a small ring, so the code behind
// those roles turns over as the run goes rather than being a fixed set installed once.
// A contract is deployed with value, since the factory forwards CALLVALUE to CREATE2 and
// a paymaster must hold at least the transaction's maximum cost.
//
// Each contract is used once and then queued for a wipe: a sender's nonce is a one-shot
// (see senderNonce), and using a paymaster once keeps the two roles symmetric and lets its
// funding be reclaimed the same way. There is no Wallet behind a contract sender at all --
// it holds no key, its nonce is known, and its one transaction is submitted raw and
// tracked by hash -- so nothing has to be registered with the transaction pool.

// accountBufferSize is how many recently deployed account contracts stay reachable for a
// role before the oldest is dropped unused.
const accountBufferSize = 8

// accountFundingGas is how much gas the value deployed with an account contract must
// cover at the current fee cap, since a paymaster is charged the transaction's maximum
// cost up front.
const accountFundingGas = 2_000_000

// senderNonce is the account nonce a contract sender's single transaction uses. A
// contract created by CREATE2 starts at nonce one under EIP-161, and nothing moves it
// before its one use.
const senderNonce = 1

// accountContract identifies a generated account contract for a caller that only needs
// its address and whether it can approve.
type accountContract struct {
	Address common.Address

	// Approves records whether the code ends in the APPROVE that lets a transaction
	// using it land. Some are generated without one so the role can fail outright.
	Approves bool
}

// accountBuffer is a bounded ring of recently deployed account contracts, plus a queue of
// used or dropped ones awaiting a wipe of their leftover balance.
type accountBuffer struct {
	mutex sync.Mutex

	// ready are deployed contracts not yet used for a role.
	ready []accountContract

	// wipeQueue holds contracts to reclaim, oldest first: ones used for a role, and
	// ones dropped from ready unused. A wipe is deferred behind wipeDelay more entries
	// so a used contract's transaction has settled -- including a paymaster's fee refund
	// -- before its balance is taken.
	wipeQueue []common.Address
}

// wipeDelay is how many further contracts must be queued before the oldest is wiped, so a
// used contract's transaction has confirmed and any fee refund has settled.
const wipeDelay = 4

// newAccountBuffer returns an empty buffer.
func newAccountBuffer() *accountBuffer {
	return &accountBuffer{}
}

// remember adds contracts a landed transaction deployed, dropping the oldest unused.
func (b *accountBuffer) remember(contracts []accountContract) {
	if len(contracts) == 0 {
		return
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.ready = append(b.ready, contracts...)

	// Contracts pushed out of the ring were never used, but they still hold their
	// funding, so they are queued for a wipe rather than dropped.
	for len(b.ready) > accountBufferSize {
		b.wipeQueue = append(b.wipeQueue, b.ready[0].Address)
		b.ready = b.ready[1:]
	}
}

// available reports whether any account contract is ready for a role.
func (b *accountBuffer) available() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return len(b.ready) > 0
}

// readyCount returns how many account contracts are ready for a role, which the account
// deployer uses to decide how hard to top the pool up.
func (b *accountBuffer) readyCount() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return len(b.ready)
}

// take removes and returns a ready account contract, queuing it to be wiped later.
//
// Both roles consume a contract: a sender because its nonce is a one-shot, a paymaster to
// keep the roles symmetric and so its funding is reclaimed the same way. The chosen index
// rotates with the caller's so concurrent transactions do not all take the same one.
func (b *accountBuffer) take(index int) (accountContract, bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.ready) == 0 {
		return accountContract{}, false
	}

	pos := index % len(b.ready)
	contract := b.ready[pos]

	b.ready = append(b.ready[:pos], b.ready[pos+1:]...)
	b.wipeQueue = append(b.wipeQueue, contract.Address)

	return contract, true
}

// takeWipe returns a contract to reclaim, once enough others are queued behind it that a
// used contract's transaction has settled. It returns false when nothing is due.
func (b *accountBuffer) takeWipe() (common.Address, bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.wipeQueue) <= wipeDelay {
		return common.Address{}, false
	}

	address := b.wipeQueue[0]
	b.wipeQueue = b.wipeQueue[1:]

	return address, true
}

// prepareKeylessFrameTx assigns the fields the canonical signature hash covers for a
// contract sender, which has no key and a known nonce.
//
// A contract sender approves in its own code, so it contributes no signature entry and
// nothing signs the transaction. The fields still have to be final before the hash is
// taken, so this runs where PrepareFrameTx would for a keyed sender.
func prepareKeylessFrameTx(frameTx *txtypes.FrameTx, sender common.Address, chainID *uint256.Int) {
	frameTx.ChainID = chainID
	frameTx.Sender = sender

	if frameTx.HasKeyedNonces() && len(frameTx.NonceKeys) == 0 {
		frameTx.NonceKeys = []*uint256.Int{new(uint256.Int)}
	}

	if frameTx.UsesLegacyNonce() {
		frameTx.NonceSeq = senderNonce
	}
}
