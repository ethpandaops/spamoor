#!/usr/bin/env bash
#
# Fails if the core packages import go-ethereum's transaction type system.
#
# The engine works on txtypes so that spamoor can send transaction types go-ethereum
# has not implemented. core/types and ethclient are allowed only in the files listed
# below, which exist to convert at the boundary.

set -euo pipefail

cd "$(dirname "$0")/.."

CORE_DIRS=(spamoor scenario txbuilder txtypes utils)

# path:package pairs exempt from the check
ALLOWED=(
	"spamoor/compat.go:core/types"
	"spamoor/wallet.go:core/types"
	"spamoor/client.go:ethclient"
	"txtypes/types.go:core/types"
	"txtypes/setcode.go:core/types"
	"txtypes/geth.go:core/types"
	"txtypes/geth_diff_test.go:core/types"
)

is_allowed() {
	local entry="$1:$2"
	for allowed in "${ALLOWED[@]}"; do
		[ "$allowed" = "$entry" ] && return 0
	done
	return 1
}

violations=0

for pkg in core/types ethclient; do
	while IFS= read -r file; do
		[ -z "$file" ] && continue
		if ! is_allowed "$file" "$pkg"; then
			echo "error: $file imports github.com/ethereum/go-ethereum/$pkg"
			violations=$((violations + 1))
		fi
	done < <(grep -rl "\"github.com/ethereum/go-ethereum/$pkg\"" --include='*.go' "${CORE_DIRS[@]}" 2>/dev/null || true)
done

if [ "$violations" -gt 0 ]; then
	cat >&2 <<-'EOF'

		The core packages must not depend on go-ethereum's transaction types.
		Use github.com/ethpandaops/spamoor/txtypes instead, and convert at the
		boundary with txtypes.FromGethTx / (*Transaction).ToGethTx if a
		go-ethereum value is genuinely required.

		If a new boundary file is justified, add it to ALLOWED in this script.
	EOF
	exit 1
fi

echo "core import check passed"
