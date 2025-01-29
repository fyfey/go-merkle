# Go Merkle Tree

[![CI Status](https://github.com/fyfey/go-merkle/actions/workflows/go.yml/badge.svg)](https://github.com/fyfey/go-merkle/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/fyfey/go-merkle)](https://goreportcard.com/report/github.com/fyfey/go-merkle)
[![Go Reference](https://pkg.go.dev/badge/github.com/fyfey/go-merkle.svg)](https://pkg.go.dev/github.com/fyfey/go-merkle)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/fyfey/go-merkle/tree/main/LICENSE)
[![Coverage Report](https://raw.githubusercontent.com/wiki/fyfey/go-merkle/coverage.svg)](https://raw.githack.com/wiki/fyfey/go-merkle/coverage.html)

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
