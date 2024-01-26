package leader

import (
	"fmt"
	"math/rand"
	"time"
)

func Run() {
	fmt.Println("Running leader election simulation...")

	WaitForSignal()
}

func WaitForSignal() {
	ch := make(chan bool)

	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		ch <- true
		fmt.Println("start leading")
	}()

	p := <-ch
	fmt.Println("election confirmed: recv'd signal :", p)
	time.Sleep(time.Second)
	fmt.Println("=========")
}
