package codeadopt

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/txtypes"
)

type caseFn func(ctx context.Context, run *caseRun) error

var caseTable = []struct {
	name string
	fn   func(s *Scenario) caseFn
}{
	{"migrate", func(s *Scenario) caseFn { return s.caseMigrate }},
	{"adopt-runtime", func(s *Scenario) caseFn { return s.caseAdoptRuntime }},
	{"adopt-initcode", func(s *Scenario) caseFn { return s.caseAdoptInitCode }},
	{"adopt-same-tx-source", func(s *Scenario) caseFn { return s.caseAdoptSameTxSource }},
	{"sendall-runtime", func(s *Scenario) caseFn { return s.caseSendAllRuntime }},
	{"sendall-initcode", func(s *Scenario) caseFn { return s.caseSendAllInitCode }},
}

var ecrecoverAddress = common.BytesToAddress([]byte{0x01})

var errNoClient = errors.New("no client available")

func allCaseNames() []string {
	names := make([]string, 0, len(caseTable))
	for _, c := range caseTable {
		names = append(names, c.name)
	}

	return names
}

func selectCases(s *Scenario, spec string) ([]caseFn, []string, error) {
	want := map[string]bool{}
	if spec != "" && spec != "all" {
		for _, name := range strings.Split(spec, ",") {
			want[strings.TrimSpace(name)] = true
		}
	}

	var (
		fns   []caseFn
		names []string
	)

	filter := len(want) > 0

	for _, c := range caseTable {
		if filter && !want[c.name] {
			continue
		}

		delete(want, c.name)
		fns = append(fns, c.fn(s))
		names = append(names, c.name)
	}

	for name := range want {
		return nil, nil, fmt.Errorf("unknown case %q (known: %s)", name, strings.Join(allCaseNames(), ", "))
	}

	return fns, names, nil
}

type caseRun struct {
	s      *Scenario
	name   string
	idx    uint64
	logger *logrus.Entry
	client *spamoor.Client
	wallet *spamoor.Wallet
	passed int
	failed int
}

func (s *Scenario) newCaseRun(name string, idx uint64) *caseRun {
	client := s.client()
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(idx))

	logger := s.logger.WithField("case", name)
	if client != nil {
		logger = logger.WithField("rpc", client.GetName())
	}

	if wallet != nil {
		logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
	}

	return &caseRun{s: s, name: name, idx: idx, logger: logger, client: client, wallet: wallet}
}

func (r *caseRun) expect(check string, ok bool, format string, args ...any) {
	if ok {
		r.passed++
		r.s.checksPassed.Add(1)

		return
	}

	r.failed++
	r.s.checksFailed.Add(1)
	r.logger.Warnf("check failed: %s: %s", check, fmt.Sprintf(format, args...))
}

func (r *caseRun) ready() error {
	if r.client == nil || r.wallet == nil {
		return errNoClient
	}

	return nil
}

func (r *caseRun) send(ctx context.Context, to *common.Address, value *uint256.Int, data []byte, auths []txtypes.SetCodeAuthorization) (*txtypes.Receipt, error) {
	return r.s.send(ctx, r.client, r.wallet, to, value, data, auths)
}

func (s *Scenario) send(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, to *common.Address, value *uint256.Int, data []byte, auths []txtypes.SetCodeAuthorization) (*txtypes.Receipt, error) {
	if value == nil {
		value = new(uint256.Int)
	}

	gas, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:              wallet.GetAddress(),
		To:                to,
		Value:             value.ToBig(),
		Data:              data,
		AuthorizationList: auths,
	})
	if err != nil {
		return nil, fmt.Errorf("estimating gas: %w", err)
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)

	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("getting fees: %w", err)
	}

	metadata := &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gas*3/2 + 10_000,
		To:        to,
		Value:     value,
		Data:      data,
		AuthList:  auths,
	}

	var tx *txtypes.Transaction

	if len(auths) > 0 {
		txData, err := txbuilder.SetCodeTx(metadata)
		if err != nil {
			return nil, err
		}

		tx, err = wallet.BuildSetCodeTx(txData)
		if err != nil {
			return nil, err
		}
	} else {
		txData, err := txbuilder.DynFeeTx(metadata)
		if err != nil {
			return nil, err
		}

		tx, err = wallet.BuildDynamicFeeTx(txData)
		if err != nil {
			return nil, err
		}
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: true,
	})
	if err != nil {
		return nil, err
	}

	if receipt == nil {
		return nil, fmt.Errorf("transaction %s was not confirmed", tx.Hash())
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return receipt, fmt.Errorf("transaction %s failed in block %v", tx.Hash(), receipt.BlockNumber)
	}

	return receipt, s.awaitBlock(ctx, client, receipt.BlockNumber.Uint64())
}

