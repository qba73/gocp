package prodline2

import (
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
	Items   int
}

// StartStage is a first stage in the production line.
// At this stage items (products) are initially prepared
// for furthe processing. This about this stage as a stage
// where raw items are produced and sent to the next stage.
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

// Stage represents an intermittent or the final stage in an imaginary
// production line. It processes items that are sent by initial stage
// on the incomming channel (in), processes the items and forwards
// them to the next stage by placing them onto the outging channel (out).
func (pl *ProductionLine) Stage(name string, in <-chan item, out chan<- item, fn Work, d, stddev time.Duration) {
	for i := range in {
		if pl.Verbose {
			pl.Logger.Printf("stage %s, %v", name, i)
		}
		fn(d, stddev)
		out <- i
	}
}

// Run simulates running all stages of the production line.
// It defines an initial and further stages of the line.
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
