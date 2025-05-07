# Contract Compilation

This directory contains the Solidity contract and its compiled artifacts. To compile the contract and generate Go bindings:

1. Install solc (Solidity compiler):
```bash
# On macOS
brew install solidity

# On Ubuntu/Debian
sudo add-apt-repository ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get install solc
```

2. Install abigen (ABI generator):
```bash
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

3. Compile the contract:
```bash
solc --abi StateBloatToken.sol -o . --overwrite
solc --bin StateBloatToken.sol -o . --overwrite
```

4. Generate Go bindings:
```bash
abigen --bin=StateBloatToken.bin --abi=StateBloatToken.abi --pkg=contract --out=StateBloatToken.go
```

The generated files will be:
- `StateBloatToken.abi` - Contract ABI
- `StateBloatToken.bin` - Contract bytecode
- `StateBloatToken.go` - Go bindings 