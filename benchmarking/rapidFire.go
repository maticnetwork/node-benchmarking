package benchmarking

import (
	"context"
	"fmt"
	"log"
	"math/big"
	// "os"
	"strconv"
	"time"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/bor/crypto"
	"github.com/maticnetwork/bor/ethclient"
)

// RapidFire fires txs in bulk
// NEEDS PRIVATE_KEY env variable - without 0x prefix
func RapidFire() {
	numTxs := 10000 // Number of txs to fire
	numClients := 4 // Number of nodes that will be hit with txs randomly
	// strategy := "regular"
	strategy := "reverse"

	clients := make(map[int]*ethclient.Client)
	for i := 0; i < numClients; i++ {
		// client, err := ethclient.Dial("http://localhost:" + strconv.Itoa(8545 + i))
		// client, err := ethclient.Dial("ws://localhost:" + strconv.Itoa(8585 + i))
		client, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(8545 + i))
		if err != nil {
			panic(err)
		}
		clients[i] = client
	}

	// key, _ := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY")) // PRIVATE_KEY env var should not be 0x prefixed
	key, _ := crypto.HexToECDSA(randomHex()) // PRIVATE_KEY env var should not be 0x prefixed
	fromAddress := crypto.PubkeyToAddress(key.PublicKey)
	toAddress := common.HexToAddress("0x9fB29AAc15b9A4B7F17c3385939b007540f4d791") // any address is fine

	client := clients[0]
	chainID, _ := client.NetworkID(context.Background())
	nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
	startBlock, _ := client.BlockByNumber(context.Background(), nil)
	fmt.Println("startBlock", startBlock.Number())

	type Data struct {
		clientNumber int
		tx *types.Transaction
	}

	txs := make([]*Data, numTxs)
	for i := 0; i < numTxs; i++ {
		tx := types.NewTransaction(
			uint64(i) + nonce,
			toAddress,
			big.NewInt(0), // value in wei,
			uint64(21000),  // gasLimit
			new(big.Int),
			nil)
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
		if err != nil {
			log.Fatal(err)
		}
		// txs[i] = &Data{rand.Intn(numClients), signedTx}
		// fmt.Println(i, i / (numTxs / numClients))
		txs[i] = &Data{i / (numTxs / numClients), signedTx}
	}
	fmt.Printf("constructed %d txs, now firing them...\n", len(txs))

	if strategy == "reverse" {
		for i := numTxs-1; i >= 0; i-- {
			data := txs[i]
			txErr := clients[data.clientNumber].SendTransaction(context.Background(), data.tx)
			if txErr != nil {
				panic(txErr)
			}
		}
	} else {
		for i := 0; i < numTxs; i++ {
			data := txs[i]
			txErr := clients[data.clientNumber].SendTransaction(context.Background(), data.tx)
			if txErr != nil {
				panic(txErr)
			}
		}
	}

	_nonce, _ := client.NonceAt(context.Background(), fromAddress, nil)
	pendingTransactions := nonce + uint64(numTxs) - _nonce
	for pendingTransactions > 0 {
		fmt.Printf("%d pending txs; Sleeping for 3 seconds...\n", pendingTransactions)
		time.Sleep(3 * time.Second)
		_nonce, _ = client.NonceAt(context.Background(), fromAddress, nil)
		pendingTransactions = nonce + uint64(numTxs) - _nonce
	}
	// time.Sleep(5 * time.Second)

	// All transactions confirmed
	latestBlock, _ := client.BlockByNumber(context.Background(), nil /* latest */)
	fmt.Println("latestBlock", latestBlock.Number())

	var txCount uint
	var gasUsed uint64
	var blockCount uint
	for i := latestBlock.Header().Number.Int64(); i > startBlock.Header().Number.Int64(); i-- {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			panic(err)
		}

		// for block == nil {
		// }
		// fmt.Println(block)
		count, _ := client.TransactionCount(context.Background(), block.Hash())
		fmt.Printf("%d txs and %d gas used in block: %d\n", count, block.GasUsed(), i)
		if count == 0 {
			continue
		}
		txCount += count
		gasUsed += block.GasUsed()
		blockCount++
	}
	fmt.Printf("Total tx count: %d\n", txCount)
	fmt.Printf("Average tx count: %d\n", txCount/blockCount)
	fmt.Printf("Average gas used: %d\n", gasUsed/uint64(blockCount))
}
