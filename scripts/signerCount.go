package scripts

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/bor/crypto"
	"github.com/maticnetwork/bor/ethclient"
	"github.com/maticnetwork/bor/rlp"
	"github.com/maticnetwork/bor/rpc"
	"golang.org/x/crypto/sha3"
)

var (
	extraSeal           = 65
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")
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
	fmt.Printf("Last block %v\n", latestBlock.Header().Number.Int64())
	for i := latestBlock.Header().Number; i.Int64() > 0; i.Sub(i, big.NewInt(1)) {
		block, err := client.BlockByNumber(context.Background(), i)
		if err != nil {
			log.Fatal(err)
		}

		signer, _ := ecrecover(block.Header())
		signerCount[signer.String()]++
	}
	for key, element := range signerCount {
		fmt.Println(key, element)
	}
}

func ecrecover(header *types.Header) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(SealHash(header).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	return signer, nil
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header)
	hasher.Sum(hash[:0])
	return hash
}

func encodeSigHeader(w io.Writer, header *types.Header) {
	err := rlp.Encode(w, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}
