package main

import (
	"context"
	"fmt"
	"log"

	"github.com/moly-space/molylibs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = "localhost:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error dialing %v", err)
	}
	defer conn.Close()
	c := pb.NewGreetServiceClient(conn)

	doGreet(c)
}

func doGreet(c pb.GreetServiceClient) {
	log.Println("doGreet was invoked")
	res, err := c.Greet(context.Background(), &pb.GreetRequest{
		FirstName: "alex kim",
	})
	if err != nil {
		log.Fatalf("error::: %s", err)
	}

	fmt.Println("res:", res)

}
