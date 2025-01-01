package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fyfey/merkle/internal/merkle"
	"github.com/fyfey/merkle/internal/proto"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Connecting to host...")

	var serverAddr string
	flag.StringVar(&serverAddr, "addr", "127.0.0.1:8080", "Server address")

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewMerkleClient(conn)
	hasher := &merkle.SHA256Hasher{}

	metadata, err := client.GetMetadata(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(filepath.Join("in", metadata.Filename))
	if err != nil {
		log.Fatalf("Failed to open file")
	}
	defer file.Close()

	var wg sync.WaitGroup
	wg.Add(int(metadata.Parts))
	for i := 0; i < int(metadata.Parts); i++ {
		go func() {
			part, err := client.GetPart(context.Background(), &proto.PartRequest{Idx: int32(i)})
			if err != nil {
				log.Fatal(err)
			}
			if !prove(part.Proof, hasher.Hash(part.Data), hasher) {
				log.Fatalf("Part %d failed merkle proof check\n", part.Idx)
			}
			offset := int64(int(part.Idx) * int(metadata.ChunkSize))
			log.Printf("Writing %d bytes @ %d", len(part.Data), offset)
			_, err = file.WriteAt(part.Data, offset)
			if err != nil {
				log.Fatal("Failed writing data to file")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func prove(p *proto.Proof, ha []byte, hasher merkle.Hasher) bool {
	rootHash := p.Nodes[len(p.Nodes)-1].Hash
	for i := 0; i < len(p.Nodes)-1; i++ {
		if p.Nodes[i].Side == proto.Proof_ProofNode_LEFT {
			ha = hasher.Hash(append(ha, p.Nodes[i].Hash...))
		} else {
			ha = hasher.Hash(append(p.Nodes[i].Hash, ha...))
		}
	}

	return bytes.Equal(ha, rootHash)
}
