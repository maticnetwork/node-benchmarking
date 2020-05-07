package scripts

import (
	"context"
	"fmt"

	"github.com/maticnetwork/bor/ethclient"
	"github.com/maticnetwork/bor/rpc"
)

// RootHash returns the merkle root of start to end block headers
func RootHash(start uint64, end uint64) {
	RPCClient, err := rpc.Dial("ws://127.0.0.1:8585") // websocket port of a node started from bor-devnet directory
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(RPCClient)
	ctx, _ := context.WithCancel(context.Background())

	root, err := client.GetRootHash(ctx, start, end)
	if err != nil {
		panic(err)
	}
	fmt.Println(root)
}
