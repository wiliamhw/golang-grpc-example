package main

import (
	"context"
	"fmt"
	"github.com/wiliamhw/golang-grpc-example/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const serverPort = 50051

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.Dial(fmt.Sprintf("localhost:%d", serverPort), opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer func(cc *grpc.ClientConn) {
		err := cc.Close()
		if err != nil {
			log.Fatalf("Could not close connection: %v", err)
		}
	}(cc)

	c := blogpb.NewBlogServiceClient(cc)

	// Create blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", createBlogRes)
}
