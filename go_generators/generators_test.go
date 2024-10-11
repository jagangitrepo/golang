package main

import (
	"io/ioutil"
	"iter"
	"log"
	"math/rand"
	"testing"
)

const N int = 1000

func RandomNumberStream() iter.Seq[int] {
	var random_numbers []int

	// Generate the random numbers and store it inside the slice
	for i := 1; i <= N; i++ {
		random_numbers = append(random_numbers, rand.Intn(1000))
	}

	// Return the yeild function, so that while caller calls it will take from the stream and push to the caller sequentially
	return func(yeild func(int) bool) {
		for _, random_number := range random_numbers {
			if !yeild(random_number) {
				return
			}
		}
	}
}

/*
GoGenerateUsingIterators - Using the iter to generrate the stream of numbers
*/
func GoGenerateUsingIterators() {
	// Get the random number from the stream one by one.
	for random_num := range RandomNumberStream() {
		log.Println("Random No: ", random_num)
	}
}

/*
GoGenerateUsingChannels - Using the Go channels to generate the stream of numbers
*/
func GoGenerateUsingChannels() {
	random_num_stream := make(chan int)
	random_num_stream_close_chn := make(chan bool)

	// Random Number stream generator go function
	go func(random_num_stream chan int, random_num_stream_close_chn chan bool) {
		defer close(random_num_stream)
		defer close(random_num_stream_close_chn)

		for i := 1; i <= N; i++ {
			random_num_stream <- rand.Intn(10000)
		}
		random_num_stream_close_chn <- true

	}(random_num_stream, random_num_stream_close_chn)

	exit_loop := false
	for !exit_loop {
		select {
		case random_num := <-random_num_stream:
			log.Println("Random No: ", random_num)
		case is_exit := <-random_num_stream_close_chn:
			exit_loop = is_exit
			break
		}
	}
}

/*
TestGoGenerateUsingChannels - Using the Go channels  to generate the stream of numbers
*/
func BenchmarkGoGenerateUsingChannels(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	b.Run("", func(b *testing.B) {
		GoGenerateUsingChannels()
	})
}

func BenchmarkGoGenerateUsingIterators(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	b.Run("", func(b *testing.B) {
		GoGenerateUsingIterators()
	})
}
