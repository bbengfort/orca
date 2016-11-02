package orca

import (
	"fmt"
	"time"
)

// Parse a Unix timestamp from an echo.Time message.
func (ts *Time) Parse() time.Time {
	if ts != nil {
		secs := ts.Seconds
		nsecs := ts.Nanoseconds
		return time.Unix(secs, nsecs)
	}
	return time.Time{}
}

// GetSentTime parses the sent time on an Echo message to a time.Time
func (m *Request) GetSentTime() time.Time {
	ts := m.GetSent()
	return ts.Parse()
}

// LogRecord returns the echo request as a string in loggable format.
func (m *Request) LogRecord() string {

	var sender string
	if device := m.GetSender(); device != nil {
		sender = fmt.Sprintf("%s (%s)", device.Name, device.IPAddr)
	} else {
		sender = "Unknown (N/A)"
	}

	delta := time.Now().Sub(m.GetSentTime())

	output := "Echo %d bytes from %s seq=%d ttl=%d time=%s"
	return fmt.Sprintf(output, len(m.Payload), sender, m.Sequence, m.TTL, delta)
}

// GetReceivedTime parses the received time on an Reply message to a time.Time
func (m *Reply) GetReceivedTime() time.Time {
	ts := m.GetReceived()
	return ts.Parse()
}

// LogRecord returns the echo reply as a string in loggable format.
func (m *Reply) LogRecord() string {

	var remote string
	if device := m.GetReceiver(); device != nil {
		remote = fmt.Sprintf("%s (%s)", device.Name, device.IPAddr)
	} else {
		remote = "Unknown (N/A)"
	}

	echo := m.GetEcho()
	delta := time.Now().Sub(echo.GetSentTime())

	output := "received %d bytes from %s order=%d seq=%d time=%s"
	return fmt.Sprintf(output, len(echo.Payload), remote, echo.Sequence, m.Sequence, delta)
}
