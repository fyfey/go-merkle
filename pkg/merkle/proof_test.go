package merkle

import "testing"

//         abcde
//        /    \
//      abcd    e
//     /    \
//   ab     cd
//  /  \   /  \
// a    b c    d

func TestProve(t *testing.T) {
	tree := NewTree()
	tree.
		Add([]byte("a")).
		Add([]byte("b")).
		Add([]byte("c")).
		Add([]byte("d")).
		Add([]byte("e")).
		Build()

	hasher := &SHA256Hasher{}

	for i := 0; i < 4; i++ {
		nodeToProve := tree.GetLeaf(i)
		proof := nodeToProve.GetProof()
		ok := proof.Prove(nodeToProve.hash, hasher)

		if !ok {
			t.Errorf("Expected true, got false")
		}
	}
}
