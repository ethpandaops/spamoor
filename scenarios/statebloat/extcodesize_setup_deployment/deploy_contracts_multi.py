#!/usr/bin/env python3
"""
Deploy multiple contracts via CREATE2 factories for different sizes for BloatNet benchmarks.

This script deploys contracts of specific sizes using the corresponding CREATE2 factories.

USAGE:
    1. First deploy initcode contracts:
       python3 deploy_initcode_multi.py

    2. Then deploy the factories:
       python3 deploy_factory_multi.py

    3. Finally deploy contracts:
       python3 deploy_contracts_multi.py --size 5 --count 1000
       python3 deploy_contracts_multi.py --size 24 --count 350

    4. Or deploy remaining contracts for all incomplete sizes (mainnet):
       python3 deploy_contracts_multi.py --deploy-remaining --count 18100

    5. Or use a specific factory address directly:
       python3 deploy_contracts_multi.py --size 24 --count 18100 --factory 0xa4a1aF502114DAA5856b6FBD849e14a535A69eE8

    6. Run EEST benchmarks:
       uv run execute remote --fork Prague \\
         --rpc-endpoint http://127.0.0.1:8545 \\
         --address-stubs stubs.json \\
         -- --gas-benchmark-values 30 \\
         tests/benchmark/stateful/bloatnet/test_extcodesize_bytecode_sizes.py -v

REQUIREMENTS:
    - web3.py
    - eth-utils
    - Local geth instance running on http://127.0.0.1:8545
    - stubs.json from deploy_factory_multi.py (or use --factory flag)
"""

import sys
import json
import time
from eth_utils import keccak
from web3 import Web3
from eth_account import Account

# Mainnet factory addresses (from extcodesize_setup_mainnet deployment)
MAINNET_FACTORIES = {
    0.5: "0xE19DA6B7c3716e79dccb84a591B2F36aF03487D5",
    1: "0x6bB03b0bB8b6528243912834D520FA71d77D420B",
    2: "0x61b25103ff58A81c3079C23e57352a62f05d38dd",
    5: "0xa35C3fDdCa307786E74d657A6b4c8952db440076",
    10: "0x401ef5eBd3F749B839e13363aDD39bb85Bb8B204",
    24: "0xa4a1aF502114DAA5856b6FBD849e14a535A69eE8",
    32: "0xFc7c3eda39Aa275fA0483942ba7e56195C7aE321",
    64: "0xd7238974c4d3c3241a974b8e1A7C0b594a5d1e14",
}


def estimate_gas_for_size(size_kb: float, block_gas_limit: int) -> int:
    """Estimate gas needed to deploy a contract of given size, respecting network limits.

    Based on actual measurements:
    - 0.5KB deployment used ~183,163 gas

    Gas costs breakdown:
    - Transaction intrinsic: 21,000
    - Factory execution: ~1,000
    - CREATE2 overhead: ~32,000
    - Contract bytecode storage: 200 gas per byte
    - Init code execution: varies by size
    """
    size_bytes = int(size_kb * 1024)

    # Precise gas calculation based on actual measurements
    if size_kb <= 0.5:
        # 0.5KB used 183,163 gas = ~21K intrinsic + ~32K CREATE2 + ~102K storage (512*200) + ~28K execution
        base_gas = 21_000 + 32_000 + (size_bytes * 200) + 30_000
    elif size_kb <= 1:
        # 1KB: similar but with more storage cost
        base_gas = 21_000 + 32_000 + (size_bytes * 200) + 35_000
    elif size_kb <= 5:
        # 5KB uses While loop for init, more execution cost
        base_gas = 21_000 + 32_000 + (size_bytes * 200) + 150_000
    elif size_kb <= 10:
        # 10KB: more While iterations
        base_gas = 21_000 + 32_000 + (size_bytes * 200) + 300_000
    elif size_kb <= 24:
        # 24KB: maximum standard contract size
        base_gas = 21_000 + 32_000 + (size_bytes * 200) + 500_000
    elif size_kb <= 32:
        # 32KB: exceeds EIP-170 limit, needs more execution
        base_gas = 21_000 + 32_000 + (size_bytes * 200) + 700_000
    else:
        # 64KB: very large contracts
        base_gas = 21_000 + 32_000 + (size_bytes * 200) + 1_000_000

    # Add 10% buffer for safety
    final_gas = int(base_gas * 1.1)

    # Cap at 80% of block gas limit to ensure inclusion
    max_safe_gas = int(block_gas_limit * 0.8)

    return min(final_gas, max_safe_gas)


