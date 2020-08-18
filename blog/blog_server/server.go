package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yurianxdev/grpc-course/blog/blogpb"
)

type server struct{}

func (s server) CreateBlog(_ context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	log.Println("CreateBlog RPC called...")
	data := req.GetBlog()
	blog := Blog{
		AuthorID: data.GetAuthorId(),
		Title:    data.GetTitle(),
		Content:  data.GetContent(),
	}

	result, err := collection.InsertOne(context.Background(), blog)
	if err != nil {
		log.Printf("Error inserting blog on collection: %v\n", err)
		// Return error throw gRPC
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error getting OID: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot conver to OID: %v", err))
	}

	log.Printf("Blog created: %v\n", oid.Hex())
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.AuthorID,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}, nil
}

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
