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

	points := []Point{}

	for _, line := range inputLines {
		var x, y int
		_, err := fmt.Sscanf(line, "%d,%d", &x, &y)
		if err != nil {
			log.Println("Error parsing point:", err)
			os.Exit(1)
		}

		points = append(points, Point{X: x, Y: y})
	}

	part1(points)
	part2(points)
}

func part1(points []Point) {
	begin := time.Now()

	var maxA, maxB Point
	var maxArea int

	for i := 0; i < len(points)-1; i++ {
		a := points[i]

		for j := i + 1; j < len(points); j++ {
			b := points[j]

			aDiff := a.X - b.X
			if aDiff < 0 {
				aDiff = -aDiff
			}

			bDiff := a.Y - b.Y
			if bDiff < 0 {
				bDiff = -bDiff
			}

			area := (aDiff + 1) * (bDiff + 1)

			if area > maxArea {
				maxArea = area
				maxA = a
				maxB = b
			}
		}
	}

	fmt.Printf("%v - %v: %d\n", maxA, maxB, maxArea)
	fmt.Printf("Part 1: %d (%v)\n", maxArea, time.Since(begin))
}

func part2(points []Point) {
	begin := time.Now()

	// Find all of the edges, categorized into horizontal and vertical
	cw := false
	horizontalEdges, verticalEdges := calculateEdges(points, cw)

	// Build a rectangle for each pair of points
	var maxA, maxB Point
	var maxArea int

	for i := 0; i < len(points)-1; i++ {
		a := points[i]

		for j := i + 1; j < len(points); j++ {
			b := points[j]

			topLeft := Point{X: min(a.X, b.X), Y: min(a.Y, b.Y)}
			bottomRight := Point{X: max(a.X, b.X), Y: max(a.Y, b.Y)}

			good := containedSimplified(topLeft, bottomRight, horizontalEdges, verticalEdges)
			// fmt.Printf("%v %v: %v\n", topLeft, bottomRight, good)
			// if contained(topLeft, bottomRight, horizontalEdges, verticalEdges) {
			if good {
				// Rectangle is fully contained
				area := (bottomRight.X - topLeft.X + 1) * (bottomRight.Y - topLeft.Y + 1)

				if area > maxArea {
					maxArea = area
					maxA = a
					maxB = b
				}
			}
		}
	}

	// fmt.Println(horizontalEdges)
	// fmt.Println(verticalEdges)
	fmt.Printf("%v - %v: %d\n", maxA, maxB, maxArea)
	fmt.Printf("Part 2: %d (%v)\n", maxArea, time.Since(begin))
}

func shoelace(points []Point) int {
	var area int

	for i := range points {
		var a, b Point
		a = points[i]

		if i == len(points)-1 {
			b = points[0]
		} else {
			b = points[i+1]
		}

		area += (a.X * b.Y) - (a.Y * b.X)
	}

	return area / 2
}

func calculateEdges(points []Point, cw bool) ([]Edge, []Edge) {
	// Work out each edge. Organize them based on horizontal vs vertical. For each edge, record if the left side is "inside"
	horizontalEdges, verticalEdges := []Edge{}, []Edge{}

	for i := range points {
		var a, b Point
		a = points[i]

		if i == len(points)-1 {
			b = points[0]
		} else {
			b = points[i+1]
		}

		if a.X == b.X {
			// Vertical edge
			if b.Y > a.Y {
				verticalEdges = append(verticalEdges, Edge{A: a, B: b, left: cw})
			} else {
				verticalEdges = append(verticalEdges, Edge{A: b, B: a, left: !cw})
			}
		}

		if a.Y == b.Y {
			// Horizontal edge
			if b.X > a.X {
				horizontalEdges = append(horizontalEdges, Edge{A: a, B: b, left: !cw})
			} else {
				horizontalEdges = append(horizontalEdges, Edge{A: b, B: a, left: cw})
			}
		}
	}

	return horizontalEdges, verticalEdges
}

// This is a simplified version of contained() that only checks if edges fall fully within the rectangle. This probably
// works for the specific input data, but is not a general solution.
func containedSimplified(topLeft, bottomRight Point, horizontalEdges, verticalEdges []Edge) bool {
	// Check horizontal edges
	for _, edge := range horizontalEdges {
		// Does the edge horizontally line up with the rectangle?
		if edge.A.X >= bottomRight.X || edge.B.X <= topLeft.X {
			// Falls completely outside horizontally, so ignore it
			continue
		}

		if edge.A.Y > topLeft.Y && edge.A.Y < bottomRight.Y {
			// Edge falls within the rectangle, so it's not contained
			return false
		}
	}

	// Check vertical edges
	for _, edge := range verticalEdges {
		// Does the edge vertically line up with the rectangle?
		if edge.A.Y >= bottomRight.Y || edge.B.Y <= topLeft.Y {
			// Falls completely outside vertically, so ignore it
			continue
		}

		if edge.A.X > topLeft.X && edge.A.X < bottomRight.X {
			// Edge falls within the rectangle, so it's not contained
			return false
		}
	}

	return true
}

func contained(topLeft, bottomRight Point, horizontalEdges, verticalEdges []Edge) bool {
	// Check horizontal edges
	for _, edge := range horizontalEdges {
		// Does the edge horizontally line up with the rectangle?
		if edge.A.X >= bottomRight.X || edge.B.X <= topLeft.X {
			// Falls completely outside horizontally, so ignore it
			continue
		}

		// Potentialy involved, so perform additional checks

		// Does the edge fall within the rectangle vertically?
		if edge.A.Y > topLeft.Y && edge.A.Y < bottomRight.Y {
			// Edge falls within the rectangle, so it's not contained
			return false
		}

		// Now check to see if the rectangle is "outside" the edge
		if (edge.left && bottomRight.Y > edge.A.Y) || (!edge.left && topLeft.Y < edge.A.Y) {
			return false
		}
	}

	// Check vertical edges
	for _, edge := range verticalEdges {
		// Does the edge vertically line up with the rectangle?
		if edge.A.Y >= bottomRight.Y || edge.B.Y <= topLeft.Y {
			// Falls completely outside vertically, so ignore it
			continue
		}

		// Potentialy involved, so perform additional checks

		// Does the edge fall within the rectangle horizontally?
		if edge.A.X > topLeft.X && edge.A.X < bottomRight.X {
			// Edge falls within the rectangle, so it's not contained
			return false
		}

		// Now check to see if the rectangle is "outside" the edge
		if (edge.left && bottomRight.X > edge.A.X) || (!edge.left && topLeft.X < edge.A.X) {
			return false
		}
	}

	return true
}

type Point struct {
	X int
	Y int
}

type Edge struct {
	A    Point
	B    Point
	left bool
}
