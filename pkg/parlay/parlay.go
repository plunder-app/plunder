package parlay

type actionType string

const (
	//upload - defines that this action will upload a file to a remote system
	upload   actionType = "upload" //
	download actionType = "download"
	command  actionType = "command"
	pkg      actionType = "package"
)

// KeyMap

// Keys are used to store information between sessions and deployments
var Keys map[string]string

func init() {
	// Initialise the map
	Keys = make(map[string]string)
}
