package orca

import (
	"fmt"
	"time"

	"github.com/bbengfort/orca/echo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Timeout is the amount of time sonar will wait for a reply
const Timeout = time.Duration(30) * time.Second

// Generate is long running function that initializes pings then sleeps.
func (app *App) Generate() error {

	// Temporarily just ping the local machine
	addr, err := ResolveAddr("")
	device := &Device{
		Name:   "apollo",
		IPAddr: addr,
	}

	reply, err := Ping(device)
	if err != nil {
		return err
	}

	fmt.Println(reply)
	return nil
}

// Ping sends an echo request to another device
func Ping(device *Device) (*echo.Reply, error) {

	// Connect to the remote node
	conn, err := grpc.Dial(device.IPAddr, grpc.WithInsecure(), grpc.WithTimeout(Timeout))
	if err != nil {
		return nil, err
	}

	// Defer closing the connection and create an Echo client.
	defer conn.Close()
	client := echo.NewOrcaClient(conn)

	// Create an EchoRequest to send to the node
	request := &echo.Request{
		Sequence: 1,
		Sender:   nil,
		Sent:     &echo.Time{Nanoseconds: time.Now().UnixNano()},
		TTL:      30,
		Payload:  []byte("Clutter to be replaced with random or actual data."),
	}

	// Send the Echo request to the remote reflector and return
	return client.Echo(context.Background(), request)
}
