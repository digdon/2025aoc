package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	// fmt.Println(turn(11, RIGHT, 8))
	// fmt.Println(turn(0, LEFT, 1))
	// fmt.Println(turn(99, RIGHT, 1))
	// fmt.Println(turn(5, LEFT, 10))
	pos := 50
	zeroCount := 0
	zeroPassCount := 0

	for _, line := range inputLines {
		var dir rune
		var count int
		var zp int

		n, err := fmt.Sscanf(line, "%c%d", &dir, &count)
		if err != nil || n != 2 {
			log.Fatalf("Failed to parse line: %s", line)
		}
		fmt.Println(line)
		pos, zp = turn(pos, dir, count)

		if pos == 0 {
			zeroCount++
		}

		zeroPassCount += zp
	}

	fmt.Println("Part 1:", zeroCount)
	fmt.Println("Part 2:", zeroCount+zeroPassCount)
}

// type Direction int

// const (
// 	LEFT Direction = iota
// 	RIGHT
// )

func turn(origPos int, dir rune, count int) (int, int) {
	zpc := count / 100
	count %= 100
	newPos := origPos

	if dir == 'L' {
		newPos -= count

		if newPos < 0 {
			newPos += 100

			if newPos != 0 && origPos != 0 {
				zpc++
			}
		}
	} else {
		newPos += count

		if newPos > 99 {
			newPos -= 100

			if newPos != 0 && origPos != 0 {
				zpc++
			}
		}
	}

	return newPos, zpc
}
