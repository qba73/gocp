package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/qba73/gocp/log"
)

func main() {
	http.HandleFunc("/", log.Decorate(handlerWithCTX))
	panic(http.ListenAndServe("127.0.0.1:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(context.Background(), "handler started")
	defer log.Println(context.Background(), "handler ended")

	time.Sleep(5 * time.Second)
	fmt.Fprintln(w, "sensor status: OK")
}

func handlerWithCTX(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println(ctx, "handler with ctx started")
	defer log.Println(ctx, "handler with ctx ended")

	select {
	case <-time.After(5 * time.Second):
		fmt.Fprintln(w, "sensor status: OK")
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(ctx, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