def deploy_contracts(w3: Web3, factory_address: str, count: int, size_kb: float,
                     private_key: str = None, sender_address: str = None) -> int:
    """Deploy contracts via a CREATE2 factory.

    If private_key is provided, transactions are signed locally.
    Otherwise, uses w3.eth.accounts[0] (requires unlocked account on node).
    """
    if private_key:
        account = sender_address
        use_local_signing = True
    else:
        account = w3.eth.accounts[0]
        use_local_signing = False

    # Get network parameters dynamically
    latest_block = w3.eth.get_block('latest')
    block_gas_limit = latest_block.gasLimit
    print(f"Network block gas limit: {block_gas_limit:,}")

    # Get current counter
    current = int.from_bytes(w3.eth.get_storage_at(factory_address, 0), 'big')
    print(f"Factory has already deployed {current} contracts")

    if count <= current:
        print(f"✅ Already have {current} contracts (target: {count})")
        return current

    remaining = count - current
    print(f"Deploying {remaining} more contracts of {size_kb}KB...")

    # Estimate gas needed based on network's block gas limit
    gas_limit = estimate_gas_for_size(size_kb, block_gas_limit)
    print(f"Using gas limit: {gas_limit:,} per deployment")
    print(f"  (Network allows up to {int(block_gas_limit * 0.8):,} per transaction)")

    # Calculate optimal batch size based on gas costs
    # Use actual network block gas limit to maximize throughput
    # Each deployment is a separate transaction with our current factory
    # Leave some margin for safety (use 95% of block gas limit)
    usable_gas = int(block_gas_limit * 0.95)
    deployments_per_block = usable_gas // gas_limit

    print(f"Optimal batch size: {deployments_per_block} deployments per block")
    print(f"  (Each deployment uses {gas_limit:,} gas)")
    print(f"  (Network block limit is {block_gas_limit:,}, using {usable_gas:,})")

    # Deploy contracts in optimized batches
    batch_size = deployments_per_block
    deployed = 0
    failed = 0
    start_time = time.time()

    print(f"\nDeploying {remaining} contracts in batches of {batch_size}...")
    print(f"Expected blocks needed: {(remaining + batch_size - 1) // batch_size}")

    for batch_start in range(0, remaining, batch_size):
        batch_end = min(batch_start + batch_size, remaining)
        batch_count = batch_end - batch_start

        # Get fresh nonce for this batch including pending txs to avoid "already known" errors
        nonce = w3.eth.get_transaction_count(account, "pending")

        # Send batch of transactions rapidly to fill a block
        tx_hashes = []
        batch_time = time.time()

        print(f"\nBatch {batch_start//batch_size + 1}: Sending {batch_count} transactions rapidly to fill block...")
        print(f"  Starting nonce: {nonce}")

        # Send all transactions as fast as possible with pre-calculated nonces
        for i in range(batch_count):
            try:
                tx = {
                    "from": account,
                    "to": factory_address,
                    "data": "0x01",  # Non-empty data to trigger CREATE2
                    "gas": gas_limit,
                    "gasPrice": w3.to_wei(22, 'gwei'),  # base_fee(20) + tip(2) to match spamoor config
                    "nonce": nonce + i,  # Pre-calculate nonce for speed
                    "chainId": w3.eth.chain_id,
                }

                if use_local_signing:
                    # Sign locally and send raw transaction
                    signed_tx = Account.sign_transaction(tx, private_key)
                    tx_hash = w3.eth.send_raw_transaction(signed_tx.raw_transaction)
                else:
                    # Use node's unlocked account
                    tx_hash = w3.eth.send_transaction(tx)

                tx_hashes.append(tx_hash)
            except Exception as e:
                print(f"  Error sending tx {batch_start + i + 1}: {e}")
                # Try to recover and continue
                break

        send_time = time.time() - batch_time
        print(f"  Sent {len(tx_hashes)} transactions in {send_time:.2f}s ({len(tx_hashes)/send_time:.1f} tx/s)")
        print(f"  Waiting for confirmations...")

        # Wait for batch receipts
        batch_deployed = 0
        for j, tx_hash in enumerate(tx_hashes):
            try:
                receipt = w3.eth.wait_for_transaction_receipt(tx_hash, timeout=60)
                if receipt['status'] == 1:
                    deployed += 1
                    batch_deployed += 1
                else:
                    failed += 1
                    print(f"  Transaction {j+1} failed")
            except Exception as e:
                print(f"  Transaction {j+1} failed or timed out: {e}")
                failed += 1

        # Progress update
        counter = int.from_bytes(w3.eth.get_storage_at(factory_address, 0), 'big')
        elapsed = time.time() - start_time
        rate = deployed / elapsed if elapsed > 0 else 0
        eta = (remaining - deployed) / rate if rate > 0 else 0
        batch_elapsed = time.time() - batch_time

        print(f"  Batch complete: {batch_deployed}/{batch_count} deployed in {batch_elapsed:.1f}s")
        print(f"  Gas used per deployment: {gas_limit:,} ({batch_deployed * gas_limit:,} total)")
        print(f"  Overall: {counter}/{count} contracts ({deployed}/{remaining} new)")
        print(f"  Rate: {rate:.1f} contracts/sec, ETA: {eta:.0f}s")

        if failed > 20:
            print("\n⚠️ Too many failures, stopping...")
            break

    # Final check
    final_counter = int.from_bytes(w3.eth.get_storage_at(factory_address, 0), 'big')
    elapsed = time.time() - start_time
    print(f"\n✅ Deployment complete in {elapsed:.1f} seconds")
    print(f"Total contracts deployed: {final_counter}")

    return final_counter


