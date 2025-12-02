# AI Transaction Generator (aitx)

The AI Transaction Generator scenario leverages OpenRouter's API to generate diverse, AI-powered Ethereum transaction payloads for comprehensive stress testing and validation.

## Overview

This scenario uses AI models to create dynamic transaction patterns that go beyond static predefined scenarios. It supports multiple generation modes and includes a feedback loop to improve payload diversity over time.

## Features

- **AI-Powered Generation**: Uses OpenRouter API with configurable AI models (default: Claude 3.5 Sonnet)
- **Multiple Generation Modes**: Support for geas assembly, calldata, transfers, or mixed mode
- **Feedback Loop**: AI learns from transaction execution results to generate better payloads
- **Batch Processing**: Generates multiple payloads per API call to minimize costs
- **Placeholder System**: Dynamic parameter substitution for transaction variations
- **Safety Validation**: Prevents malicious code generation with built-in safety checks
- **Cost Management**: Configurable limits on API calls and token consumption

## Configuration

### Required Configuration

```yaml
openrouter_api_key: "your-api-key-here"  # OpenRouter API key (or set OPENROUTER_API_KEY env var)
```

### AI Configuration Options

```yaml
# AI Model Settings
model: "anthropic/claude-3.5-sonnet"      # AI model to use
test_direction: "focus on gas optimization" # Directional guidance for AI
payloads_per_request: 50                  # Payloads generated per API call
max_ai_calls: 10                          # Maximum API calls limit
max_tokens: 100000                        # Maximum token consumption limit

# Generation Settings
generation_mode: "mixed"                  # "geas", "calldata", "transfer", "mixed"

# Feedback Settings
feedback_batch_size: 20                   # Transaction results included in feedback
enable_feedback_loop: true                # Enable AI learning from results

# Debug Settings
log_ai_conversations: false               # Enable detailed AI conversation logging
```

### Standard Transaction Options

```yaml
total_count: 1000                         # Total transactions to send (0 = unlimited)
throughput: 10                            # Transactions per slot
max_pending: 100                          # Maximum pending transactions
max_wallets: 50                           # Maximum child wallets
rebroadcast: 1                            # Enable transaction rebroadcast
basefee: 20                               # Base fee in gwei
tipfee: 2                                 # Tip fee in gwei
gaslimit: 1000000                         # Gas limit for transactions
timeout: "30m"                            # Scenario timeout
client_group: ""                          # Client group preference
log_txs: false                            # Log individual transactions
```

## Usage

### Basic Usage

```bash
./spamoor aitx --openrouter-api-key="your-key" --count=100 --throughput=5
```

### Advanced Configuration

```bash
./spamoor aitx \
  --openrouter-api-key="your-key" \
  --model="anthropic/claude-3.5-sonnet" \
  --generation-mode="mixed" \
  --test-direction="focus on complex contract interactions" \
  --payloads-per-request=30 \
  --max-ai-calls=5 \
  --enable-feedback-loop=true \
  --count=500 \
  --throughput=10 \
  --gaslimit=2000000
```

### YAML Configuration

```yaml
# aitx-config.yaml
scenarios:
  aitx:
    # AI Configuration
    openrouter_api_key: "your-api-key-here"
    model: "anthropic/claude-3.5-sonnet"
    test_direction: "explore edge cases and gas optimization patterns"
    generation_mode: "mixed"
    payloads_per_request: 40
    max_ai_calls: 8
    max_tokens: 80000
    
    # Feedback Configuration
    feedback_batch_size: 25
    enable_feedback_loop: true
    
    # Debug Configuration
    log_ai_conversations: true
    
    # Transaction Configuration
    total_count: 2000
    throughput: 15
    max_pending: 150
    basefee: 25
    tipfee: 3
    gaslimit: 1500000
    rebroadcast: 1
    log_txs: true
```

## Generation Modes

### Geas Mode (`generation_mode: "geas"`)
**PREFERRED MODE** - Provides the most interesting and comprehensive EVM testing.

