package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/chunlinwang/grpc-demo/notification"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedNotificationServer
	client pb.NotificationClient
	cc     *grpc.ClientConn
}

func (s *server) UnaryNotify(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	fmt.Printf("UnaryNotify request: %v \n", req)

	return &pb.NotificationResponse{Content: req.Content, RequestId: req.RequestId}, nil
}

func (s *server) BidirectionalStreamingNotify(stream pb.Notification_BidirectionalStreamingNotifyServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return status.Error(codes.InvalidArgument, "request message not received")
		}
		if err != nil {
			return err
		}

		if req.GetPropagate() {
			res, err := s.client.UnaryNotify(stream.Context(), &pb.NotificationRequest{Content: req.GetContent(), RequestId: req.GetRequestId()})
			if err != nil {
				return err
			}

			stream.Send(res)
		}

		fmt.Printf("BidirectionalStreamingNotify request: %v, #%v.\n", req.GetContent(), req.GetRequestId())

		stream.Send(&pb.NotificationResponse{Content: req.GetContent(), RequestId: req.GetRequestId() * 2})
	}
}

func (s *server) Close() {
	s.cc.Close()
}

func newNotificationServer() *server {
	target := fmt.Sprintf("localhost:%v", *port)
	cc, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return &server{client: pb.NewNotificationClient(cc), cc: cc}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	notificationServer := newNotificationServer()
	defer notificationServer.Close()

	s := grpc.NewServer()
	pb.RegisterNotificationServer(s, notificationServer)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
