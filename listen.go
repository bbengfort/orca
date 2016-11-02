package orca

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Echo implements the echo.EchoServer interface on the App
func (app *App) Echo(ctx context.Context, in *Request) (*Reply, error) {

	// Store the RECV timestamp before any work
	recv := time.Now()

	// Log the echo request
	if !app.Silent {
		log.Println(in.LogRecord())
	}

	// Return the Reply
	return &Reply{
		Sequence: 0,
		Receiver: nil,
		Received: &Time{Nanoseconds: recv.UnixNano()},
		Echo:     in,
	}, nil

}

// Reflect listens for EchoRequests and Replies to them.
func (app *App) Reflect() error {
	// Look up the address to listen on
	addr, err := app.GetListenAddr()
	if err != nil {
		return err
	}

	// Create the socket to listen on
	sock, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// Log the fact that we are listening on the address
	if !app.Silent {
		log.Printf("Listening for Echo Requests on %s\n", addr)
	}

	// Create the grpc server, handler, and listen
	server := grpc.NewServer()
	RegisterOrcaServer(server, app)
	server.Serve(sock)

	// Serve until finished
	return nil
}
