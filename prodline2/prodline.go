package prodline2

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type item int

type Work func(d, stddev time.Duration)

// ProductionLine represents an imaginary production line
// that process N number of items per run.
type ProductionLine struct {
	Logger  log.Logger
	Verbose bool
	Stages  []workerFn

	output <-chan item
	ctx    context.Context
}

type Stage struct{}

type workerFn func(context.Context, <-chan item, chan<- item)

func (pl *ProductionLine) AddStage(name string, worker workerFn) {
	pl.Stages = append(pl.Stages, worker)
}

// Run simulates running all stages of the production line.
// It defines an initial and further stages of the line.
func (pl *ProductionLine) Start() {
	prev := make(chan item, 1)

	go func(ctx context.Context, ch chan<- item) {
		i := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Println("first stage cancelled!")
				close(ch)
				return
			default:
				ch <- item(i)
				i++
			}
		}
	}(pl.ctx, prev)

	for _, stage := range pl.Stages {
		out := make(chan item, 1)
		go stage(pl.ctx, prev, out)
		prev = out
	}
	pl.output = prev
}

func (pl *ProductionLine) Items() <-chan item {
	return pl.output
}

// have a function that registers stages and run the full pipeline

func newDummyStage(t, stddev time.Duration) workerFn {
	return func(ctx context.Context, in <-chan item, out chan<- item) {
		for item := range in {
			select {
			case <-ctx.Done():
				fmt.Println("worker cancelled!")
				close(out)
				return
			default:
			}

			delay := t + time.Duration(rand.NormFloat64()*float64(stddev))
			time.Sleep(delay)
			out <- item
		}
	}
}

func Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pl := ProductionLine{
		Logger:  *log.Default(),
		Verbose: true,
		ctx:     ctx,
	}

	pl.AddStage("baking", newDummyStage(time.Second, 200*time.Millisecond))
	pl.AddStage("icing", newDummyStage(time.Second, 200*time.Millisecond))
	pl.AddStage("inscribing", newDummyStage(time.Second, 200*time.Millisecond))
	pl.AddStage("packaging", newDummyStage(time.Second, 200*time.Millisecond))

	pl.Start()

	for item := range pl.Items() {
		fmt.Println(item)
	}
}

func work(d, stddev time.Duration) {
	delay := d + time.Duration(rand.NormFloat64()*float64(stddev))
	time.Sleep(delay)
}
