package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBidirectionalStreaming(c)
}

func doBidirectionalStreaming(c calculatorpb.CalculatorClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Printf("Error creating streaming: %v\n", err)
	}

	requests := []*calculatorpb.FindMaximumRequest{
		{
			Number: 1,
		},
		{
			Number: 5,
		},
		{
			Number: 3,
		},
		{
			Number: 6,
		},
		{
			Number: 2,
		},
		{
			Number: 20,
		},
	}

	wg := make(chan struct{})

	go func() {
		for _, req := range requests {
			log.Printf("Sending %v...\n", err)
			err := stream.Send(req)
			if err != nil {
				log.Printf("Error sending request: %v", err)
			}

			time.Sleep(1000 * time.Millisecond)
		}

		_ = stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Erorr reciving response: \n%v\n", err)
			}

			fmt.Printf("The current maximum number is: %d\n", res.GetMaximum())
		}
		close(wg)
	}()

	<-wg
}

func doClientStreaming(c calculatorpb.CalculatorClient) {
	log.Printf("Starting streaming...\n")

	numbers := []int32{1, 3, 3, 5, 2}
	stream, err := c.ComputeAverage(context.Background())
	for err != nil {
		log.Printf("Error creating stream: %v\n", err)
		return
	}

	for _, val := range numbers {
		log.Printf("Sending: %d\n", val)
		err := stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: val,
		})

		time.Sleep(time.Millisecond * 1000)

		if err != nil {
			log.Printf("Error sending request: %v\n", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Error reciving the response from server: %v", err)
	}

	fmt.Printf("Response: %v", res)
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
