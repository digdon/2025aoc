package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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

	begin := time.Now()

	// Start by parsing out the present shapes
	presentRE := regexp.MustCompile(`^\d+:`)
	regionRE := regexp.MustCompile(`^\d+x\d+:`)

	presents := map[int]Present{}
	lineNum := 0

	for ; lineNum < len(inputLines); lineNum++ {
		line := inputLines[lineNum]

		if presentRE.MatchString(line) {
			lineNum++
			presentNum := 0
			_, err := fmt.Sscanf(line, "%d:", &presentNum)
			if err != nil {
				log.Println("Error parsing present number:", err)
				os.Exit(1)
			}

			// Parse out this present
			present := Present{}

			for ; lineNum < len(inputLines); lineNum++ {
				line = inputLines[lineNum]

				if line == "" {
					// End of this present. Add to map and break
					presents[presentNum] = present
					break
				}

				row := []bool{}
				for _, char := range line {
					if char == '#' {
						row = append(row, true)
						present.complexArea++
					} else {
						row = append(row, false)
					}
				}
				present.simpleArea += len(row)
				present.shape = append(present.shape, row)
			}
		} else if regionRE.MatchString(line) {
			// Found the start of the region definitions - we're done with presents
			break
		}
	}

	// Now pull out the regions and see if they can hold the specified shapes
	regions := []Region{}
	countRe := regexp.MustCompile(`\s+`)

	for ; lineNum < len(inputLines); lineNum++ {
		line := inputLines[lineNum]

		var x, y int
		idx := strings.Index(line, ": ")
		_, err := fmt.Sscanf(line[:idx], "%dx%d", &x, &y)
		if err != nil {
			log.Println("Error parsing region dimensions:", err)
			os.Exit(1)
		}

		countParts := countRe.Split(line[idx+2:], -1)
		counts := []int{}
		for _, part := range countParts {
			v, err := strconv.Atoi(part)
			if err != nil {
				log.Println("Error parsing present count:", err)
				os.Exit(1)
			}
			counts = append(counts, v)
		}

		region := Region{
			width:             x,
			length:            y,
			presentQuantities: counts,
		}
		regions = append(regions, region)
	}

	// Proccess each region - can they hold the specified presents?
	canFitCount := 0

	for _, region := range regions {
		fits := canFit(region, presents)
		// fmt.Printf("Region %dx%d can fit: %v\n", region.width, region.length, fits)

		if fits {
			canFitCount++
		}
	}

	fmt.Printf("Part 1: %d (%v)\n", canFitCount, time.Since(begin))
}

func canFit(region Region, presents map[int]Present) bool {
	regionArea := region.width * region.length

	// As a first sanity check, calculate a couple approximate areas needed for the specified presents and
	// see if they fit (or don't)
	totalSimpleArea := 0
	totalComplexArea := 0

	for presentNum, count := range region.presentQuantities {
		present := presents[presentNum]
		totalSimpleArea += present.simpleArea * count
		totalComplexArea += present.complexArea * count
	}

	// fmt.Printf("region area: %d, simple present area: %d, complex present area: %d\n", regionArea, totalSimpleArea, totalComplexArea)

	if totalSimpleArea <= regionArea {
		// Room for all of the presents, sitting side-by-side
		return true
	}

	if totalComplexArea > regionArea {
		// Not ernough room even if we could pack them perfectly
		return false
	}

	// More complex fitting logic to be added here

	return true
}

type Present struct {
	shape       [][]bool
	simpleArea  int
	complexArea int
}

type Region struct {
	width             int
	length            int
	presentQuantities []int
}
