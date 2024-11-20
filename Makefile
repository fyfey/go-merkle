.PHONY: test test-coverage run-server proto
test:
	go test -v -cover ./...
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
run-server:
	go run cmd/server/*.go -c 1024 -f arrival_in_nara.txt
proto:
	# go get -u -v google.golang.org/protobuf/cmd/protoc-gen-go
	protoc \
		--go_out=./internal/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=./internal/proto \
		--go-grpc_opt=paths=source_relative \
		merkle.proto
