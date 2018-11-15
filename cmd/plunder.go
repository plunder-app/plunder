package cmd

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Release - this struct contains the release information populated when building plunder
var Release struct {
	Version string
	Build   string
}

var plunderCmd = &cobra.Command{
	Use:   "plunder",
	Short: "This is a tool for finding gold amongst bare-metal (and provisioning kubernetes)",
}

var logLevel int

func init() {
	// Global flag across all subcommands
	plunderCmd.PersistentFlags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
	plunderCmd.AddCommand(plunderVersion)
}

// Execute - starts the command parsing process
func Execute() {
	if os.Getenv("PLUNDER_LOGLEVEL") != "" {
		i, err := strconv.ParseInt(os.Getenv("PLUNDER_LOGLEVEL"), 10, 8)
		if err != nil {
			log.Fatalf("Error parsing environment variable [PLUNDER_LOGLEVEL")
		}
		// We've only parsed to an 8bit integer, however i is still a int64 so needs casting
		logLevel = int(i)
	} else {
		// Default to logging anything Info and below
		logLevel = int(log.InfoLevel)
	}

	log.SetLevel(log.Level(logLevel))
	if err := plunderCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var plunderVersion = &cobra.Command{
	Use:   "version",
	Short: "Version and Release information about the plunder tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Plunder Release Information\n")
		fmt.Printf("Version:  %s\n", Release.Version)
		fmt.Printf("Build:    %s\n", Release.Build)
	},
}
