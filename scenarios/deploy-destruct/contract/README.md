## Commands

```
solc --standard-json ./DeployTest.input.json | jq '.contracts["DeployTest.sol"].DeployTest.abi' > DeployTest.abi
solc --standard-json ./DeployTest.input.json | jq -r '.contracts["DeployTest.sol"].DeployTest.evm.bytecode.object' > DeployTest.bin
abigen --bin=./DeployTest.bin --abi=./DeployTest.abi --pkg=contract --out=DeployTest.go
```