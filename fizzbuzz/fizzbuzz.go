package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	// Create processing/filtering pipeline
	in := createPipeline(
		makeGenerator(),
		makeLimitFilter(15),
		makeValueFilter(3, "fizz"),
		makeValueFilter(5, "buzz"),
		makeValueFilter(6, "boom"),
		makeValueFilter(9, "bang"),
	)

	// Display all values that are passed through the pipeline
	for v := range in {
		fmt.Println(v)
	}
}

// Pipeline factory
func createPipeline(gen Generator, filters ...Filter) <-chan Val {
	ch := gen()

	for _, f := range filters {
		// Add filter to the chain
		ch = f(ch)
	}

	return ch
}

// Val is a supporting type to hold and represent enumerated value
type Val struct {
	i int
	s []string
}

// String - outputs either stringify number or accumulated phrase
func (v Val) String() string {
	s := strconv.Itoa(v.i)

	if len(v.s) > 0 {
		s = strings.Join(v.s, " ")
	}

	return s
}

// Generator type
type Generator func() <-chan Val

// Create Val generator
func makeGenerator() Generator {
	return func() <-chan Val {
		ch := make(chan Val)

		go func() {
			for i := 1; ; i++ {
				ch <- Val{i: i, s: make([]string, 0, 2)}
			}
		}()

		return ch
	}
}

// Filter is generic interface to a filtering function
type Filter func(in <-chan Val) <-chan Val

// Logic function that does the actual filterign
type FilterLogic func(out chan<- Val, in <-chan Val)

func genericFilterFuction(filter FilterLogic) Filter {
	return func(in <-chan Val) <-chan Val {
		out := make(chan Val)
		go func() {
			defer close(out)
			filter(out, in)
		}()
		return out
	}
}

// Limit filter generator
func makeLimitFilter(limit int) Filter {
	return genericFilterFuction(func(out chan<- Val, in <-chan Val) {
		for i := 0; i < limit; i++ {
			v, ok := <-in
			if !ok {
				break
			}

			// passthrough
			out <- v
		}
	})
}

// Value fliter generator
func makeValueFilter(div int, msg string) Filter {
	return genericFilterFuction(func(out chan<- Val, in <-chan Val) {
		for {
			v, ok := <-in
			if !ok {
				break
			}

			// Check if number if divisible by our divizor
			if v.i%div == 0 {
				v.s = append(v.s, msg)
			}
			out <- v
		}
	})
}
