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

	begin := time.Now()

	// Find the start position
	var startCol int

	for col := range inputLines[0] {
		if inputLines[0][col] == 'S' {
			startCol = col
			break
		}
	}

	// Place the first beam
	beams := map[Point]int{}
	beams[Point{X: startCol, Y: 1}] = 1
	splitterCount := 0

	// Now we sart moving down the grid, starting from row 2 (0-based)
	for y := 2; y < len(inputLines); y++ {
		for x := range inputLines[y] {
			switch inputLines[y][x] {
			case '.':
				// If there's a beam directly above, it continues down
				if beams[Point{X: x, Y: y - 1}] > 0 {
					beams[Point{X: x, Y: y}] += beams[Point{X: x, Y: y - 1}]
				}
			case '^':
				// Found a splitter. If there's a beam above, split it left and right
				if beams[Point{X: x, Y: y - 1}] > 0 {
					splitterCount++

					if inputLines[y][x-1] == '.' {
						beams[Point{X: x - 1, Y: y}] += beams[Point{X: x, Y: y - 1}]
					}

					if inputLines[y][x+1] == '.' {
						beams[Point{X: x + 1, Y: y}] += beams[Point{X: x, Y: y - 1}]
					}
				}
			}
		}
	}

	pathCount := 0
	for x := range inputLines[len(inputLines)-1] {
		pathCount += beams[Point{X: x, Y: len(inputLines) - 1}]
	}

	fmt.Println("Part 1:", splitterCount)
	fmt.Println("Part 2:", pathCount)
	fmt.Println("Time:", time.Since(begin))
}

type Point struct {
	X int
	Y int
}
