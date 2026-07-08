package ensnames

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethpandaops/spamoor/scenarios/ens-names/contract"
)

// Well-known ENS registry nodes used throughout the scenario.
var (
	rootNode        = [32]byte{}
	ethNode         = namehash("eth")
	reverseNode     = namehash("reverse")
	addrReverseNode = namehash("addr.reverse")
)

// namehash implements the ENS namehash algorithm (EIP-137): the empty name is
// the zero hash, and every label is folded in right-to-left as
// keccak256(namehash(rest) . keccak256(label)).
func namehash(name string) [32]byte {
	node := [32]byte{}
	if name == "" {
		return node
	}

	labels := strings.Split(name, ".")
	for i := len(labels) - 1; i >= 0; i-- {
		labelHash := crypto.Keccak256([]byte(labels[i]))
		copy(node[:], crypto.Keccak256(node[:], labelHash))
	}

	return node
}

// labelhash returns keccak256(label), the registrar token id of a .eth label.
func labelhash(label string) [32]byte {
	hash := [32]byte{}
	copy(hash[:], crypto.Keccak256([]byte(label)))
	return hash
}

// ethNameNode returns namehash(label + ".eth") without string round-trips.
func ethNameNode(label string) [32]byte {
	node := [32]byte{}
	lh := labelhash(label)
	copy(node[:], crypto.Keccak256(ethNode[:], lh[:]))
	return node
}

// nameState tracks one .eth name owned (or being acquired) by a child wallet,
// including the commit-reveal pipeline state. Name state is fresh per run: the
// scenario never reconstructs names from previous runs, it only registers new
// ones (labels embed a per-run random id, so they cannot collide with earlier
// runs).
type nameState struct {
	label   string
	node    [32]byte
	tokenID [32]byte

	// commit-reveal pipeline
	registration contract.IETHRegistrarControllerRegistration
	committed    bool      // commit tx confirmed
	commitTime   time.Time // local time the commit confirmation was observed
	registerTry  uint64    // failed register attempts

	// post-registration state
	registered  bool
	wrapped     bool // held by the NameWrapper as ERC1155
	needReclaim bool // ERC721 received via transfer, registry owner not yet updated

	inFlight bool // a tx for this name is currently pending
}

// walletState tracks the ENS state of one child wallet.
type walletState struct {
	mu      sync.Mutex
	names   []*nameState
	pending *nameState // name in the commit-reveal pipeline

	wrapperApproved bool // base registrar setApprovalForAll(NameWrapper) confirmed
	approveInFlight bool
	churnSeq        uint64
	nameSeq         uint64
}

// walletState returns (creating on first use) the state of the given wallet.
func (s *Scenario) getWalletState(addr common.Address) *walletState {
	s.walletStatesMtx.Lock()
	defer s.walletStatesMtx.Unlock()

	state := s.walletStates[addr]
	if state == nil {
		state = &walletState{}
		s.walletStates[addr] = state
	}

	return state
}

// newLabel returns a fresh, collision-safe .eth label for the given wallet.
// Labels embed the per-run random id, the wallet index and a sequence number,
// so they are unique across wallets, runs and root keys.
func (s *Scenario) newLabel(state *walletState, walletIdx uint64) string {
	state.nameSeq++
	return fmt.Sprintf("sp%s-%d-%d", s.runID, walletIdx, state.nameSeq)
}

// newChurnLabel returns a fresh label for a short-lived churn registration.
func (s *Scenario) newChurnLabel(state *walletState, walletIdx uint64) string {
	state.churnSeq++
	return fmt.Sprintf("c%s-%d-%d", s.runID, walletIdx, state.churnSeq)
}
