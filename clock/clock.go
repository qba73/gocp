package clock

import (
	"io"
	"log"
	"net"
	"time"
)

// =========
// Concurrent clock server
// =========

func RunClock() error {
	l, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		handleConnectionClock(conn)
	}
}

func handleConnectionClock(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // it will disconnect a client
		}
		time.Sleep(5 * time.Second)
	}
}
