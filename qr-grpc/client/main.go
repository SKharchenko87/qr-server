package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "qr-grpc/proto"
	"time"
)

func main() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewQRGenerateServiceClient(conn)
	req := &pb.GenerateRequest{
		Text: "1234567890",
		Kind: pb.KindType_Medium,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := client.Generate(ctx, req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("QR: %v", res)
}