def verify_contracts(w3: Web3, factory_address: str, count: int, size_kb: float) -> bool:
    """Verify that contracts exist at expected CREATE2 addresses."""
    print(f"\n--- Verifying CREATE2 Addresses for {size_kb}KB contracts ---")

    # Get init code hash from factory storage
    stored_hash = w3.eth.get_storage_at(factory_address, 1)

    # Verify a sample of contracts
    sample_size = min(5, count)
    verified = 0

    for salt in range(sample_size):
        # Calculate CREATE2 address
        create2_input = (
            b"\xff" +
            bytes.fromhex(factory_address[2:].lower()) +
            salt.to_bytes(32, "big") +
            stored_hash
        )
        expected_addr = Web3.to_checksum_address("0x" + keccak(create2_input)[-20:].hex())

        # Check if contract exists
        code = w3.eth.get_code(expected_addr)
        if len(code) > 0:
            print(f"  Salt {salt}: ✅ Found at {expected_addr} ({len(code)} bytes)")
            verified += 1
        else:
            print(f"  Salt {salt}: ❌ Not found at {expected_addr}")

    return verified == sample_size


def get_deployment_status(w3: Web3, target_count: int) -> dict:
    """Get current deployment status for all mainnet factories."""
    status = {}
    for size_kb, factory_addr in MAINNET_FACTORIES.items():
        current = int.from_bytes(w3.eth.get_storage_at(factory_addr, 0), 'big')
        remaining = max(0, target_count - current)
        status[size_kb] = {
            'factory': factory_addr,
            'deployed': current,
            'target': target_count,
            'remaining': remaining,
            'complete': current >= target_count
        }
    return status


