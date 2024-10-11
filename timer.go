package main

import (
	"log"
	"time"
)

func main() {

	// Timer
	timerOne := time.NewTimer(time.Second * 2)
	log.Println("Before Wait")
	<-timerOne.C
	log.Println("After Wait")

	// Ticker
	tickerOne := time.NewTicker(time.Millisecond * 500)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-tickerOne.C:
				log.Println("Tick at ", t)
			}
		}
	}()

	time.Sleep(time.Second * 10)
	tickerOne.Stop()
	done <- true
}
