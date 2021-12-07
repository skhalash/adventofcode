package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	result, err := run(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Println(result)
}

func run(filepath string) (int, error) {
	crabPositions, err := loadCrabPositions(filepath)
	if err != nil {
		return 0, err
	}

	crabCountByPosition := make(map[int]int)
	for _, pos := range crabPositions {
		crabCountByPosition[pos]++
	}

	return alignCrabs(crabCountByPosition), nil
}

func alignCrabs(crabCountByPosition map[int]int) int {
	minPos, maxPos := minMaxPosition(crabCountByPosition)
	minFuel := math.MaxInt
	for i := minPos; i < maxPos; i++ {
		fuel := 0
		for pos, count := range crabCountByPosition {
			fuel += count * fuelSpent(abs(pos-i))
		}
		minFuel = min(minFuel, fuel)
	}

	return minFuel
}

func minMaxPosition(crabCountByPosition map[int]int) (int, int) {
	min, max := math.MaxInt, math.MinInt
	for pos := range crabCountByPosition {
		if max < pos {
			max = pos
		}
		if min > pos {
			min = pos
		}
	}
	return min, max
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func fuelSpent(distance int) int {
	return distance * (distance + 1) / 2
}

func loadCrabPositions(filepath string) ([]int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}

	strs := strings.Split(strings.TrimSpace(string(content)), ",")

	var result []int
	for _, s := range strs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number: %v", err)
		}
		result = append(result, n)
	}
	return result, nil
}
