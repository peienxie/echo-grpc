
protobuf:
	protoc -I proto/ proto/*.proto --go_out=plugins=grpc:.

server:
	go build -o bin/server cmd/server/main.go

client:
	go build -o bin/client cmd/client/main.go

build: server client
	
.PHONY: protobuf server client build