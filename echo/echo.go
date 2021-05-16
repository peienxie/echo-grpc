package echo

import (
	"context"
	"echo-grpc/pb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc/peer"
)

type EchoServer struct {}


func NewEchoServer() *EchoServer {
	return &EchoServer{}
}

func getClientAddr(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClientAddr] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientAddr] peer.Addr is nil")
	}

	return pr.Addr.String(), nil
}

func (s *EchoServer) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	addr, err := getClientAddr(ctx)
	if err == nil {
		log.Printf("receive: '%s' from %s\n", req.Message, addr)
	}

	var resp pb.EchoResponse
	resp.Message = req.Message
	return &resp, nil
}

