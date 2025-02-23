package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/vinicius77/coffee-shop-grpc/coffee_shop_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Just example, do not do it with insecure credentials on production environments
	var host string = "localhost:8085"
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at %v. %v \n", host, err)
	}

	defer conn.Close()
	log.Printf("server running at http://%v \n", host)
	c := pb.NewCoffeeShopClient(conn)

	// Create a new context that cancels after a 1-second timeout.
	// Useful for limiting the duration of operations that rely on this context.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Stream the menu
	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("error getting the menu stream %v \n", err)
	}

	done := make(chan bool)
	var items []*pb.Item

	go func() {
		for {
			res, err := menuStream.Recv()

			// end of file
			if err == io.EOF {
				done <- true
			}

			if err != nil {
				log.Fatalf("could not receive from menu stream: %v \n", err)
			}

			items = res.Items
			log.Printf("response from menu stream received: %v \n", items)
		}
	}()

	<-done

	receipt, err := c.PlaceOrder(ctx, &pb.Order{Items: items})
	if err != nil {
		log.Fatalf("could not receive the receipt %v \n", err)
	}

	status, err := c.GetOrderStatus(ctx, receipt)
	if err != nil {
		log.Fatalf("could not get order status %v \n", err)
	}

	log.Printf("order status: %v \n", status)

}
