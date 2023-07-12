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

func RunSearch() {
	start := time.Now()
	results := GoogleLinear("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
