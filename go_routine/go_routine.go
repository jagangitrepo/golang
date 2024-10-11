package main

import (
	"log"
	"math/rand/v2"
	"time"
)

func main() {
	// Wait Group
	// localvar := 100
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// func(localvar int) {
	// 	log.Printf("GO localvar: %v\n", localvar)
	// 	wg.Done()
	// }(199)
	// go func() {
	// 	localvar = 299
	// 	log.Printf("GO localvar: %v\n", localvar)
	// 	wg.Done()
	// }()
	// log.Println("Hello!.. Go from main.....")
	// log.Printf("Main localvar: %v\n", localvar)
	// wg.Wait()

	// Generator
	numberStream := make(chan int)
	numberStreamCloseChn := make(chan bool)
	// Ananymous function
	go func() {
		defer close(numberStream)
		time.AfterFunc(time.Second*2, func() {
			numberStreamCloseChn <- true
		})
		for {
			numberStream <- rand.IntN(100)
			// time.Sleep(time.Millisecond * 100)  
		}
	}()

	exit_loop := false
	for !exit_loop {
		select {
		case num := <-numberStream:
			log.Printf("From numberstream: %v\n", num)
		case cl := <-numberStreamCloseChn:
			log.Printf("Closing %b\n", cl)
			exit_loop = true
			break
		default:
			// log.Printf("Default")
		}
	}

	// for num := range numberStream {
	// 	log.Printf("From numberstream: %v\n", num)
	// }
}
