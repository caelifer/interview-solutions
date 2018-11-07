// This program implements "fizz/buzz" solution in a rather overrengineered sort of way.
// Running example can be found here:
//
//     https://play.golang.org/p/budLggnTMXH
//
// This example studies the way to implement functional style programming in Go.
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Start of the program
func main() {
	// Create filtering pipeline and display resuluts
	makePipeline().
		AddFilter(makeLimitFilter(45)).        // Limit number of generated values in a stream.
		AddFilter(makeValueFilter(3, "fizz")). // Add tag "fizz" to all values devisible by 3.
		AddFilter(makeValueFilter(5, "buzz")). // Add tag "buzz" to all values devisible by 5.
		AddFilter(makeValueFilter(7, "zang")). // Add tag "boom" to all values devisible by 7.
		AddFilter(makeValueFilter(9, "bang")). // Add tag "bang" to all values devisible by 9.
		Run(func(v Val) { fmt.Println(v) })    // Collect and display filtered values.
}

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

type Processor func(Val)

type Sink chan<- Val

type Pipeline <-chan Val

// Pipeline factory
func makePipeline() Pipeline {
	ch := make(chan Val)
	go func() {
		for i := 1; ; i++ {
			ch <- Val{i, make([]string, 0, 1)}
		}
	}()
	return ch
}

func (p Pipeline) Run(proc Processor) {
	for v := range p {
		proc(v)
	}
}

func (p Pipeline) AddFilter(filter Filter) Pipeline {
	return filter(p)
}

// Filter is connecting pipelines
type Filter func(in Pipeline) Pipeline

// FilterFn is a function that does the actual filtering
type FilterFn func(out Sink, in Pipeline)

// apply creates Filter by applying provided FilterFn. It takes care of resources by closing
// channels at the appropriate time. It also handles panics in filtering funcitons.
func apply(filterFn FilterFn) Filter {
	return func(in Pipeline) Pipeline {
		out := make(chan Val)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("pipeline paniced: %v", err)
				}
				close(out) // Always close channel, even if filterFn() panics
			}()
			filterFn(out, in)
		}()
		return out
	}
}

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
				// Check if number if divisible by our divizor
				if v.i%div == 0 {
					v.s = append(v.s, tag)
				}
				out <- v
			}
		}
	}(num, tag))
}
