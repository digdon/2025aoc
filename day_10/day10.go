package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
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

	begin := time.Now()

	machines := []Machine{}

	for _, line := range inputLines {
		machine := parseMachine(line)
		machines = append(machines, machine)
	}

	fmt.Printf("Parsing time: %v\n", time.Since(begin))

	// Part 1 stuff
	begin = time.Now()

	totalPresses := 0

	for _, machine := range machines {
		// buttons := solvePart1(machine)
		// totalPresses += len(buttons)
		totalPresses += solvePart1New(machine)
	}

	fmt.Printf("Part 1: %d (%v)\n", totalPresses, time.Since(begin))
}

type Machine struct {
	origString string
	lights     rune
	buttons    []rune
	joltages   []int
}

func parseMachine(line string) Machine {
	// find lights and joltages positions
	lightStart, lightEnd := -1, -1
	joltageStart, joltageEnd := -1, -1

	for i, char := range line {
		switch char {
		case '[':
			lightStart = i
		case ']':
			lightEnd = i
		case '{':
			joltageStart = i
		case '}':
			joltageEnd = i
		}
	}

	// Build lights value
	var lights rune
	lightsLength := lightEnd - lightStart - 1
	for i := range lightsLength {
		if line[lightStart+1+i] == '#' {
			lights |= 1 << i
		}
	}

	// Parse out joltage values
	parts := strings.Split(line[joltageStart+1:joltageEnd], ",")
	joltages := make([]int, len(parts))
	for i, part := range parts {
		value, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			log.Println("Error parsing joltage:", err)
			os.Exit(1)
		}
		joltages[i] = value
	}

	// Parse out the buttons
	buttons := []rune{}
	buttonsString := strings.TrimSpace(line[lightEnd+1 : joltageStart])
	buttonParts := strings.SplitSeq(buttonsString, " ")

	for part := range buttonParts {
		wiring := part[1 : len(part)-1]
		bits := strings.Split(wiring, ",")
		var buttonValue rune

		for _, bit := range bits {
			bitValue, err := strconv.Atoi(bit)
			if err != nil {
				log.Println("Error parsing button bit:", err)
				os.Exit(1)
			}

			buttonValue |= 1 << bitValue
		}

		buttons = append(buttons, buttonValue)
	}

	return Machine{origString: line, lights: lights, buttons: buttons, joltages: joltages}
}

// And here's yet another way to do this - just brute-force all combinations of button presses. Turns out that, given the
// input size, this is actually much faster than the BFS appriach.
func solvePart1New(machine Machine) int {
	minPresses := math.MaxInt

	for i := 0; i < (1 << len(machine.buttons)); i++ {
		currentLights := rune(0)
		presses := 0

		for j := 0; j < len(machine.buttons); j++ {
			if (i & (1 << j)) != 0 {
				currentLights ^= machine.buttons[j]
				presses++
			}
		}
		if currentLights == machine.lights {
			minPresses = min(minPresses, presses)
		}
	}

	return minPresses
}

// My original attempt at this involved BFS and tracking light states and which buttons were pressed.
// Each button can be pressed at most once, so by tracking which buttons were pressed, we could avoid
// trying to push a button a second time. This greatly reduced the search space, but I was treating
// A, B, C button presses as different from B, A, C, even though they result in the same final state.
// This required a lot of extra memory copying and resulted in, ultimately, redundant work.
// The new approach tracks only the current light state, and for each new state reached, tracks
// which buttons were pressed (regardless of order) to reach that state.
type QueueItem struct {
	currentLights rune
	nextButton    rune
}

func solvePart1(machine Machine) []rune {
	queue := []QueueItem{}
	visited := map[rune]map[rune]bool{}
	visited[0] = map[rune]bool{}

	for _, button := range machine.buttons {
		queue = append(queue, QueueItem{currentLights: 0, nextButton: button})
	}

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		currentPresses, _ := visited[item.currentLights]
		newLights := item.currentLights ^ item.nextButton

		if newLights == machine.lights {
			// Found a solution
			buttons := []rune{}

			for button := range currentPresses {
				buttons = append(buttons, button)
			}

			buttons = append(buttons, item.nextButton)
			return buttons
		}

		_, seen := visited[newLights]

		if seen {
			// Already seen this lights state, so no further processing needed
			continue
		}

		// A new lights state - record the state and what buttons were pressed to get here
		newPresses := make(map[rune]bool)
		maps.Copy(newPresses, currentPresses)
		newPresses[item.nextButton] = true
		visited[newLights] = newPresses

		// Now queue up further button presses, skipping buttons that have already been pressed
		for _, button := range machine.buttons {
			if !newPresses[button] {
				queue = append(queue, QueueItem{currentLights: newLights, nextButton: button})
			}
		}
	}

	return nil
}
