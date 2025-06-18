## Commands

```
solc --standard-json ./CREATE2Factory.input.json | jq '.contracts["CREATE2Factory.sol"].CREATE2Factory.abi' > CREATE2Factory.abi
solc --standard-json ./CREATE2Factory.input.json | jq -r '.contracts["CREATE2Factory.sol"].CREATE2Factory.evm.bytecode.object' > CREATE2Factory.bin
abigen --bin=./CREATE2Factory.bin --abi=./CREATE2Factory.abi --pkg=contract --out=CREATE2Factory.go
```