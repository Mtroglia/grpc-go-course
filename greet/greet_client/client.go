package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/mtroglia/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	//THis is creating the client on the port 50051
	fmt.Println("Hello Im a client")
	conn_c, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cound not connect %v", err)
	}
	//defer means to performe this at the end of the function
	defer conn_c.Close()
	//create a new service client with the connection conn_c
	c := greetpb.NewGreetServiceClient(conn_c)

	//fmt.Printf("created client %f", c)

	//This fucntion will run the Unary Greet API
	//doUnary(c)

	// client to request streaming
	//doServerStreaming(c)

	// client streaming to server
	//doClientStreaming(c)

	// Bidirectional Streaming
	doBidirectional(c)

}

func doBidirectional(c greetpb.GreetServiceClient) {

	requests := []*greetpb.GreetEveroneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Matt",
				LastName:  "Trogdor",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Konni",
				LastName:  "Wilson",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Susan",
				LastName:  "Trog",
			},
		},
	}

	// create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {

		log.Fatalf("Error creating stream %v\n", err)
		return
	}

	//use a wait channel for blocking down below..
	// make a channel of struct, which doesnt take any data.
	waitc := make(chan struct{})
	// send a bunch of messages to the server go routine
	go func() { //this runs in its own go routine
		//function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Client Sending message to server::: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive a bunch of messages from the server go routine
	go func() {
		//function to receive lots of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving  %v\n", err)
				break
			}
			fmt.Printf("Client Recv from server :::%v \n", res.GetResult())
		}
		//if we get an EOF from the server, we can now close the wait channel
		close(waitc)
		// closing the wait channel, unblocks everything

	}()

	// block until everything is done
	//waits for the channel to be closed
	<-waitc
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a client streaming RPC")

	// look at the api calls stream has
	//want to use Send and CloseAndRecieve
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Matt",
				LastName:  "Trogdor",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Konni",
				LastName:  "Wilson",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Susan",
				LastName:  "Trog",
			},
		},
	}
	// dont have to pass the request in client streaming
	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error calling LongGreet streamer :: %v", err)
	}
	// iterage over slice and send each message one by one
	for iter, req := range requests {
		fmt.Printf("Iterated over requests. req:%v \n", req)
		fmt.Printf("Iterated over requests. iter:%v \n", iter)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	// Now we need to close and resepce the server response.
	// check return of stream.CloseAndRecv() which is greetpb.GreetService_LongGreetClient.CloseAndRecv
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Problem, error while recieving repsonse from longgreet %v \n", err)

	}
	fmt.Printf("LongGreet Response: %v \n", res)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Server Streaminig RPC Client")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Matt",
			LastName:  "Troglia",
		},
	}
	// Observe what GreetManyTimes returns from a GreetServiceClient.
	// Returns a "(GreetService_GreetManyTimesClient, error)"
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC::: %v ", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//reached end of the stream. Server sends EOF, break from for loop
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream:: %v ", err)

		}
		// if no errors, get the results from msg.
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC")
	//c can now invoke Greet from the GreetServiceClient interface.
	//which was created in mains c := greetpb.NewGreetServiceClient(conn_c)
	//navigate to the greet.pb.go generated code and determine what the Greet fuction needs.
	//i.e. Greet(ctx context.Context, in *GreetRequest, opts ...grpc.CallOption)

	//for the context, we can use context.Background(). We need to create the *GreetRequest object...
	//constist of a greeting with FirstName and LN
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Matt",
			LastName:  "Trog",
		},
	}
	//This will return the GreetResponse "res" and error "err" from the server
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet::: %v ", res.Result)

}
