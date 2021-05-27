package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/mtroglia/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Client Calculator ")
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	defer cc.Close()
	//create new service client with the dialed connection
	c := calcpb.NewCalculatorServiceClient(cc) // NewSummingServiceClient(cc)

	//doUnary(c)
	doUnarySqrt(c)
	//doStreamingServer(c)

	//doStreamingClient(c)

	//doBidirectionalStreaming(c)

}

func doBidirectionalStreaming(c calcpb.CalculatorServiceClient) {
	fmt.Println("Strating FIND MAX BiDi stremaing RPC")
	requests := []int32{2, 5, 8, 4, 9, 5, 7, 6, 10}
	stream, err := c.FindMax(context.Background())

	if err != nil {
		log.Fatalf("Error setting up client stream ")
	}
	// set up a blocking function
	waitChan := make(chan struct{})
	// set up go function routine for sending messages to server
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending req val::%v to server \n", req)
			stream.Send(&calcpb.FindMaxRequest{
				Number: req,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// setup go fucntion routine for receiving messages from server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Printf("Sever is done sending stream, reached EOF \n")
				break
			}
			if err != nil {
				log.Fatalf("Error receiving stream from server :: %v\n", err)
				break
			}
			fmt.Printf(" Got new max value from Server. Value is::: %v \n", res.GetMaxVal())

		}
		close(waitChan)
	}()
	//this is the blocking function waiting for the server to unblock the go routines
	<-waitChan
}

func doStreamingClient(c calcpb.CalculatorServiceClient) {
	fmt.Println("Streaming Client RPC for Averaging ")
	requests := []*calcpb.ComputeAverageRequest{
		{
			Values: 8,
		},
		{
			Values: 12,
		},
		{
			Values: 7,
		},
		{
			Values: 6,
		},
	}
	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error calling ComputeAverage client streamer :: %v \n", err)
	}

	for iter, req := range requests {
		fmt.Printf("Iter:%v \n", iter)
		fmt.Printf("Req:%v \n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Problem, error while recieving repsonse from ComputeAverage %v \n", err)
	}
	fmt.Printf("ComputeAverage Response from server: %v \n", res)

}

func doStreamingServer(c calcpb.CalculatorServiceClient) {
	fmt.Println("Server Streaming  RPC for Prime Number Decomp")
	req := &calcpb.PrimeNumberDecompositionRequest{
		Number: 70,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error creating request for  PrimeDecomposition RPC:: %v ", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			//reached end of the stream. Server sends EOF, break from for loop
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream:: %v ", err)

		}
		// if no errors, get the results from res.
		log.Printf("Response from PrimeFactorDecomposition: %v", res.GetPrimeFactor())
	}

}

func doUnary(c calcpb.CalculatorServiceClient) {
	fmt.Println("Unary RPC for Summing")

	//save summing request object with values defined into req
	req := &calcpb.SummingRequest{
		Summing: &calcpb.Summing{
			NumOne: 6,
			NumTwo: 9,
		},
	}

	//use the SummingServiceClient
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling the Sum RPC: %v ", err)
	}

	log.Printf("Response from Sum::::: %v ", res.Result)

}
func doUnarySqrt(c calcpb.CalculatorServiceClient) {
	fmt.Println("Unary RPC for SquareRoot")

	// Sqrt request object with values defined into req
	requests := []*calcpb.SquareRootRequest{
		{
			Number: 9,
		},
		{
			Number: -3,
		},
	}

	for _, req := range requests {
		//use the SquareRootServiceClient
		res, err := c.SquareRoot(context.Background(), req)
		if err != nil {
			respErr, ok := status.FromError(err)
			if ok {
				//actual error from gRPC (user error)
				fmt.Printf("Error message from server %v \n ", respErr.Message())
				fmt.Println(respErr.Code())
				if respErr.Code() == codes.InvalidArgument {
					fmt.Println("The client must have send a negative number ")
					return
				}

			} else {
				log.Fatalf("Major Error while calling the SquareRoot RPC: %v ", err)
				return
			}

		}
		log.Printf("Response from SquareRoot of %v ::::: %v \n", req.GetNumber(), res.GetNumberRoot())
		time.Sleep(1000 * time.Microsecond)
	}

}
