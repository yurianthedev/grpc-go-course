package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBidirectionalStreaming(c)
}

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting streaming...")

	stream, err := c.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v\n", err)
	}

	requests := []*greetpb.GreetEveryOneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Julian",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bryan",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Carlos",
			},
		},
	}

	wc := make(chan struct{})

	go func() {
		for _, req := range requests {
			log.Printf("Sending: %v\n", req)
			err := stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
			if err != nil {
				log.Printf("Error sending request: %v\n", err)
			}
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
				log.Printf("Error while reciving response: %v\n", err)
				break
			}
			fmt.Printf("Recived: %v\n", res.GetResponse())
		}
		close(wc)
	}()

	<-wc
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting streaming...")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bryan",
				LastName:  "Garzon",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Juanes",
				LastName:  "Martinez",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Nicolas",
				LastName:  "Vera",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Printf("Error creating streaming: %v\n", err)
	}

	for _, req := range requests {
		log.Printf("Sending request: %v...\n", req)
		err := stream.Send(req)
		if err != nil {
			log.Printf("There was an error sending stream: %v", err)
		}
		time.Sleep(time.Millisecond * 1000)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Error recieving streams: %v", err)
	}

	fmt.Printf("LongGreet response: %v\n", res)
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
