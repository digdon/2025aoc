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

	var batteryBanks [][]int

	for _, line := range inputLines {
		var bank []int

		for _, char := range line {
			bank = append(bank, int(char-'0'))
		}

		batteryBanks = append(batteryBanks, bank)
	}

	begin := time.Now()

	fmt.Println("Part 1:", processBanks(batteryBanks, 2))
	fmt.Println("Part 2:", processBanks(batteryBanks, 12))

	fmt.Println("Elapsed time:", time.Since(begin))
}

func processBanks(batteryBanks [][]int, digitCount int) int {
	totalJoltage := 0

	for _, bank := range batteryBanks {
		joltage := process(bank, 0, 0, digitCount)
		totalJoltage += joltage
	}

	return totalJoltage
}

func process(bank []int, currentJoltage int, pos int, digitCount int) int {
	if digitCount == 0 {
		// Collected enough digits - we're done
		return currentJoltage
	}

	// From the starting position, find the highest number that still has enough remaining digits
	var maxValue, maxPos int

	for i := pos; i < len(bank)-digitCount+1; i++ {
		if bank[i] > maxValue {
			maxValue = bank[i]
			maxPos = i
		}
	}

	joltage := (currentJoltage * 10) + maxValue
	return process(bank, joltage, maxPos+1, digitCount-1)
}

func partOne(batteryBanks [][]int) {
	var totalJoltage int

	for _, bank := range batteryBanks {
		var maxTens, maxJoltage int

		for i := 0; i < len(bank)-1; i++ {
			tens := bank[i] * 10
			if tens <= maxTens {
				continue
			}

			for j := i + 1; j < len(bank); j++ {
				joltage := tens + bank[j]

				if joltage > maxJoltage {
					maxJoltage = joltage
					maxTens = tens
				}
			}
		}

		totalJoltage += maxJoltage
	}

	fmt.Println("Part 1:", totalJoltage)
}
