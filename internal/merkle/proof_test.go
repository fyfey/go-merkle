package merkle

import "testing"

//         abcde
//        /    \
//      abcd    e
//     /    \
//   ab     cd
//  /  \   /  \
// a    b c    d

func TestProveLeft(t *testing.T) {
	tree := NewTree()
	tree.
		Add([]byte("a")).
		Add([]byte("b")).
		Add([]byte("c")).
		Add([]byte("d")).
		Add([]byte("e")).
		Build()

	hasher := &SHA256Hasher{}
	nodeToProve := tree.GetLeaf(0)
	proof := nodeToProve.GetProof()
	ok := proof.Prove(nodeToProve.hash, hasher)

	if !ok {
		t.Errorf("Expected true, got false")
	}
}

func TestProveRight(t *testing.T) {
	tree := NewTree()
	tree.
		Add([]byte("a")).
		Add([]byte("b")).
		Add([]byte("c")).
		Add([]byte("d")).
		Add([]byte("e")).
		Build()

	hasher := &SHA256Hasher{}
	nodeToProve := tree.GetLeaf(3)
	proof := nodeToProve.GetProof()
	ok := proof.Prove(nodeToProve.hash, hasher)

	if !ok {
		t.Errorf("Expected true, got false")
	}
}
