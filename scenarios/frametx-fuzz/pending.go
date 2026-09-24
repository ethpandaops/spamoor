package frametxfuzz

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Pending account contracts: generated account code that has *not* been deployed yet.
//
// This is what makes EIP-8141's deploy-led validation prefixes reachable. A CREATE2
// address is a function of the code, so the account's address is known before it exists,
// and a deploy frame leading the validation prefix installs the code at `tx.sender` in
// the very transaction that uses it. The account is therefore created, validated and used
// in one transaction, which is the flagship native-account-abstraction flow and the only
// shape that puts a contract creation inside the part of a transaction the public mempool
// simulates.
//
// Two ways to pay for such a transaction, and both are drawn:
//
//   - `[deploy, self_verify]` — the account pays for itself, so an earlier transaction has
//     to credit its address first. That is what this buffer holds: addresses a landed
//     transaction funded, waiting for the transaction that deploys them.
//   - `[deploy, only_verify, pay]` — a paymaster pays, so the account needs no balance at
//     all and is minted on the spot rather than drawn from here.
//
// Funding an address that holds no code is safe: CREATE2 refuses an address with code or
// a non-zero nonce, but balance alone does not block it.

// pendingReadyTarget is how many funded, undeployed accounts the buffer is kept topped up
// to.
const pendingReadyTarget = 4

// pendingBufferSize bounds the buffer, so a run that never draws a deploy prefix does not
// keep funding addresses it will not use.
const pendingBufferSize = 8

// pendingFundingsPerTx caps how many addresses one transaction credits.
const pendingFundingsPerTx = 2

// pendingFundingIndexBase keeps funding frames in a frame-index range of their own, so
// the salts they derive never collide with a body deployment's or an account deployment's.
const pendingFundingIndexBase = 200

// deployedSenderNonce is the account nonce a deploy-led transaction uses. The account does
// not exist when the nonce is checked -- the deploy frame in the same transaction creates
// it -- so the nonce is zero, unlike a contract deployed by an earlier transaction
// (see senderNonce).
const deployedSenderNonce = 0

// pendingAccount is an account contract whose code is known but not yet deployed.
type pendingAccount struct {
	// Address is the CREATE2 address the deploy frame installs the code at, which is
	// also the transaction's sender.
	Address common.Address

	// Factory is the CREATE2 factory the address was derived against, and the deploy
	// frame's target.
	Factory common.Address

	// Salt and InitCode are what the deploy frame passes to the factory.
	Salt     [32]byte
	InitCode []byte

	// Approves records whether the code ends in the APPROVE that lets a transaction
	// using it land. Some are generated without one so the role can fail outright.
	Approves bool

	// Funded records whether an earlier transaction credited the address, which decides
	// whether the account can pay for its own transaction or needs a paymaster.
	Funded bool
}

// deployData renders the factory calldata that deploys the account: the salt followed by
// the init code, which is what the canonical CREATE2 factory expects.
func (a pendingAccount) deployData() []byte {
	data := make([]byte, 0, len(a.Salt)+len(a.InitCode))
	data = append(data, a.Salt[:]...)
	data = append(data, a.InitCode...)

	return data
}

// newPendingAccount derives an account contract from a seed and a frame index, without
// touching the chain: the address follows from the code and the salt.
func newPendingAccount(seed string, factory common.Address, txID uint64, index int, gasLimit uint64, approves bool) (pendingAccount, error) {
	runtime := GenerateAccountCode(seed, txID, index, gasLimit, approves)

	initCode, err := DeploymentInitCode(runtime)
	if err != nil {
		return pendingAccount{}, err
	}

	salt := deploySalt(txID, index)

	return pendingAccount{
		Address:  crypto.CreateAddress2(factory, salt, crypto.Keccak256(initCode)),
		Factory:  factory,
		Salt:     salt,
		InitCode: initCode,
		Approves: approves,
	}, nil
}

// pendingBuffer is a bounded ring of funded, undeployed account contracts.
type pendingBuffer struct {
	mutex sync.Mutex
	ready []pendingAccount
}

// newPendingBuffer returns an empty buffer.
func newPendingBuffer() *pendingBuffer {
	return &pendingBuffer{}
}

// remember adds accounts a landed transaction funded, dropping the oldest when the buffer
// is full. A dropped address keeps its funding: it is an ordinary address with a balance
// and no code, so a later run with the same seed finds it again.
func (b *pendingBuffer) remember(accounts []pendingAccount) {
	if len(accounts) == 0 {
		return
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.ready = append(b.ready, accounts...)

	if len(b.ready) > pendingBufferSize {
		b.ready = b.ready[len(b.ready)-pendingBufferSize:]
	}
}

// available reports whether a funded account is ready to be deployed.
func (b *pendingBuffer) available() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return len(b.ready) > 0
}

// readyCount returns how many funded accounts are waiting, which the funder uses to decide
// how hard to top the buffer up.
func (b *pendingBuffer) readyCount() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return len(b.ready)
}

// take removes and returns a funded account. The chosen index rotates with the caller's so
// concurrent transactions do not all take the same one.
func (b *pendingBuffer) take(index int) (pendingAccount, bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.ready) == 0 {
		return pendingAccount{}, false
	}

	pos := index % len(b.ready)
	account := b.ready[pos]

	b.ready = append(b.ready[:pos], b.ready[pos+1:]...)

	return account, true
}
