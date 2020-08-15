package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	"github.com/yurianxdev/grpc-course/calculator/calculatorpb"
)

func main() {
	log.Println("Starting client...")

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error dialing to server: %v\n", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorClient(conn)

	doUnary(c)

	doServerStreaming(c)
}

func doServerStreaming(c calculatorpb.CalculatorClient) {
	log.Printf("Starting streaming...\n")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Printf("Something failed in the streaming: %v\n", err)
		return
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			log.Printf("End of the stream\n")
			break
		}
		if err != nil {
			log.Printf("Error reading stream: %v\n", err)
		}
		fmt.Printf("Response from PrimeNumberDecomposition: %d\n", msg.GetPrimeNumber())
	}
}

func doUnary(c calculatorpb.CalculatorClient) {
	req := &calculatorpb.CalculatorRequest{NumberOne: 1, NumberTwo: 2}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error requesting for a Sum: %v\n", err)
	}

	fmt.Printf("The sum of those two numbers is: %d\n", res.Result)
}
