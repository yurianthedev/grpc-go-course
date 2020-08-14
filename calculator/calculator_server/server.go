package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/yurianxdev/grpc-course/calculator/calculatorpb"
)

type server struct{}

func (s *server) Sum(_ context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Println("Request for sum accepted")
	result := req.NumberOne + req.NumberTwo
	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	log.Println("Starting server...")
	li, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error listening server: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServer(s, &server{})

	if err := s.Serve(li); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
