package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Just example, do not do it with insecure credentials on production environments
	var host string = "localhost:8085"
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server at %v. Error: %v", host, err)
	}

	defer conn.Close()
	fmt.Printf("Server running at http://%v \n", host)

}
