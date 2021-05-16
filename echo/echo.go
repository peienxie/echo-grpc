package echo

import (
	"context"
	"echo-grpc/pb"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/peer"
)

type EchoServer struct{}

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

func (s *EchoServer) EchoStreaming(req *pb.EchoStreamingRequest, stream pb.EchoService_EchoStreamingServer) error {
	log.Printf("receive request: %+v\n", req)
	for i := 0; i < int(req.Count); i++ {
		resp := &pb.EchoStreamingResponse{}
		resp.Message = req.Message
		resp.Remaining = int32(int(req.Count) - i - 1)
		log.Printf("sending response: %+v\n", resp)
		if err := stream.Send(resp); err != nil {
			return err
		}
		time.Sleep(time.Duration(req.Interval) * time.Millisecond)
	}
	log.Printf("send complete\n")
	return nil
}
