package merkle

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type ProofNode struct {
	Left bool
	Hash []byte
}

func (n ProofNode) String() string {
	return fmt.Sprintf("Left: %v, Hash: %s", n.Left, hex.EncodeToString(n.Hash))
}

type MerkleProof []ProofNode

// GetProof returns a merkleProof (list of hashes and whether to place your computed hash on the left)
// The first item in the hash is for the sibling at height 0, then for the sibling of the computed hash
// The last item in the hash is the root hash and should be compared against the computed root hash.
func (n *Node) GetProof() MerkleProof {
	proof := MerkleProof{}
	nextProof := n.Sibling()
	for {
		left := bytes.Equal(nextProof.parent.right.hash, nextProof.hash)
		proof = append(proof, ProofNode{left, nextProof.hash})
		if nextProof.Uncle() == nil {
			proof = append(proof, ProofNode{left, nextProof.parent.hash})
			break
		}
		nextProof = nextProof.Uncle()
	}
	return proof
}

// Prove proves that the a hash is correct for the given proof
func (p MerkleProof) Prove(h []byte, hasher Hasher) bool {
	root := p[len(p)-1].Hash
	for i := 0; i < len(p)-1; i++ {
		if p[i].Left {
			h = hasher.Hash(append(h, p[i].Hash...))
		} else {
			h = hasher.Hash(append(p[i].Hash, h...))
		}
	}

	return bytes.Equal(h, root)
}
