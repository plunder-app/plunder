package parlay

import (
	"fmt"
)

// ValidateAction will parse an action to ensure it is valid
func (action *Action) ValidateAction() error {
	switch action.ActionType {
	case "upload":
		// Validate the upload action
		if action.Source == "" {
			return fmt.Errorf("The Source field can not be blank")
		}

		if action.Destination == "" {
			return fmt.Errorf("The Destination field can not be blank")
		}
		return nil
	case "download":
		// Validate the download action
		if action.Source == "" {
			return fmt.Errorf("The Source field can not be blank")
		}

		if action.Destination == "" {
			return fmt.Errorf("The Destination field can not be blank")
		}
		return nil
	case "command":
		// Validate the Command action
		if action.Command == "" && action.KeyName == "" {
			return fmt.Errorf("Neither a command or a key has been specified to execute")
		}
		if action.Command != "" && action.KeyName != "" {
			return fmt.Errorf("Unable to use both a Command and a Command Key")
		}

		return nil
	case "pkg":
		// Validate the Package action
		if action.PkgManager == "" {
			return fmt.Errorf("The Package Manager field can not be blank")
		} else if action.PkgManager != "apt" && action.PkgManager != "yum" {
			return fmt.Errorf("Unknown Package Manager [%s]", action.PkgManager)
		}

		if action.PkgOperation == "" {
			return fmt.Errorf("The Package Operation field can not be blank")
		} else if action.PkgOperation != "install" && action.PkgOperation != "remove" {
			return fmt.Errorf("Unknown Package Operation [%s]", action.PkgOperation)
		}

		if action.Packages == "" {
			return fmt.Errorf("The Packages field can not be blank")
		}
		return nil
	case "key":
		// Validate the Key action
		if action.KeyFile == "" {
			return fmt.Errorf("The KeyField field can not be blank")
		}
		return nil
	default:
		return fmt.Errorf("Unknown Action [%s]", action.ActionType)
	}
}
