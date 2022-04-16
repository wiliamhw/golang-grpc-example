package main

import (
	"context"
	"fmt"
	"github.com/frain8/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

const serverPort = 50051
const tls = true

func main() {
	fmt.Println("Hello, I'm a client")

	creds := insecure.NewCredentials()
	if tls {
		certFile := "ssl/ca.crt" // Certificate Authority Trust
		credsBuffer, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
		}
		creds = credsBuffer
	}

	opts := grpc.WithTransportCredentials(creds)
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

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created client: %f\n", c)

	doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	//doUnaryWithDeadline(c, 5*time.Second) // should complete
	//doUnaryWithDeadline(c, 1*time.Second) // should timeout
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Marek",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet Unary RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Mareek",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet Server Streaming RPC: %v", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", res.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := [...]*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephane",
			},
		}, {
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		}, {
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		}, {
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
			},
		}, {
			Greeting: &greetpb.Greeting{
				FirstName: "Jonathan",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		err := stream.Send(req)
		if err != nil {
			fmt.Printf("Error on sending request: %v\n", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	requests := [...]*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephane",
			},
		}, {
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		}, {
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		}, {
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
			},
		}, {
			Greeting: &greetpb.Greeting{
				FirstName: "Jonathan",
			},
		},
	}

	// Sender
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			err := stream.Send(req)
			if err != nil {
				log.Printf("Error while sending data to stream: %v\n", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("Error while closing send stream: %v\n", err)
		}
	}()

	// Receiver
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error while receiving: %v\n", err)
			break
		}
		fmt.Printf("Received: %v\n", res.GetResult())
	}
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	log.Println("Starting to do a UnaryWithDeadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephane",
			LastName:  "Marek",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)

		// Golang/framework error
		if !ok {
			log.Fatalf("Big Error while calling SquareRoot RPC: %v", err)
		}

		// gRPC error
		if statusErr.Code() == codes.DeadlineExceeded {
			log.Fatalf("Timeout was hit! Deadline was exceeded")
		}
		log.Fatalf("Unexpected gRPC error: %v\n", statusErr)
	}
	log.Printf("Response from Greet: %v", res.Result)
}
