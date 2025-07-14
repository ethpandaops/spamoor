package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethpandaops/spamoor/scenarios/storagespam/contract"
)

type EventSummary struct {
	GasLimit *big.Int
	Loops    *big.Int
	BlockNum uint64
	TxHash   common.Hash
}

type GasSummary struct {
	GasLimit   *big.Int
	TotalLoops *big.Int
	Count      int
	Events     []EventSummary
}

func main() {
	var rpcHost string
	var contractAddr string
	var batchSize int

	flag.StringVar(&rpcHost, "rpc", "", "RPC host URL")
	flag.StringVar(&contractAddr, "contract", "", "Contract address")
	flag.IntVar(&batchSize, "batch", 100, "Batch size for block queries")
	flag.Parse()

	if rpcHost == "" || contractAddr == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -rpc <RPC_URL> -contract <CONTRACT_ADDRESS> [-batch <SIZE>]\n", os.Args[0])
		os.Exit(1)
	}

	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, rpcHost)
	if err != nil {
		log.Fatalf("Failed to connect to RPC: %v", err)
	}
	defer client.Close()

	address := common.HexToAddress(contractAddr)
	fmt.Printf("Connected to RPC: %s\n", rpcHost)
	fmt.Printf("Analyzing contract: %s\n", address.Hex())
	fmt.Printf("Batch size: %d blocks\n\n", batchSize)

	// Get current block
	latestBlock, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to get latest block: %v", err)
	}

	currentBlock := latestBlock.Number().Uint64()
	fmt.Printf("Starting from block: %d\n", currentBlock)

	// Get contract ABI
	contractABI, err := abi.JSON(strings.NewReader(contract.StorageSpamMetaData.ABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Topic for RandomForGas event
	randomForGasTopic := contractABI.Events["RandomForGas"].ID

	allEvents := make([]EventSummary, 0)
	emptyBatchCount := 0
	maxEmptyBatches := 10 // Stop after 10 consecutive batches with no logs

	// Walk backwards from head using log filtering
	for currentBlock > 0 && emptyBatchCount < maxEmptyBatches {
		fromBlock := currentBlock
		if fromBlock > uint64(batchSize) {
			fromBlock = currentBlock - uint64(batchSize) + 1
		} else {
			fromBlock = 0
		}

		// Query logs for this batch using FilterLogs
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlock)),
			ToBlock:   big.NewInt(int64(currentBlock)),
			Addresses: []common.Address{address},
			Topics:    [][]common.Hash{{randomForGasTopic}},
		}

		logs, err := client.FilterLogs(ctx, query)
		if err != nil {
			log.Printf("Error getting logs for blocks %d-%d: %v", fromBlock, currentBlock, err)
			currentBlock = fromBlock - 1
			continue
		}

		batchEventCount := len(logs)
		if batchEventCount == 0 {
			emptyBatchCount++
		} else {
			emptyBatchCount = 0 // Reset counter when we find events
		}

		// Process each log
		for _, vLog := range logs {
			event := &contract.StorageSpamRandomForGas{}
			err := contractABI.UnpackIntoInterface(event, "RandomForGas", vLog.Data)
			if err != nil {
				log.Printf("Failed to unpack event: %v", err)
				continue
			}

			// Get indexed gas parameter
			if len(vLog.Topics) < 2 {
				log.Printf("Invalid event topics")
				continue
			}
			gasLimit := new(big.Int).SetBytes(vLog.Topics[1].Bytes())

			allEvents = append(allEvents, EventSummary{
				GasLimit: gasLimit,
				Loops:    event.Loops,
				BlockNum: vLog.BlockNumber,
				TxHash:   vLog.TxHash,
			})
		}

		if batchEventCount > 0 {
			fmt.Printf("Blocks %d-%d: %d events\n", fromBlock, currentBlock, batchEventCount)
		}

		if fromBlock == 0 {
			break
		}
		currentBlock = fromBlock - 1
	}

	fmt.Printf("\n=== EVENT SUMMARY ===\n")
	fmt.Printf("Total RandomForGas events found: %d\n", len(allEvents))

	if len(allEvents) == 0 {
		fmt.Println("\nNo RandomForGas events found")
		return
	}

	// Sort events by block number
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].BlockNum < allEvents[j].BlockNum
	})

	// Summarize by gas limit
	gasSummaries := make(map[string]*GasSummary)
	for _, event := range allEvents {
		key := event.GasLimit.String()
		if summary, exists := gasSummaries[key]; exists {
			summary.TotalLoops.Add(summary.TotalLoops, event.Loops)
			summary.Count++
			summary.Events = append(summary.Events, event)
		} else {
			gasSummaries[key] = &GasSummary{
				GasLimit:   new(big.Int).Set(event.GasLimit),
				TotalLoops: new(big.Int).Set(event.Loops),
				Count:      1,
				Events:     []EventSummary{event},
			}
		}
	}

	// Print summary
	fmt.Printf("\n=== EVENTS SUMMARY ===\n")
	fmt.Printf("Total events: %d\n", len(allEvents))
	fmt.Printf("Block range: %d - %d\n", allEvents[0].BlockNum, allEvents[len(allEvents)-1].BlockNum)

	fmt.Printf("\n=== GAS LIMIT SUMMARY ===\n")
	fmt.Printf("%-20s | %-10s | %-20s | %-20s\n", "Gas Limit", "Count", "Total Loops", "Avg Loops")
	fmt.Printf("%-20s-+-%-10s-+-%-20s-+-%-20s\n", "--------------------", "----------", "--------------------", "--------------------")

	// Sort summaries by gas limit
	var sortedGasLimits []*big.Int
	for _, summary := range gasSummaries {
		sortedGasLimits = append(sortedGasLimits, summary.GasLimit)
	}
	sort.Slice(sortedGasLimits, func(i, j int) bool {
		return sortedGasLimits[i].Cmp(sortedGasLimits[j]) < 0
	})

	totalGas := new(big.Int)
	totalLoops := new(big.Int)
	for _, gasLimit := range sortedGasLimits {
		summary := gasSummaries[gasLimit.String()]
		avgLoops := new(big.Int).Div(summary.TotalLoops, big.NewInt(int64(summary.Count)))

		fmt.Printf("%-20s | %-10d | %-20s | %-20s\n",
			gasLimit.String(),
			summary.Count,
			summary.TotalLoops.String(),
			avgLoops.String())

		// Calculate total gas burned
		gasForThisLimit := new(big.Int).Mul(gasLimit, big.NewInt(int64(summary.Count)))
		totalGas.Add(totalGas, gasForThisLimit)
		totalLoops.Add(totalLoops, summary.TotalLoops)
	}

	fmt.Printf("\n=== TOTALS ===\n")
	fmt.Printf("Total gas burned: %s\n", totalGas.String())
	fmt.Printf("Total loops executed: %s\n", totalLoops.String())

	// Storage calculations - each loop creates one storage slot (32 byte key + 32 byte value)
	storagePerLoop := int64(64) // 32 bytes key + 32 bytes value
	totalStorageBytes := new(big.Int).Mul(totalLoops, big.NewInt(storagePerLoop))
	totalStorageKB := new(big.Float).SetInt(totalStorageBytes)
	totalStorageKB.Quo(totalStorageKB, big.NewFloat(1024))
	totalStorageMB := new(big.Float).Quo(totalStorageKB, big.NewFloat(1024))
	totalStorageGB := new(big.Float).Quo(totalStorageMB, big.NewFloat(1024))

	fmt.Printf("\n=== STORAGE USAGE ===\n")
	fmt.Printf("Storage per loop: 64 bytes (32 byte key + 32 byte value)\n")
	fmt.Printf("Total storage slots created: %s\n", totalLoops.String())
	fmt.Printf("Total storage used: %s bytes\n", totalStorageBytes.String())
	kbStr, _ := totalStorageKB.Float64()
	mbStr, _ := totalStorageMB.Float64()
	gbStr, _ := totalStorageGB.Float64()
	fmt.Printf("Total storage used: %.2f KB / %.2f MB / %.2f GB\n", kbStr, mbStr, gbStr)

	if len(allEvents) > 0 {
		avgGasPerEvent := new(big.Int).Div(totalGas, big.NewInt(int64(len(allEvents))))
		avgLoopsPerEvent := new(big.Int).Div(totalLoops, big.NewInt(int64(len(allEvents))))
		avgStoragePerEvent := new(big.Int).Mul(avgLoopsPerEvent, big.NewInt(storagePerLoop))
		fmt.Printf("\n=== AVERAGES PER EVENT ===\n")
		fmt.Printf("Average gas per event: %s\n", avgGasPerEvent.String())
		fmt.Printf("Average loops per event: %s\n", avgLoopsPerEvent.String())
		fmt.Printf("Average storage per event: %s bytes\n", avgStoragePerEvent.String())

		// Gas efficiency
		if totalLoops.Sign() > 0 {
			gasPerLoop := new(big.Float).SetInt(totalGas)
			gasPerLoop.Quo(gasPerLoop, new(big.Float).SetInt(totalLoops))
			gasPerLoopVal, _ := gasPerLoop.Float64()
			fmt.Printf("\n=== EFFICIENCY ===\n")
			fmt.Printf("Average gas per storage slot: %.2f\n", gasPerLoopVal)
			fmt.Printf("Average gas per KB of storage: %.2f\n", gasPerLoopVal*1024/64)
		}
	}

	// Optional: Print detailed events
	var printDetails string
	fmt.Printf("\nPrint detailed event list? (y/N): ")
	fmt.Scanln(&printDetails)

	if printDetails == "y" || printDetails == "Y" {
		fmt.Printf("\n=== DETAILED EVENTS ===\n")
		fmt.Printf("%-10s | %-20s | %-20s | %-66s\n", "Block", "Gas Limit", "Loops", "Tx Hash")
		fmt.Printf("%-10s-+-%-20s-+-%-20s-+-%-66s\n", "----------", "--------------------", "--------------------", "------------------------------------------------------------------")

		for _, event := range allEvents {
			fmt.Printf("%-10d | %-20s | %-20s | %-66s\n",
				event.BlockNum,
				event.GasLimit.String(),
				event.Loops.String(),
				event.TxHash.Hex())
		}
	}
}
