package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/yurianxdev/grpc-course/calculator/calculatorpb"
)

type server struct{}

func (s *server) ComputeAverage(averageServer calculatorpb.Calculator_ComputeAverageServer) error {
	log.Printf("Request for ComputeAverage was accepted...\n")

	var average float64
	var acc int32
	var counter int32
	for {
		res, err := averageServer.Recv()
		if err == io.EOF {
			return averageServer.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Printf("Something goes wrong while reciving client stream: %v", err)
		}

		counter++
		number := res.GetNumber()
		acc += number

		average += float64(acc) / float64(counter)
	}
}

func (s *server) Sum(_ context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Println("Request for sum accepted")
	result := req.NumberOne + req.NumberTwo
	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}

	return res, nil
}

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, decompositionServer calculatorpb.Calculator_PrimeNumberDecompositionServer) error {
	log.Printf("Request for PrimeDecomposition was accepted...\n")
	prime := int32(2)
	n := req.Number
	for n > 1 {
		if n%prime == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeNumber: prime,
			}
			n /= prime

			err := decompositionServer.Send(res)
			if err != nil {
				return err
			}
			time.Sleep(time.Millisecond * 1000)
		} else {
			prime++
		}
	}

	return nil
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
