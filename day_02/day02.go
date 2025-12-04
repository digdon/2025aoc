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
	// input := "11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124"
	// inputLines = append(inputLines, input)

	var ranges []Range

	for _, line := range inputLines {
		parts := strings.SplitSeq(line, ",")

		for part := range parts {
			var start, end int
			n, err := fmt.Sscanf(part, "%d-%d", &start, &end)
			if err != nil || n != 2 {
				fmt.Printf("Failed to parse range: %s\n", part)
				return
			}

			ranges = append(ranges, Range{Start: start, End: end})
		}
	}

	part1Total := 0
	part2Total := 0

	begin := time.Now()

	for _, r := range ranges {
		for i := r.Start; i <= r.End; i++ {
			str := strconv.Itoa(i)
			fullLen, halfLen := len(str), len(str)/2

			// Part 1 stuff
			if len(str)%2 == 0 {
				if str[0:halfLen] == str[halfLen:] {
					// fmt.Println(str)
					part1Total += i
				}
			}

			// Part 2 stuff
			for j := range halfLen {
				if fullLen%(j+1) != 0 {
					// Skip over anything we know can't match based on length
					continue
				}

				// Build comparison string
				re := strings.Repeat(str[:j+1], fullLen/(j+1))

				if re == str {
					// fmt.Println(str)
					part2Total += i
					break
				}
			}
		}
	}

	fmt.Println("Part 1:", part1Total)
	fmt.Println("Part 2:", part2Total)
	fmt.Println("Elapsed time:", time.Since(begin))
}

type Range struct {
	Start int
	End   int
}
