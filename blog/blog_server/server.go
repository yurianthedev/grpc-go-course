package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"github.com/yurianxdev/grpc-course/blog/blogpb"
)

type server struct{}

type Blog struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

var collection *mongo.Collection

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := "50051"
	log.Println("Starting server...")
	// Listen tcp connections
	li, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error listening server: %v", err)
	}

	// MongoDB client
	log.Printf("Connecting to database client on port %s...", "27017")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed creating database client: %v\n", err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Failed connecting to database: %v\n", err)
	}

	// Create database or connection
	collection = client.Database("mydb").Collection("blog")

	// Create new server
	s := grpc.NewServer()
	// Append implementations of methods defined on
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		log.Printf("Server listening on port %s", port)
		// Accept connections
		if err := s.Serve(li); err != nil {
			log.Fatalf("Error serving: %v", err)
		}
	}()

	// Wait for control c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Exit gracefully
	<-ch
	log.Printf("Stopping the database client...\n")
	_ = client.Disconnect(context.TODO())
	log.Printf("Stopping the server...\n")
	s.Stop()
	log.Printf("Closing the listener...\n")
	li.Close()
	log.Printf("Server stopped\n")
}
