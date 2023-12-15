package prodline

import (
	"fmt"
	"math/rand"
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
func (b *Bakery) baker(baked chan<- cake) {
	for i := 0; i < b.Cakes; i++ {
		c := cake(i)
		if b.Verbose {
			fmt.Println("baking", c)
		}
		work(b.BakeTime, b.BakeStdDev)
		baked <- c
	}
	close(baked)
}

// icer represents a person or machine for puting ising
// on the top of the cake. Icer sits in the middle of
// the 'production' cakes chain between the baker and the inscriber.
func (b *Bakery) icer(iced chan<- cake, baked <-chan cake) {
	// range over the channel
	for c := range baked {
		if b.Verbose {
			fmt.Println("icing", c)
		}
		work(b.IceTime, b.IceStdDev)
		iced <- c
	}
}

// inscriber represents an entity (machine or a person) that
// takes cakes from the icer and decorate them (inscribe).
// inscriber is a last stage of the cake production line.
func (b *Bakery) inscriber(inscribed chan<- cake, iced <-chan cake) {
	// for i := 0; i < b.Cakes; i++ {
	// 	c := <-iced
	// 	if b.Verbose {
	// 		fmt.Println("icing", c)
	// 	}
	// 	work(b.InscribeTime, b.InscribeStdDev)
	// 	if b.Verbose {
	// 		fmt.Println("finished", c)
	// 	}
	// }
	for i := 0; i < b.Cakes; i++ {
		c := <-iced
		if b.Verbose {
			fmt.Println("inscribing", c)
		}
		work(b.InscribeTime, b.InscribeStdDev)
		inscribed <- c
	}
}

func (b *Bakery) packer(inscribed <-chan cake) {
	for i := 0; i < b.Cakes; i++ {
		c := <-inscribed
		if b.Verbose {
			fmt.Println("packaging", c)
		}
		work(b.PackingTime, b.PackingStdDev)
		if b.Verbose {
			fmt.Println("finished packaging", c)
		}
	}
}

// Work runs the baking simulation 'runs' times.
func (b *Bakery) Work(runs int) {
	for run := 0; run < runs; run++ {
		baked := make(chan cake, b.BakeBuf)
		iced := make(chan cake, b.IceBuf)
		inscribed := make(chan cake, b.InscribeBuf)

		// start baking
		go b.baker(baked)

		// start icing cakes
		for i := 0; i < b.NumIcers; i++ {
			go b.icer(iced, baked)
		}

		// start inscribing cakes
		b.inscriber(iced, inscribed)

		// start packaging cakes for storage
		// Packaging is the last step in the production line.
		b.packer(inscribed)
	}
}

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
		BakeBuf:    3,

		NumIcers:  2,
		IceTime:   3 * time.Second,
		IceStdDev: 500 * time.Millisecond,
		IceBuf:    3,

		InscribeTime:   3 * time.Second,
		InscribeStdDev: 1 * time.Second,
		InscribeBuf:    2,

		NumPackers:    1,
		PackingTime:   200 * time.Millisecond,
		PackingStdDev: 100 * time.Millisecond,
		PackerBuf:     1,
	}
	b.Work(1)
}
