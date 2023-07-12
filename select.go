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

// RunTimeAfterEntireConversation illustrates timing out
// the entire conversation (for loop).
// Create the timer once, outside of the loop, to time out
// the entire conversation. Note that in the func RunTimeAfter
// we have a timeout for each message.
func RunTimeAfterEntireConversation() {
	c := boringSelect("Bolek")             // Create a generator
	timeout := time.After(5 * time.Second) // Create a func scoped timeout (chan)
	for {
		select {
		case s := <-c: // keep getting values from the channel
			fmt.Println(s)
		case <-timeout: // when the value is available on this channel take it and return
			fmt.Println("You talk too much")
			return
		}
	}
}

// RunQuit shows how to quit generator by sending a "signal" to the quit channel.
// boring generator stops putting values on the c channel when quit chan receives
// a value. Once value is placed on the quit channel select - case <- quit
// executes and exits the for loop.
func RunQuit() {

	// boring is a generator that puts strings on the returned
	// channel until receiving `quit` signal on the quit channel.
	boring := func(msg string, quit chan bool) chan string {
		c := make(chan string)
		go func() {
			for i := 0; ; i++ {
				select {
				case c <- fmt.Sprintf("%s %d", msg, i):
					// do nothing
					time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
				case <-quit:
					return
				}
			}
		}()
		return c
	}

	quit := make(chan bool)    // chan we will send a message to the generator to stop producing strings and putting them on the chan
	c := boring("Lolek", quit) // boring returns a chan containing strings being generated by the boring func
	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-c) // printing string off the channel
	}
	quit <- true // sent a signal to the generator to stop
}

// RunQuitWithCleanup illustrates two-way communication. We send a signal to quit producing and putting on the channel data.
// Then we call clanup func and then we send a confirmation message on the same channel.
func RunQuitWithCleanup() {
	cleanup := func() { fmt.Println("got a signal to exit. cleaning up!") }

	boring := func(msg string, quit chan string) chan string {
		c := make(chan string)
		go func() {
			for i := 0; ; i++ {
				select {
				case c <- fmt.Sprintf("%s %d", msg, i):
					// do nothing
					time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
				case <-quit:
					cleanup() // run possible cleanup
					quit <- "see you!"
					return
				}
			}
		}()
		return c
	}

	quit := make(chan string)
	c := boring("Jonny", quit)
	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-c)
	}
	quit <- "bye!"
	fmt.Printf("Jonny says %s\n", <-quit)
}
