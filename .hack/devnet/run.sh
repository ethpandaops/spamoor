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

## Generate spamoor hosts
ENCLAVE_UUID=$(kurtosis enclave inspect "$ENCLAVE_NAME" --full-uuids | grep 'UUID:' | awk '{print $2}')

EXECUTION_NODES=$(docker ps -aq -f "label=kurtosis_enclave_uuid=$ENCLAVE_UUID" \
              -f "label=com.kurtosistech.app-id=kurtosis" \
              -f "label=com.kurtosistech.custom.ethereum-package.client-type=execution" | tac)

cat <<EOF > "${__dir}/generated-hosts.txt"
$(for node in $EXECUTION_NODES; do
  name=$(docker inspect -f "{{ with index .Config.Labels \"com.kurtosistech.id\"}}{{.}}{{end}}" $node)
  ip=$(echo '127.0.0.1')
  port=$(docker inspect --format='{{ (index (index .NetworkSettings.Ports "8545/tcp") 0).HostPort }}' $node)
  if [ -z "$port" ]; then
    port="65535"
  fi

  IFS='-' read -r -a name_parts <<< "$name"
  trimmed_name="${name_parts[1]}-${name_parts[2]}-${name_parts[3]}"
  
  echo "group(${name_parts[1]},${name_parts[2]},${name_parts[3]})name($trimmed_name)http://$ip:$port"
done)
EOF


cat <<EOF
============================================================================================================
Spamoor hosts at ${__dir}/generated-hosts.txt
============================================================================================================
EOF
