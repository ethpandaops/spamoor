package main

import (
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/holiman/uint256"
	flag "github.com/spf13/pflag"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethpandaops/spamoor/txbuilder"
)

type CliArgs struct {
	verbose       bool
	rpchost       string
	privkey       string
	randomPrivKey bool

	txCount      uint
	txTo         string
	txValue      uint64
	txData       string
	txBlobs      []string
	gaslimit     uint64
	maxfeepergas float32
	maxpriofee   float32
	maxblobfee   float32

	output   bool
	chainid  uint64
	nonce    uint64
	addnonce int64
}

func main() {
	cliArgs := CliArgs{}

	flag.BoolVarP(&cliArgs.verbose, "verbose", "v", false, "Run the script with verbose output")
	flag.StringVarP(&cliArgs.rpchost, "rpchost", "r", "http://127.0.0.1:8545", "The RPC host to send transactions to.")
	flag.StringVarP(&cliArgs.privkey, "privkey", "p", "", "The private key of the wallet to send blobs from.\n(Special: \"env\" to read from BLOBSENDER_PRIVKEY environment variable)")
	flag.BoolVar(&cliArgs.randomPrivKey, "random-privkey", false, "Use random private key if no privkey supplied")

	flag.UintVarP(&cliArgs.txCount, "count", "n", 1, "The number of transactions to send.")
	flag.StringVarP(&cliArgs.txTo, "to", "t", "", "The transaction to address.")
	flag.Uint64VarP(&cliArgs.txValue, "value", "a", 0, "The transaction value.")
	flag.StringVarP(&cliArgs.txData, "data", "d", "", "The transaction calldata.")
	flag.StringArrayVarP(&cliArgs.txBlobs, "blobs", "b", []string{}, "The blobs to reference in the transaction (binary file).")
	flag.Uint64Var(&cliArgs.gaslimit, "gaslimit", 500000, "The gas limit for transactions.")
	flag.Float32Var(&cliArgs.maxfeepergas, "maxfeepergas", 20, "The gas limit for transactions.")
	flag.Float32Var(&cliArgs.maxpriofee, "maxpriofee", 1.2, "The maximum priority fee per gas in gwei.")
	flag.Float32Var(&cliArgs.maxblobfee, "maxblobfee", 10, "The maximum blob fee per chunk in gwei.")

	flag.BoolVarP(&cliArgs.output, "output", "o", false, "Output signed transactions to stdout instead of broadcasting them (offline mode).")
	flag.Uint64Var(&cliArgs.chainid, "chainid", 0, "ChainID of the network (For offline mode in combination with --output or to override transactions)")
	flag.Uint64Var(&cliArgs.nonce, "nonce", 0, "Current nonce of the wallet (For offline mode in combination with --output)")
	flag.Int64Var(&cliArgs.addnonce, "addnonce", 0, "Nonce offset to use for transactions (useful for replacement transactions)")

	flag.Parse()

	var client *txbuilder.Client
	wallet, err := loadPrivkey(&cliArgs)
	if err != nil {
		panic(err)
	}

	if !cliArgs.output {
		client, err = txbuilder.NewClient(cliArgs.rpchost)
		if err != nil {
			panic(err)
		}
		err = client.UpdateWallet(wallet)
		if err != nil {
			panic(err)
		}
	}
	if cliArgs.chainid != 0 {
		wallet.SetChainId(big.NewInt(int64(cliArgs.chainid)))
	}
	if cliArgs.nonce != 0 {
		wallet.SetNonce(cliArgs.nonce)
	}
	if cliArgs.addnonce != 0 {
		nonce := int64(wallet.GetNonce()) + cliArgs.addnonce
		if nonce < 0 {
			panic(fmt.Errorf("cannot use negative nonce"))
		}
		wallet.SetNonce(uint64(nonce))
	}

	toAddr := common.HexToAddress(cliArgs.txTo)
	txMetadata := txbuilder.TxMetadata{
		GasTipCap:  uint256.NewInt(uint64(cliArgs.maxpriofee * 1000000000)),
		GasFeeCap:  uint256.NewInt(uint64(cliArgs.maxfeepergas * 1000000000)),
		BlobFeeCap: uint256.NewInt(uint64(cliArgs.maxblobfee * 1000000000)),
		Gas:        cliArgs.gaslimit,
		To:         &toAddr,
		Value:      uint256.NewInt(cliArgs.txValue),
		Data:       common.FromHex(cliArgs.txData),
	}
	blobRefs := [][]string{}
	for _, blobRef := range cliArgs.txBlobs {
		blobRefs = append(blobRefs, []string{blobRef})
	}

	for idx := 0; idx < int(cliArgs.txCount); idx++ {
		txData, err := txbuilder.BuildBlobTx(&txMetadata, blobRefs)
		if err != nil {
			panic(err)
		}
		tx, err := wallet.BuildBlobTx(txData)
		if err != nil {
			panic(err)
		}
		txBytes, err := tx.MarshalBinary()
		if err != nil {
			panic(err)
		}

		if cliArgs.output {
			fmt.Println(common.Bytes2Hex(txBytes))
		} else {
			txHash := client.SubmitTransaction(txBytes)
			if txHash != nil {
				fmt.Printf("TX %v: %v\n", idx+1, txHash.String())
			}
		}
	}

}

func loadPrivkey(cliArgs *CliArgs) (*txbuilder.Wallet, error) {
	privkey := cliArgs.privkey
	if privkey == "env" {
		privkey = os.Getenv("BLOBSENDER_PRIVKEY")
	}
	if privkey == "" && !cliArgs.randomPrivKey {
		return nil, errors.New("No private key specified.")
	}
	return txbuilder.NewWallet(privkey)
}