def deploy_remaining_contracts(w3: Web3, target_count: int,
                               private_key: str = None, sender_address: str = None) -> None:
    """Deploy remaining contracts for all incomplete sizes."""
    print("=" * 60)
    print("DEPLOYING REMAINING CONTRACTS FOR ALL INCOMPLETE SIZES")
    print("=" * 60)

    status = get_deployment_status(w3, target_count)

    # Show current status
    print("\nCurrent deployment status:")
    for size_kb in sorted(status.keys()):
        info = status[size_kb]
        pct = (info['deployed'] / info['target'] * 100) if info['target'] > 0 else 0
        marker = "✅" if info['complete'] else "🔄"
        print(f"  {size_kb:>5}KB: {info['deployed']:>6} / {info['target']:>6} ({pct:5.1f}%) {marker}")

    # Deploy for each incomplete size
    incomplete = [(size_kb, info) for size_kb, info in status.items() if not info['complete']]

    if not incomplete:
        print("\n✅ All sizes are complete!")
        return

    print(f"\n{len(incomplete)} sizes need more contracts:")
    for size_kb, info in incomplete:
        print(f"  {size_kb}KB: {info['remaining']} remaining")

    for size_kb, info in incomplete:
        print(f"\n{'=' * 60}")
        print(f"DEPLOYING {size_kb}KB CONTRACTS")
        print(f"{'=' * 60}")
        deploy_contracts(w3, info['factory'], target_count, size_kb, private_key, sender_address)

    # Final status
    print(f"\n{'=' * 60}")
    print("FINAL DEPLOYMENT STATUS")
    print("=" * 60)
    final_status = get_deployment_status(w3, target_count)
    for size_kb in sorted(final_status.keys()):
        info = final_status[size_kb]
        pct = (info['deployed'] / info['target'] * 100) if info['target'] > 0 else 0
        marker = "✅" if info['complete'] else "❌"
        print(f"  {size_kb:>5}KB: {info['deployed']:>6} / {info['target']:>6} ({pct:5.1f}%) {marker}")


