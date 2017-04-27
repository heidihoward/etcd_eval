// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// the file is borrowed from github.com/rakyll/boom/boomer/print.go

package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)


type result struct {
	errStr   string
	duration time.Duration
	happened time.Time
}

type report struct {
	results chan result
	errorDist map[string]int
	lats      []time.Duration
	times			[]time.Time
}

func printReport(results chan result) <-chan struct{} {
	return wrapReport(func() {
		r := &report{
			results:   results,
			errorDist: make(map[string]int),
		}
		r.finalize()
		r.print()
	})
}

func printRate(results chan result) <-chan struct{} {
	return wrapReport(func() {
		r := &report{
			results:   results,
			errorDist: make(map[string]int),
		}
		r.finalize()
	})
}

func wrapReport(f func()) <-chan struct{} {
	donec := make(chan struct{})
	go func() {
		defer close(donec)
		f()
	}()
	return donec
}

func (r *report) finalize() {
	//st := time.Now()
	for res := range r.results {
		if res.errStr != "" {
			r.errorDist[res.errStr]++
		} else {
			r.lats = append(r.lats, res.duration)
			r.times = append(r.times, res.happened)
		}
	}
}

func (r *report) print() {
		r.printLatencies()

	if len(r.errorDist) > 0 {
		r.printErrors()
	}
}

// Prints percentile latencies.
func (r *report) printLatencies() {

	filename := csvfile
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	stats := csv.NewWriter(file)
	defer stats.Flush()

	for i := 0; i < len(r.lats); i++ {
		latency := strconv.FormatInt(r.lats[i].Nanoseconds(), 10)
		stats.Write([]string{strconv.FormatInt(r.times[i].UnixNano(),10), strconv.Itoa(i), latency, strconv.Itoa(1)})

	}
	stats.Flush()
}



func (r *report) printErrors() {
	fmt.Printf("\nError distribution:\n")
	for err, num := range r.errorDist {
		fmt.Printf("  [%d]\t%s\n", num, err)
	}
}
