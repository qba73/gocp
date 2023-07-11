package main

import "github.com/qba73/gocp"

func main() {
	// ==========
	//gocp.RunExamples()

	// ==========
	// Example with time sleep to prevent from
	// main exiting and bringing goroutine down with it.
	// ==========
	// gocp.RunMain()

	// ==========
	// Example with a goroutine feeding the channel with string values
	// ==========
	//gocp.RunMainChannels()

	// ==========
	// Examples from the Patterns section:
	// ==========

	// Run generator - one function
	// gocp.RunMainGenerator()

	// run main generator service - multiple functions - return channel pattern
	// gocp.RunMainGeneratorService()

	// Run fan-in multiplex (file)
	//gocp.RunMultiplex()

	// run fan-in example pattern
	// gocp.RunFanIn()

	// Run fan-in pattern with sequence managed by a chan sent inside chan
	gocp.RunSequence()
}
