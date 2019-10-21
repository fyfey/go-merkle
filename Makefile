proto:
	go get -u -v github.com/golang/protobuf/protoc-gen-go
	protoc --go_out=plugins=grpc:$(shell pwd) merkle.proto
