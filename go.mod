module github.com/mtroglia/grpc-go-course

go 1.16

replace github.com/mtroglia/grpc-go-course/greet/greetpb => /Users/trog/go/src/grpc-go-course/greet/greetpb

require (
	go.mongodb.org/mongo-driver v1.5.2 // indirect
	google.golang.org/grpc v1.37.1
	google.golang.org/protobuf v1.26.0
)
