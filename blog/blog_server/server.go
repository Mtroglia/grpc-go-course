package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/mtroglia/grpc-go-course/blog/blogpb"
)

type server struct {
}

var collection *mongo.Collection

func main() {
	//if go code crahses, we get the file name and line number using the below
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Starting Blog Main() Server")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error connecting %v \n", err)
		return
	}
	// Start a collection which is like a table
	collection = client.Database("mydb").Collection("blog")

	// Start listenter on default port for gRPC 50051
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
	blogpb.RegisterBlogServiceServer(s, &server{})

	// To use reflection service, install gRPC cleint i.e. evans
	//https://github.com/ktr0731/evans
	// Once installed, connect gRPC client CLI to server
	// e.g. evans -p 50051 -r
	// Register reflection service on gRPC server

	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server Now...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("End of Program")
}
