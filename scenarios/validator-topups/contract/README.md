## Commands

```
solc --standard-json ./BatchTopUps.input.json | jq '.contracts["contracts/BatchTopUps.sol"].BatchTopUps.abi' > BatchTopUps.abi
solc --standard-json ./BatchTopUps.input.json | jq -r '.contracts["contracts/BatchTopUps.sol"].BatchTopUps.evm.bytecode.object' > BatchTopUps.bin
abigen --bin=./BatchTopUps.bin --abi=./BatchTopUps.abi --pkg=contract --out=BatchTopUps.go
```