package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/fyfey/go-merkle/internal/proto"
	"github.com/fyfey/go-merkle/internal/server"
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
			file, err := os.Stat(filepath.Join("out", filename))
			if err != nil {
				log.Fatal(err)
			}

			r, err := os.Open(filepath.Join("out", filename))
			if err != nil {
				log.Fatal(err)
			}
			defer r.Close()
			s := server.NewServer(c, 8080, r, chunkSize, filename)
			s.Start()

			log.Printf("file %s size %d\n", filename, file.Size())

			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			grpcServer := grpc.NewServer()
			proto.RegisterMerkleServer(grpcServer, s)

			log.Printf("Server listening on 0.0.0.0:%d\n", s.Port)
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
