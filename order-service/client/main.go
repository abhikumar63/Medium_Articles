package main

import (
	"context"
	"log"
	"time"

	orderpb "order-service/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := orderpb.NewOrderServiceClient(conn)

	req := &orderpb.OrderRequest{
		ItemId:   "ITEM42",
		Quantity: 3,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreateOrder(ctx, req)
	if err != nil {
		log.Fatalf("Error calling CreateOrder: %v", err)
	}

	log.Printf("Order Response: ID=%s, Status=%s", res.OrderId, res.Status)
}
