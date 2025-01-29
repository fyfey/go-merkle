package merkle

import (
	"errors"
	"io"
)

type OddLeafStrategy int

const (
	DuplicateOddLeaves OddLeafStrategy = iota + 1
	IgnoreOddLeaves
)

type Tree struct {
	data            [][]*Node
	hasher          Hasher
	oddLeafStrategy OddLeafStrategy
}

type TreeOpt func(*Tree)

func WithSHA256Hasher() TreeOpt {
	return func(t *Tree) {
		t.hasher = SHA256Hasher{}
	}
}

func WithDoubleSHA256Hasher() TreeOpt {
	return func(t *Tree) {
		t.hasher = DoubleSHA256Hasher{}
	}
}

func WithDuplicateOddLeaves() TreeOpt {
	return func(t *Tree) {
		t.oddLeafStrategy = DuplicateOddLeaves
	}
}

func NewTree(opts ...TreeOpt) *Tree {
	data := make([][]*Node, 0)
	data = append(data, make([]*Node, 0))
	t := &Tree{data: data}

	for _, opt := range opts {
		opt(t)
	}

	if t.hasher == nil {
		t.hasher = SHA256Hasher{}
	}
	if t.oddLeafStrategy == 0 {
		t.oddLeafStrategy = IgnoreOddLeaves
	}

	return t
}

func ReadTree(r io.Reader, chunkSize int) (*Tree, error) {
	if chunkSize <= 0 {
		return nil, errors.New("invalid chunk size. Must be greater than 0")
	}

	buf := make([]byte, chunkSize)
	Tree := NewTree()
	for {
		read, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		newData := make([]byte, read)
		copy(newData, buf[:read])

		Tree.Add(newData)
	}
	Tree.Build()

	return Tree, nil
}

func (t *Tree) Add(data []byte) *Tree {
	base := t.data
	base[0] = append(base[0], NewNode(data, t.hasher))
	t.data = base

	return t
}

func (t *Tree) AddRaw(hash []byte) *Tree {
	base := t.data
	base[0] = append(base[0], NewRawNode(hash, t.hasher))
	t.data = base

	return t
}

// GetLeaf returns the leaf node at the given index
func (t *Tree) GetLeaf(index int) *Node {
	return t.data[0][index]
}

// GetHeight returns the nodes at the given height
func (t *Tree) GetHeight(index int) []*Node {
	return t.data[index]
}

// GetLeaves returns a slice of leaf nodes
func (t *Tree) GetLeaves() []*Node {
	return t.data[0]
}

// Height returns the height of the tree (leaf nodes are at height 0)
func (t *Tree) Height() int {
	return len(t.data) - 1
}

// Root returns the root hash of the tree
func (t *Tree) Root() []byte {
	return t.data[len(t.data)-1][0].hash
}

// Build builds the tree from the leaf nodes
func (t *Tree) Build() error {
	nodes := t.data

	if len(nodes) != 1 {
		return errors.New("Tree already built")
	}
	if len(nodes[0]) == 0 {
		return errors.New("No nodes to build")
	}

	height := 0
	for {
		if len(nodes[height]) == 1 {
			break
		}
		nextHeight := make([]*Node, 0)
		for i := 0; i < int(len(nodes[height])/2)*2; i += 2 {
			newNode := NewParent(nodes[height][i], nodes[height][i+1])
			nextHeight = append(nextHeight, newNode)
		}
		if len(nodes[height])%2 != 0 {
			switch t.oddLeafStrategy {
			case DuplicateOddLeaves:
				// create a parent with the same node as both children
				newNode := NewParent(nodes[height][len(nodes[height])-1], nodes[height][len(nodes[height])-1])
				nextHeight = append(nextHeight, newNode)
			case IgnoreOddLeaves:
				nextHeight = append(nextHeight, nodes[height][len(nodes[height])-1])
			default:
				return errors.New("Invalid odd leaf strategy")
			}
		}
		nodes = append(nodes, nextHeight)
		height++
	}
	t.data = nodes

	return nil
}
