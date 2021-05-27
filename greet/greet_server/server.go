package main

import (
	//"github.com/mtroglia/grpc-go-course/greet/greetpb"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/mtroglia/grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

//This is the server object. Add services for the implementation type
type server struct {
	greetpb.UnimplementedGreetServiceServer
}

//implement the interface for GreetServiceServer
//function on pointer to server
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	//to implement return  hello and first name
	//the GreetRequest struct "req" has a Greeting struct. Greeting struct holds first name and last name
	fmt.Printf("Greet function was invoked with %v  \n", req)
	firstName := req.GetGreeting().GetFirstName()
	//form a GreetResponse
	result := "Hello " + firstName
	//include & to dereference the pointer to the GreetResponse
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil

}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes Function invoked::: %v", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		// stream object to use from the function type "stream greetpb.GreetService_GreetManyTimesServer"
		// stream Send will update the GreetService_GreetManyTimesServer object
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

//Look for server interface
// LongGreet(GreetService_LongGreetServer) error
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	//This server takes a stream
	fmt.Printf("LongGreet server function was invoked with streaming \n")
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("We have finished reading the client streaming")
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error reading client requeset ::: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "

	}

}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("Greet Everyone  server function was invoked with BiDi streaming \n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Print("Done got and EOF from client  \n")
			return err
		}
		if err != nil {
			log.Fatalf("Error recieving from client :: %v \n", err)
			return err
		}
		fName := req.GetGreeting().GetFirstName()
		result := "Hello " + fName

		sendErr := stream.Send(&greetpb.GreetEveroneResponse{
			Result: result,
		})

		if sendErr != nil {
			log.Fatalf("Error sending result")
			return err
		}
	}

}

func main() {
	fmt.Println("Starting Main() Server")
	// default port for gRPC 50051
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen to: %v", err)

	}

	s := grpc.NewServer()
	//greetpb is the full path of the import statement at the top of script
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
