package netcat3

import (
	"io"
	"log"
	"net"
	"os"
)

func Run() {
	conn, err := net.Dial("tcp", "localhost:9000")
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
	<-done // wait for the goroutine that operates in the background to finish
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalln(err)
	}
}
