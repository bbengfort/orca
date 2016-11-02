package orca

import (
	"errors"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Interval is the amount of time between pings
const Interval = time.Duration(1) * time.Second

// Timeout is the amount of time sonar will wait for a reply
const Timeout = time.Duration(30) * time.Second

var sequence int64

// Generate is long running function that initializes pings then sleeps.
func (app *App) Generate(addr string, count int) error {

	// Loop forever with a delay between the interval
	for i := 0; i < count; i++ {

		// Wait for the specified interval
		select {
		case <-time.After(Interval):
			// This breaks out of select not the for loop
			break
		}

		if err := app.Ping(addr); err != nil {
			return err
		}

	}

	// Compute the statistics
	app.ComputeStats()
	return nil
}

// Ping sends an echo request to a device and handles the response
func (app *App) Ping(addr string) error {
	// Send the ping out and get a reply (blocking)
	reply, err := app.SendPing(addr)
	if err != nil {
		return err
	}

	// Get the echo request from the payload
	echo := reply.GetEcho()
	if echo == nil {
		return errors.New("No echo payload to compute latency!")
	}

	// Compute the latency and append to latencies
	recv := time.Now()
	msecs := float64(recv.Sub(echo.GetSentTime()).Seconds()) * 1000.0
	app.Latencies = append(app.Latencies, msecs)

	// Log the echo reply
	if !app.Silent {
		log.Println(reply.LogRecord())
	}

	return nil
}

// SendPing sends an echo request to another device
func (app *App) SendPing(addr string) (*Reply, error) {

	// Connect to the remote node
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(Timeout))
	if err != nil {
		return nil, err
	}

	// Defer closing the connection and create an Echo client.
	defer conn.Close()
	client := NewOrcaClient(conn)

	sequence++

	// Create an EchoRequest to send to the node
	request := &Request{
		Sequence: sequence,
		Sender:   &Device{Name: "orca", IPAddr: app.IPAddr},
		Sent:     &Time{Nanoseconds: time.Now().UnixNano()},
		TTL:      int64(Timeout.Seconds()),
		Ping:     sequence,
		Payload:  []byte("Clutter to be replaced with random or actual data."),
	}

	// Send the Echo request to the remote reflector and return
	return client.Echo(context.Background(), request)
}
