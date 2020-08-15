package main

import (
	"context"
	"fmt"
	"io"
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

	doServerStreaming(c)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting streaming...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bryan",
			LastName:  "Rincon",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
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
		fmt.Printf("Response from GreetManyTimes: %s\n", msg.GetResult())
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	log.Println("Starting to do an Unary RPC...")
	req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Julian"}}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Printf("Error while calling GreetRPC: %v\n", err)
	}

	fmt.Printf("Response from greet: %v\n", res.Result)
}
