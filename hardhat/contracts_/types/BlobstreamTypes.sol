// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @notice A tuple of data root with metadata. Each data root is associated
///  with a Celestia block height.
/// @dev `availableDataRoot` in
///  https://github.com/celestiaorg/celestia-specs/blob/master/src/specs/data_structures.md#header
struct DataRootTuple {
    // Celestia block height the data root was included in.
    // Genesis block is height = 0.
    // First queryable block is height = 1.
    uint256 Height;
    // Data root.
    bytes32 DataRoot;
}

/// @notice Merkle Tree Proof structure.
struct BinaryMerkleProof {
    // List of side nodes to verify and calculate tree.
    bytes32[] SideNodes;
    // The key of the leaf to verify.
    uint256 Key;
    // The number of leaves in the tree
    uint256 NumLeaves;
}

/// @notice A representation of the Celestia-app namespace ID and its version.
/// See: https://celestiaorg.github.io/celestia-app/specs/namespace.html
struct Namespace {
    // The namespace version.
    bytes1 Version;
    // The namespace ID.
    bytes28 Id;
}

/// @notice Namespace Merkle Tree node.
struct NamespaceNode {
    // Minimum namespace.
    Namespace Min;
    // Maximum namespace.
    Namespace Max;
    // Node value.
    bytes32 Digest;
}

/// @notice Namespace Merkle Tree Multiproof structure. Proves multiple leaves.
struct NamespaceMerkleMultiproof {
    // The beginning key of the leaves to verify.
    uint256 BeginKey;
    // The ending key of the leaves to verify.
    uint256 EndKey;
    // List of side nodes to verify and calculate tree.
    NamespaceNode[] SideNodes;
}

/// @notice Contains the necessary parameters needed to verify that a data root tuple
/// was committed to, by the Blobstream smart contract, at some specif nonce.
struct AttestationProof {
    // the attestation nonce that commits to the data root tuple.
    uint256 TupleRootNonce;
    // the data root tuple that was committed to.
    DataRootTuple Tuple;
    // the binary merkle proof of the tuple to the commitment.
    BinaryMerkleProof Proof;
}

/// @notice Contains the necessary parameters to prove that some shares, which were posted to
/// the Celestia network, were committed to by the Blobstream smart contract.
struct SharesProof {
    // The shares that were committed to.
    bytes[] Data;
    // The shares proof to the row roots. If the shares span multiple rows, we will have multiple nmt proofs.
    NamespaceMerkleMultiproof[] ShareProofs;
    // The namespace of the shares.
    Namespace Namespace;
    // The rows where the shares belong. If the shares span multiple rows, we will have multiple rows.
    NamespaceNode[] RowRoots;
    // The proofs of the rowRoots to the data root.
    BinaryMerkleProof[] RowProofs;
    // The proof of the data root tuple to the data root tuple root that was posted to the Blobstream contract.
    AttestationProof AttestationProof;
}
