package echo

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func Run() error {
	listener, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", toTitle(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", "--- --- ---")
}

func toTitle(s string) string {
	if s == "" {
		return s
	}
	first := strings.ToUpper(string(s[0]))
	return first + strings.ToLower(s[1:])
}

func handleConn(c net.Conn) {
	defer c.Close()
	input := bufio.NewScanner(c)
	for input.Scan() {
		go echo(c, input.Text(), 1*time.Second)
	}
	if err := input.Err(); err != nil {
		log.Println(err)
	}
}
