package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dannwee/dbc-go/instructions"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetPool() {
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")

	poolAddressStr := "YOUR_POOL_ADDRESS"

	fmt.Println("Getting pool...")
	poolAddress := solana.MustPublicKeyFromBase58(poolAddressStr)

	ctx := context.Background()

	pool, err := instructions.GetPool(ctx, poolAddress, rpcClient)
	if err != nil {
		log.Fatalf("Failed to get pool: %v", err)
	}

	// Marshal the pool to JSON
	jsonData, err := json.MarshalIndent(pool, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal pool to JSON: %v", err)
	}

	fmt.Printf("Pool JSON: %s\n", string(jsonData))

	// Print some key metrics
	fmt.Printf("Base Reserve: %d\n", pool.BaseReserve)
	fmt.Printf("Quote Reserve: %d\n", pool.QuoteReserve)
	fmt.Printf("SqrtPrice: %s\n", pool.SqrtPrice.String())
	fmt.Printf("Total Trading Base Fee: %d\n", pool.Metrics.TotalTradingBaseFee)
	fmt.Printf("Total Trading Quote Fee: %d\n", pool.Metrics.TotalTradingQuoteFee)
}

// func main() {
// 	GetPool()
// }
