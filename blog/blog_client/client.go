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
	readBlog(c, "somerandomid")
	readBlog(c, "5f3b6340b0039109d78d521e")
}

func readBlog(c blogpb.BlogServiceClient, id string) {
	log.Println("Calling ReadBlog RPC...")

	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: id,
	})
	if err != nil {
		log.Printf("Error reading blog: %v\n", err)
		return
	}
	fmt.Printf("Blog found: %v", res.GetBlog())
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
		return
	}

	fmt.Printf("Blog created: %v\n", res.GetBlog())
}
