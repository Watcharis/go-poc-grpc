package main

import (
	"context"
	"log"
	"time"
	"watcharis/go-poc-grpc/go-poc-grpc/pb-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {

	// Set up a connection to the server.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	creds := insecure.NewCredentials().Clone()

	// Create a new gRPC client connection (ต่อไปยังที่เดียวกับ gRPC server)
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer conn.Close()

	userClient := proto.NewUserServiceClient(conn)

	// Contact the server and print out its response.
	// name := "Mike"
	personRequest := proto.PersonRequest{Id: 2}

	r, err := userClient.GetUser(ctx, &personRequest)
	if err != nil {
		log.Fatalf("could not getUser: %v", err)

		st, ok := status.FromError(err)
		if ok {
			// แสดง code และข้อความ error ของ gRPC
			log.Printf("gRPC error code: %v, error message: %v", st.Code(), st.Message())
		} else {
			// แสดงข้อความ error ที่ไม่คาดคิด
			log.Fatalf("Unexpected error: %v", err)
		}
		return

	}
	log.Printf("Greeting: %s", r.ProtoReflect().Descriptor())

}
