package gocp

import (
	"fmt"
	"math/rand"
	"time"
)

// =========
// Fan-in pattern - sending a chan inside a chan
// =========

//  This pattern is used for signaling. We use it here to restore the order of print execution.
//
// Restoring sequencing - sending a chan inside a chan - pattern for signalling!

type Message struct {
	str  string
	wait chan bool // signaller
}

func boringSequenceGenerator(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
		}
	}()
	return c
}

func fanInSequence(input1, input2 <-chan string) <-chan Message {
	waitFor := make(chan bool)
	c := make(chan Message)
	go func() {
		for {
			c <- Message{<-input1, waitFor}
			<-waitFor
		}
	}()
	go func() {
		for {
			c <- Message{<-input2, waitFor}
			<-waitFor
		}
	}()
	return c
}

func RunSequence() {
	c := fanInSequence(boringSequenceGenerator("Bolek"), boringSequenceGenerator("Lolek"))
	for i := 0; i < 10; i++ {
		msg1 := <-c
		fmt.Println(msg1.str)

		msg2 := <-c
		fmt.Println(msg2.str)

		msg1.wait <- true
		msg2.wait <- true
	}
}
