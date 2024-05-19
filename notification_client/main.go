package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/chunlinwang/grpc-demo/notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func unaryCall(c pb.NotificationClient, requestId uint64, message string, propagate bool) {
	// Creates a context with a one second deadline for the RPC.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.NotificationRequest{Content: message, RequestId: requestId, Propagate: propagate}

	r, err := c.UnaryNotify(ctx, req)

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	statusCode := status.Code(err)
	fmt.Printf("#%v Content = [%v], statusCode = %v\n", r.GetRequestId(), r.GetContent(), statusCode)
}

func streamingCall(c pb.NotificationClient, requestId uint64, message string, propagate bool) {
	// Creates a context with a one second deadline for the RPC.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.BidirectionalStreamingNotify(ctx)
	if err != nil {
		log.Printf("Stream err: %v", err)
		return
	}

	err = stream.Send(&pb.NotificationRequest{Content: message, RequestId: requestId, Propagate: propagate})
	if err != nil {
		log.Printf("Send error: %v", err)
		return
	}

	r, err := stream.Recv()

	statusCode := status.Code(err)
	fmt.Printf("#%v Content = [%v], statusCode = %v\n", r.GetRequestId(), r.GetContent(), statusCode)
}

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("cannt connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewNotificationClient(conn)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()
	var requestId uint64 = 1
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			//unaryCall(c, requestId, fmt.Sprintf("Current time: ", t), false)
			streamingCall(c, requestId, fmt.Sprintf("Current time: ", t), true)
		}

		requestId++
	}
}
