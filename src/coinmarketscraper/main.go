package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"coinmarketscraper/engine"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	stopChan := make(chan struct{})
	go func() {
		<-signalChan
		stopChan <- struct{}{}
		close(stopChan)
		return
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	coins := engine.Run(stopChan)
	for coin := range coins {
		fmt.Println(fmt.Sprintf("%+v", coin))
	}
}
