package main

import (
	"context"
	"fmt"
	"github.com/frain8/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Received Sum RPC: %v\n", req)
	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()
	sum := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(
	req *calculatorpb.PrimeNumberDecompositionRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer,
) error {
	fmt.Printf("Received PrimeNumberDecomposition RPC: %v\n", req)
	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor != 0 {
			divisor++
			fmt.Printf("Divisor has increased to %v\n", divisor)
			continue
		}
		res := &calculatorpb.PrimeNumberDecompositionResponse{
			PrimeNumber: divisor,
		}
		err := stream.Send(res)
		if err != nil {
			log.Printf("Erorr while sending data to client: %v\n", err)
			return err
		}
		number /= divisor
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage function was invoked with a streaming request")

	var sum int32 = 0
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break // Finish reading the client stream
		}
		if err != nil {
			log.Printf("Erorr while reading client stream: %v", err)
			return err
		}
		number := req.GetNumber()
		sum += number
		count++
	}
	average := float64(sum) / float64(count)
	err := stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
		Average: average,
	})
	if err != nil {
		log.Printf("Erorr while sending data to client: %v\n", err)
	}
	return err
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum function was invoked with a streaming request")
	var maximum int32 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Erorr while reading client stream: %v\n", err)
			return err
		}
		number := req.GetNumber()

		if number <= maximum {
			continue
		}
		maximum = number
		sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
			Maximum: maximum,
		})
		if sendErr != nil {
			log.Printf("Erorr while sending data to client: %v\n", sendErr)
			return sendErr
		}
	}
}

func (*server) SquareRoot(
	ctx context.Context,
	req *calculatorpb.SquareRootRequest,
) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

const port = 50051

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
