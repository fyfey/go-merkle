package server

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/fyfey/go-merkle/pkg/merkle"
	"github.com/fyfey/go-merkle/internal/proto"
)

type Server struct {
	context   context.Context
	Port      int
	File      io.Reader
	ChunkSize int
	tree      *merkle.Tree
	Filename  string
	data      [][]byte
	proto.UnimplementedMerkleServer
}

func NewServer(c context.Context, port int, file io.Reader, chunkSize int, filename string) *Server {
	return &Server{
		context:   c,
		Port:      port,
		File:      file,
		ChunkSize: chunkSize,
		Filename:  filename,
	}
}

func (s *Server) Start() {
	log.Println("Starting server...")
	log.Println("Chunk size:", s.ChunkSize)
	log.Println("Creating merkle tree...")

	tree, err := merkle.ReadTree(s.File, s.ChunkSize)
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
		Filename:  s.Filename,
		Parts:     int32(len(s.tree.GetLeaves())),
		ChunkSize: int32(s.ChunkSize),
	}, nil
}

// GetPart returns a given part
func (s *Server) GetPart(ctx context.Context, in *proto.PartRequest) (*proto.Part, error) {
	log.Printf("GetPart %d\n", in.Idx)
	data := s.tree.GetLeaf(int(in.Idx)).Data
	if len(data) == 0 {
		return nil, errors.New("Part does not exist")
	}
	proof := s.tree.GetLeaf(int(in.Idx)).GetProof()
	pbProof := &proto.Proof{
		Nodes: make([]*proto.Proof_ProofNode, len(proof)),
	}
	for i, node := range proof {
		side := proto.Proof_ProofNode_RIGHT
		if node.Left {
			side = proto.Proof_ProofNode_LEFT
		}
		pbProof.Nodes[i] = &proto.Proof_ProofNode{
			Hash: node.Hash,
			Side: side,
		}
	}
	part := &proto.Part{
		Idx:   in.Idx,
		Data:  data,
		Proof: pbProof,
	}

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