Two geas deployment methods are supported:

#### Simple Method (for arbitrary opcode/precompile testing)
- Generates standalone geas assembly code
- Code gets deployed as a contract and can be called with specific calldata
- Best for: testing specific opcodes, precompiles, edge cases
- Example: Testing modular exponentiation precompile, specific memory operations

#### Init/Run Method (for performance testing)
- Generates separate initialization code and run code
- Init code executes once during contract deployment
- Run code executes in a loop until gas is consumed (like gasburnertx)
- Best for: performance benchmarks, gas burning, stress testing
- Example: Testing complex loops, storage operations, cryptographic operations

Both methods automatically validate for security (blocks dangerous operations like selfdestruct, delegatecall, create2)

### Calldata Mode (`generation_mode: "calldata"`)
- Generates raw calldata for contract interactions
- Creates function calls with ABI-encoded parameters
- Supports common patterns like transfers, approvals, and complex calls
- Includes both simple and complex contract interaction patterns

### Transfer Mode (`generation_mode: "transfer"`)
- Generates simple ETH transfers between addresses
- Focuses on value transfer patterns and address variations
- Uses placeholder system for dynamic amounts and recipients

### Mixed Mode (`generation_mode: "mixed"`)
- Combines all generation modes in a single scenario
- **Prioritizes geas generation** (70-80% of payloads will be geas-based)
- Provides maximum diversity in transaction patterns
- AI automatically balances between geas methods and other transaction types
- Recommended for comprehensive testing scenarios

## Placeholder System

The AI can use placeholders that are dynamically substituted during transaction building:

- `${WALLET_ADDRESS}`: Random wallet from pool
- `${RANDOM_ADDRESS}`: Randomly generated address
- `${ETH_AMOUNT_SMALL/MEDIUM/LARGE}`: Dynamic ETH amounts
- `${GAS_LIMIT_LOW/MEDIUM/HIGH}`: Dynamic gas limits
- `${RANDOM_UINT256}`: Random 256-bit integer
- `${RANDOM_BYTES32}`: Random 32-byte value
- `${LOOP_COUNT_SMALL/MEDIUM/LARGE}`: Loop iteration counts

## Feedback Loop

When enabled, the scenario:
1. Collects transaction execution results (success/failure, gas usage, errors)
2. Provides statistical analysis to the AI for subsequent generations
3. Encourages the AI to avoid failing patterns and explore successful variations
4. Builds context across multiple AI calls for progressive improvement

## Cost Management

- **API Call Limits**: Prevents runaway costs with `max_ai_calls`
- **Token Limits**: Controls total token consumption with `max_tokens`
- **Batch Processing**: Generates multiple payloads per API call
- **Efficient Caching**: Reuses generated payloads until cache is exhausted

## Debugging

### AI Conversation Logging

Enable detailed logging of AI conversations with `--log-ai-conversations=true` or in YAML:

```yaml
log_ai_conversations: true
```

This will log:
- Full conversation history for each AI request
- All retry attempts with error feedback
- AI responses before and after parsing
- Truncated content for readability (messages over 2000 chars)

Use with debug logging level for maximum detail:
```bash
./spamoor aitx --log-ai-conversations=true --log-level=debug
```

### Geas Compilation Validation

The scenario automatically validates all geas code during the AI conversation:

**Real-time Compilation**: Every geas payload is compiled using the geas compiler before being accepted
**Immediate Feedback**: Compilation errors are sent back to the AI in the same conversation for immediate correction
**Comprehensive Error Guidance**: AI receives specific guidance on:
- Valid EVM opcodes and syntax
- Proper formatting requirements (newlines, hex format)
- Common compilation issues and fixes
- Stack management requirements

