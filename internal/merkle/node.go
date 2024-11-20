package merkle

import (
	"encoding/hex"
	"encoding/json"
)

type Node struct {
	hasher      Hasher
	parent      *Node
	left, right *Node
	hash        []byte
	// Data        []byte
}

// NewNode creates a new node and hashes the data
func NewNode(data []byte, hasher Hasher) *Node {
	return &Node{
		hash:   hasher.Hash(data),
		hasher: hasher,
		// Data:   data,
	}
}

// NewNode creates a new node and hashes the data
func NewRawNode(hash []byte, hasher Hasher) *Node {
	return &Node{hash: hash, hasher: hasher}
}

// NewParent creates a new node and sets the left and right children
func NewParent(left *Node, right *Node) *Node {
	parent := &Node{hasher: left.hasher}

	return parent.SetChildren(left, right)
}

// SetChildren sets the left and right children of a node and calculates the hash
func (n *Node) SetChildren(left *Node, right *Node) *Node {
	n.left = left
	n.right = right
	left.parent = n
	right.parent = n
	n.hash = n.hasher.Hash(append(left.hash, right.hash...))
	return n
}

// MarshalJSON marshals the node to JSON
func (n Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Left  *Node  `json:"left"`
		Right *Node  `json:"right"`
		Hash  string `json:"hash"`
	}{
		n.left,
		n.right,
		hex.EncodeToString(n.hash),
	})
}

// Sibling returns the sibling of a node
func (n *Node) Sibling() *Node {
	if n.parent == nil {
		return nil
	}
	if n.parent.left == n {
		return n.parent.right
	}
	return n.parent.left
}

// Uncle returns the uncle of a node
func (n *Node) Uncle() *Node {
	if n.parent == nil {
		return nil
	}
	return n.parent.Sibling()
}
