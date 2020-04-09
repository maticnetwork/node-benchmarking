package scripts

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/bor/crypto"
	"github.com/maticnetwork/bor/ethclient"
)

// Stress fires txs in bulk
// NEEDS PRIVATE_KEY env variable - without 0x prefix
func RapidFire() {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		panic(err)
	}

	key, _ := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	fromAddress := crypto.PubkeyToAddress(key.PublicKey)
	toAddress := common.HexToAddress("0x9fB29AAc15b9A4B7F17c3385939b007540f4d791") // any address is fine

	chainID, _ := client.NetworkID(context.Background())
	nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
	startBlock, _ := client.BlockByNumber(context.Background(), nil)
	fmt.Println("startBlock", startBlock.Number())

	const numTxs = 5000
	fmt.Printf("Firing %d txs\n", numTxs)
	for i := nonce; i < nonce+numTxs; i++ {
		tx := types.NewTransaction(
			i,
			toAddress,
			big.NewInt(10), // in wei,
			uint64(21000),  // gasLimit
			new(big.Int),
			nil)
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
		if err != nil {
			log.Fatal(err)
		}
		txErr := client.SendTransaction(context.Background(), signedTx)
		if txErr != nil {
			panic(txErr)
		}
		// fmt.Printf("%d send success tx.hash=%s\n", i, signedTx.Hash().String())
	}

	pendingTransactions, _ := client.PendingTransactionCount(context.Background())
	for pendingTransactions > 0 {
		fmt.Printf("%d pending txs; Sleeping for 3 seconds...\n", pendingTransactions)
		time.Sleep(3 * time.Second)
		pendingTransactions, _ = client.PendingTransactionCount(context.Background())
	}

	// All transactions confirmed
	latestBlock, _ := client.BlockByNumber(context.Background(), nil /* latest */)
	fmt.Println("latestBlock", latestBlock.Number())

	var txCount uint
	var gasUsed uint64
	var blockCount uint
	for i := latestBlock.Header().Number; i.Int64() > startBlock.Header().Number.Int64(); i.Sub(i, big.NewInt(1)) {
		block, _ := client.BlockByNumber(context.Background(), i)
		count, _ := client.TransactionCount(context.Background(), block.Hash())
		fmt.Printf("%d txs and %d gas used in block: %d\n", count, block.GasUsed(), i)
		if count == 0 {
			continue
		}
		txCount += count
		gasUsed += block.GasUsed()
		blockCount++
	}
	fmt.Printf("Average tx count: %d\n", txCount/blockCount)
	fmt.Printf("Average gas used: %d\n", gasUsed/uint64(blockCount))
}