**Example Error Feedback to AI**:
```
GEAS COMPILATION ERROR DETECTED:
Error: geas compilation failed: unknown opcode 'keccak256'

Your geas assembly code failed to compile. Please fix the following issues:

GEAS CODE REQUIREMENTS:
1. Use VALID EVM opcodes only (e.g., push1, add, mul, sstore, sload, etc.)
2. Format: ONE opcode per line, separated by \n
3. Use correct syntax: 'push1 0x20' not 'push1(0x20)' or 'PUSH1 0x20'
4. Use 'sha3' instead of 'keccak256' for the EVM opcode
...
```

### Common Debug Scenarios

**Parsing Errors**: When AI generates invalid JSON, conversation logs show exactly what was sent and received
**Geas Compilation Errors**: Invalid assembly code is caught immediately with specific error feedback to guide AI fixes
**Retry Logic**: Track how error feedback is provided and how AI responds to corrections
**Token Usage**: Monitor conversation length and token consumption patterns
**Response Quality**: Analyze AI output quality and prompt effectiveness

## Security

- **Code Validation**: Blocks dangerous operations (selfdestruct, delegatecall, create2)
- **Geas Compilation**: Validates assembly code before execution
- **Payload Validation**: Ensures all generated payloads meet safety requirements
- **API Key Security**: Supports environment variable configuration

## Troubleshooting

### Common Issues

#### API Key Not Found
```
Error: OpenRouter API key is required
```
**Solution**: Set the API key using `--openrouter-api-key` flag or `OPENROUTER_API_KEY` environment variable.

#### AI Generation Failures
```
Error: AI payload generation failed: HTTP request failed
```
**Solutions**: 
- Check internet connectivity
- Verify API key validity
- Ensure OpenRouter service is available
- Try reducing `payloads_per_request` if hitting rate limits

#### Invalid Generation Mode
```
Error: invalid generation mode 'invalid', must be one of: geas, calldata, transfer, mixed
```
**Solution**: Use a valid generation mode: `geas`, `calldata`, `transfer`, or `mixed`.

#### Geas Compilation Errors
```
Error: geas compilation failed: syntax error
```
**Solutions**:
- AI-generated geas code may be invalid
- Check the AI model configuration
- Try different `test_direction` guidance
- Verify the AI is generating valid assembly syntax

### Performance Optimization

- **Throughput**: Start with low throughput (5-10) and increase gradually
- **Wallet Count**: Ensure adequate wallets for transaction distribution
- **Gas Limits**: Adjust based on payload complexity
- **Batch Size**: Increase `payloads_per_request` for better cost efficiency
- **Feedback**: Enable feedback loop for improved AI generation over time

### Monitoring

- Monitor token consumption to stay within budget
- Track API call usage against limits
- Watch transaction success rates for payload quality
- Review gas usage patterns for optimization opportunities

## Examples

### Gas Optimization Focus
```bash
./spamoor aitx \
  --openrouter-api-key="your-key" \
  --test-direction="create gas-efficient transaction patterns" \
  --generation-mode="mixed" \
  --enable-feedback-loop=true \
  --count=1000
```

### Complex Contract Testing
```bash
./spamoor aitx \
  --openrouter-api-key="your-key" \
  --test-direction="generate complex contract interactions with edge cases" \
  --generation-mode="calldata" \
  --gaslimit=3000000 \
  --count=500
```

### Assembly Code Generation
```bash
./spamoor aitx \
  --openrouter-api-key="your-key" \
  --test-direction="create diverse EVM assembly patterns for stress testing" \
  --generation-mode="geas" \
  --gaslimit=5000000 \
  --payloads-per-request=20
```

### Performance Testing with Init/Run Geas
```bash
./spamoor aitx \
  --openrouter-api-key="your-key" \
  --test-direction="generate performance benchmarks using init/run pattern" \
  --generation-mode="geas" \
  --gaslimit=8000000 \
  --count=200
```

### Opcode Testing with Simple Geas
```bash
./spamoor aitx \
  --openrouter-api-key="your-key" \
  --test-direction="test specific opcodes and precompiles using simple deployment" \
  --generation-mode="geas" \
  --gaslimit=3000000 \
  --payloads-per-request=15
```