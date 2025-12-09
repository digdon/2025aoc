package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"sort"
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

	points := []Point{}

	// Parse all of the 3D points from the input
	for _, line := range inputLines {
		var x, y, z int
		_, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		if err != nil {
			log.Println("Error parsing point:", err)
			os.Exit(1)
		}

		points = append(points, Point{X: x, Y: y, Z: z})
	}

	// Calculate the distances between every pair of points
	edges := []Edge{}

	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			edges = append(edges, Edge{A: points[i], B: points[j], distance: distance(points[i], points[j])})
		}
	}

	// Sort edges by distance, shortest to longest
	slices.SortFunc(edges, func(a, b Edge) int {
		return int(a.distance - b.distance)
	})

	// Now build circuit groups by adding edges in order of shortest distance
	circuitList := []map[Point]bool{}

	for i := range edges {
		if i == 1000 {
			// Part 1 stuff
			sort.Slice(circuitList, func(i, j int) bool {
				return len(circuitList[i]) > len(circuitList[j])
			})

			total := 1
			for i := range 3 {
				total *= len(circuitList[i])
			}

			fmt.Printf("Part 1: %d (%v)\n", total, time.Since(begin))
		}

		a, b := edges[i].A, edges[i].B

		// Check to see if a and/or b are part of any existing circuit
		existAIdx := findCircuitIndex(circuitList, a)
		existBIdx := findCircuitIndex(circuitList, b)

		if existAIdx == -1 && existBIdx == -1 {
			// Neither point is part of an existing circuit, so create a new one
			newCircuit := map[Point]bool{}
			newCircuit[a] = true
			newCircuit[b] = true
			circuitList = append(circuitList, newCircuit)
		} else if existAIdx != -1 && existBIdx == -1 {
			// Only a is part of an existing circuit, so add b to it
			circuitList[existAIdx][b] = true
		} else if existAIdx == -1 && existBIdx != -1 {
			// Only b is part of an existing circuit, so add a to it
			circuitList[existBIdx][a] = true
		} else if existAIdx != existBIdx {
			// Both points are part of different circuits, so merge them
			for p := range circuitList[existBIdx] {
				circuitList[existAIdx][p] = true
			}

			// Remove B circuit from the active list
			circuitList = append(circuitList[:existBIdx], circuitList[existBIdx+1:]...)
		}

		// Part 2 stuff
		if len(circuitList) == 1 && len(circuitList[0]) == len(points) {
			// All points are now connected
			xMult := a.X * b.X
			fmt.Printf("Part 2: %d (%v)\n", xMult, time.Since(begin))
			break
		}
	}
}

func findCircuitIndex(circuitList []map[Point]bool, p Point) int {
	for i, circuit := range circuitList {
		if circuit[p] {
			return i
		}
	}

	return -1
}

func distance(a, b Point) float64 {
	dx := float64(b.X - a.X)
	dy := float64(b.Y - a.Y)
	dz := float64(b.Z - a.Z)

	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

type Edge struct {
	A        Point
	B        Point
	distance float64
}

type Point struct {
	X int
	Y int
	Z int
}
