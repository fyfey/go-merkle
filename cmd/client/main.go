package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fyfey/go-merkle/internal/merkle"
	"github.com/fyfey/go-merkle/internal/proto"
	"github.com/gosuri/uilive"
	"google.golang.org/grpc"
)

type Status int

const (
	Downloading Status = iota
	Completed
)

type PartStatus struct {
	Idx     int
	Message string
	Status  Status
}

const ProgressChar = "█"
const BarWidth = 20.00

var messages []string
var hasher merkle.Hasher = &merkle.SHA256Hasher{}

func main() {
	var serverAddr string
	var workerCount int
	flag.StringVar(&serverAddr, "addr", "127.0.0.1:8080", "Server address")
	flag.IntVar(&workerCount, "workers", 10, "Worker count")
	flag.Parse()

	messages = make([]string, workerCount)

	log.Println("Connecting to host...")
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := proto.NewMerkleClient(conn)

	metadata, err := client.GetMetadata(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(filepath.Join("in", metadata.Filename))
	if err != nil {
		log.Fatalf("Failed to open file")
	}
	defer file.Close()

	in := make(chan int)
	ret := make(chan error)
	status := make(chan PartStatus)

	writer := uilive.New()
	writer.RefreshInterval = time.Millisecond
	writer.Start()

	go statusWriter(writer, status, metadata.Filename, metadata.Parts, metadata.ChunkSize, workerCount)

	wg := sync.WaitGroup{}
	wg.Add(int(metadata.Parts))
	for i := 0; i < workerCount; i++ {
		go worker(i, in, ret, status, client, metadata, file, hasher, &wg)
	}
	go func() {
		for i := 0; i < int(metadata.Parts); i++ {
			in <- i
		}
	}()
	go func() {
		for err := range ret {
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}()

	wg.Wait()
	writer.Stop()
	fmt.Println("Done!")
	close(in)
	close(ret)
}

// prove verifies the merkle proof of a part
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

// woerkers downloads a part from the server and writes it to the file. sending the status of the download to the status channel
func worker(idx int, in chan int, ret chan error, status chan PartStatus, client proto.MerkleClient, metadata *proto.Metadata, file *os.File, hasher merkle.Hasher, wg *sync.WaitGroup) {
	for x := range in {
		status <- PartStatus{Idx: idx, Status: Downloading, Message: fmt.Sprintf("%d) #%d...\n", idx, x)}
		part, err := client.GetPart(context.Background(), &proto.PartRequest{Idx: int32(x)})
		if err != nil {
			ret <- err
		}
		if !prove(part.Proof, hasher.Hash(part.Data), hasher) {
			ret <- fmt.Errorf("Part %d failed merkle proof check\n", part.Idx)
		}
		offset := int64(int(part.Idx) * int(metadata.ChunkSize))
		_, err = file.WriteAt(part.Data, offset)
		if err != nil {
			ret <- fmt.Errorf("Failed writing data to file")
		} else {
			ret <- nil
		}
		time.Sleep(time.Duration(rand.Int63n(500)+50) * time.Millisecond)
		status <- PartStatus{Idx: idx, Status: Completed, Message: fmt.Sprintf("%d) #%d... OK\n", idx, x)}
		time.Sleep(500 * time.Millisecond)
		wg.Done()
	}
}

// statusWriter updates the messages slice with the latest status of each worker's part download
// and each time a status changes, all statuses are redrawn to the terminal
func statusWriter(writer *uilive.Writer, status chan PartStatus, filename string, totalParts int32, chunkSize int32, workerCount int) {
	partsCompleted := 0
	for msg := range status {
		if msg.Status == Completed {
			partsCompleted++
		}
		messages[msg.Idx] = msg.Message
		fmt.Fprintf(writer, "Downloading %s with %d workers...\n", filename, workerCount)
		for _, m := range messages {
			fmt.Fprintf(writer.Newline(), m)
		}
		percent := float64(partsCompleted) / float64(totalParts) * 100
		bars := percent * BarWidth / 100
		fmt.Fprintf(
			writer.Newline(),
			"Progress: %s%s %d/%d\n",
			strings.Repeat(ProgressChar, int(bars)),
			strings.Repeat(" ", int(BarWidth-bars)),
			int32(partsCompleted)*chunkSize,
			int32(totalParts)*chunkSize,
		)
	}
}
