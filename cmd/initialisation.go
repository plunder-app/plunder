package cmd

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	plunderCmd.AddCommand(PlunderInit)
}

// PlunderInit - This is for intialising a blank or partial configuration
var PlunderInit = &cobra.Command{
	Use:   "init",
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
