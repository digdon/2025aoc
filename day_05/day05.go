package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

	rangeMode := true
	ranges := []Range{}
	idList := []int{}

	for _, line := range inputLines {
		if line == "" {
			rangeMode = false
			continue
		}

		if rangeMode {
			var r Range
			_, err := fmt.Sscanf(line, "%d-%d", &r.Start, &r.End)
			if err != nil {
				log.Println("Error parsing range:", err)
				os.Exit(1)
			}
			ranges = append(ranges, r)
		} else {
			id, err := strconv.Atoi(line)
			if err != nil {
				log.Println("Error parsing ID:", err)
				os.Exit(1)
			}
			idList = append(idList, id)
		}
	}

	// Part 1 stuff
	begin := time.Now()

	part1FreshCount := 0

	for _, id := range idList {
		for _, r := range ranges {
			if id >= r.Start && id <= r.End {
				part1FreshCount++
				break
			}
		}
	}

	fmt.Printf("Part 1: %d (%v)\n", part1FreshCount, time.Since(begin))

	// Part 2 stuff
	begin = time.Now()

	mergedRanges := []Range{}

	for _, r := range ranges {
		newRanges := []Range{}

		contained := false

		for _, mr := range mergedRanges {
			if r.Start >= mr.Start && r.Start <= mr.End && r.End >= mr.Start && r.End <= mr.End {
				// Fully-encompassed by merged range - split to make room
				// newRanges = append(newRanges, Range{Start: mr.Start, End: r.Start - 1})
				// newRanges = append(newRanges, Range{Start: r.End + 1, End: mr.End})
				// Fully-encompassed by merged range - keep the merged range and ignore the new one
				newRanges = append(newRanges, mr)
				contained = true
			} else if r.Start <= mr.Start && r.End >= mr.End {
				// Fully-encompasses merged range - ignore the merged range
				continue
			} else if r.Start < mr.Start && r.End >= mr.Start && r.End <= mr.End {
				// Overlaps at the start of merged range
				newRanges = append(newRanges, Range{Start: r.End + 1, End: mr.End})
			} else if r.Start >= mr.Start && r.Start <= mr.End && r.End > mr.End {
				// Overlaps at the end of merged range
				newRanges = append(newRanges, Range{Start: mr.Start, End: r.Start - 1})
			} else {
				// No overlap
				newRanges = append(newRanges, mr)
			}
		}

		if !contained {
			newRanges = append(newRanges, r)
		}

		mergedRanges = newRanges
	}

	part2FreshCount := 0

	for _, r := range mergedRanges {
		part2FreshCount += r.End - r.Start + 1
	}

	fmt.Printf("Part 2: %d (%v)\n", part2FreshCount, time.Since(begin))
}

type Range struct {
	Start int
	End   int
}
