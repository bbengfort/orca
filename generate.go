package orca

import (
	"database/sql"
	"log"
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
	device := new(Device)
	if err := device.GetByName("apollo", app.db); err != nil {
		return err
	}

	addr, err := ResolveAddr(device.IPAddr)
	if err != nil {
		return err
	}
	device.IPAddr = addr

	if err := app.Ping(device); err != nil {
		return err
	}

	return nil
}

// Ping sends an echo request to a device and handles the response
func (app *App) Ping(device *Device) error {
	// Send the ping out and get a reply (blocking)
	reply, err := app.SendPing(device)
	if err != nil {
		return err
	}

	// Store the recv timestamp before any work.
	recv := time.Now()

	// Log the echo reply
	if app.Config.Debug {
		log.Println(reply.LogRecord())
	}

	// Fetch the Ping from the database
	ping := new(Ping)
	echo := reply.GetEcho()
	if err := ping.Get(echo.Ping, app.db); err != nil {
		return err
	}

	// Update the ping information
	ping.Recv = recv
	ping.Response = reply.Sequence

	// Compute the latency and set it
	msecs := float64(recv.Sub(echo.GetSentTime()).Seconds()) * 1000.0
	ping.Latency = sql.NullFloat64{Float64: msecs, Valid: true}

	// Save the ping to the database
	ping.Save(app.db)

	return nil
}

// SendPing sends an echo request to another device
func (app *App) SendPing(device *Device) (*echo.Reply, error) {

	// Connect to the remote node
	conn, err := grpc.Dial(device.IPAddr, grpc.WithInsecure(), grpc.WithTimeout(Timeout))
	if err != nil {
		return nil, err
	}

	// Defer closing the connection and create an Echo client.
	defer conn.Close()
	client := echo.NewOrcaClient(conn)

	// Refresh the current location of the source
	// NOTE: location errors are ignored
	app.SyncLocation()

	// Create a Ping record for experimental metrics
	ping := new(Ping)

	// Set the source with the current information
	ping.Source = app.GetDevice()
	ping.Source.IPAddr = app.ExternalIP
	ping.Location = app.Location

	// Set the target as the passed in device and increment the sequence
	ping.Target = device
	ping.Target.Sequence++
	ping.Target.Save(app.db)

	// Set the sent timestamp - note this is not the same as the timestamp
	// in the echo.Request - in order to eliminate database access latency.
	// The latency saved on the record should be computed from the request.
	ping.Sent = time.Now()
	ping.Request = ping.Target.Sequence

	// Save the ping to the database (updates the ID)
	if _, err := ping.Save(app.db); err != nil {
		return nil, err
	}

	// Create an EchoRequest to send to the node
	request := &echo.Request{
		Sequence: ping.Target.Sequence,
		Sender:   ping.Source.Echo(),
		Sent:     &echo.Time{Nanoseconds: time.Now().UnixNano()},
		TTL:      int64(Timeout.Seconds()),
		Ping:     ping.ID,
		Payload:  []byte("Clutter to be replaced with random or actual data."),
	}

	// Send the Echo request to the remote reflector and return
	return client.Echo(context.Background(), request)
}
