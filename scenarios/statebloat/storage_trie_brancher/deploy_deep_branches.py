#!/usr/bin/env python3
"""
Deploy worst-case attack contracts to an Ethereum node.

This script deploys X WorstCaseERC20.sol contracts using Nick's factory method,
funds auxiliary accounts, and saves all deployed addresses to a JSON file.

Usage:
    python deploy_worst_case_contracts.py --rpc-url http://localhost:8546 \
        --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
        --storage-depth 9 --account-depth 3 --num-contracts 15000

The script will:
1. Load CREATE2 salt data from s9_acc3.json (or similar)
2. Deploy all contracts via Nick's factory
3. Fund auxiliary accounts for trie depth
4. Save deployed addresses to deployed_contracts.json
"""

import argparse
import json
import sys
import time
from pathlib import Path
from typing import Dict, List, Any
from eth_account import Account
from web3 import Web3
from web3.types import TxParams, HexBytes
from eth_utils import keccak, to_checksum_address
import subprocess

# Nick's factory address
NICK_FACTORY = "0x4e59b44847b379578588920ca78fbf26c0b4956c"

# Gas settings
MAX_GAS_PER_TX = 16_000_000  # Fusaka limit
DEPLOY_GAS_LIMIT = 3_000_000  # Gas for deploying one contract
FUND_GAS_LIMIT = 21_000  # Gas for funding accounts

def compile_contract(storage_depth: int, account_depth: int) -> bytes:
    """Compile the depth contract and return init bytecode."""
    sol_path = Path(f"depth_{storage_depth}.sol")
    if not sol_path.exists():
        raise FileNotFoundError(f"Contract file not found: {sol_path}")

    print(f"Compiling {sol_path}...")

    # Compile to get init code (with constructor)
    # MUST use same flags as pre-mining to get matching bytecode
    result = subprocess.run(
        [
            "solc",
            "--bin",
            "--optimize",
            "--optimize-runs", "200",
            "--metadata-hash", "none",  # Critical for reproducible bytecode
            str(sol_path)
        ],
        capture_output=True,
        text=True
    )

    if result.returncode != 0:
        raise Exception(f"Failed to compile: {result.stderr}")

    # Extract bytecode from solc output
    lines = result.stdout.split('\n')
    init_code_hex = None
    in_binary_section = False
    for line in lines:
        if "Binary:" in line:
            in_binary_section = True
        elif in_binary_section and line.strip() and not line.startswith("="):
            init_code_hex = line.strip()
            break

    if init_code_hex.startswith("0x"):
        init_code_hex = init_code_hex[2:]

    return bytes.fromhex(init_code_hex)

def load_create2_data(storage_depth: int, account_depth: int) -> Dict:
    """Load pre-mined CREATE2 salt data."""
    filename = f"s{storage_depth}_acc{account_depth}.json"
    filepath = Path(filename)

    if not filepath.exists():
        raise FileNotFoundError(f"CREATE2 data file not found: {filename}")

    with open(filepath, 'r') as f:
        return json.load(f)

def calculate_create2_address(deployer: str, salt: bytes, init_code: bytes) -> str:
    """Calculate CREATE2 address."""
    init_code_hash = keccak(init_code)
    create2_input = (
        bytes.fromhex("ff") +
        bytes.fromhex(deployer[2:]) +
        salt +
        init_code_hash
    )
    addr = keccak(create2_input)[-20:]
    return to_checksum_address("0x" + addr.hex())

def deploy_via_nick_factory(
    w3: Web3,
    account: Account,
    init_code: bytes,
    salt: bytes,
    nonce: int
) -> tuple[str, int]:
    """Deploy contract via Nick's factory."""
    # Call Nick's factory with salt + init_code
    factory_data = salt + init_code

    tx: TxParams = {
        'from': account.address,
        'to': to_checksum_address(NICK_FACTORY),
        'data': factory_data,
        'gas': DEPLOY_GAS_LIMIT,
        'gasPrice': w3.eth.gas_price,
        'nonce': nonce,
    }

    signed_tx = account.sign_transaction(tx)
    tx_hash = w3.eth.send_raw_transaction(signed_tx.raw_transaction)

    # Calculate deployed address
    deployed_addr = calculate_create2_address(NICK_FACTORY, salt, init_code)

    return deployed_addr, nonce + 1

