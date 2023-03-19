package main

import (
	"echo-grpc/echo"
	"echo-grpc/pb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	host = "localhost"
	port = 9001
)

func main() {
	apiListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on %s\n", apiListener.Addr().String())

	echoServer := echo.NewEchoServer()

	grpc := grpc.NewServer()
	pb.RegisterEchoServiceServer(grpc, echoServer)
	grpc.Serve(apiListener)
}
