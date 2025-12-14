package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	connections := map[string][]string{}

	for _, line := range inputLines {
		idx := strings.Index(line, ":")
		device := line[:idx]
		connectedDevices := strings.Split(strings.TrimSpace(line[idx+1:]), " ")
		connections[device] = connectedDevices
	}

	// Part 1 stuff
	pathCount := findPart1Paths(connections, "you")
	fmt.Printf("Part 1: %d (%v)\n", pathCount, time.Since(begin))

	// Part 2 stuff
	begin = time.Now()
	pathCount = findPart2Paths(connections, "svr", &map[string]bool{})
	fmt.Printf("Part 2: %d (%v)\n", pathCount, time.Since(begin))
}

func findPart1Paths(connections map[string][]string, device string) int {
	if device == "out" {
		return 1
	}

	totalPaths := 0

	for _, connectedDevice := range connections[device] {
		totalPaths += findPart1Paths(connections, connectedDevice)
	}

	return totalPaths
}

var cache = map[string]int{}

func findPart2Paths(connections map[string][]string, device string, visited *map[string]bool) int {
	// Check the cache to see if we've already computed this
	cachekey := fmt.Sprintf("%s-%t-%t", device, (*visited)["dac"], (*visited)["fft"])
	cachedVal, found := cache[cachekey]

	if found {
		return cachedVal
	}

	if device == "out" {
		if (*visited)["dac"] && (*visited)["fft"] {
			return 1
		} else {
			return 0
		}
	}

	totalPaths := 0

	for _, connectedDevice := range connections[device] {
		(*visited)[connectedDevice] = true
		totalPaths += findPart2Paths(connections, connectedDevice, visited)
		(*visited)[connectedDevice] = false
	}

	cache[cachekey] = totalPaths

	return totalPaths
}
