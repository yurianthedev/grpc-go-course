package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/yurianxdev/grpc-course/greet/greetpb"
)

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed dialing: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do an Unary RPC...\n")
	req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Julian"}}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Printf("Error while calling GreetRPC: %v", err)
	}

	fmt.Printf("Response from greet: %v", res.Result)
}
