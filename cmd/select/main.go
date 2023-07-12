package main

import "github.com/qba73/gocp"

func main() {
	// example 1
	// gocp.RunTimeAfter()

	// example 2 - time After - timeout for the entire conversation
	//gocp.RunTimeAfterEntireConversation()

	// example 3 - run Quit - explicitely send quit message to the channel
	// and make it return - meaning stop generating data.
	//gocp.RunQuit()

	// example 4
	gocp.RunQuitWithCleanup()
}
