package benchmarking

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"encoding/hex"
	"sync"

	"strconv"
	"time"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/bor/crypto"
	"github.com/maticnetwork/bor/ethclient"
)

const maxFileDescriptors = 50

// RapidFire2 fires txs in bulk
func RapidFire2() {
	numTxs := 100000 // Number of txs to fire
	numClients := 4 // Number of nodes that will be hit with txs randomly

	clients := make(map[int]*ethclient.Client)
	for i := 0; i < numClients; i++ {
		// client, err := ethclient.Dial("http://localhost:" + strconv.Itoa(8545 + i))
		client, err := ethclient.Dial("ws://localhost:" + strconv.Itoa(8585 + i))
		if err != nil {
			panic(err)
		}
		clients[i] = client
	}

	toAddress := common.HexToAddress("0x9fB29AAc15b9A4B7F17c3385939b007540f4d791") // any address is fine

	ctx := context.Background()
	client, _ := ethclient.Dial("http://localhost:8546")
	chainID, _ := client.NetworkID(ctx)
	startBlock, _ := client.BlockByNumber(ctx, nil)
	fmt.Println("startBlock", startBlock.Number())

	type Data struct {
		clientNumber int
		tx *types.Transaction
	}

	rand.Seed(time.Now().Unix())
	// rand.Seed(int64(rand.Intn(100)))
	txs := make([]*Data, numTxs)
	for i := 0; i < numTxs; i++ {
		key, _ := crypto.HexToECDSA(randomHex())
		tx := types.NewTransaction(
			0,
			toAddress,
			big.NewInt(0), // value in wei,
			uint64(21000),  // gasLimit
			new(big.Int),
			nil)
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
		if err != nil {
			log.Fatal(err)
		}
		txs[i] = &Data{rand.Intn(numClients), signedTx}
	}
	fmt.Printf("constructed %d txs, now firing them...\n", len(txs))


	maxChan := make(chan bool, maxFileDescriptors)
	wg := new(sync.WaitGroup)
	txHashes := make([]common.Hash, 0, numTxs)
	for i := numTxs-1; i >= 0; i-- {
		wg.Add(1)
		maxChan <- true
		go func(data *Data) {
			// defer func(maxChan chan bool) { <-maxChan }(maxChan)
			txErr := clients[data.clientNumber].SendTransaction(context.Background(), data.tx)
			if txErr != nil {
				panic(txErr)
			}
			txHashes = append(txHashes, data.tx.Hash())
			<-maxChan
			wg.Done()
		}(txs[i])
	}
	wg.Wait()
	close(maxChan)

	unconfirmed := areAllTxsConfirmed(ctx, client, txHashes)
	for unconfirmed > 0 {
		fmt.Printf("%d pending txs; Sleeping for 3 seconds...\n", unconfirmed)
		time.Sleep(3 * time.Second)
		unconfirmed = areAllTxsConfirmed(ctx, client, txHashes)
	}

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

func areAllTxsConfirmed(ctx context.Context, client *ethclient.Client, txHashes []common.Hash) int {
	var unconfirmed int
	wg := new(sync.WaitGroup)
	maxChan := make(chan bool, maxFileDescriptors)
	for _, h := range txHashes {
		wg.Add(1)
		maxChan <- true
		go func(ctx context.Context, hash common.Hash) {
			defer wg.Done()
			_, err := client.TransactionReceipt(ctx, hash)
			// receipt, err := client.TransactionReceipt(ctx, hash)
			if err != nil {
				unconfirmed++
				if err.Error() != "not found" {
					fmt.Println(err)
				}
			}
			<-maxChan
		}(ctx, h)
	}
	wg.Wait()
	close(maxChan)
	return unconfirmed
}

func randomHex() string {
  bytes := make([]byte, 32)
	rand.Read(bytes)
	// fmt.Println(hex.EncodeToString(bytes))
  return hex.EncodeToString(bytes)
}
