package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/mtroglia/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	//calcpb.UnimplementedSummingServiceServer
	calcpb.UnimplementedCalculatorServiceServer
}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {
	fmt.Printf("Received SquareRoot RPC::: %v", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Oooops, Received a negative number: %v", number),
		)
	}
	return &calcpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil

}

//implement the interface for CalculatorServiceServer
//Find this in calc.pb.go
func (*server) Sum(ctx context.Context, req *calcpb.SummingRequest) (*calcpb.SummingResponse, error) {
	fmt.Printf("Sum fuction is being implemented with: %v", req)
	firstNum := req.GetSumming().GetNumOne()
	secondNum := req.GetSumming().GetNumTwo()

	//for the SummingRepsonse
	result := firstNum + secondNum
	res := &calcpb.SummingResponse{
		Result: result,
	}
	return res, nil //res, nil
}

//	PrimeNumberDecomposition(*PrimeNumberDecompositionRequest, CalculatorService_PrimeNumberDecompositionServer) error
func (*server) PrimeNumberDecomposition(req *calcpb.PrimeNumberDecompositionRequest, stream calcpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf(" Invoking PrimeNumberDecomposition server Function RPC ::: %v\n", req)

	number := req.GetNumber()

	k := int64(2)
	for number > 1 {

		if number%k == 0 {
			stream.Send(&calcpb.PrimeNumberDecompositionResponse{
				PrimeFactor: k,
			})
			number = number / k

		} else {
			fmt.Printf("k:%v is not a divisor. Increased\n", k)
			k = k + 1

		}

	}

	return nil

}

// implement server interface
// FindMax(CalculatorService_FindMaxServer) error
func (*server) FindMax(stream calcpb.CalculatorService_FindMaxServer) error {
	fmt.Println(" FIND MAX RPC")
	maxVal := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Received EOF from client")

			return nil
		}
		if err != nil {
			fmt.Printf("Error receiving stream from client \n")
			return err
		}
		clientNumber := req.GetNumber()

		if maxVal < clientNumber {
			fmt.Printf("Max Value Update. Sending to Client::: %v \n", clientNumber)
			maxVal = clientNumber
			sendErr := stream.Send(&calcpb.FindMaxResponse{
				MaxVal: maxVal,
			})
			if sendErr != nil {
				log.Fatalf("Failed to send FindMaxResponse")
				return err
			}

		} else {
			fmt.Printf("Max Value Remains::: %v \n", maxVal)
		}

	}
}

func (*server) ComputeAverage(stream calcpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Server ComputeAverage Invoked. Waiting for client stream...")

	summingAvg := 0
	count := 0
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			fmt.Printf("We have finished recieving the client stream. %v \n", req)
			return stream.SendAndClose(&calcpb.ComputeAverageResponse{
				Average: float32(summingAvg) / float32(count),
			})

		}
		if err != nil {
			log.Fatalf("Problem with recieving client stream::: %v ", err)
		}
		summingAvg += int(req.GetValues())
		count++

	}

}

//Main function to set up TCP listener
func main() {
	fmt.Println("Running calculator service server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listin to address:port: %v", err)
	}

	//gRPC server
	s := grpc.NewServer()
	//calcpb.RegisterSummingServiceServer(s, &server{})
	calcpb.RegisterCalculatorServiceServer(s, &server{}) //SummingServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
