package main

import (
	"context"
	"fmt"
	"github.com/frain8/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doErrorUnary(c)
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

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a FindMaximum BiDi Streaming RPC...")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream and calling FindMaximum: %v", err)
	}
	numbers := [...]int32{4, 7, 2, 19, 4, 6, 32}

	// Sender
	go func() {
		for _, number := range numbers {
			fmt.Printf("Sending number: %v\n", number)
			err := stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
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
			log.Fatalf("Error while receiving: %v\n", err)
		}
		fmt.Printf("Received a new maximum of...: %v\n", res.GetMaximum())
	}
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")

	// Correct call
	doErrorCall(c, 10)

	// Error Call
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})
	if err != nil {
		respErr, ok := status.FromError(err)

		// Default error
		if !ok {
			log.Fatalf("Big Error while calling SquareRoot RPC: %v", err)
		}

		// Custom error
		fmt.Printf("Error message from server: %v\n", respErr.Message())
		fmt.Println(respErr.Code())
		if respErr.Code() == codes.InvalidArgument {
			fmt.Println("We probably sent a negative number!")
			return
		}
	}
	fmt.Printf("Result of square root of %v: %v\n", number, res.GetNumberRoot())
}
