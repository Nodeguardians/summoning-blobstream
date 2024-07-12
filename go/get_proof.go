package main

import (
	"context"
	"encoding/json"
	"encoding/hex"
	"fmt"
	"github.com/celestiaorg/celestia-app/pkg/square"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	blobstreamxwrapper "github.com/succinctlabs/blobstreamx/bindings"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"math/big"
	"os"
)

// INPUTS:
// PayForBlob Transaction Hash (w/o 0x prefix)
const TX_HASH = "5EC9E71F71C40A16B6F1E71CE4E8ADA205EE07072052DF25E0A44A1ED21DF17A"
// Searches last N commitments in BlobstreamX
const SEARCH_RANGE = 5000;

func main() {
	err := getProofs(TX_HASH)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// A function that queries and pretty-prints the Blobstream proofs
// of a given PayForBlob transaction.
func getProofs(txhash string) error {
	ctx := context.Background()

	// Start the tendermint RPC client (w/ Celestia Mocha Full Node)
	trpc, err := http.New("tcp://consensus-full-mocha-4.celestia-mocha.com:26657", "/websocket")
	if err != nil {
		return err
	}
	err = trpc.Start()
	if err != nil {
		return err
	}

	// Retrieve the target PayForBlob transaction
	txhashbytes, err := hex.DecodeString(txhash)
	tx, err := trpc.Tx(ctx, txhashbytes, true)
	if err != nil {
		return err
	}

	// Find BlobstreamX commitment corresponding to the block height 
	// that contains the PayForBlob transaction.
	// This can be found by scanning through the emitted events of the BlobstreamX contract.

	// Connect to an EVM RPC endpoint
	ethClient, err := ethclient.Dial("https://rpc-sepolia-eth.nodeguardians.io")
	if err != nil {
		return err
	}
	defer ethClient.Close()

	// Use the BlobstreamX contract binding
	wrapper, err := blobstreamxwrapper.NewBlobstreamX(
		ethcmn.HexToAddress("0xF0c6429ebAB2e7DC6e05DaFB61128bE21f13cb1e"), 
		ethClient,
	)
	if err != nil {
		return err
	}

	LatestBlockNumber, err := ethClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	// Scan the latest blocks for a BlobstreamX commitment event
	// that has a range covering our required transaction height.
	eventsIterator, err := wrapper.FilterDataCommitmentStored(
		&bind.FilterOpts{
			Context: ctx,
			Start: LatestBlockNumber - SEARCH_RANGE,
			End: &LatestBlockNumber,
		},
		nil,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	var event *blobstreamxwrapper.BlobstreamXDataCommitmentStored
	for eventsIterator.Next() {
		e := eventsIterator.Event
		if int64(e.StartBlock) <= tx.Height && tx.Height < int64(e.EndBlock) {
			event = &blobstreamxwrapper.BlobstreamXDataCommitmentStored{
				ProofNonce:     e.ProofNonce,
				StartBlock:     e.StartBlock,
				EndBlock:       e.EndBlock,
				DataCommitment: e.DataCommitment,
			}
			break
		}
	}
	if err := eventsIterator.Error(); err != nil {
		return err
	}
	err = eventsIterator.Close()
	if err != nil {
		return err
	}
	if event == nil {
		return fmt.Errorf("couldn't find range containing the transaction height")
	}

	// Query the Celestia full node for the block data root inclusion proof 
	// to the data root tuple root
	dcProof, err := trpc.DataRootInclusionProof(
		ctx, 
		uint64(tx.Height), 
		event.StartBlock, 
		event.EndBlock,
	)
	if err != nil {
		return err
	}

	blockRes, err := trpc.Block(ctx, &tx.Height)
	if err != nil {
		return err
	}

	// Get the start and end shares used to store the message
	blobShareRange, err := square.BlobShareRange(
		blockRes.Block.Txs.ToSliceOfBytes(), 
		int(tx.Index), 
		0, 
		blockRes.Block.Header.Version.App,
	)
	if err != nil {
		return err
	}

	sharesProof, err := trpc.ProveShares(
		ctx, 
		uint64(tx.Height), 
		uint64(blobShareRange.Start), 
		uint64(blobShareRange.End),
	)
	if err != nil {
		return err
	}

	printProofs(
		event.ProofNonce,
		tx.Height,
		blockRes.Block.DataHash,
		dcProof,
		sharesProof,
	)
    return nil
}

// A function that pretty-prints a given Blobstream proof
// In a format expected by the BlobstreamX contract.
func printProofs(
	proofNonce *big.Int,
	height int64,
	blockDataRoot tmbytes.HexBytes,
	dcProof *coretypes.ResultDataRootInclusionProof,
	sharesProof types.ShareProof,
) {
	fmt.Println("Height: ", height)
	sp := toSharesProof(
		proofNonce,
		height,
		blockDataRoot,
		dcProof,
		sharesProof,
	)

	b, err := json.MarshalIndent(sp, "", "  ")
    if err != nil {
        fmt.Println("JSON failed")
    }

	fmt.Println("SharesProof: ", string(b))


	fmt.Println(
		"DataHash: ", 
		"0x" + hex.EncodeToString(blockDataRoot.Bytes()),
	)

	f, err := os.Create("data/proof.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(string(b))
	fmt.Println(l, "bytes written successfully at data/proof.json")
	if err != nil {
		fmt.Println(err)
        f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}
