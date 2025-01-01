.PHONY: test test-coverage run-server run-client proto
test:
	go test -v -cover ./...
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
.PHONY: run-server
run-server:
	go run cmd/server/main.go -c 512 -f arrival_in_nara.txt
.PHONY: run-client
run-client:
	go run cmd/client/main.go
.PHONY: proto
proto:
	# go get -u -v google.golang.org/protobuf/cmd/protoc-gen-go
	protoc \
		--go_out=./internal/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=./internal/proto \
		--go-grpc_opt=paths=source_relative \
		merkle.proto