def main():
    """Main deployment script."""
    import argparse

    # Anvil's default account private key
    ANVIL_DEFAULT_PRIVATE_KEY = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

    parser = argparse.ArgumentParser(description='Deploy BloatNet contracts via CREATE2 factory')
    parser.add_argument('--size', type=float,
                        choices=[0.5, 1, 2, 5, 10, 24, 32, 64],
                        help='Contract size in KB')
    parser.add_argument('--count', type=int,
                        help='Total number of contracts to deploy')
    parser.add_argument('--rpc-url', default='http://127.0.0.1:8545',
                        help='RPC URL')
    parser.add_argument('--stubs', default='stubs.json',
                        help='Path to stubs JSON file')
    parser.add_argument('--factory', type=str,
                        help='Factory address (overrides stubs.json lookup)')
    parser.add_argument('--use-mainnet-factories', action='store_true',
                        help='Use hardcoded mainnet factory addresses')
    parser.add_argument('--deploy-remaining', action='store_true',
                        help='Deploy remaining contracts for all incomplete sizes')
    parser.add_argument('--status', action='store_true',
                        help='Show deployment status and exit')
    parser.add_argument('--private-key', '-p', type=str,
                        help='Private key for signing transactions (signs locally instead of using node account)')
    parser.add_argument('--anvil', action='store_true',
                        help='Use Anvil default account (0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266)')
    args = parser.parse_args()

    # Validate arguments
    if not args.deploy_remaining and not args.status and not args.size:
        parser.error("--size is required unless using --deploy-remaining or --status")
    if not args.deploy_remaining and not args.status and not args.count:
        parser.error("--count is required unless using --deploy-remaining or --status")

    # Connect to local geth instance
    w3 = Web3(Web3.HTTPProvider(args.rpc_url))
    if not w3.is_connected():
        print(f"❌ Failed to connect to {args.rpc_url}")
        sys.exit(1)

    print(f"Connected to: {args.rpc_url}")

    # Handle private key / account setup
    private_key = None
    sender_address = None

    if args.anvil:
        private_key = ANVIL_DEFAULT_PRIVATE_KEY
        sender_address = Account.from_key(private_key).address
        print(f"Using Anvil default account: {sender_address}")
    elif args.private_key:
        private_key = args.private_key
        sender_address = Account.from_key(private_key).address
        print(f"Using account from private key: {sender_address}")
    else:
        try:
            sender_address = w3.eth.accounts[0]
            print(f"Using node account: {sender_address}")
        except Exception:
            print("❌ No accounts available on node. Use --private-key or --anvil")
            sys.exit(1)

    # Show account balance
    balance = w3.eth.get_balance(sender_address)
    print(f"Account balance: {w3.from_wei(balance, 'ether'):.4f} ETH")

    # Handle --status mode
    if args.status:
        target = args.count or 18100
        status = get_deployment_status(w3, target)
        print(f"\nDeployment status (target: {target}):")
        for size_kb in sorted(status.keys()):
            info = status[size_kb]
            pct = (info['deployed'] / info['target'] * 100) if info['target'] > 0 else 0
            marker = "✅" if info['complete'] else "🔄"
            print(f"  {size_kb:>5}KB: {info['deployed']:>6} / {info['target']:>6} ({pct:5.1f}%) {marker} {info['factory']}")
        return

    # Handle --deploy-remaining mode
    if args.deploy_remaining:
        target = args.count or 18100
        deploy_remaining_contracts(w3, target, private_key, sender_address)
        return

    # Single size deployment mode
    factory_address = None

    # Priority: --factory > --use-mainnet-factories > stubs.json
    if args.factory:
        factory_address = args.factory
        print(f"Using provided factory address: {factory_address}")
    elif args.use_mainnet_factories:
        factory_address = MAINNET_FACTORIES.get(args.size)
        if not factory_address:
            print(f"❌ No mainnet factory for {args.size}KB")
            print(f"Available sizes: {list(MAINNET_FACTORIES.keys())}")
            sys.exit(1)
        print(f"Using mainnet factory: {factory_address}")
    else:
        # Load factory address from stubs
        try:
            with open(args.stubs, 'r') as f:
                stubs = json.load(f)
        except FileNotFoundError:
            # Fall back to mainnet factories
            print(f"⚠️ Stubs file not found: {args.stubs}, using mainnet factories")
            factory_address = MAINNET_FACTORIES.get(args.size)
            if not factory_address:
                print(f"❌ No mainnet factory for {args.size}KB")
                sys.exit(1)
            print(f"Using mainnet factory: {factory_address}")

        if not factory_address:
            # Find the appropriate factory from stubs
            # Handle size key naming: 0.5 -> "0_5kb", 1.0 -> "1kb", 24.0 -> "24kb"
            if args.size == int(args.size):
                size_key = f"{int(args.size)}kb"
            else:
                size_key = f"{args.size}kb".replace(".", "_")
            factory_key = f"bloatnet_factory_{size_key}"
            factory_address = stubs.get(factory_key)

            if not factory_address:
                print(f"❌ Factory not found for {args.size}KB contracts")
                print(f"Looking for key: {factory_key}")
                print(f"Available factories: {list(stubs.keys())}")
                print(f"\nTip: Use --use-mainnet-factories to use hardcoded addresses")
                sys.exit(1)

            print(f"Using factory from stubs: {factory_address}")

    # Deploy contracts
    final_count = deploy_contracts(w3, factory_address, args.count, args.size, private_key, sender_address)

    # Verify deployment
    if verify_contracts(w3, factory_address, final_count, args.size):
        print(f"\n✅ Successfully verified {args.size}KB contracts")

    # Calculate gas requirements for testing
    cost_per_contract = 2660  # Approximate gas per EXTCODESIZE with CREATE2
    test_gas_30m = 30_000_000
    max_contracts_30m = (test_gas_30m - 21000 - 1000) // cost_per_contract

    print(f"\n=== Ready for Testing ===")
    print(f"Contract size: {args.size}KB")
    print(f"Contracts deployed: {final_count}")
    print(f"Factory address: {factory_address}")
    print(f"Max contracts for 30M gas: ~{max_contracts_30m}")

    if final_count < max_contracts_30m:
        print(f"⚠️ Consider deploying {max_contracts_30m - final_count} more contracts "
              f"to fully utilize 30M gas")

    print("\nTo run the benchmark test:")
    print("uv run execute remote --fork Prague \\")
    print(f"  --rpc-endpoint {args.rpc_url} \\")
    print(f"  --address-stubs {args.stubs} \\")
    print("  -- --gas-benchmark-values 30 \\")
    print("  tests/benchmark/stateful/bloatnet/test_extcodesize_bytecode_sizes.py \\")
    print(f"  -k '{args.size}KB' -v")

if __name__ == "__main__":
    main()