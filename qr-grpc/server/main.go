package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "qr-grpc/proto"

	qr "github.com/SKharchenko87/qr"
)

// Сервер, реализующий генератор
type server struct {
	pb.UnimplementedQRGenerateServiceServer
}

func (s *server) Generate(ctx context.Context, req *pb.GenerateRequest) (*pb.GenerateResponse, error) {
	log.Printf("Received: %v", req)
	text := req.GetText()
	levelCorrection := req.GetLevelCorrection()
	qr := qr.GenerateQR(text, qr.LevelCorrection(levelCorrection))
	size := len(qr)
	res := pb.GenerateResponse{Qr: make([]*pb.GenerateResponseRow, size)}
	for i, row := range qr {
		res.Qr[i] = &pb.GenerateResponseRow{V: row}
	}
	return &res, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterQRGenerateServiceServer(grpcServer, &server{})
	log.Printf("Listening on %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
