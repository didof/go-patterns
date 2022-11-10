package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	exit := make(chan struct{})

	b := new(barrack)
	swordmenCh, _ := b.swordmen(10)
	archersCh, _ := b.archers(5)

	unitsCh := b.mergeProductions(swordmenCh, archersCh)
	i := 0
	go func() {
		for unit := range unitsCh {
			fmt.Println(*unit)
			i++
		}
		close(exit)
	}()

	<-exit

	fmt.Printf("%d units produced\n", i)

	os.Exit(0)
}

type barrack struct{}

type soldier struct {
	t        string
	atk, dfe int
}

func (s soldier) String() string {
	var sound string
	switch s.t {
	case "swordman":
		sound = "eis Machin!"
	case "archer":
		sound = "eisvoli"
	}
	return fmt.Sprintf("%10s: %s", s.t, sound)
}

func (b *barrack) swordmen(amount int) (<-chan *soldier, error) {
	ch := make(chan *soldier)

	go func() {
		defer close(ch)
		for i := 0; i < amount; i++ {
			ch <- &soldier{"swordman", 6, 70}
		}
	}()

	return ch, nil
}

func (b *barrack) archers(amount int) (<-chan *soldier, error) {
	ch := make(chan *soldier)

	go func() {
		defer close(ch)
		for i := 0; i < amount; i++ {
			ch <- &soldier{"archer", 8, 40}
		}
	}()

	return ch, nil
}

func (b *barrack) mergeProductions(chs ...<-chan *soldier) <-chan *soldier {
	var wg sync.WaitGroup
	wg.Add(len(chs))

	out := make(chan *soldier)

	for _, ch := range chs {
		go func(ch <-chan *soldier) {
			defer wg.Done()
			for unit := range ch {
				out <- unit
			}
		}(ch)
	}

	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}
