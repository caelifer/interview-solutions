package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Val is a supporting type to hold and represent enumerated value
type Val struct {
	i int
	s []string
}

// String pretty printer for Val type. It outputs either space-joined V.s,
// or V.i in string form.
func (v Val) String() string {
	s := strconv.Itoa(v.i)
	if len(v.s) > 0 {
		s = strings.Join(v.s, " ")
	}
	return s
}

// Start of the program
func main() {
	// Create filtering pipeline
	in := createPipeline(
		makeGenerator(),
		makeLimitFilter(45),
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
	// Start generator
	ch := gen()
	for _, f := range filters {
		// Add filter to the chain
		ch = f(ch)
	}
	// Return resulting channel
	return ch
}

// Generator type
type Generator func() <-chan Val

// Create Val generator
func makeGenerator() Generator {
	return func() <-chan Val {
		ch := make(chan Val)
		go func() {
			for i := 1; ; i++ {
				ch <- Val{i, make([]string, 0, 1)}
			}
			close(ch)
		}()
		return ch
	}
}

// Filter is generic interface to a filtering function
type Filter func(in <-chan Val) <-chan Val

// Logic function that does the actual filtering
type FilterLogicFn func(out chan<- Val, in <-chan Val)

// apply creats Filter by applying specified FilterLogicFn
func apply(fl FilterLogicFn) Filter {
	return func(in <-chan Val) <-chan Val {
		out := make(chan Val)
		go func() {
			defer close(out) // Always close channel, even if fl() panics
			fl(out, in)
		}()
		return out
	}
}

// Limit filter generated implementation.
// Must use function to capture parameter in a closure
func makeLimitFilter(limit int) Filter {
	return apply(func(limit int) FilterLogicFn {
		return func(out chan<- Val, in <-chan Val) {
			for i := 0; i < limit; i++ {
				v, ok := <-in
				if !ok {
					break
				}
				// passthrough
				out <- v
			}
		}
	}(limit))
}

// Value filter generator implementaion
// Must use function to capture parameter in a closure
func makeValueFilter(num int, msg string) Filter {
	return apply(func(div int, msg string) FilterLogicFn {
		return func(out chan<- Val, in <-chan Val) {
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
		}
	}(num, msg))
}
