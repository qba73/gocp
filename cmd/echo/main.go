package main

import (
	"log"

	"github.com/qba73/gocp/echo"
)

func main() {
	if err := echo.Run(); err != nil {
		log.Fatalln(err)
	}
}
