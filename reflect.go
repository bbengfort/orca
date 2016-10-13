package orca

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/context"

	pb "github.com/bbengfort/orca/echo"
	"google.golang.org/grpc"
)

// Echo implements the echo.EchoServer interface on the App
func (app *App) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoReply, error) {

	// Return the EchoReply
	return &pb.EchoReply{
		Received: &pb.Time{
			Seconds:     0,
			Nanoseconds: time.Now().UnixNano(),
		},
		Echo: in,
	}, nil

}

// Reflect listens for EchoRequests and Replies to them.
func (app *App) Reflect() error {
	// Look up the address to listen on
	addr, err := app.GetAddr()
	if err != nil {
		return err
	}

	// Create the socket to listen on
	sock, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// Log the fact that we are listening on the address
	log.Printf("Listening for Echo Requests on %s\n", addr)

	// Create the grpc server, handler, and listen
	server := grpc.NewServer()
	pb.RegisterEchoServiceServer(server, app)
	server.Serve(sock)

	// Serve until finished
	return nil
}
