package main

import (
	"bufio"
	"context"
	"echo-grpc/pb"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
)

var (
	host = "localhost"
	port = 9001
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Printf("Connection has been successfully established with %s:%d\n", host, port)

	client := pb.NewEchoServiceClient(conn)

	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for {
		fmt.Printf("Enter your message: ")

		scanner.Scan()
		text = scanner.Text()
		if text == "q" {
			fmt.Printf("bye!\n")
			os.Exit(2)
		}

		resp, err := client.Echo(context.Background(), &pb.EchoRequest{Message: text})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Reply message: '%s'\n", resp.Message)
	}
}
