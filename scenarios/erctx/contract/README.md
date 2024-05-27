## Commands

```
solc --standard-json ./TestToken.input.json | jq '.contracts["contracts/TestToken.sol"].TestToken.abi' > TestToken.abi
solc --standard-json ./TestToken.input.json | jq -r '.contracts["contracts/TestToken.sol"].TestToken.evm.bytecode.object' > TestToken.bin
abigen --bin=./TestToken.bin --abi=./TestToken.abi --pkg=contract --out=TestToken.go
```