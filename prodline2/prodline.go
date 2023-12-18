package prodline2

import (
	"log"
	"math/rand"
	"time"
)

type item int

type ProductionLine struct {
	Logger  log.Logger
	Verbose bool
	Items   int
}

type Work func(d, stddev time.Duration)

func (pl *ProductionLine) StartStage(name string, out chan<- item, fn Work, d, stddev time.Duration) {
	for i := 0; i < pl.Items; i++ {
		c := item(i)
		if pl.Verbose {
			pl.Logger.Printf("starting stage %v, %d", c, i)
		}
		fn(d, stddev)
		out <- c
	}
}

func (pl *ProductionLine) Stage(name string, in <-chan item, out chan<- item, fn Work, d, stddev time.Duration) {
	for i := range in {
		if pl.Verbose {
			pl.Logger.Printf("stage %s, %v", name, i)
		}
		fn(d, stddev)
		out <- i
	}
}

func (pl *ProductionLine) Run() {
	baked := make(chan item, 1)
	iced := make(chan item, 1)
	inscribed := make(chan item, 1)
	packed := make(chan item, 1)

	go pl.StartStage("baking", baked, work, time.Second, 200*time.Millisecond)
	go pl.Stage("icing", baked, iced, work, time.Second, 200*time.Millisecond)
	go pl.Stage("inscribing", iced, inscribed, work, time.Second, 300*time.Millisecond)
	go pl.Stage("packing", inscribed, packed, work, time.Second, 50*time.Millisecond)

	// drain the packed channel
	for i := 0; i < pl.Items; i++ {
		<-packed
	}
}

func Run() {
	pl := ProductionLine{
		Logger:  *log.Default(),
		Verbose: true,
		Items:   3,
	}
	pl.Run()
}

func work(d, stddev time.Duration) {
	delay := d + time.Duration(rand.NormFloat64()*float64(stddev))
	time.Sleep(delay)
}
