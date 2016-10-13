package echo

import "time"

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
func (m *EchoRequest) GetSentTime() time.Time {
	ts := m.GetSent()
	return ts.Parse()
}

// GetReceivedTime parses the received time on an EchoReply message to a time.Time
func (m *EchoReply) GetReceivedTime() time.Time {
	ts := m.GetReceived()
	return ts.Parse()
}
