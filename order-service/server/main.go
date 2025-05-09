package main

import (
	"context"
	"log"
	"net"

	orderpb "order-service/proto"

	"google.golang.org/grpc"
)

type server struct {
	orderpb.UnimplementedOrderServiceServer
}

func (s *server) CreateOrder(ctx context.Context, req *orderpb.OrderRequest) (*orderpb.OrderResponse, error) {
	log.Printf("Received order: ItemID=%s, Quantity=%d", req.ItemId, req.Quantity)
	return &orderpb.OrderResponse{
		OrderId: "ORD123",
		Status:  "CONFIRMED",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(s, &server{})
	log.Println("gRPC server running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
