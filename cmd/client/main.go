package main

import (
	"bufio"
	"context"
	"echo-grpc/pb"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

var (
	host = "localhost"
	port = 9001
)

func getMessage() string {
	fmt.Printf("Enter your message: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func getCountValue() (int, error) {
	fmt.Printf("Enter how many time message should reply: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	if text == "" {
		return 1, nil
	}
	count, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func doNormalEcho(client pb.EchoServiceClient) error {
	fmt.Printf("Start processing echo message with grpc server\n")
	for {
		text := getMessage()
		if text == "q" {
			fmt.Printf("bye!\n")
			break
		}

		resp, err := client.Echo(context.Background(), &pb.EchoRequest{Message: text})
		if err != nil {
			return err
		}
		fmt.Printf("Reply message: '%s'\n", resp.Message)
	}
	return nil
}

func doStreamingEcho(client pb.EchoServiceClient) error {

	for {
		text := getMessage()
		if text == "q" {
			fmt.Printf("bye!\n")
			break
		}
		count, err := getCountValue()
		if err != nil || count == 0 {
			fmt.Printf("invalid count value: %d\n", count)
			continue
		}

		fmt.Printf("Start processing streaming echo message with grpc server\n")
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		stream, err := client.EchoStreaming(ctx, &pb.EchoStreamingRequest{
			Message:  text,
			Count:    int32(count),
			Interval: 1000,
		})
		if err != nil {
			return err
		}

		for {
			resp, err := stream.Recv()
			if err != nil && err != io.EOF {
				return err
			}

			log.Printf("response from server: %+v\n", resp)
			if resp.Remaining == 0 {
				break
			}
		}
	}
	log.Printf("receive complete\n")
	return nil
}

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Printf("Connection has been successfully established with %s:%d\n", host, port)

	client := pb.NewEchoServiceClient(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Select Processing Mode:\n")
		fmt.Printf("	S: Stream Echo\n")
		fmt.Printf("	N: Normal Echo\n")
		fmt.Printf("Enter your chioce: ")
		scanner.Scan()
		text := scanner.Text()
		if text == "S" || text == "s" {
			if err := doStreamingEcho(client); err != nil {
				log.Fatal(err)
			}
		} else if text == "N" || text == "n" {
			if err := doNormalEcho(client); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Invalid option")
		}
	}
}
