package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var originalPatterns = []string{
	"abcefg",
	"cf",
	"acdeg",
	"acdfg",
	"bcdf",
	"abdfg",
	"abdefg",
	"acf",
	"abcdefg",
	"abcdfg",
}

func main() {
	result, err := run(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Println(result)
}

func run(filepath string) (int, error) {
	heightmap, err := loadHeightmap(filepath)
	if err != nil {
		return 0, err
	}

	return riskLevel(heightmap), nil
}

func riskLevel(heightmap [][]int) int {
	var lowPoints []int
	for i := 0; i < len(heightmap); i++ {
		for j := 0; j < len(heightmap[i]); j++ {
			val := heightmap[i][j]

			upperLower := i > 0 && heightmap[i-1][j] <= val
			lowerLower := i < len(heightmap)-1 && heightmap[i+1][j] <= val
			leftLower := j > 0 && heightmap[i][j-1] <= val
			rightLower := j < len(heightmap[i])-1 && heightmap[i][j+1] <= val
			if upperLower || lowerLower || leftLower || rightLower {
				continue
			}

			lowPoints = append(lowPoints, val)

		}
	}

	riskLevel := 0
	for _, lp := range lowPoints {
		riskLevel += (lp + 1)
	}
	return riskLevel
}

func loadHeightmap(filepath string) ([][]int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var result [][]int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var raw []int
		line := scanner.Text()
		for _, char := range line {
			n, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, fmt.Errorf("failed to parse number: %v", err)
			}
			raw = append(raw, n)
		}

		result = append(result, raw)
	}

	return result, nil
}
