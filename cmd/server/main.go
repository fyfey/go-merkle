package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fyfey/merkle/internal/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var filename string
var chunkSize int

func main() {
	c := context.Background()

	rootCmd := &cobra.Command{
		Use:   "merkleserver",
		Short: "File server using merkle",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := os.Stat(filename)
			if err != nil {
				log.Fatal(err)
			}

			r, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer r.Close()
			s := NewServer(c, 8080, r, chunkSize, filename)
			s.Start()

			log.Printf("file %s size %d\n", filename, file.Size())

			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			grpcServer := grpc.NewServer()
			proto.RegisterMerkleServer(grpcServer, s)

			log.Printf("Server listening on 0.0.0.0:%d\n", s.port)
			grpcServer.Serve(lis)
		},
	}
	rootCmd.Flags().StringVarP(&filename, "filename", "f", "", "File to serve")
	rootCmd.Flags().IntVarP(&chunkSize, "chunksize", "c", 1024, "Chunk size in bytes to split the file")
	rootCmd.MarkFlagRequired("filename")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c.Done()
}
