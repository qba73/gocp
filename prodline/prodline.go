package prodline

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Bakery struct {
	Verbose bool
	Cakes   int

	BakeTime   time.Duration
	BakeStdDev time.Duration
	BakeBuf    int

	NumIcers  int
	IceTime   time.Duration
	IceStdDev time.Duration
	IceBuf    int

	NumInscribers  int
	InscribeTime   time.Duration
	InscribeStdDev time.Duration
	InscribeBuf    int

	NumPackers    int
	PackingTime   time.Duration
	PackingStdDev time.Duration
	PackerBuf     int
}

type cake int

// baker represents a baking machine or a baker - person
// responsible for baking cakes in the bakery.
// baker bakes cakes and places them on a table. cakes at
// this point are ready to ice - a job for the icer.
func (b *Bakery) baker(out chan<- cake) {
	for i := 0; i < b.Cakes; i++ {
		c := cake(i)
		if b.Verbose {
			fmt.Println("baking", c)
		}
		work(b.BakeTime, b.BakeStdDev)
		out <- c
	}
	fmt.Println("baker done, closing")
}

// icer represents a person or machine for puting ising
// on the top of the cake. Icer sits in the middle of
// the 'production' cakes chain between the baker and the inscriber.
func (b *Bakery) icer(in <-chan cake, out chan<- cake) {
	// range over the channel
	for c := range in {
		if b.Verbose {
			fmt.Println("icing", c)
		}
		work(b.IceTime, b.IceStdDev)
		out <- c
	}
	fmt.Println("icer done, closing")
}

// inscriber represents an entity (machine or a person) that
// takes cakes from the icer and decorate them (inscribe).
// inscriber is a last stage of the cake production line.
func (b *Bakery) inscriber(in <-chan cake, out chan<- cake) {
	for c := range in {
		if b.Verbose {
			fmt.Println("inscribing", c)
		}
		work(b.InscribeTime, b.InscribeStdDev)
		out <- c
	}
	fmt.Println("inscriber done, closing!")
}

func (b *Bakery) packer(in <-chan cake, out chan<- cake) {
	for c := range in {
		if b.Verbose {
			fmt.Println("packaging", c)
		}
		work(b.PackingTime, b.PackingStdDev)
		if b.Verbose {
			fmt.Println("finished packaging", c)
		}
		fmt.Println("packed", c)
		out <- c
	}
	fmt.Println("packer done, closing!")
}

// Work runs the baking simulation 'runs' times.
func (b *Bakery) Work(runs int) {
	for run := 0; run < runs; run++ {
		baked := make(chan cake, b.BakeBuf)
		iced := make(chan cake, b.IceBuf)
		inscribed := make(chan cake, b.InscribeBuf)
		packed := make(chan cake, b.PackerBuf)

		// start baking using one baker (machine or human)
		go b.baker(baked)

		// start icing cakes - using 1 or more icers (machines or humans)
		for i := 0; i < b.NumIcers; i++ {
			go b.icer(baked, iced) // from backed -> to iced output
		}

		// start inscribing cakes - using 1 or more inscribers (machines or humans)
		for i := 0; i < b.NumInscribers; i++ {
			go b.inscriber(iced, inscribed)
		}

		// start packaging cakes for storage - using n packagers
		// Packaging is the last step in the production line.
		for i := 0; i < b.NumPackers; i++ {
			go b.packer(inscribed, packed)
		}

		// drain the queue - loop n times where n it the number of cakes in the batch.
		for i := 0; i < b.Cakes; i++ {
			<-packed
		}
	}
}

// work simulates a work items like baking, icing, packaging etc.
func work(d, stddev time.Duration) {
	delay := d + time.Duration(rand.NormFloat64()*float64(stddev))
	time.Sleep(delay)
}

func RunBakery() {
	b := Bakery{
		Verbose: true,

		Cakes:      3,
		BakeTime:   2 * time.Second,
		BakeStdDev: 500 * time.Millisecond,
		BakeBuf:    1,

		NumIcers:  2,
		IceTime:   3 * time.Second,
		IceStdDev: 500 * time.Millisecond,
		IceBuf:    1,

		NumInscribers:  2,
		InscribeTime:   3 * time.Second,
		InscribeStdDev: 1 * time.Second,
		InscribeBuf:    2,

		NumPackers:    1,
		PackingTime:   200 * time.Millisecond,
		PackingStdDev: 100 * time.Millisecond,
		PackerBuf:     3,
	}
	b.Work(1)
}

// concurent fetch emulation
func fetch(urls []string) {
	results := make(chan string, 5)

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, u := range urls {
		go func(u string) {
			defer wg.Done()
			time.Sleep(time.Second)
			results <- "readings"
		}(u)
	}

	go func() {
		for res := range results {
			// do sth with res
			fmt.Println(res)
		}
	}()

	wg.Wait()
}
