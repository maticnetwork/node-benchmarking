package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	RPCClient, err := rpc.Dial("http://localhost:8545")
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(RPCClient)
	fmt.Println("client generated", client)
}
