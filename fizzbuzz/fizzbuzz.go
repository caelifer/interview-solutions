// This program implements "fizz/buzz" solution in an overrengineered sort of way.
//
// Running example can be found here:
//
//     https://play.golang.org/p/budLggnTMXH
//
// This example studies the way to implement functional programming in Go.

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const DefaultUppperLimit = 25

// Start of the program
func main() {
	upperLimit := DefaultUppperLimit

	if len(os.Args) > 1 {
		// Take first argument
		if n, err := strconv.Atoi(os.Args[1]); err != nil {
			log.Printf("Failed to parse %q: %v", os.Args[1], err)
			log.Printf("Using default upper limit: %d", upperLimit)
		} else {
			upperLimit = n
		}
	}

	// Create filtering pipeline and display resuluts
	makePipeline().
		AddFilter(makeLimitFilter(upperLimit)). // Limit number of generated values in a stream.
		AddFilter(makeValueFilter(3, "fizz")).  // Add tag "fizz" to all values devisible by 3.
		AddFilter(makeValueFilter(5, "buzz")).  // Add tag "buzz" to all values devisible by 5.
		AddFilter(makeValueFilter(7, "bang")).  // Add tag "bang" to all values devisible by 7.
		AddFilter(makeValueFilter(9, "zang")).  // Add tag "zang" to all values devisible by 9.
		Run(func(v Val) { fmt.Println(v) })     // Run pipeline and display filtered values.
}

// Val is a supporting type to hold and represent enumerated value
type Val struct {
	i int
	s []string
}

// String pretty printer for Val type. It outputs either space-joined V.s or V.i in string form.
func (v Val) String() string {
	s := strconv.Itoa(v.i)
	if len(v.s) > 0 {
		s = strings.Join(v.s, " ")
	}
	return s
}

// Pipeline helper types

// Processor is a functor object used to process each value in a filtered pipeline.
type Processor func(Val)

// Sink is a helper type to denote collector in the pipeline chain.
type Sink chan<- Val

// Pipeline is a helper object to represent generator in the pipeline chain.
type Pipeline <-chan Val

// makePipeline is a Pipeline factory function. It creates generator responsible for
// creating initial stream of Val values.
func makePipeline() Pipeline {
	ch := make(chan Val)
	go func() {
		for i := 1; ; i++ {
			ch <- Val{i, make([]string, 0, 1)}
		}
	}()
	return ch
}

// Run method wraps final collector.
func (p Pipeline) Run(proc Processor) {
	for v := range p {
		proc(v)
	}
}

// AddFilter method inserts new Filter node into a Pipeline.
func (p Pipeline) AddFilter(filter Filter) Pipeline {
	return filter(p)
}

// Filter is a helper type that create new generator, based on the processed stream from `in` Pipeline, effectively
// creating a new node in the chain of connected Pipeline objects.
type Filter func(in Pipeline) Pipeline

// FilterFn is a function that does the actual filtering by recieving a Val object from `in`, processing it, and
// optionally passing it out to `out` Sink.
type FilterFn func(out Sink, in Pipeline)

// Limiting filter generator implemented as closure to capture provided parameter.
func makeLimitFilter(limit int) Filter {
	return apply(func(limit int) FilterFn {
		return func(out Sink, in Pipeline) {
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

// Value filter generator implemented as closure to capture provided parameter.
func makeValueFilter(num int, tag string) Filter {
	return apply(func(div int, tag string) FilterFn {
		return func(out Sink, in Pipeline) {
			for {
				v, ok := <-in
				if !ok {
					break
				}
				// Check if number if divisible by our divizor.
				if v.i%div == 0 {
					v.s = append(v.s, tag)
				}
				out <- v
			}
		}
	}(num, tag))
}

// apply constructs a new Filter by applying provided FilterFn. It takes care of resources by closing
// `out` channels at the appropriate time. It also handles panics in filtering funcitons.
func apply(filterFn FilterFn) Filter {
	return func(in Pipeline) Pipeline {
		out := make(chan Val)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("pipeline paniced: %v", err)
				}
				close(out) // Always close channel, even if filterFn() panics.
			}()
			filterFn(out, in)
		}()
		return out
	}
}
