package log

import (
	"context"
	"log"
	"math/rand"
	"net/http"
)

type reqKeyID uint64

const reqID reqKeyID = 1

func Println(ctx context.Context, msg string) {
	id, ok := ctx.Value(reqID).(uint64)
	if !ok {
		log.Println("could not find request ID in context")
		return
	}
	log.Printf("[%d] %s", id, msg)
}

func Decorate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := rand.Uint64()
		ctx = context.WithValue(ctx, reqID, id)
		f(w, r.WithContext(ctx))
	}
}
