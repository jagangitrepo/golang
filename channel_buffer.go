package main

import (
	"log"
	"strconv"
)

func channel_buffer_static_generic() {

	buffer_count := 5
	// String Type
	messages := make(chan interface{}, buffer_count)

	go func() {
		// Store/push the data into the buffered channel
		messages <- 23
		messages <- 23.0
		messages <- "Jagan"
		messages <- 'c'
		messages <- byte(1)

		messages <- byte(1) //In buffered channel, when try to push extra data then it will block the untill the receiver reads it.
		messages <- byte(1)
		messages <- byte(1)
		messages <- byte(1)
		messages <- byte(1)
	}()
	// Retrive the data from the buffered channel for processing
	for i := range buffer_count {
		_ = i
		msg := <-messages
		log.Printf("Received: %T: %#v\n", msg, msg)
	}
}

func channel_buffer_static_type() {

	// String Type
	buffer_count := 5
	messages := make(chan string, buffer_count)

	// Store/push the data into the buffered channel
	for i := range buffer_count {
		messages <- strconv.Itoa(i)
	}

	// Retrive the data from the buffered channel for processing
	for i := range buffer_count {
		_ = i
		log.Println("Received: ", <-messages)
	}
}

func main() {
	channel_buffer_static_generic()
	channel_buffer_static_type()
}
