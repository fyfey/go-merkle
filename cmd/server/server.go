package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/fyfey/merkle/internal/merkle"
	"github.com/fyfey/merkle/internal/proto"
)

type Server struct {
	context   context.Context
	port      int
	file      io.Reader
	chunkSize int
	tree      *merkle.Tree
	filename  string
	data      [][]byte
	proto.UnimplementedMerkleServer
}

func NewServer(c context.Context, port int, file io.Reader, chunkSize int, filename string) *Server {
	return &Server{
		context:   c,
		port:      port,
		file:      file,
		chunkSize: chunkSize,
		filename:  filename,
	}
}

func (s *Server) Start() {
	log.Println("Starting server...")
	log.Println("Chunk size:", s.chunkSize)
	log.Println("Creating merkle tree...")

	tree, err := merkle.ReadTree(s.file, s.chunkSize)
	if err != nil {
		log.Fatal(err)
	}

	printTree(tree)
	log.Println("Merkle tree created with root", hex.EncodeToString(tree.Root()))
	s.tree = tree
}

// GetMetedata gets the file's metadata
func (s *Server) GetMetadata(context.Context, *proto.Empty) (*proto.Metadata, error) {
	log.Println("GetMetadata")
	return &proto.Metadata{
		Filename:  s.filename,
		Parts:     int32(len(s.tree.GetLeaves())),
		ChunkSize: int32(s.chunkSize),
	}, nil
}

// GetPart returns a given part
func (s *Server) GetPart(ctx context.Context, in *proto.PartRequest) (*proto.Part, error) {
	data := []byte("test") //s.tree.GetLeaf(int(in.Idx)).Data
	if len(data) == 0 {
		return nil, errors.New("Part does not exist")
	}
	proof := s.tree.GetLeaf(int(in.Idx)).GetProof()
	pbProof := &proto.Proof{
		Nodes: make([]*proto.Proof_ProofNode, len(proof)),
		// MerkleRoot: s.tree.Root(),
	}
	for _, node := range proof {
		side := proto.Proof_ProofNode_RIGHT
		if node.Left {
			side = proto.Proof_ProofNode_LEFT
		}
		pbProof.Nodes = append(pbProof.Nodes, &proto.Proof_ProofNode{
			Hash: node.Hash,
			Side: side,
		})
	}
	log.Println("GetPart")
	part := &proto.Part{
		Idx:   in.Idx,
		Data:  data,
		Proof: pbProof,
	}

	log.Printf("GetPart: %v\n", part)

	return part, nil
}

// printTree prints the merkle tree stats in a visual way
func printTree(tree *merkle.Tree) {
	height := strconv.Itoa(tree.Height())
	hexRoot := hex.EncodeToString(tree.Root())
	leaves := strconv.Itoa(len(tree.GetLeaves()))
	fmt.Printf(
		"\n"+
			" ╔═══════════════════════════╗ \n"+
			" ║   root: %s ║                \n"+
			" ║    / \\    ↑               ║ \n"+
			" ║   /   \\   height: %s%s ║ \n"+
			" ║  / \\ / \\  ↓               ║ \n"+
			" ║ .  . .  .                 ║ \n"+
			" ║ %s leaves%s       ║ \n"+
			" ╚═══════════════════════════╝ \n"+
			"\n",
		hexRoot[:8]+"…"+hexRoot[len(hexRoot)-8:],
		height,
		strings.Repeat(" ", 7-len(height)),
		leaves,
		strings.Repeat(" ", 12-len(leaves)),
	)
}
