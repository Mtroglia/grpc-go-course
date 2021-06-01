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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

//This is the server object. Add services for the implementation type
type server struct {
	greetpb.UnimplementedGreetServiceServer
}

//implement the interface for GreetServiceServer
//function on pointer to server
func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	//to implement return  hello and first name
	//the GreetWithDeadlineRequest struct "req" has a Greeting struct. Greeting struct holds first name and last name
	fmt.Printf("GreetWithDeadline function was invoked with %v  \n", req)

	for i := 0; i < 3; i++ {
		fmt.Println("Sleeping ....")
		if ctx.Err() == context.Canceled {
			fmt.Println("Client canceled request")
			return nil, status.Error(codes.Canceled, "Client canceled request ")
		}
		time.Sleep(1 * time.Second)

	}
	firstName := req.GetGreeting().GetFirstName()
	//form a GreetWithDeadlineResponse
	result := "Hello " + firstName
	//include & to dereference the pointer to the GreetWithDeadlineResponse
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil

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
	opts := []grpc.ServerOption{}

	// gRPC ssl
	//https://grpc.io/docs/guides/auth/
	tls := false

	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed to loading certificates to: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	//greetpb is the full path of the import statement at the top of script
	greetpb.RegisterGreetServiceServer(s, &server{})

	// To use reflection service, install gRPC cleint i.e. evans
	//https://github.com/ktr0731/evans
	// Once installed, connect gRPC client CLI to server
	// e.g. evans -p 50051 -r
	// Register reflection service on gRPC server

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
