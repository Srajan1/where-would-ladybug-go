//go:build parallel

package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var MILLION = 1000000

func getNextPosition(currentPosition int, clockwise bool) int {
	nextPosition := -1
	if clockwise {
		nextPosition = currentPosition + 1
		if nextPosition > 12 {
			nextPosition = 1
		}
	} else {
		nextPosition = currentPosition - 1
		if nextPosition < 1 {
			nextPosition = 12
		}
	}
	return nextPosition
}

func getLastPosition(r *rand.Rand) int {
	currentPosition := 12
	remainingPosition := make(map[int]bool, 12)

	for i := 1; i <= 12; i++ {
		remainingPosition[i] = true
	}

	for {
		delete(remainingPosition, currentPosition)
		if len(remainingPosition) == 0 {
			return currentPosition
		}

		moveClockWise := r.Intn(2) == 0
		currentPosition = getNextPosition(currentPosition, moveClockWise)
	}
}

func main() {
	fmt.Println("Starting the simulation concurrently. Running for a million iterations.")

	// I'll divide the million iterations equally between CPU threads which will be hanlded by workers.s
	numWorkers := runtime.NumCPU()
	chunk := MILLION / numWorkers

	fmt.Printf("Number of workers: %d\n", numWorkers)

	var wg sync.WaitGroup
	results := make([]map[int]int, numWorkers)

	start := time.Now()

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			// Turns out rand.Intn(2) aquires cpu lock, which would cause all go routines
			// to fight with each other to aquire it. So I'm passing the random from here.
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))

			localCount := make(map[int]int)
			start := workerID * chunk
			end := start + chunk

			for i := start; i < end; i++ {
				pos := getLastPosition(r)
				localCount[pos]++
			}

			results[workerID] = localCount
		}(w)
	}

	wg.Wait()

	// Merge results
	finalCombinedPositions := make(map[int]int)
	for _, m := range results {
		for k, v := range m {
			finalCombinedPositions[k] += v
		}
	}

	t := time.Now()
	elapsed := t.Sub(start)

	for i := 1; i <= 11; i++ {
		fmt.Printf("Number of times ladybug ended up last at %d is: %d\n", i, finalCombinedPositions[i])
	}

	fmt.Printf("\n\nTime taken to run all simulations in series is %dms\n", elapsed.Milliseconds())

}
