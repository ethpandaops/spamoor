#!/bin/bash

compile_contract() {
    local workdir=$1
    local solc_version=$2
    local solc_args=$3
    local contract_file=$4
    local contract_name=$5
    
    if [ -z "$contract_name" ]; then
        contract_name="$contract_file"
    fi

    #echo "docker run --rm -v $workdir:/contracts ethereum/solc:$solc_version /contracts/$contract_file.sol --combined-json abi,bin $solc_args"
    local contract_json=$(docker run --rm -v $workdir:/contracts ethereum/solc:$solc_version /contracts/$contract_file.sol --combined-json abi,bin $solc_args)

    local contract_abi=$(echo "$contract_json" | jq -r '.contracts["/contracts/'$contract_file'.sol:'$contract_name'"].abi')
    if [ "$contract_abi" == "null" ]; then
        contract_abi=$(echo "$contract_json" | jq -r '.contracts["contracts/'$contract_file'.sol:'$contract_name'"].abi')
    fi

    local contract_bin=$(echo "$contract_json" | jq -r '.contracts["/contracts/'$contract_file'.sol:'$contract_name'"].bin')
    if [ "$contract_bin" == "null" ]; then
        contract_bin=$(echo "$contract_json" | jq -r '.contracts["contracts/'$contract_file'.sol:'$contract_name'"].bin')
    fi

    echo "$contract_abi" > $contract_name.abi
    echo "$contract_bin" > $contract_name.bin
    docker run --rm -u $(id -u):$(id -g) -v $(pwd):/workspace ethereum/client-go:alltools-latest abigen --bin=/workspace/$contract_name.bin --abi=/workspace/$contract_name.abi --pkg=contract --out=/workspace/$contract_file.go --type $contract_name
    rm $contract_name.bin $contract_name.abi
    echo "$contract_json" | jq > $contract_file.output.json
}
