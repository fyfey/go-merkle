package merkle

import (
	"encoding/hex"
	"encoding/json"
	"testing"
)

const (
	HashLeft      = "360f84035942243c6a36537ae2f8673485e6c04455a0a85a0db19690f2541480"
	HashRight     = "27042f4e6eca7d0b2a7ee4026df2ecfa51d3339e6d122aa099118ecd8563bad9"
	HashLeftRight = "2a9870f5b7eb1cd732d95224cfea825a7b8772136cb497b20d2e3c612dfc90fe"
	ExpectedJSON  = `{"left":{"left":null,"right":null,"hash":"360f84035942243c6a36537ae2f8673485e6c04455a0a85a0db19690f2541480"},"right":{"left":null,"right":null,"hash":"27042f4e6eca7d0b2a7ee4026df2ecfa51d3339e6d122aa099118ecd8563bad9"},"hash":"2a9870f5b7eb1cd732d95224cfea825a7b8772136cb497b20d2e3c612dfc90fe"}`
)

func TestNode(t *testing.T) {
	hasher := &SHA256Hasher{}
	left := NewNode([]byte("left"), hasher)
	right := NewNode([]byte("right"), hasher)

	parent := NewParent(left, right)

	// expect the sha256 hash of the concatenation of the left and right hashes
	if hex.EncodeToString(parent.hash) != HashLeftRight {
		t.Errorf("Expected %s, got %s", HashLeftRight, hex.EncodeToString(parent.hash))
	}
}

func TestRawNode(t *testing.T) {
	hasher := &SHA256Hasher{}
	hLeft, _ := hex.DecodeString(HashLeft)
	hRight, _ := hex.DecodeString(HashRight)
	left := NewRawNode(hLeft, hasher)
	right := NewRawNode(hRight, hasher)

	parent := NewParent(left, right)

	// expect the sha256 hash of the concatenation of the left and right hashes
	if hex.EncodeToString(parent.hash) != HashLeftRight {
		t.Errorf("Expected %s, got %s", HashLeftRight, hex.EncodeToString(parent.hash))
	}
}

func TestSibling(t *testing.T) {
	hasher := &SHA256Hasher{}
	left := NewNode([]byte("left"), hasher)
	right := NewNode([]byte("right"), hasher)

	leftsSibling := left.Sibling()
	if leftsSibling != nil {
		t.Errorf("Expected nil, got %p", leftsSibling)
	}

	NewParent(left, right)

	leftsSibling = left.Sibling()

	if leftsSibling != right {
		t.Errorf("Expected sibling to be %p, got %p", right, leftsSibling)
	}

	rightsSibling := right.Sibling()

	if rightsSibling != left {
		t.Errorf("Expected sibling to be %p, got %p", left, rightsSibling)
	}

}

func TestUncle(t *testing.T) {
	/*
	 *      grandparent
	 *       |      |
	 *    parent  uncle
	 *     |  |
	 *  left right
	 *
	 */

	hasher := &SHA256Hasher{}
	parent := NewNode([]byte("parent"), hasher)
	uncle := NewNode([]byte("uncle"), hasher)
	NewParent(parent, uncle) // grandparent

	left := NewNode([]byte("left"), hasher)
	right := NewNode([]byte("right"), hasher)

	leftUncle := left.Uncle()
	if leftUncle != nil {
		t.Errorf("Expected nil, got %p", leftUncle)
	}

	parent.SetChildren(left, right)

	leftUncle = left.Uncle()
	if leftUncle != uncle {
		t.Errorf("Expected uncle to be %p, got %p", uncle, leftUncle)
	}

	rightUncle := right.Uncle()
	if rightUncle != uncle {
		t.Errorf("Expected uncle to be %p, got %p", uncle, rightUncle)
	}
}

func TestMarshaLJSON(t *testing.T) {
	hasher := &SHA256Hasher{}
	left := NewNode([]byte("left"), hasher)
	right := NewNode([]byte("right"), hasher)

	parent := NewParent(left, right)

	json, err := json.Marshal(parent)
	if err != nil {
		t.Errorf("Error marshalling json: %s", err)
	}

	// expect the sha256 hash of the concatenation of the left and right hashes
	if string(json) != ExpectedJSON {
		t.Errorf("Expected %s, got %s", ExpectedJSON, json)
	}

}