func (s *Scenario) awaitBlock(ctx context.Context, client *spamoor.Client, number uint64) error {
	for {
		head, err := client.GetBlockNumber(ctx)
		if err == nil && head >= number {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}

func (s *Scenario) call(ctx context.Context, client *spamoor.Client, to common.Address, data []byte) ([]byte, error) {
	return client.GetEthClient().CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
}

func word(out []byte) uint64 {
	if len(out) != 32 {
		return ^uint64(0)
	}

	return new(big.Int).SetBytes(out).Uint64()
}

func (s *Scenario) checkSources(ctx context.Context) error {
	run := s.newCaseRun("sources", 0)
	if run.client == nil {
		return scenario.ErrNoClients
	}

	if err := s.awaitBlock(ctx, run.client, s.fixtureBlock); err != nil {
		return err
	}

	empty := randomAddress()
	cases := []struct {
		check  string
		source common.Address
		want   uint64
	}{
		{"code source", s.valueSource, 1},
		{"same code hash", s.adopter, 1},
		{"empty account", empty, 0},
		{"EOA", s.walletPool.GetWellKnownWallet("codeadopt-deployer").GetAddress(), 0},
		{"precompile", ecrecoverAddress, 0},
	}

	for _, c := range cases {
		out, err := s.call(ctx, run.client, s.adopter, common.LeftPadBytes(c.source.Bytes(), 32))
		run.expect("SETCODEFROM "+c.check, err == nil && word(out) == c.want, "source %s: got %x err %v, want %d", c.source, out, err, c.want)
	}

	out, err := s.call(ctx, run.client, s.staticCaller, common.LeftPadBytes(s.valueSource.Bytes(), 32))
	run.expect("SETCODEFROM halts under STATICCALL", err == nil && len(out) == 64 && word(out[:32]) == 0, "got %x err %v", out, err)

	s.logger.Infof("SETCODEFROM source checks: %d passed, %d failed", run.passed, run.failed)

	return nil
}

func randomAddress() common.Address {
	var address common.Address

	_, _ = rand.Read(address[:])

	return address
}

func (s *Scenario) codeHashesMatch(ctx context.Context, client *spamoor.Client, a, b common.Address) (bool, error) {
	codeA, err := client.GetCodeAt(ctx, a)
	if err != nil {
		return false, err
	}

	codeB, err := client.GetCodeAt(ctx, b)
	if err != nil {
		return false, err
	}

	return len(codeA) > 0 && bytes.Equal(codeA, codeB), nil
}

func (s *Scenario) caseAdoptRuntime(ctx context.Context, run *caseRun) error {
	if err := run.ready(); err != nil {
		return err
	}

	receipt, err := run.send(ctx, nil, nil, deployCode(adopterCode), nil)
	if err != nil {
		return err
	}

	adopter := receipt.ContractAddress

	if _, err := run.send(ctx, &adopter, nil, common.LeftPadBytes(s.valueSource.Bytes(), 32), nil); err != nil {
		return err
	}

	same, err := s.codeHashesMatch(ctx, run.client, adopter, s.valueSource)
	run.expect("adopter holds the source code", err == nil && same, "%s vs %s err %v", adopter, s.valueSource, err)

	out, err := s.call(ctx, run.client, adopter, nil)
	run.expect("adopter runs the adopted code", err == nil && word(out) == 0x42, "got %x err %v", out, err)

	return nil
}

func (s *Scenario) caseAdoptInitCode(ctx context.Context, run *caseRun) error {
	if err := run.ready(); err != nil {
		return err
	}

	receipt, err := run.send(ctx, nil, nil, adoptInInitCode(s.valueSource), nil)
	if err != nil {
		return err
	}

	same, err := s.codeHashesMatch(ctx, run.client, receipt.ContractAddress, s.valueSource)
	run.expect("creation installs the adopted code, ignoring the 0xEF return data", err == nil && same, "%s err %v", receipt.ContractAddress, err)

	return nil
}

func (s *Scenario) caseAdoptSameTxSource(ctx context.Context, run *caseRun) error {
	if err := run.ready(); err != nil {
		return err
	}

	receipt, err := run.send(ctx, nil, nil, adoptSameTxSourceInitCode(s.valueSource), nil)
	if err != nil {
		return err
	}

	created := receipt.ContractAddress

	slot0, err := run.client.GetStorageAt(ctx, created, common.Hash{})
	run.expect("source created in the same transaction is refused", err == nil && slot0 == common.Hash{}, "slot0 %s err %v", slot0, err)

	slot1, err := run.client.GetStorageAt(ctx, created, common.BigToHash(big.NewInt(1)))
	run.expect("pre-existing source is accepted", err == nil && slot1 == common.BigToHash(big.NewInt(1)), "slot1 %s err %v", slot1, err)

	same, err := s.codeHashesMatch(ctx, run.client, created, s.valueSource)
	run.expect("creation keeps the adopted code", err == nil && same, "%s err %v", created, err)

	return nil
}

func (s *Scenario) caseSendAllRuntime(ctx context.Context, run *caseRun) error {
	if err := run.ready(); err != nil {
		return err
	}

	beneficiary := randomAddress()
	value := new(uint256.Int).Mul(uint256.NewInt(s.options.SendAllValue), uint256.NewInt(1e9))

	receipt, err := run.send(ctx, nil, value, sendAllDeployCode(beneficiary, 0x2a), nil)
	if err != nil {
		return err
	}

	contract := receipt.ContractAddress

	if _, err := run.send(ctx, &contract, nil, nil, nil); err != nil {
		return err
	}

	balance, err := run.client.GetBalanceAt(ctx, beneficiary)
	run.expect("SENDALL moves the whole balance", err == nil && balance.Cmp(value.ToBig()) == 0, "beneficiary %s got %v err %v", beneficiary, balance, err)

	left, err := run.client.GetBalanceAt(ctx, contract)
	run.expect("SENDALL empties the account", err == nil && left.Sign() == 0, "contract %s holds %v err %v", contract, left, err)

	code, err := run.client.GetCodeAt(ctx, contract)
	run.expect("SENDALL keeps the code", err == nil && bytes.Equal(code, sendAllCode(beneficiary)), "code %x err %v", code, err)

	slot, err := run.client.GetStorageAt(ctx, contract, common.Hash{})
	run.expect("SENDALL keeps the storage", err == nil && slot == common.BigToHash(big.NewInt(0x2a)), "slot0 %s err %v", slot, err)

	nonce, err := run.client.GetNonceAt(ctx, contract, nil)
	run.expect("SENDALL keeps the nonce", err == nil && nonce == 1, "nonce %d err %v", nonce, err)

	return nil
}

func (s *Scenario) caseSendAllInitCode(ctx context.Context, run *caseRun) error {
	if err := run.ready(); err != nil {
		return err
	}

	beneficiary := randomAddress()
	value := new(uint256.Int).Mul(uint256.NewInt(s.options.SendAllValue), uint256.NewInt(1e9))

	receipt, err := run.send(ctx, nil, value, sendAllInitCode(beneficiary, 0x07), nil)
	if err != nil {
		return err
	}

	created := receipt.ContractAddress

	balance, err := run.client.GetBalanceAt(ctx, beneficiary)
	run.expect("SENDALL in initcode moves the balance", err == nil && balance.Cmp(value.ToBig()) == 0, "beneficiary %s got %v err %v", beneficiary, balance, err)

	nonce, err := run.client.GetNonceAt(ctx, created, nil)
	run.expect("account created and SENDALLed in one transaction survives", err == nil && nonce == 1, "nonce %d err %v", nonce, err)

	slot, err := run.client.GetStorageAt(ctx, created, common.Hash{})
	run.expect("its storage survives", err == nil && slot == common.BigToHash(big.NewInt(7)), "slot0 %s err %v", slot, err)

	return nil
}

func (s *Scenario) caseMigrate(ctx context.Context, run *caseRun) error {
	if err := run.ready(); err != nil {
		return err
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	account := crypto.PubkeyToAddress(key.PublicKey)
	recoverInput, err := ecrecoverInput(key)
	if err != nil {
		return err
	}

	out, err := s.call(ctx, run.client, ecrecoverAddress, recoverInput)
	run.expect("ecRecover returns an account without code", err == nil && bytes.Equal(out, common.LeftPadBytes(account.Bytes(), 32)), "got %x err %v", out, err)

	auth, err := txtypes.SignAuthorization(txtypes.SetCodeAuthorization{
		ChainID: *uint256.MustFromBig(s.walletPool.GetChainId()),
		Address: s.adopter,
		Nonce:   0,
	}, key)
	if err != nil {
		return err
	}

	fund := new(uint256.Int).Mul(uint256.NewInt(s.options.FundAmount), uint256.NewInt(1e9))

	if _, err := run.send(ctx, &account, fund, nil, []txtypes.SetCodeAuthorization{auth}); err != nil {
		return err
	}

	code, err := run.client.GetCodeAt(ctx, account)
	run.expect("account carries the delegation", err == nil && bytes.Equal(code, append(common.FromHex("ef0100"), s.adopter.Bytes()...)), "code %x err %v", code, err)

	out, err = s.call(ctx, run.client, ecrecoverAddress, recoverInput)
	run.expect("ecRecover returns a delegated account", err == nil && bytes.Equal(out, common.LeftPadBytes(account.Bytes(), 32)), "got %x err %v", out, err)

	out, err = s.call(ctx, run.client, s.adopter, common.LeftPadBytes(account.Bytes(), 32))
	run.expect("a delegated account is not a SETCODEFROM source", err == nil && word(out) == 0, "got %x err %v", out, err)

	if _, err := run.send(ctx, &account, nil, common.LeftPadBytes(s.approverSource.Bytes(), 32), nil); err != nil {
		return err
	}

	same, err := s.codeHashesMatch(ctx, run.client, account, s.approverSource)
	run.expect("account adopted the approver code", err == nil && same, "%s err %v", account, err)

	out, err = s.call(ctx, run.client, ecrecoverAddress, recoverInput)
	run.expect("ecRecover refuses an account with code", err == nil && len(out) == 0, "got %x err %v", out, err)

	nonce, err := run.client.GetNonceAt(ctx, account, nil)
	if err != nil {
		return err
	}

	err = s.sendECDSA(ctx, run.client, key, nonce, run.wallet.GetAddress())
	run.expect("ECDSA transaction from the migrated account is refused", err != nil, "node accepted it")

	if err != nil {
		run.logger.Debugf("ECDSA transaction refused: %v", err)
	}

	receipt, err := s.sendFrameFrom(ctx, run.client, account, nonce, run.wallet.GetAddress())
	run.expect("frame transaction approved by the adopted code lands", err == nil && receipt != nil && receipt.Status == types.ReceiptStatusSuccessful, "receipt %v err %v", receiptStatus(receipt), err)

	return nil
}

func receiptStatus(receipt *txtypes.Receipt) string {
	if receipt == nil {
		return "none"
	}

	return fmt.Sprintf("status %d in block %v", receipt.Status, receipt.BlockNumber)
}
