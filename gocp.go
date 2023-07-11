package gocp

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Second)
	}
}

func boringRandom(msg string) {
	for i := 0; ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func RunExamples() {
	// boring
	//boring("hello word")

	// boring random (with time sleep random duration)
	// boringRandom("hello Gophers!")

	// Run boring as a Goroutine
	go boring("Hello boring Gopher")

}

// When main returns, the program exists and takes the boring function down with it.
func RunMain() {
	go boring("boring")
	fmt.Println("we are listening")
	time.Sleep(2 * time.Second)
	fmt.Println("You are boring, I am leaving!")
}

// =========
// Using Channels
// =========

func boringChan(msg string, c chan string) {
	for i := 0; ; i++ {
		c <- fmt.Sprintf("%s %d", msg, i) // Expression to be sent can be any suitable value.
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func RunMainChannels() {
	// create a channel - this will be a place where the func sends values
	c := make(chan string)

	// start a goroutine - start sending data (strings) to the channel
	go boringChan("boring!", c)

	// get data from the channel
	for i := 0; i < 6; i++ {
		fmt.Printf("received from the channel: %s\n", <-c) // Received expression is a value
	}
	// we are finishing here
	fmt.Println("You are boring, I am leaving!")

}

// =========
// Patterns
// =========

// Generator pattern - function that returns a channel.

func boringGenerator(msg string) <-chan string { // Returns receive-only channel of strings.
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // Return the channel to the caller.
}

func RunMainGenerator() {
	// create a channel by calling boringGenerator func
	c := boringGenerator("Hello from boring generator!")

	// create a loop and take values from the channel and print them
	for i := 0; i < 6; i++ {
		fmt.Printf("You say: %q\n", <-c)
	}

	// finishing func by printing out the message.
	fmt.Println("You are boring; I am leaving for good!")
}

// Channels as a handle on a service.
// Our boringGenerator function returns a channel that let's us communicate with the boring service it provides.

func RunMainGeneratorService() {
	joe := boringGenerator("Joe")
	mark := boringGenerator("Mark")
	for i := 0; i < 5; i++ {
		// important note about synchronization here:
		// If joe is not ready yet, mark won't be able to send values.
		// In other words mark needs to wait for joe
		fmt.Println(<-joe)
		fmt.Println(<-mark)
	}
	fmt.Println("You are both boring. I am leaving.")
}

// Multiplexing aka "Fan-In" pattern.
// Fan-In pattern - function that takes multiple channels and return a channel.

func fanIn(input1, input2 <-chan string) <-chan string {
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

func RunFanIn() {
	// joe := boringGenerator("Joe")
	// ann := boringGenerator("Ann")
	// c := fanIn(joe, ann)

	c := fanIn(boringGenerator("Joe"), boringGenerator("Ann"))

	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}

	// finish the function
	fmt.Println("You are both boring. I am leaving!")
}
