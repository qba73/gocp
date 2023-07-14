package gocp

import (
	"fmt"
	"math/rand"
	"sync"
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

	// fan-in pattern start
	// Start a goroutine for each search and send results to the channel
	for _, s := range []Search{Web, Image, Video} {
		go func(search Search) {
			c <- search(query)
		}(s)
	}
	// fan-in pattern end

	var results []Result

	// time out pattern for all 'conversation'
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
	// end time out pattern

	return results
}

func GoogleSearchRange(query string) []Result {
	c := make(chan Result)

	// fan-in pattern start
	// Start a goroutine for each search and send results to the channel
	for _, s := range []Search{Web, Image, Video} {
		go func(search Search) {
			c <- search(query)
		}(s)
	}
	// fan-in pattern end

	var results []Result

	// time out pattern for all 'conversation'
	timeout := time.After(80 * time.Millisecond) // timeout on the entire for loop

	for i := 0; i < 3; i++ {
		select {
		case res := <-c:
			results = append(results, res)
		case <-timeout:
			fmt.Println("timeout")
			return results
		}
	}
	// end time out pattern

	return results
}

// GoogleSearchWG illustrates how to run independent searches
// and use WaitGroup
func GoogleSearchWG(query string) []Result {
	c := make(chan Result)

	var wg sync.WaitGroup

	// fan-in pattern start
	// Start a goroutine for each search and send results to the channel
	for _, s := range []Search{Web, Image, Video} {
		wg.Add(1)
		go func(search Search) {
			defer wg.Done()
			c <- search(query)
		}(s)
	}
	// fan-in pattern end

	// Wait for all searches to finish and close the channel.
	go func() {
		wg.Wait()
		close(c)
	}()

	var results []Result

	// Get results from the channel and aggregate
	for r := range c {
		results = append(results, r)
	}

	return results
}

// =========
// Example - how to avoid slow servers and use server replicas.
// =========

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c // return result from search that responded first (the fastest)
}

// =========
// Example - putting together patterns, and create server replicas with timeout.
//
// Reducing tail latency using replicated search servers.
// =========

func GoogleSearchReplicas(query string) []Result {
	c := make(chan Result) // channel for our serch results
	// 3 goroutines for searching web, video and image using fastest replicas
	go func() { c <- First(query, fakeSearch("web1"), fakeSearch("web2")) }()
	go func() { c <- First(query, fakeSearch("video1"), fakeSearch("video2")) }()
	go func() { c <- First(query, fakeSearch("image1"), fakeSearch("image2")) }()

	timeout := time.After(80 * time.Millisecond) // "global timeout" for the entire search query (for web, video and image)
	var results []Result

	for i := 0; i < 3; i++ {
		select {
		case res := <-c:
			results = append(results, res)
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
	// results := GoogleSearchWithTimeout("golang")

	// Search with Replicas
	// results := First("golang",
	// 	fakeSearch("replica 1"),
	// 	fakeSearch("replica 2"),
	// )

	// Search with replicas and time outs:
	// results := GoogleSearchReplicas("golang")

	// Closing channel example
	//results := GoogleSearchRange("golang")

	// Search with WaitGroup
	results := GoogleSearchWG("golang")

	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
