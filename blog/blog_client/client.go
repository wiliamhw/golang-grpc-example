package main

import (
	"context"
	"fmt"
	"github.com/wiliamhw/golang-grpc-example/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
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

	blog := createBlog(c)
	readBlog(c, "1dfsoijfs")
	readBlog(c, blog.GetId())
	updateBlog(c, blog.GetId())
	deleteBlog(c, blog.GetId())
	listBlog(c)
}

func createBlog(c blogpb.BlogServiceClient) *blogpb.Blog {
	fmt.Println("\nCreating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v\n", res)
	return res.GetBlog()
}

func readBlog(c blogpb.BlogServiceClient, blogId string) {
	fmt.Printf("\nReading the blog with id: %v\n", blogId)

	req := &blogpb.ReadBlogRequest{
		BlogId: blogId,
	}

	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}
	fmt.Printf("Blog was read: %v\n", res)
}

func updateBlog(c blogpb.BlogServiceClient, blogId string) *blogpb.Blog {
	fmt.Println("\nUpdating the blog")
	newBlog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "ChangeAuthor",
		Title:    "title-test",
		Content:  "content-test",
	}
	res, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been updated: %v\n", res)
	return res.GetBlog()
}

func deleteBlog(c blogpb.BlogServiceClient, blogId string) {
	fmt.Printf("\nDeleting the blog with id: %v\n", blogId)

	req := &blogpb.DeleteBlogRequest{
		BlogId: blogId,
	}

	res, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		fmt.Printf("Error happened while deleting: %v\n", err)
	}
	fmt.Printf("Blog was deleted: %v\n", res)
}

func listBlog(c blogpb.BlogServiceClient) {
	fmt.Println("\nListing all blogs")

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling ListBlog Streaming RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream %v", err)
		}
		fmt.Printf("Response from ListBlog: %v", res.GetBlog())
	}
}
