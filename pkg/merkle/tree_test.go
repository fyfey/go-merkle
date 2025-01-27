package merkle

import (
	"encoding/hex"
	"errors"
	"os"
	"testing"
)

const (
	HashA      = "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"
	HashB      = "3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d"
	HashC      = "2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6"
	HashD      = "18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4"
	HashE      = "3f79bb7b435b05321651daefd374cdc681dc06faa65e374e38337b88ca046dea"
	Root       = "d71f8983ad4ee170f8129f1ebcdd7440be7798d8e1c80420bf11f1eced610dba"
	ReaderRoot = "75894b4ffb30aa65e4624b84611dd47f17ebf8f4caeeb9a8109bd851506ae170"
)

func TestTree(t *testing.T) {
	tree := NewTree(WithSHA256Hasher())

	if tree.Height() != 0 {
		t.Errorf("Expected 0, got %d", tree.Height())
	}

	err := tree.Build()

	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	hD, _ := hex.DecodeString(HashD)

	err = tree.
		Add([]byte("a")).
		Add([]byte("b")).
		Add([]byte("c")).
		AddRaw(hD).
		Add([]byte("e")).
		Build()

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if tree.Height() != 3 {
		t.Errorf("Expected 3, got %d", tree.Height())
	}

	if len(tree.GetLeaves()) != 5 {
		t.Errorf("Expected 5, got %d", len(tree.GetLeaves()))
	}

	if hex.EncodeToString(tree.Root()) != Root {
		t.Errorf("Expected %s, got %s", HashC, hex.EncodeToString(tree.Root()))
	}

	a := tree.GetLeaf(0)
	if hex.EncodeToString(a.hash) != HashA {
		t.Errorf("Expected %s, got %s", HashA, hex.EncodeToString(a.hash))
	}

	// test build errors when already built
	err = tree.Build()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestReader(t *testing.T) {
	r, err := os.Open("../../out/arrival_in_nara.txt")
	if err != nil {
		t.Fatalf("Expected nil, got %s", err)
	}
	defer r.Close()

	_, err = ReadTree(r, 0)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	_, err = ReadTree(r, -1)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	tree, err := ReadTree(r, 512)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if hex.EncodeToString(tree.Root()) != ReaderRoot {
		t.Errorf("Expected %s, got %s", ReaderRoot, hex.EncodeToString(tree.Root()))
	}

	// Test WithDoubleSHA256Hasher

	tree = NewTree(WithDoubleSHA256Hasher())

	if tree.hasher == nil {
		t.Errorf("Expected hasher, got nil")
	}
}

// DodgeyReader implements the io.Reader interface but always returns an error
type DodgeyReader struct {
}

func (r *DodgeyReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("DodgeyReader")
}

func TestDodgeyReader(t *testing.T) {
	r := &DodgeyReader{}
	_, err := ReadTree(r, 512)
	if err == nil && err.Error() != "DodgeyReader" {
		t.Errorf("Expected nil, got %s", err)
	}
}
