package benchmarking

import (
	"context"
	"fmt"
	"log"

	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/bor/ethclient"
	"github.com/maticnetwork/bor/rpc"
)

func SubscribeBlocks() {
	RPCClient, err := rpc.Dial("ws://127.0.0.1:8585") // websocket port of a node started from bor-devnet directory
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(RPCClient)
	ctx := context.Background()

	ch := make(chan *types.Header)
	subscription, err := client.SubscribeNewHead(ctx, ch)
	if err != nil {
		panic(err)
	}

	var totalTxs uint

	for {
		select {
		case header := <-ch:
			count := printBlockTxCount(ctx, client, header)
			totalTxs += count
			fmt.Printf("Total: %d\n", totalTxs)
		case err := <-subscription.Err():
			log.Fatal(err)
			return
		case <-ctx.Done():
			return
		}
	}
}

func printBlockTxCount(ctx context.Context, client *ethclient.Client, header *types.Header) uint {
	count, err := client.TransactionCount(ctx, header.Hash())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%d txs in block %d; ", count, header.Number.Uint64())
	return count
}
