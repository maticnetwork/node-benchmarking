package scripts

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/maticnetwork/bor/ethclient"
	"github.com/maticnetwork/bor/rpc"
)

// SignerCount prints the number of blocks produced by each block producer
func SignerCount() {
	RPCClient, err := rpc.Dial("http://localhost:8545")
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(RPCClient)

	latestBlock, err := client.BlockByNumber(context.Background(), nil /* latest */)
	if err != nil {
		panic(err)
	}

	signerCount := make(map[string]int)
	for i := latestBlock.Header().Number; i.Int64() > 0; i.Sub(i, big.NewInt(1)) {
		block, err := client.BlockByNumber(context.Background(), i)
		if err != nil {
			log.Fatal(err)
		}

		signer := block.Coinbase().Hex()
		signerCount[signer]++
	}
	for key, element := range signerCount {
		fmt.Println(key, element)
	}
}
