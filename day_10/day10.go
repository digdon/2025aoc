package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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

	// Parse out all of the machine data
	machines := []Machine{}

	for _, line := range inputLines {
		machines = append(machines, parseMachine(line))
	}

	// Part 1 stuff
	begin := time.Now()

	totalPresses := 0

	for _, machine := range machines {
		totalPresses += part1MinPresses(machine)
	}

	fmt.Printf("Part 1: %d (%v)\n", totalPresses, time.Since(begin))

	// Part 2 stuff
	begin = time.Now()

	totalPresses = 0

	for _, machine := range machines {
		patterns := generatePatterns(len(machine.joltages), machine.buttons)
		cache := map[string]int{}
		totalPresses += part2MinPresses(machine.joltages, patterns, cache)
	}

	fmt.Printf("Part 2: %d (%v)\n", totalPresses, time.Since(begin))
}

type Machine struct {
	lights   rune
	buttons  [][]int
	joltages []int
}

func parseMachine(line string) Machine {
	parts := strings.Split(line, " ")

	// Parse out the target light config
	lights := rune(0)
	lightsString := parts[0][1 : len(parts[0])-1]

	for i := range lightsString {
		if lightsString[i] == '#' {
			lights |= 1 << i
		}
	}

	// Parse out the joltages
	joltages := []int{}
	joltageStrings := strings.Split(parts[len(parts)-1][1:len(parts[len(parts)-1])-1], ",")

	for i := range joltageStrings {
		val, err := strconv.Atoi(joltageStrings[i])

		if err != nil {
			log.Fatalf("Error parsing joltage value: %v", err)
		}

		joltages = append(joltages, val)
	}

	// Parse out the buttons. A button is represented as a list of light indices it toggles.
	buttons := [][]int{}

	for i := 1; i < len(parts)-1; i++ {
		buttonValues := strings.Split(parts[i][1:len(parts[i])-1], ",")
		button := []int{}

		for j := range buttonValues {
			val, err := strconv.Atoi(buttonValues[j])

			if err != nil {
				log.Fatalf("Error parsing button value: %v", err)
			}

			button = append(button, val)
		}

		buttons = append(buttons, button)
	}

	return Machine{lights: lights, buttons: buttons, joltages: joltages}
}

func part1MinPresses(machine Machine) int {
	minPresses := math.MaxInt

	for i := range 1 << len(machine.buttons) {
		// Iterate through all possible button press combinations (button numbers represented by bits in i)
		lights := rune(0)
		presses := 0

		for j := range machine.buttons {
			// Check to see if button j is pressed in this combination (button number bit in i is set)
			if (i & (1 << j)) != 0 {
				// Press the button. This involves looking at each light index the button toggles and flipping that bit in lights.
				for _, lightIndex := range machine.buttons[j] {
					lights ^= 1 << lightIndex
				}

				presses++
			}
		}

		if lights == machine.lights {
			minPresses = min(minPresses, presses)
		}
	}

	return minPresses
}

type PatternInfo struct {
	pattern []int
	cost    int
}

func generatePatterns(numLights int, buttons [][]int) map[rune]map[string]PatternInfo {
	patterns := map[rune]map[string]PatternInfo{}

	for i := range 1 << len(buttons) {
		lightCounts := make([]int, numLights)
		presses := 0

		for j := range buttons {
			if (i & (1 << j)) != 0 {
				// press button j
				for _, lightIndex := range buttons[j] {
					lightCounts[lightIndex]++
				}

				presses++
			}
		}

		parityPattern := rune(0)

		for k := range lightCounts {
			if lightCounts[k]%2 != 0 {
				parityPattern |= 1 << k
			}
		}

		// Add to patterns map
		ppMap, exists := patterns[parityPattern]

		if !exists {
			ppMap = map[string]PatternInfo{}
			patterns[parityPattern] = ppMap
		}

		ltk := fmt.Sprintf("%v", lightCounts)
		item, exists := ppMap[ltk]

		if !exists || presses < item.cost {
			ppMap[ltk] = PatternInfo{pattern: lightCounts, cost: presses}
		}
	}

	return patterns
}

func part2MinPresses(joltages []int, patterns map[rune]map[string]PatternInfo, cache map[string]int) int {
	// Check if we've reached all-zero joltages
	allZero := true

	for _, joltage := range joltages {
		if joltage != 0 {
			allZero = false
			break
		}
	}

	if allZero {
		return 0
	}

	// Create a cache key based on current joltages
	key := fmt.Sprintf("%v", joltages)

	// Check cache
	if val, ok := cache[key]; ok {
		return val
	}

	// Determine current parity pattern
	parityPattern := rune(0)

	for i, joltage := range joltages {
		if joltage%2 != 0 {
			parityPattern |= 1 << i
		}
	}

	// Look up possible patterns for this parity
	patternOptions, exists := patterns[parityPattern]

	if !exists {
		// No patterns available for this parity, return a large number
		return 1000000
	}

	minPresses := 1000000

	for _, patternInfo := range patternOptions {
		// Check to see if we can apply this particular pattern. Unlike part 1, which toggles the lights,
		// the joltage for a particular light increases every time the light is hit by a button. So we have
		// to check to see if the pattern (which represents a series of button presses) would result in a light
		// being hit too many times. Because we are subtracting the pattern from the joltages, too many
		// presses would result in a negative joltage for a particular light.
		canApply := true

		for i := range joltages {
			if joltages[i]-patternInfo.pattern[i] < 0 {
				canApply = false
				break
			}
		}

		if !canApply {
			continue
		}

		// Pattern is valid, so apply it to the joltages. At the same time, we know that after applying the pattern,
		// all resulting joltages will be even, so we can halve them for the next recursion.
		newJoltages := make([]int, len(joltages))

		for i := range joltages {
			newJoltages[i] = (joltages[i] - patternInfo.pattern[i]) / 2
		}

		// Recurse. Because we halved the joltages, we need to multiply the resulting presses by 2
		// (to account for the doubling effect of halving). Plus we add the cost of this pattern.
		minPresses = min(minPresses, patternInfo.cost+part2MinPresses(newJoltages, patterns, cache)*2)
	}

	cache[key] = minPresses

	return minPresses
}
