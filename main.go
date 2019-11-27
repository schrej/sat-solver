package main

import (
	"flag"
	"fmt"
	"math"
	"sync"
)

func main() {
	n := flag.Int("n", 0, "2^n is the number of threads")
	flag.Parse()
	f := flag.Arg(0)

	parsed, err := Parse(f)
	if err != nil {
		fmt.Printf("Failed to parse formula: %e", err)
		return
	}
	symbols := []string{}
	form := parsed.Form(&symbols)

	fmt.Printf("Formula parsed successfully, there are %d symbols.\n", len(symbols))

	res := make(chan Interpretation)
	wg := sync.WaitGroup{}
	rootInt := Interpretation{Fixed: 0, Vals: make([]bool, len(symbols))}
	for i := 0; i < int(math.Pow(2, float64(*n))); i++ {
		in := Interpretation{Fixed: *n, Vals: make([]bool, len(symbols))}
		copy(in.Vals[:*n], rootInt.Vals[:*n])
		wg.Add(1)
		go Solver(in, form, res, &wg)
		rootInt.Next()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case inter := <-res:
		fmt.Println("Formula is satisfyable")
		inter.Print(symbols)
	case <-done:
		fmt.Println("Formula is a contradiction.")
	}

}

func Solver(in Interpretation, form Formula, res chan Interpretation, wg *sync.WaitGroup) {
	for {
		if form.Val(in.Vals) {
			res <- in
			break
		}
		if !in.Next() {
			break
		}
	}
	wg.Done()
}
