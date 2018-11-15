package utils

import (
	"os"
	"os/signal"
	"sync"
)

//WaitForCtrlC - This function is the loop that will catch a Control-C keypress
func WaitForCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		endWaiter.Done()
	}()
	endWaiter.Wait()
}
