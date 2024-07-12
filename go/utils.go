package main

import (
	"encoding/hex"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/rpc/core/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"math/big"
)

// These functions are all adapted from Celestia's official tutorial:
// https://docs.celestia.org/developers/blobstream-proof-queries#converting-the-proofs-to-be-usable-in-the-daverifier-library
// They are modified such that []byte are encoded as hexstring

// They mainly convert proofs from the types in Celestia's library to types understood by the BlobstreamX library
// i.e., Go => Solidity

// https://github.com/celestiaorg/blobstream-contracts/blob/v4.1.0/src/lib/tree/Types.sol#L6
type Namespace struct {
	Version string
	Id string
}

func toNamespaceV0(namespaceID []byte) Namespace {
	return Namespace{
		Version: "0x00",
		Id: "0x" + hex.EncodeToString(namespaceID),
	}
}

// https://github.com/celestiaorg/blobstream-contracts/blob/v4.1.0/src/lib/tree/namespace/NamespaceNode.sol#L7
type NamespaceNode struct {
	Min Namespace
	Max Namespace
	Digest string
}

func minNamespace(innerNode []byte) *Namespace {
	return &Namespace{
		Version: "0x" + hex.EncodeToString(innerNode[0:1]),
		Id: "0x" + hex.EncodeToString(innerNode[1:29]),
	}
}

func maxNamespace(innerNode []byte) *Namespace {
	return &Namespace{
		Version: "0x" + hex.EncodeToString(innerNode[29:30]),
		Id: "0x" + hex.EncodeToString(innerNode[30:58]),
	}
}

func toNamespaceNode(node []byte) *NamespaceNode {
	minNs := minNamespace(node)
	maxNs := maxNamespace(node)
	digestString := "0x" + hex.EncodeToString(node[58:])
	return &NamespaceNode{
		Min:    *minNs,
		Max:    *maxNs,
		Digest: digestString,
	}
}

// https://github.com/celestiaorg/blobstream-contracts/blob/v4.1.0/src/lib/tree/namespace/NamespaceMerkleMultiproof.sol
type NamespaceMerkleMultiproof struct {
	BeginKey int64
	EndKey int64
	SideNodes []NamespaceNode
}

func toNamespaceMerkleMultiProofs(proofs []*tmproto.NMTProof) []NamespaceMerkleMultiproof {
	shareProofs := make([]NamespaceMerkleMultiproof, len(proofs))
	for i, proof := range proofs {
		sideNodes := make([]NamespaceNode, len(proof.Nodes))
		for j, node := range proof.Nodes {
			sideNodes[j] = *toNamespaceNode(node)
		}
		shareProofs[i] = NamespaceMerkleMultiproof{
			BeginKey:  int64(proof.Start),
			EndKey:    int64(proof.End),
			SideNodes: sideNodes,
		}
	}
	return shareProofs
}

func toRowRoots(roots []tmbytes.HexBytes) []NamespaceNode {
	rowRoots := make([]NamespaceNode, len(roots))
	for i, root := range roots {
		rowRoots[i] = *toNamespaceNode(root.Bytes())
	}
	return rowRoots
}

// https://github.com/celestiaorg/blobstream-contracts/blob/v4.1.0/src/lib/tree/namespace/NamespaceMerkleProof.sol
type BinaryMerkleProof struct {
	SideNodes []string
	Key int64
	NumLeaves int64
}

func toRowProofs(proofs []*merkle.Proof) []BinaryMerkleProof {
	rowProofs := make([]BinaryMerkleProof, len(proofs))
	for i, proof := range proofs {
		sideNodes := make( []string, len(proof.Aunts))
		for j, sideNode :=  range proof.Aunts {
			sideNodes[j] = "0x" + hex.EncodeToString(sideNode[:])
		}
 		rowProofs[i] = BinaryMerkleProof{
			SideNodes: sideNodes,
			Key:       proof.Index,
			NumLeaves: proof.Total,
		}
	}
	return rowProofs
}

// https://github.com/celestiaorg/blobstream-contracts/blob/v4.1.0/src/DataRootTuple.sol
type DataRootTuple struct {
	Height int64
	DataRoot string
}

// https://github.com/celestiaorg/blobstream-contracts/blob/v4.1.0/src/lib/verifier/DAVerifier.sol#L33
type AttestationProof struct {
	TupleRootNonce *big.Int
	Tuple DataRootTuple 
	Proof BinaryMerkleProof
}

func toAttestationProof(
	nonce *big.Int,
	height int64,
	blockDataRoot tmbytes.HexBytes,
	dataRootInclusionProof merkle.Proof,
) AttestationProof {
	sideNodes := make( []string, len(dataRootInclusionProof.Aunts))
	for i, sideNode :=  range dataRootInclusionProof.Aunts {
		sideNodes[i] = "0x" + hex.EncodeToString(sideNode[:])
	}

	return AttestationProof{
		TupleRootNonce: nonce,
		Tuple: DataRootTuple{
			Height: int64(height),
			DataRoot: "0x" + 
				hex.EncodeToString(blockDataRoot.Bytes()),
		},
		Proof: BinaryMerkleProof{
			SideNodes: sideNodes,
			Key: dataRootInclusionProof.Index,
			NumLeaves: dataRootInclusionProof.Total,
		},
	}
}

// https://github.com/celestiaorg/blobstream-contracts/blob/v4.1.0/src/lib/verifier/DAVerifier.sol#L16
type SharesProof struct {
	Data []string
	ShareProofs []NamespaceMerkleMultiproof
	Namespace Namespace
	RowRoots []NamespaceNode
	RowProofs []BinaryMerkleProof
	AttestationProof AttestationProof
}

func toSharesProof(
	proofNonce *big.Int,
	height int64,
	blockDataRoot tmbytes.HexBytes,
	dcProof *coretypes.ResultDataRootInclusionProof,
	sharesProof types.ShareProof,
) SharesProof {
	
	data := make([]string, len(sharesProof.Data))
	for i, chunk :=  range sharesProof.Data {
		data[i] = "0x" + hex.EncodeToString(chunk[:])
	}

	nmtProof := toNamespaceMerkleMultiProofs(sharesProof.ShareProofs)
	namespace := toNamespaceV0(sharesProof.NamespaceID)
	rowRoots := toRowRoots(sharesProof.RowProof.RowRoots)
	rowProofs := toRowProofs(sharesProof.RowProof.Proofs)
	attProof := toAttestationProof(
		proofNonce,
		height,
		blockDataRoot,
		dcProof.Proof,
	)

	return SharesProof {
		Data: data,
		ShareProofs: nmtProof,
		Namespace: namespace,
		RowRoots: rowRoots,
		RowProofs: rowProofs,
		AttestationProof: attProof,
	}
}

