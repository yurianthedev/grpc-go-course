package main

import (
	"context"
	"fmt"
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
}

func doUnary(c calculatorpb.CalculatorClient) {
	req := &calculatorpb.CalculatorRequest{NumberOne: 1, NumberTwo: 2}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error requesting for a Sum: %v", err)
	}

	fmt.Printf("The sum of those two numbers is: %d", res.Result)
}
