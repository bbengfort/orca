package orca

// Device is an entity that represents nodes in the network that can be pinged.
// Device objects are stored in the devices table.
type Device struct {
	Name   string
	IPAddr string
}
