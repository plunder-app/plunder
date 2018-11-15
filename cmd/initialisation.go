package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/thebsdbox/plunder/pkg/utils"

	log "github.com/Sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	plunderCmd.AddCommand(PlunderConfig)
	plunderCmd.AddCommand(PlunderGet)

}

// PlunderConfig - This is for intialising a blank or partial configuration
var PlunderConfig = &cobra.Command{
	Use:   "config",
	Short: "Initialise a plunder configuration",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		// Indent (or pretty-print) the configuration output
		b, err := json.MarshalIndent(controller, "", "\t")
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("\n%s\n", b)
		return
	},
}

// PlunderGet - The Get command will pull any required components (iPXE boot files)
var PlunderGet = &cobra.Command{
	Use:   "get",
	Short: "Get any components needed for bootstrapping (internet access required)",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		err := utils.PullPXEBooter()
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}
