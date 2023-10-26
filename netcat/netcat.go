package netcat

import (
	"io"
	"net"
	"os"
)

func Run() error {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		return err
	}
	defer conn.Close()
	return copy(os.Stdout, conn)
}

func copy(dst io.Writer, src io.Reader) error {
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}
