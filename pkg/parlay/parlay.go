package parlay

import "github.com/plunder-app/plunder/pkg/parlay/types"

type actionType string

const (
	//upload - defines that this action will upload a file to a remote system
	upload   actionType = "upload" //
	download actionType = "download"
	command  actionType = "command"
	pkg      actionType = "package"
)

// TreasureMap - X Marks the spot
// The treasure maps define the automation that will take place on the hosts defined
type TreasureMap struct {
	// An array/list of deployments that will take places as part of this "map"
	Deployments []Deployment `json:"deployments"`
}

// Deployment defines the hosts and the action(s) that should be performed on them
type Deployment struct {
	// Name of the deployment that is taking place i.e. (Install MySQL)
	Name string `json:"name"`
	// An array/list of hosts that these actions should be performed upon
	Hosts []string `json:"hosts"`

	// Parallel allow multiple actions across multiple hosts in parallel
	Parallel         bool `json:"parallel"`
	ParallelSessions int  `json:"parallelSessions"`

	// The actions that should be performed
	Actions []types.Action `json:"actions"`
}

// KeyMap

// Keys are used to store information between sessions and deployments
var Keys map[string]string

func init() {
	// Initialise the map
	Keys = make(map[string]string)
}
