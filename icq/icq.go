package icq

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type client chan<- string // outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
	errors   = make(chan error)
)

// messanger broadcast messages to connected clients
// and adds/removes clients from the pool.
func messanger() {
	// clients represents all connected clients to the ICQ server
	clients := make(map[client]bool)

	for {
		select {
		case msg := <-messages:
			// broadcast a message to all clients' outgoing message channels
			for cl := range clients {
				cl <- msg
			}

		// a new client connects to the server:
		//  - add it to the client pool
		case cl := <-entering:
			clients[cl] = true

		// a client disconnects from the server
		//  - remove it from the pool
		//  - close the channel (client) the client uses to communicate with the server
		case cl := <-leaving:
			delete(clients, cl)
			close(cl)
		}
	}
}

func handleConnection(conn net.Conn) {
	// outgoing client messages
	// the channel represents a new client that will be registered in the clients map in the func messanger.
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "Connected new client: " + who
	messages <- who + " has joined conversation"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	if err := input.Err(); err != nil {
		// handle error here, possible in a new, error channel?
		errors <- err
	}

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

// clientWriter writes messages comming from the channel
// to the provided connection.
func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			errors <- err
		}
	}
}

// RunServer starts a new ICQ Chat Server.
// The main job of the function is to listen for and accept new incoming
// newtwork connections from clients. For each connection the func creates
// a new handleConnection goroutine.
func RunServer() {
	icqlog := log.New(os.Stdout, "ICQ:", log.Lshortfile)
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		icqlog.Fatal(err)
	}
	icqlog.Printf("listening on " + l.Addr().String())

	go messanger()
	for {
		conn, err := l.Accept()
		if err != nil {
			icqlog.Print(err)
			continue
		}
		go handleConnection(conn)
	}
}
