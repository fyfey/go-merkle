package merkle

import (
	"encoding/hex"
	"slices"
	"testing"
)

func TestBitcoinBlock1000(t *testing.T) {

	// Note that all bitcoin txIDs and roots are stored in little-endian format,
	// so we need to reverse them before using them
	txIDs := []string{
		"8c14f0db3df150123e6f3dbbf30f8b955a8249b62ac1d1ff16284aefa3d06d87",
		"fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4",
		"6359f0868171b1d194cbee1af2f16ea598ae8fad666d9b012c8ed2b79a236ec4",
		"e9a66845e05d5abc0ad04ec80f774a7e585c6e8db975962d069a522137b80c1d",
	}
	expectedRoot := "f3e94742aca4b5ef85488dc37c06c3282295ffec960994b2c0d5ac2a25a95766"

	hex.DecodeString(txIDs[0])

	b0, _ := hex.DecodeString(txIDs[0])
	b1, _ := hex.DecodeString(txIDs[1])
	b2, _ := hex.DecodeString(txIDs[2])
	b3, _ := hex.DecodeString(txIDs[3])

	slices.Reverse(b0)
	slices.Reverse(b1)
	slices.Reverse(b2)
	slices.Reverse(b3)

	tree := NewTree(WithDoubleSHA256Hasher(), WithDuplicateOddLeaves()).
		AddRaw(b0).
		AddRaw(b1).
		AddRaw(b2).
		AddRaw(b3)

	tree.Build()

	root := tree.Root()
	slices.Reverse(root)

	if hex.EncodeToString(tree.Root()) != expectedRoot {
		t.Errorf("Expected %s, got %s", expectedRoot, hex.EncodeToString(tree.Root()))
	}
}

// Adding this test to test odd number of leaves
func TestBitcoinBlock1018(t *testing.T) {

	txIDs := []string{
		"a335b243f5e343049fccac2cf4d70578ad705831940d3eef48360b0ea3829ed4",
		"d5fd11cb1fabd91c75733f4cf8ff2f91e4c0d7afa4fd132f792eacb3ef56a46c",
		"0441cb66ef0cbf78c9ecb3d5a7d0acf878bfdefae8a77541b3519a54df51e7fd",
		"1a8a27d690889b28d6cb4dacec41e354c62f40d85a7f4b2d7a54ffc736c6ff35",
		"1d543d550676f82bf8bf5b0cc410b16fc6fc353b2a4fd9a0d6a2312ed7338701",
	}
	expectedRoot := "5766798857e436d6243b46b5c1e0af5b6806aa9c2320b3ffd4ecff7b31fd4647"

	hex.DecodeString(txIDs[0])

	b0, _ := hex.DecodeString(txIDs[0])
	b1, _ := hex.DecodeString(txIDs[1])
	b2, _ := hex.DecodeString(txIDs[2])
	b3, _ := hex.DecodeString(txIDs[3])
	b4, _ := hex.DecodeString(txIDs[4])

	slices.Reverse(b0)
	slices.Reverse(b1)
	slices.Reverse(b2)
	slices.Reverse(b3)
	slices.Reverse(b4)

	tree := NewTree(WithDoubleSHA256Hasher(), WithDuplicateOddLeaves()).
		AddRaw(b0).
		AddRaw(b1).
		AddRaw(b2).
		AddRaw(b3).
		AddRaw(b4)

	tree.Build()

	root := tree.Root()
	slices.Reverse(root)

	if hex.EncodeToString(tree.Root()) != expectedRoot {
		t.Errorf("Expected %s, got %s", expectedRoot, hex.EncodeToString(tree.Root()))
	}
}
