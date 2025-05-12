package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "streaming-practice/pb/pb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	// Unary
	unaryResp, _ := client.SayHelloUnary(context.Background(), &pb.HelloRequest{Name: "Unary"})
	log.Println("Unary:", unaryResp.Message)

	// Server Streaming
	stream1, _ := client.SayHelloServerStream(context.Background(), &pb.HelloRequest{Name: "ServerStream"})
	for {
		msg, err := stream1.Recv()
		if err == io.EOF {
			break
		}
		log.Println("ServerStream:", msg.Message)
	}

	// Client Streaming
	stream2, _ := client.SayHelloClientStream(context.Background())
	for _, name := range []string{"Alice", "Bob", "Carol"} {
		stream2.Send(&pb.HelloRequest{Name: name})
	}
	resp2, _ := stream2.CloseAndRecv()
	log.Println("ClientStream:", resp2.Message)

	// Bidirectional Streaming
	stream3, _ := client.SayHelloBiDiStream(context.Background())
	go func() {
		for _, name := range []string{"X", "Y", "Z"} {
			stream3.Send(&pb.HelloRequest{Name: name})
			time.Sleep(time.Millisecond * 500)
		}
		stream3.CloseSend()
	}()
	for {
		resp, err := stream3.Recv()
		if err == io.EOF {
			break
		}
		log.Println("BiDiStream:", resp.Message)
	}
}
