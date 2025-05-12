package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "streaming-practice/pb/pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHelloUnary(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func (s *server) SayHelloServerStream(req *pb.HelloRequest, stream pb.HelloService_SayHelloServerStreamServer) error {
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Hello %s - %d", req.Name, i)
		if err := stream.Send(&pb.HelloResponse{Message: msg}); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (s *server) SayHelloClientStream(stream pb.HelloService_SayHelloClientStreamServer) error {
	var names []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			msg := fmt.Sprintf("Hello to: %v", names)
			return stream.SendAndClose(&pb.HelloResponse{Message: msg})
		}
		if err != nil {
			return err
		}
		names = append(names, req.Name)
	}
}

func (s *server) SayHelloBiDiStream(stream pb.HelloService_SayHelloBiDiStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("Hello, %s", req.Name)
		if err := stream.Send(&pb.HelloResponse{Message: msg}); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, &server{})
	log.Println("Server running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
