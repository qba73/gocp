package main

import (
	"log"
	"os"

	"github.com/qba73/gocp/clock"
)

func main() {
	if err := clock.RunClock(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
