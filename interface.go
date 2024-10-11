package main

import (
	"fmt"
	"iter"
	"log"
)

// Writable interface
type IWriteable interface {
	Write(data ...interface{}) (int, error)
}

// Readable interface
type IReadable interface {
	Read() iter.Seq[string]
}

type ConsoleMessage struct {
	str_msg []string
}

// Variadic function
func (c *ConsoleMessage) Write(msgs ...interface{}) (int, error) {
	var ret = 0

	for _, msg := range msgs {
		log.Println("Msg: ", msg)
		c.str_msg = append(c.str_msg, fmt.Sprintf("%v", msg))
	}

	return ret, nil
}

// Iterator function
func (c *ConsoleMessage) Read() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, msg := range c.str_msg {
			if !yield(msg) {
				return
			}
		}
	}
}

func main() {
	consolemsg := ConsoleMessage{}
	consolemsg.Write("Vaishu", "Pattu", 1, 23.1)
	for msg := range consolemsg.Read() {
		fmt.Println(msg)
	}
}
