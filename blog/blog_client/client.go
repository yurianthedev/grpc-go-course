package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/yurianxdev/grpc-course/blog/blogpb"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Starting rpc client...")
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed dialing: %v", err)
	}
	defer conn.Close()
	log.Println("Successfully dial with RPC server on port 50051")

	c := blogpb.NewBlogServiceClient(conn)
	createBlog(c)
}

func createBlog(c blogpb.BlogServiceClient) {
	log.Println("Calling CreateBlog RPC...")

	blog := &blogpb.Blog{
		AuthorId: "Julian",
		Title:    "My first blog",
		Content:  "Some content",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Printf("Error at CreateBlog RPC: %v\n", err)
	}

	fmt.Printf("Blog created: %v\n", res.GetBlog())
}
