package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
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

	begin = time.Now()

	totalPresses := 0

	for _, machine := range machines {
		buttons := solvePart1New(machine)
		totalPresses += len(buttons)
	}

	fmt.Printf("Part 1: %d (%v)\n", totalPresses, time.Since(begin))
}

type QueueItem struct {
	currentLights rune
	nextButton    rune
}

// My original attempt at this involved BFS and tracking light states and which buttons were pressed.
// Each button can be pressed at most once, so by tracking which buttons were pressed, we could avoid
// trying to push a button a second time. This greatly reduced the search space, but I was treating
// A, B, C button presses as different from B, A, C, even though they result in the same final state.
// This requird a lot of extra memory copying and resulted in, ultimately, redundant work.
// The new approach tracks only the current light state, and for each new state reached, tracks
// which buttons were pressed to reach that state.
func solvePart1New(machine Machine) []rune {
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

	return Machine{lights: lights, buttons: buttons, joltages: joltages}
}

type Machine struct {
	lights   rune
	buttons  []rune
	joltages []int
}
