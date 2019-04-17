package types

import "encoding/json"

// Action defines what the instructions that will be executed
type Action struct {
	Name       string `json:"name"`
	ActionType string `json:"type"`
	Timeout    int    `json:"timeout"`

	// File based operations
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
	FileMove    bool   `json:"fileMove,omitempty"`

	// Package manager operations
	PkgManager   string `json:"packageManager,omitempty"`
	PkgOperation string `json:"packageOperation,omitempty"`
	Packages     string `json:"packages,omitempty"`

	// Command operations
	Command          string   `json:"command,omitempty"`
	Commands         []string `json:"commands,omitempty"`
	CommandLocal     bool     `json:"commandLocal,omitempty"`
	CommandSaveFile  string   `json:"commandSaveFile,omitempty"`
	CommandSaveAsKey string   `json:"commandSaveAsKey,omitempty"`
	CommandSudo      string   `json:"commandSudo,omitempty"`

	// Piping commands, read in a file and send over stdin, or capture stdout from a local command
	CommandPipeFile string `json:"commandPipeFile,omitempty"`
	CommandPipeCmd  string `json:"commandPipeCmd,omitempty"`

	// Ignore any failures
	IgnoreFailure bool `json:"ignoreFail,omitempty"`

	// Key operations
	KeyFile string `json:"keyFile,omitempty"`
	KeyName string `json:"keyName,omitempty"`

	//Plugin Spec
	Plugin json.RawMessage `json:"plugin,omitempty"`
}
