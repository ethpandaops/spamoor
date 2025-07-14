# Spamoor - Ethereum Transaction Spammer

Spamoor is a powerful tool for generating various types of random transactions on Ethereum testnets. It provides a modular, scenario-based architecture for stress testing, network validation, and continuous transaction testing.

## Project Structure
Claude MUST read the `.cursor/rules/project_architecture.mdc` file before making any structural changes to the project.

## Code Standards  
Claude MUST read the `.cursor/rules/code_standards.mdc` file before writing any code in this project.

## Development Workflow
Claude MUST read the `.cursor/rules/development_workflow.mdc` file before making changes to build, test, or deployment configurations.

## Component Documentation
Individual components have their own CLAUDE.md files with component-specific rules. Always check for and read component-level documentation when working on specific parts of the codebase.

## User Documentation
Comprehensive documentation is available in the `docs/` directory for different user types:

- **`docs/app-users.md`**: Complete guide for CLI usage, daemon mode, configuration, and troubleshooting
- **`docs/api-consumers.md`**: REST API documentation with bash/curl examples for all endpoints  
- **`docs/scenario-developers.md`**: Comprehensive guide for implementing custom transaction scenarios

**Before working on any scenario development or API integration, Claude MUST read the relevant documentation files to understand:**
- Critical development rules (no root wallet usage, proper nonce management, etc.)
- Wallet management patterns and selection strategies
- Transaction building and submission methods
- Contract interaction patterns using abigen with `BuildBoundTx`
- Balance management for internal transfers
- Best practices for context handling, logging, and configuration

## Key Features
- **Multiple Scenarios**: 12+ different transaction types including EOA, ERC20, blob transactions, contract deployments, and more
- **Daemon Mode**: Web-based interface for managing multiple concurrent spammers  
- **Modular Architecture**: Easy to extend with new scenario types
- **Docker Support**: Containerized deployment and development
- **DevNet Integration**: Quick development environment with `make devnet-run`
- **Metrics Collection**: Prometheus integration for monitoring
- **Flexible Configuration**: YAML files and command-line flags

## Build Requirements
- Go 1.24+
- CGO enabled (for SQLite and cryptographic operations)
- Build tags: `with_blob_v1,ckzg` for blob transaction support

## Testing & Development
- Run `make test` for all tests
- Use `make devnet-run` for full development environment with Ethereum testnet
- Individual scenarios can be tested independently
- Web interface available at http://localhost:8080 in daemon mode

## Transaction Submission Constraints
When working with transaction submission:
- **SendTransactionBatch**: All transactions must originate from the same sender wallet
- **SendMultiTransactionBatch**: Efficiently handles transactions from multiple wallets with advanced concurrency control
- **Multi-wallet batching**: Use `SendMultiTransactionBatch` for optimal performance across multiple wallets
- **Rebroadcast mechanism**: Automatic exponential backoff for stuck transactions
- **Balance management**: Automatic for ETH transfers, manual updates required for internal transfers

## Critical Development Rules
🚨 **NEVER use the root wallet directly in scenarios** - causes nonce conflicts
🚨 **NEVER use go-ethereum's bound contracts for transactions** - use `BuildBoundTx` pattern  
🚨 **ALWAYS spread transactions across multiple wallets** - respect client pending limits
🚨 **ALWAYS call onComplete() in ProcessNextTxFn** - required for scenario.RunTransactionScenario transaction counting
🚨 **ALWAYS respect context cancellation** - scenarios must stop when context is cancelled

## Code Quality Requirements
**Claude MUST automatically run these checks after ANY Go code changes:**
1. `go fmt ./...` - Format all Go code
2. `go vet ./...` - Check for Go code issues  
3. `staticcheck ./...` - Run static analysis

These checks match the GitHub Actions CI pipeline and must pass before any code is considered complete. Do NOT ask the user to run these - they should be executed automatically after every code modification.