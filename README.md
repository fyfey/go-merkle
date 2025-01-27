# Go Merkle Tree

![Example](https://github.com/fyfey/go-merkle/actions/workflows/go.yml/badge.svg)

This is a simple implementation of a Merkle Tree in Go.

There's a simple server/client transfer protocol that allows
the client to request a file from the server and verify the integrity of the file using the Merkle Tree.

## Run Server

```shell
go run ./cmd/server/... -c 64 -f arrival_in_nara.txt
```

## Run Client

```shell
go run ./cmd/client/... --workers 10
```
