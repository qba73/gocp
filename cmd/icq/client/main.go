package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	port := flag.String("port", "8000", "port to connect to")
	flag.Parse()

	conn, err := net.Dial("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)
	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- true
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
