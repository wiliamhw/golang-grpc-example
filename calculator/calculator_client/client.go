package main

import (
	"context"
	"fmt"
	"github.com/frain8/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
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

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//fmt.Printf("Created client: %f\n", c)

	doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber:  25,
		SecondNumber: 5,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculator RPC: %v", err)
	}
	log.Printf("Response from Calculator: %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 124538982,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculator Server Streaming RPC: %v", err)
	}
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream %v", err)
		}
		log.Printf("Response from PrimeNumberDecomposition: %v", res.GetPrimeNumber())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v", err)
	}

	numbers := [...]int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		req := &calculatorpb.ComputeAverageRequest{
			Number: number,
		}
		err := stream.Send(req)
		if err != nil {
			fmt.Printf("Error on sending request: %v\n", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("The average is: %v\n", res.GetAverage())
}
