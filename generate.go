package orca

import (
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// Timeout is the amount of time sonar will wait for a reply
const Timeout = time.Duration(30) * time.Second

// Generate is long running function that initializes pings then sleeps.
func (app *App) Generate() error {
	return nil
}

// Ping sends an echo request to another device
func Ping(device *Device) (*EchoReply, error) {

	// Connect to the remote node
	conn, err := grpc.Dial(device.IPAddr, grpc.WithInsecure(), grpc.WithTimeout(Timeout))
	if err != nil {
		return nil, err
	}

	// Defer closing the connection and create an Echo client.
	defer conn.Close()
	client := NewEchoServiceClient(conn)

	// Create an EchoRequest to send to the node
	request := &EchoRequest{
		Sequence: 1,
		Sender:   nil,
		Sent:     &Time{Nanoseconds: time.Now().UnixNano()},
		TTL:      30,
		Payload:  []byte("Clutter to be replaced with random or actual data."),
	}

	// Send the Echo request to the remote reflector and return
	return client.Echo(context.Background(), request)
}
