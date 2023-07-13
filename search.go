package gocp

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = fakeSearch("Web")
	Image = fakeSearch("Image")
	Video = fakeSearch("Video")
)

type Result string

type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

// Google searches lineary for query in web, images and videos.
func GoogleLinear(query string) []Result {
	var results []Result
	results = append(results, Web(query))
	results = append(results, Image(query))
	results = append(results, Video(query))
	return results
}

// GoogleGoroutines takes a query and run search for web, images and video
// independently. Each gorouting runs search and put results on a channel c.
// In the loop we run 3 itertions to get results 3x from the channel.
func GoogleGoroutines(query string) []Result {
	c := make(chan Result) // on this channel we will send query results

	// launch 3 independent searches - each in its own goroutine!
	go func() {
		c <- Web(query)
	}()
	go func() {
		c <- Image(query)
	}()
	go func() {
		c <- Video(query)
	}()

	var results []Result
	for i := 0; i < 3; i++ {
		result := <-c
		results = append(results, result)
	}
	return results
}

func GoogleSearchWithTimeout(query string) []Result {
	c := make(chan Result)

	// Start a goroutine for each search and send results to the channel
	for _, s := range []Search{Web, Image, Video} {
		go func(search Search) {
			c <- search(query)
		}(s)
	}

	var results []Result
	timeout := time.After(80 * time.Millisecond) // timeout on the entire for loop
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return results
		}
	}
	return results
}

func RunSearch() {
	start := time.Now()
	//results := GoogleLinear("golang")
	//results := GoogleGoroutines("golang")
	results := GoogleSearchWithTimeout("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
