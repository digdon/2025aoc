package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var inputLines []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		inputLines = append(inputLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	points := map[Point]bool{}

	for y, line := range inputLines {
		for x, char := range line {
			if char == '@' {
				points[Point{X: x, Y: y}] = true
			}
		}
	}

	begin := time.Now()

	// Part 1 stuff
	moveablePoints := getMoveablePoints(points)
	fmt.Printf("Part 1: %d (%v)\n", len(moveablePoints), time.Since(begin))

	// Part 2 stuff
	begin = time.Now()
	totalMoved := 0

	for len(moveablePoints) > 0 {
		totalMoved += len(moveablePoints)

		for _, p := range moveablePoints {
			delete(points, p)
		}

		moveablePoints = getMoveablePoints(points)
	}

	fmt.Printf("Part 2: %d (%v)\n", totalMoved, time.Since(begin))
}

func getMoveablePoints(points map[Point]bool) []Point {
	moveablePoints := []Point{}

	for p := range points {
		rollcount := 0

		for _, d := range Directions {
			nx, ny := p.X+d[0], p.Y+d[1]

			if points[Point{X: nx, Y: ny}] {
				rollcount++
			}
		}

		if rollcount < 4 {
			moveablePoints = append(moveablePoints, p)
		}
	}

	return moveablePoints
}

var Directions = [][]int{
	{-1, -1},
	{0, -1},
	{1, -1},
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
}

type Point struct {
	X int
	Y int
}
