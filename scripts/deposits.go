package scripts

import (
	"context"
	"fmt"
	"log"

	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/bor/ethclient"
	"github.com/maticnetwork/bor/rpc"
)

// Deposits subscribes to StateSync events on Bor
func Deposits() {
	RPCClient, err := rpc.Dial("ws://127.0.0.1:8585") // websocket port of a node started from bor-devnet directory
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(RPCClient)
	ctx, _ := context.WithCancel(context.Background())

	depositsChannel := make(chan *types.StateData)
	subscription, err := client.SubscribeNewDeposit(ctx, depositsChannel)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case stateData := <-depositsChannel:
			printStateData(stateData)
		case err := <-subscription.Err():
			log.Fatal(err)
			return
		case <-ctx.Done():
			return
		}
	}
}

func printStateData(s *types.StateData) {
	fmt.Println(s.Did, s.Contract.Hex(), s.Data, s.TxHash.Hex())
}
