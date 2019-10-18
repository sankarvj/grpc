package main

import (
	"context"
	"log"
	"net"

	pb "github.com/sankarvj/grpc/pb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var counter int64

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedStatusServiceServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) CheckStatus(ctx context.Context, in *pb.StatusRequest) (*pb.Status, error) {
	counter++
	log.Printf("Received: %v for %d time", in.GetCheck(), counter)
	return &pb.Status{Health: "Health is good for " + in.GetCheck()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("main : API listening on %s", port)
	s := grpc.NewServer()
	pb.RegisterStatusServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
