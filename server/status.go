package server

import "fmt"

// Status ...
type Status int

// Status list
const (
	StateUndefined Status = iota
	StateInitialized
	StateStarting
	StateStarted
	StateRunning
	StateStopping
	StateStopped
	StateError
)

const statusListing = "StateUndefined" +
	"StateInitialized" +
	"StateStarting" +
	"StateStarted" +
	"StateRunning" +
	"StateStopping" +
	"StateStopped" +
	"StateError"

var statusIndex = [...]uint8{0, 14, 30, 43, 55, 67, 80, 92, 102}

func (i Status) String() string {
	if i < 0 || i >= Status(len(statusIndex)-1) {
		return fmt.Sprintf("Status(%d)", i)
	}
	return statusListing[statusIndex[i]:statusIndex[i+1]]
}
