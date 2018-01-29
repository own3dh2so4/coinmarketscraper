package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
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
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			for coin := range coins {
				fmt.Println(fmt.Sprintf("%+v", coin))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
