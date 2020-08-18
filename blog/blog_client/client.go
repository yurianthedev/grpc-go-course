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
	id := createBlog(c, &blogpb.Blog{
		AuthorId: "Julian",
		Title:    "My first blog",
		Content:  "Some content",
	})
	readBlog(c, "SomeRandomId") // This will fail
	readBlog(c, id)             // Found blog created just before
	updateBlog(c, id)
	readBlog(c, id)
	deleteBlog(c, id) // Delete the only blog

	// Create some record for testing
	createBlog(c, &blogpb.Blog{
		AuthorId: "Julian",
		Title:    "Some blog",
		Content:  "Some content",
	})
	createBlog(c, &blogpb.Blog{
		AuthorId: "Bryan",
		Title:    "My third blog",
		Content:  "Some content",
	})
	createBlog(c, &blogpb.Blog{
		AuthorId: "Rincon",
		Title:    "Blog",
		Content:  "Some content",
	})
	createBlog(c, &blogpb.Blog{
		AuthorId: "Garzon",
		Title:    "Another blog",
		Content:  "Some content",
	})
}

func deleteBlog(c blogpb.BlogServiceClient, id string) {
	log.Println("Calling DeleteBlog RPC...")

	deleteRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: id,
	})
	if err != nil {
		log.Printf("Error deleting blog, %v\n", err)
		return
	}

	fmt.Printf("Blog deleted: %s\n", deleteRes.GetBlogId())
}

func updateBlog(c blogpb.BlogServiceClient, id string) {
	log.Println("Calling UpdateBlog RPC...")

	res, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       id,
			Title:    "Another title",
			AuthorId: "Another author",
			Content:  "AnotherContent",
		},
	})
	if err != nil {
		log.Printf("Error updating blog, %v\n", err)
		return
	}

	fmt.Printf("Updated Blog: %v\n", res.GetBlog())
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
	fmt.Printf("Blog found: %v\n", res.GetBlog())
}

func createBlog(c blogpb.BlogServiceClient, blog *blogpb.Blog) string {
	log.Println("Calling CreateBlog RPC...")

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})
	if err != nil {
		log.Printf("Error at CreateBlog RPC: %v\n", err)
		return ""
	}

	fmt.Printf("Blog created: %v\n", res.GetBlog())
	return res.GetBlog().GetId()
}
