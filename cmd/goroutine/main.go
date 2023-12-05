package main

import "fmt"

func main() {

	ch := make(chan struct{}, 10)

	go func() {
		fmt.Println("hello")
		<-ch
	}()
	ch <- struct{}{}

	fmt.Println("world")
}
