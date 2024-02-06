package leader

import (
	"fmt"
	"math/rand"
	"time"
)

func Run() {
	fmt.Println("Running leader election simulation...")

	//WaitForSignal()
	FanOut()
	fmt.Println("=== === === === ===")
	WaitForTask()
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

func FanOut() {
	//nsNames := []string{"default", "nsA", "nsB", "nsC"}
	//nsCount := len(nsNames)

	nsCount := 20

	ch := make(chan int)

	for n := 0; n < nsCount; n++ {
		go func(n int) {
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			ch <- rand.Intn(100)
		}(n)
	}

	totalObj := 0
	for nsCount > 0 {
		nsCount--

		p := <-ch
		fmt.Println("VS obj in ns :", p)

		totalObj += p
		fmt.Println("Total obj :", totalObj)
	}

	time.Sleep(time.Second)
	fmt.Println("=========")
}

func WaitForTask() {
	ch := make(chan string)

	go func() {
		p := <-ch
		fmt.Println("minion : recv'd signal :", p)
	}()

	ch <- "label added"
	fmt.Println("master : sent signal")

	time.Sleep(time.Second)
	fmt.Println("=========")
}
