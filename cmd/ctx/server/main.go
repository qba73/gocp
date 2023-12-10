package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handlerWithCTX)
	//http.HandleFunc("/ctx", handlerWithCTX)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler started")
	defer log.Printf("handler ended")

	time.Sleep(5 * time.Second)
	fmt.Fprintln(w, "sensor status: OK")
}

func handlerWithCTX(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Printf("handler with ctx started")
	defer log.Printf("handler with ctx ended")

	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintln(w, "sensor status: OK")
	case <-ctx.Done():
		err := ctx.Err()
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
