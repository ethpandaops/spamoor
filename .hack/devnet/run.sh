#!/bin/bash
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ -f "${__dir}/custom-kurtosis.devnet.config.yaml" ]; then
  config_file="${__dir}/custom-kurtosis.devnet.config.yaml"
else
  config_file="${__dir}/kurtosis.devnet.config.yaml"
fi

## Run devnet using kurtosis
ENCLAVE_NAME="${ENCLAVE_NAME:-spamoor}"
ETHEREUM_PACKAGE="${ETHEREUM_PACKAGE:-github.com/ethpandaops/ethereum-package}"
if kurtosis enclave inspect "$ENCLAVE_NAME" > /dev/null; then
  echo "Kurtosis enclave '$ENCLAVE_NAME' is already up."
else
  kurtosis run "$ETHEREUM_PACKAGE" \
  --image-download always \
  --enclave "$ENCLAVE_NAME" \
  --args-file "${config_file}"

  # Stop spamoor instance within ethereum-package if running
  kurtosis service stop "$ENCLAVE_NAME" spamoor > /dev/null || true
fi

# Get chain config
kurtosis files inspect "$ENCLAVE_NAME" el_cl_genesis_data ./config.yaml | tail -n +2 > "${__dir}/generated-chain-config.yaml"

# Get spamoor hosts
kurtosis files inspect "$ENCLAVE_NAME" spamoor-config rpc-hosts.txt | tail -n +2 > "${__dir}/generated-hosts.txt"


cat <<EOF
============================================================================================================
Spamoor hosts at ${__dir}/generated-hosts.txt
============================================================================================================
EOF
