package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	//reportSensorStatus1(ctx, 5*time.Second, "status checked: OK")
	//reportSensorStatus2(ctx, 5*time.Second, "status checked: OK")
	reportSensorStatus3(ctx, 5*time.Second, "OK")
	//reportSensorStatus(ctx, 5*time.Second, "sensor OK")
}

// ctx has no impact. The func does not use it.
// func reportSensorStatus1(ctx context.Context, d time.Duration, msg string) {
// 	time.Sleep(d)
// 	fmt.Println("Temp sensor: " + msg)
// }

// select ch and read a message from the time package - function After.
//
//	note that select blocks the execution as we do not specify the `default` case.
func reportSensorStatus2(ctx context.Context, d time.Duration, msg string) {
	select {
	case <-time.After(d):
		fmt.Println("Temp sensor: " + msg)
	}
}

func reportSensorStatus3(ctx context.Context, d time.Duration, msg string) {
	select {
	case <-time.After(d):
		fmt.Println("Temp sensor: " + msg)
	case <-ctx.Done():
		log.Println(ctx.Err())
	}
}

// func reportSensorStatus(ctx context.Context, d time.Duration, msg string) {
// 	select {
// 	case <-ctx.Done():

// 		log.Println("canceleld")
// 	default:
// 		time.Sleep(d)
// 		fmt.Println(msg)
// 	}
// }
