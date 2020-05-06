package benchmarking

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"math/big"
	"math/rand"

	"strconv"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/bor/crypto"
	"github.com/maticnetwork/bor/ethclient"
)

// RapidFire fires txs in bulk
func RapidFire(numTxs int, numClients int, seed int64, delay int) {
	clients := make(map[int]*ethclient.Client)
	home, _ := os.UserHomeDir()
	for i := 0; i < numClients; i++ {
		client, err := ethclient.Dial(home + "/matic/testnets/bor-devnet/node" + strconv.Itoa(i+1) + "/bor.ipc")
		if err != nil {
			panic(err)
		}
		clients[i] = client
	}

	rand.Seed(seed)
	toAddress := common.HexToAddress(randomHex())

	ctx := context.Background()
	client, _ := clients[0]
	chainID, _ := client.NetworkID(ctx)
	startBlock, _ := client.BlockByNumber(ctx, nil)
	fmt.Println("startBlock", startBlock.Number())

	type Data struct {
		clientNumber int
		tx *types.Transaction
	}

	key, _ := crypto.HexToECDSA(randomHex())
	txs := make([]*Data, numTxs)
	for i := 0; i < numTxs; i++ {
		tx := types.NewTransaction(
			uint64(i),
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
	fmt.Printf("constructed %d txs, sleeping for %d seconds before firing them...\n", len(txs), delay)
	time.Sleep(time.Duration(delay) * time.Second);

	txHashes := make([]common.Hash, 0, numTxs)
	for i := numTxs-1; i >= 0; i-- {
		data := txs[i]
		txErr := clients[data.clientNumber].SendTransaction(context.Background(), data.tx)
		if txErr != nil {
			panic(txErr)
		}
		txHashes = append(txHashes, data.tx.Hash())
	}
	fmt.Printf("fired %d txs\n", len(txHashes))
	endBlock, _ := client.BlockByNumber(ctx, nil)
	fmt.Println("endBlock", endBlock.Number())
}

func randomHex() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
