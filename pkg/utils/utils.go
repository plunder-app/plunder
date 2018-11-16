package utils

import (
	"encoding/hex"
	"io/ioutil"
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

//FileToHex, this is a helper function to allow embedding files into .go files
func FileToHex(filePath string) (sl string, err error) {

	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	sl = hex.EncodeToString(bs)
	return

}
