package main

import (
	"fmt"
	"github.com/frain8/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const serverPort = 50051

func main() {
	fmt.Println("Hello, I'm a client")

	cc, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", serverPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer func(cc *grpc.ClientConn) {
		err := cc.Close()
		if err != nil {
			log.Fatalf("Could not close connection: %v", err)
		}
	}(cc)

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Created client: %f\n", c)
}