def fund_account(
    w3: Web3,
    account: Account,
    recipient: str,
    amount: int,
    nonce: int
) -> int:
    """Fund an account with specified amount."""
    tx: TxParams = {
        'from': account.address,
        'to': to_checksum_address(recipient),
        'value': amount,
        'gas': FUND_GAS_LIMIT,
        'gasPrice': w3.eth.gas_price,
        'nonce': nonce,
    }

    signed_tx = account.sign_transaction(tx)
    w3.eth.send_raw_transaction(signed_tx.raw_transaction)

    return nonce + 1

def batch_transactions(transactions: List[Any], batch_size: int = 100):
    """Batch transactions to avoid overwhelming the RPC."""
    for i in range(0, len(transactions), batch_size):
        yield transactions[i:i + batch_size]

def main():
    parser = argparse.ArgumentParser(description='Deploy worst-case attack contracts')
    parser.add_argument('--rpc-url', default='http://localhost:8546', help='RPC endpoint URL (can include auth like https://user:pass@host)')
    parser.add_argument('--private-key', required=True, help='Private key for deployment')
    parser.add_argument('--storage-depth', type=int, default=9, help='Storage depth')
    parser.add_argument('--account-depth', type=int, default=3, help='Account depth')
    parser.add_argument('--num-contracts', type=int, default=15000, help='Number of contracts to deploy')
    parser.add_argument('--output', default='deployed_contracts.json', help='Output JSON file')
    parser.add_argument('--skip-contracts', action='store_true', help='Skip contract deployment (only fund EOAs)')
    parser.add_argument('--skip-funding', action='store_true', help='Skip EOA funding (only deploy contracts)')

    args = parser.parse_args()

    # Validate arguments
    if args.skip_contracts and args.skip_funding:
        print("Error: Cannot skip both contracts and funding. Nothing to do!")
        sys.exit(1)

    # Connect to node - HTTPProvider automatically handles basic auth in URL
    print(f"Connecting to RPC endpoint...")
    # Mask credentials in output
    display_url = args.rpc_url
    if '@' in display_url:
        # Extract and mask the auth part
        protocol_end = display_url.find('://') + 3
        auth_end = display_url.find('@', protocol_end)
        if auth_end > protocol_end:
            masked_url = display_url[:protocol_end] + '***:***@' + display_url[auth_end+1:]
            print(f"  URL: {masked_url}")
    else:
        print(f"  URL: {display_url}")

    w3 = Web3(Web3.HTTPProvider(args.rpc_url))

    if not w3.is_connected():
        print("Failed to connect to Ethereum node")
        sys.exit(1)

    # Setup account
    account = Account.from_key(args.private_key)
    print(f"Deploying from account: {account.address}")

    balance = w3.eth.get_balance(account.address)
    print(f"Account balance: {Web3.from_wei(balance, 'ether')} ETH")

    if balance == 0:
        print("\nWARNING: Account has 0 ETH balance!")
        print("You need to fund this account before deploying contracts or funding EOAs.")
        print(f"Send ETH to: {account.address}")
        response = input("\nDo you want to continue anyway? (y/n): ")
        if response.lower() != 'y':
            print("Exiting...")
            sys.exit(1)

    # Load CREATE2 data
    print(f"\nLoading CREATE2 data for depth {args.storage_depth}, account depth {args.account_depth}...")
    create2_data = load_create2_data(args.storage_depth, args.account_depth)
    contracts = create2_data.get("contracts", [])[:args.num_contracts]

    if len(contracts) < args.num_contracts:
        print(f"Warning: Only {len(contracts)} contracts available in data file")

    # We always need to compile or know the init code to calculate addresses
    if not args.skip_contracts:
        print(f"Will deploy {len(contracts)} contracts")
    else:
        print("Skipping contract deployment (--skip-contracts flag set)")

    # Compile contract (needed for address calculation even if skipping deployment)
    init_code = compile_contract(args.storage_depth, args.account_depth)
    print(f"Contract init code size: {len(init_code)} bytes")

    # Get starting nonce
    nonce = w3.eth.get_transaction_count(account.address)
    print(f"Starting nonce: {nonce}")

    # Deployed contracts info
    deployed_info = {
        "storage_depth": args.storage_depth,
        "account_depth": args.account_depth,
        "deployer": account.address,
        "nick_factory": NICK_FACTORY,
        "contracts": []
    }

    # Phase 1: Fund auxiliary accounts
    funded_count = 0
    total_aux_accounts = 0

    if not args.skip_funding:
        print("\n=== Phase 1: Funding Auxiliary Accounts ===")

        # Collect all unique auxiliary accounts to fund
        aux_accounts_to_fund = set()
        for contract_data in contracts:
            auxiliary_accounts = contract_data.get("auxiliary_accounts", [])
            aux_accounts_to_fund.update(auxiliary_accounts)

        total_aux_accounts = len(aux_accounts_to_fund)
        print(f"Total unique auxiliary accounts to fund: {total_aux_accounts}")

        # Calculate batch size for funding transactions
        # Simple transfers use about 21,000 gas each
        FUND_GAS_LIMIT = 21000
        latest_block = w3.eth.get_block('latest')
        network_gas_limit = latest_block['gasLimit']
        print(f"Network gas limit: {network_gas_limit:,}")

        # Target 50% of network gas limit for funding batches
        target_gas_per_batch = network_gas_limit // 2
        funding_batch_size = max(1, target_gas_per_batch // FUND_GAS_LIMIT)
        print(f"Calculated funding batch size: {funding_batch_size} accounts per batch")

        # Convert set to list for batching
        aux_accounts_list = list(aux_accounts_to_fund)

        # Process funding in batches
        for batch_num, batch_start in enumerate(range(0, len(aux_accounts_list), funding_batch_size)):
            batch = aux_accounts_list[batch_start:batch_start + funding_batch_size]
            print(f"\nFunding batch {batch_num + 1} ({len(batch)} accounts)...")

            tx_hashes = []
            batch_start_nonce = nonce

            # Send all funding transactions in the batch
            for aux_account in batch:
                try:
                    tx = {
                        'from': account.address,
                        'to': to_checksum_address(aux_account),
                        'value': 1,  # 1 wei to create the account
                        'gas': FUND_GAS_LIMIT,
                        'gasPrice': w3.eth.gas_price,
                        'nonce': nonce,
                        'chainId': w3.eth.chain_id,  # Add EIP-155 replay protection
                    }
                    signed = account.sign_transaction(tx)
                    tx_hash = w3.eth.send_raw_transaction(signed.raw_transaction)
                    tx_hashes.append(tx_hash)
                    nonce += 1
                except Exception as e:
                    print(f"  Error funding {aux_account}: {e}")
                    # Try to get more details about the error
                    if "invalid fields" in str(e):
                        print(f"    Debug: Original address: {aux_account}")
                        print(f"    Debug: Checksummed: {to_checksum_address(aux_account)}")
                        # Only show first few errors in detail
                        if batch.index(aux_account) < 3:
                            print(f"    Debug: Full tx: {tx}")

            # Wait for all transactions in batch to be mined
            if tx_hashes:
                print(f"  Waiting for {len(tx_hashes)} funding transactions to be mined...")
                for tx_hash in tx_hashes:
                    try:
                        receipt = w3.eth.wait_for_transaction_receipt(tx_hash, timeout=60)
                        if receipt.status == 1:
                            funded_count += 1
                        else:
                            print(f"    Funding transaction failed: {tx_hash.hex()}")
                    except Exception as e:
                        print(f"    Error waiting for funding tx receipt: {e}")

                print(f"  Funded {funded_count}/{total_aux_accounts} accounts so far")

        print(f"\nTotal auxiliary accounts funded: {funded_count}")
    else:
        print("\n=== Phase 1: Skipping EOA Funding (--skip-funding flag set) ===")

    # Phase 2: Deploy contracts via Nick's factory
    if not args.skip_contracts:
        print("\n=== Phase 2: Deploying Contracts ===")

        # Get the network gas limit
        latest_block = w3.eth.get_block('latest')
        network_gas_limit = latest_block['gasLimit']
        print(f"Network gas limit: {network_gas_limit:,}")

        # Calculate batch size: target 50% of network gas limit
        # Each deployment uses approximately DEPLOY_GAS_LIMIT
        target_gas_per_batch = network_gas_limit // 2  # Use 50% of gas limit
        batch_size = max(1, target_gas_per_batch // DEPLOY_GAS_LIMIT)
        print(f"Calculated batch size: {batch_size} contracts per batch (targeting 50% of gas limit)")

        for batch_num, batch in enumerate(batch_transactions(contracts, batch_size)):
            print(f"\nDeploying batch {batch_num + 1} ({len(batch)} contracts)...")

            batch_start_nonce = nonce
            deployed_in_batch = []
            tx_hashes = []

            # Send all transactions in the batch
            for contract_data in batch:
                salt = contract_data["salt"]

                # Convert salt to bytes
                if isinstance(salt, int):
                    salt_bytes = salt.to_bytes(32, 'big')
                else:
                    salt_bytes = bytes.fromhex(salt[2:] if salt.startswith("0x") else salt)

                # Deploy via Nick's factory
                try:
                    # Call Nick's factory with salt + init_code
                    factory_data = salt_bytes + init_code
                    tx: TxParams = {
                        'from': account.address,
                        'to': to_checksum_address(NICK_FACTORY),
                        'data': factory_data,
                        'gas': DEPLOY_GAS_LIMIT,
                        'gasPrice': w3.eth.gas_price,
                        'nonce': nonce,
                        'chainId': w3.eth.chain_id,  # Add EIP-155 replay protection
                    }

                    signed_tx = account.sign_transaction(tx)
                    tx_hash = w3.eth.send_raw_transaction(signed_tx.raw_transaction)
                    tx_hashes.append(tx_hash)

                    # Calculate deployed address
                    deployed_addr = calculate_create2_address(NICK_FACTORY, salt_bytes, init_code)

                    # Store deployed info
                    contract_info = {
                        "address": deployed_addr,
                        "salt": "0x" + salt_bytes.hex(),
                        "auxiliary_accounts": contract_data.get("auxiliary_accounts", [])
                    }
                    deployed_info["contracts"].append(contract_info)
                    deployed_in_batch.append(deployed_addr)

                    nonce += 1

                except Exception as e:
                    print(f"  Error sending deployment tx: {e}")
                    continue

            # Wait for all transactions in batch to be mined
            if tx_hashes:
                print(f"  Waiting for {len(tx_hashes)} transactions to be mined...")
                for tx_hash in tx_hashes:
                    try:
                        receipt = w3.eth.wait_for_transaction_receipt(tx_hash, timeout=60)
                        if receipt.status == 0:
                            print(f"    Transaction failed: {tx_hash.hex()}")
                    except Exception as e:
                        print(f"    Error waiting for tx receipt: {e}")

            # Verify deployments
            verified = 0
            for addr in deployed_in_batch:
                code = w3.eth.get_code(addr)
                if len(code) > 0:
                    verified += 1

            print(f"  Verified {verified}/{len(deployed_in_batch)} contracts deployed")

            # Progress update
            total_deployed = len(deployed_info["contracts"])
            print(f"Total progress: {total_deployed}/{len(contracts)} contracts deployed")
    else:
        print("\n=== Phase 2: Skipping Contract Deployment (--skip-contracts flag set) ===")
        # When skipping contracts, populate deployed_info with pre-calculated addresses
        for contract_data in contracts:
            salt = contract_data["salt"]
            # Convert salt to bytes
            if isinstance(salt, int):
                salt_bytes = salt.to_bytes(32, 'big')
            else:
                salt_bytes = bytes.fromhex(salt[2:] if salt.startswith("0x") else salt)

            # Calculate what the deployed address would be
            deployed_addr = calculate_create2_address(NICK_FACTORY, salt_bytes, init_code)

            contract_info = {
                "address": deployed_addr,
                "salt": "0x" + salt_bytes.hex(),
                "auxiliary_accounts": contract_data.get("auxiliary_accounts", [])
            }
            deployed_info["contracts"].append(contract_info)

    # Save deployed contracts info
    print(f"\n=== Saving Deployment Info ===")
    output_path = Path(args.output)
    with open(output_path, 'w') as f:
        json.dump(deployed_info, f, indent=2)

    print(f"Deployment info saved to: {output_path}")

    # Also save just the contract addresses to stubs.json
    stubs_path = Path("stubs.json")
    contract_addresses = [contract["address"] for contract in deployed_info["contracts"]]
    with open(stubs_path, 'w') as f:
        json.dump(contract_addresses, f, indent=2)

    print(f"Contract addresses saved to: {stubs_path}")
    print(f"\nOperation complete!")
    if not args.skip_contracts:
        print(f"  Total contracts deployed: {len(deployed_info['contracts'])}")
    else:
        print(f"  Total contract addresses calculated: {len(deployed_info['contracts'])}")
    if not args.skip_funding:
        print(f"  Total auxiliary accounts funded: {funded_count}")
    print(f"  Final nonce: {nonce}")

if __name__ == "__main__":
    main()