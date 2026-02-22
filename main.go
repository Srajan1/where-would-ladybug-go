package main

import (
	"fmt"
	"math/rand"
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

func getLastPosition() int {

	finalPosition := -1

	// I'll create a set of numbers between 1 and 12 and make the ladybug traverse the clock.
	// The moment it reaches a number I'll remove it from the set.
	// Last number to be removed will be last position.
	currentPosition := 12
	remainingPosition := make(map[int]bool)

	for i := 1; i <= 12; i++ {
		remainingPosition[i] = true
	}

	for {
		delete(remainingPosition, currentPosition)
		if len(remainingPosition) == 0 {
			finalPosition = currentPosition
			break
		}

		moveClockWise := rand.Intn(2) == 0

		nextPosition := getNextPosition(currentPosition, moveClockWise)
		currentPosition = nextPosition

	}
	return finalPosition

}

func main() {
	fmt.Println("Starting the sim.")

	lastPosition := make(map[int]int)
	for i := 0; i < MILLION; i++ {
		lastPosition[getLastPosition()]++
	}
	fmt.Println(lastPosition)
}
