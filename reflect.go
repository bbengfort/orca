package orca

import (
	"log"
	"net"
	"time"

	"github.com/bbengfort/orca/echo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Echo implements the echo.EchoServer interface on the App
func (app *App) Echo(ctx context.Context, in *echo.Request) (*echo.Reply, error) {

	// Log the echo request
	if app.Config.Debug {
		log.Println(in.LogRecord())
	}

	// Return the Reply
	return &echo.Reply{
		Receiver: app.GetDevice().Echo(),
		Received: &echo.Time{Nanoseconds: time.Now().UnixNano()},
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
	if app.Config.Debug {
		log.Printf("Listening for Echo Requests on %s\n", addr)
	}

	// Create the grpc server, handler, and listen
	server := grpc.NewServer()
	echo.RegisterOrcaServer(server, app)
	server.Serve(sock)

	// Serve until finished
	return nil
}
