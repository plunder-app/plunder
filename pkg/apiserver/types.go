package apiserver

import "encoding/json"

//Response - This is the wrapper for responses back to a client, if any errors are created then the payload isn't guarenteed
type Response struct {
	FriendlyError string `json:"friendlyError,omitempty"`
	Error         string `json:"error,omitempty"`   // when it goes wrong
	Success       string `json:"success,omitempty"` // when it goes correct

	Payload json.RawMessage `json:"payload,omitempty"`
}
