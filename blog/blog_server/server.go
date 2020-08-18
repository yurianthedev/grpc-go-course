package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yurianxdev/grpc-course/blog/blogpb"
)

type server struct{}

func (s server) DeleteBlog(_ context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	log.Println("DeleteBlog RPC called...")

	blogId := req.GetBlogId()
	// Parse string to Mongo ObjectId
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		log.Printf("Error parsing id: %v\n", err)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error parsing id: %v", err))
	}

	filter := bson.M{"_id": oid} // Mongo formatted filter
	deleteRes, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Error deleting blog: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error deleting blog: %v", err))
	}
	if deleteRes.DeletedCount == 0 {
		log.Printf("Blog not found: %v\n", err)
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Blog not found: %v", err))
	}

	log.Printf("Removed %d elements: %s\n", deleteRes.DeletedCount, blogId)
	return &blogpb.DeleteBlogResponse{
		BlogId: blogId,
	}, nil
}

func (s server) UpdateBlog(_ context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	log.Println("UpdateBlog RPC called...")

	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.Id)
	if err != nil {
		log.Printf("Error parsing id: %v\n", err)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error parsing id: %v", err))
	}

	filter := bson.M{"_id": oid} // Mongo formatted filter
	blogObject := Blog{
		AuthorID: blog.AuthorId,
		Title:    blog.Title,
		Content:  blog.Content,
	}

	updateRes, err := collection.ReplaceOne(context.Background(), filter, &blogObject)
	if err != nil {
		log.Printf("Error updating blog: %v\n", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Error updating blog: %v", err))
	}
	if updateRes.ModifiedCount == 0 {
		log.Printf("Blog not found: %v\n", err)
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Blog not found: %v", err))
	}

	log.Printf("Updated %d elements: %v\n", updateRes.ModifiedCount, updateRes.UpsertedID)
	return &blogpb.UpdateBlogResponse{
		Blog: blog,
	}, nil
}

func (s server) ReadBlog(_ context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	log.Println("ReadBlog RPC called...")

	blogId := req.GetBlogId()
	// Parse string to Mongo ObjectId
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		log.Printf("Error parsing id: %v\n", err)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error parsing id: %v", err))
	}

	blog := &Blog{}              // Object model to parse in
	filter := bson.M{"_id": oid} // Mongo formatted filter

	dbResult := collection.FindOne(context.Background(), filter)
	// Decode response into Golang native object of type Blog
	if err := dbResult.Decode(blog); err != nil {
		log.Printf("Error finding blog: %v\n", err)
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Error finding blog: %v", err))
	}

	log.Printf("Blog found: %v\n", blog.ID.Hex())
	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       blog.ID.Hex(),
			AuthorId: blog.AuthorID,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}, nil
}

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
