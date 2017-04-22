package astilectron

// Event names
const (
	EventNameElectronLog    = "electron.log"
	EventNameElectronStop   = "electron.stop"
	EventNameProvisionStart = "provision.start"
	EventNameProvisionStop  = "provision.stop"
)

// Other constants
const (
	mainTargetID = 0
)

// Event represents a go-astilectron event
type Event struct {
	Name     string
	Payload  interface{}
	TargetID int
}
