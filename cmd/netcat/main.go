package main

import (
	"log"

	"github.com/qba73/gocp/netcat"
)

func main() {
	if err := netcat.Run(); err != nil {
		log.Fatalln(err)
	}
}
