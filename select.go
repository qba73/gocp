package gocp

import (
	"fmt"
	"math/rand"
	"time"
)

// Select pattern
//
// Select statement provides another way to handle multiple channels.
//  - all channels are evaluated
//  - selection blocks until one communication can proceed, which then does
//  - if multiple channels can proceed, select chooses pseudo-randomly
//  - a default clause, if present, executes immediately if no channel is ready

// fnaInOrig is the original function that we used in the first example
// of fan-in pattern. two go routines process data from two channels: input1 & 2
// and send the data to one channel c. The channel c is returned from the func.
func fanInOrig(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- <-input1
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()
	return c
}

// fanInNew usess the select statement and one go routine to
// take values off from channels input1 & input2 and push values
// to the third channel that is returned.
func fanInNew(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()
	return c
}

// Timeout using select
// The time.After func returns a channel that blocks for the specified duration.
// After the interval, the channel delivers the current time, once.

func boringSelect(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(2e3)) * time.Millisecond)
		}
	}()
	return c
}

// RunTimeAfter illustrates how to use for - select case
// statements for timing out slow sources of data.
// In this example if value from chan c is not delivered to s within
// the declared time, the second case triggeres return.
func RunTimeAfter() {
	c := boringSelect("Bolek")
	for {
		select {
		case s := <-c:
			fmt.Println(s)
		case <-time.After(1 * time.Second):
			fmt.Println("You are too slow")
			return
		}
	}
}
