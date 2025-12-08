package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	part1(inputLines)
	part2(inputLines)
}

func part1(inputLines []string) {
	begin := time.Now()

	grid := [][]int{}
	operations := []rune{}

	for _, line := range inputLines {
		fields := strings.Fields(line)

		if fields[0][0] == '*' || fields[0][0] == '+' {
			for _, value := range fields {
				operations = append(operations, rune(value[0]))
			}
		} else {
			row := []int{}

			for _, value := range fields {
				value, err := strconv.Atoi(value)
				if err != nil {
					log.Println("Error parsing grid value:", err)
					os.Exit(1)
				}
				row = append(row, value)
			}

			grid = append(grid, row)
		}
	}

	colValues := []int{}

	for col, op := range operations {
		var value int

		if op == '*' {
			value = 1
		}

		for _, row := range grid {
			switch op {
			case '*':
				value *= row[col]
			case '+':
				value += row[col]
			}
		}

		colValues = append(colValues, value)
	}

	total := 0

	for _, v := range colValues {
		total += v
	}

	fmt.Printf("Part 1: %d (%v)\n", total, time.Since(begin))
}

func part2(inputLines []string) {
	begin := time.Now()

	// Start by finding the start and end positions of each column, based on location of the operations
	colRanges := []Range{}
	operations := []rune{}
	start := len(inputLines[len(inputLines)-1]) - 1
	end := start
	line := inputLines[len(inputLines)-1]
	for pos := start; pos >= 0; pos-- {
		if line[pos] != ' ' {
			operations = append(operations, rune(line[pos]))
			colRanges = append(colRanges, Range{Start: pos, End: end})
			end = pos - 2
		}
	}

	maxDigits := len(inputLines) - 1
	grandTotal := 0

	// Now we process each column range. For each column, we extract the digits from each row,
	// build the numbers, then apply the operation to get the column total. Finally, we sum
	// all the column totals to get the grand total.
	for i, cr := range colRanges {
		colValues := []int{}

		for col := cr.Start; col <= cr.End; col++ {
			value := 0

			for row := range maxDigits {
				char := inputLines[row][col]

				if char == ' ' {
					continue
				}

				digit := int(char - '0')
				value = value*10 + digit
			}

			colValues = append(colValues, value)
		}

		op := operations[i]
		colTotal := 0

		if op == '*' {
			colTotal = 1
		}

		for _, v := range colValues {
			switch op {
			case '*':
				colTotal *= v
			case '+':
				colTotal += v
			}
		}

		grandTotal += colTotal
	}

	fmt.Printf("Part 2: %d (%v)\n", grandTotal, time.Since(begin))
}

type Range struct {
	Start int
	End   int
}
