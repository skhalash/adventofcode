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
	riskLevel := 0
	for _, lp := range lowPoints(heightmap) {
		riskLevel += (lp + 1)
	}
	return riskLevel
}

func lowPoints(heightmap [][]int) []int {
	var lowPoints []int
	for i := 0; i < len(heightmap); i++ {
		for j := 0; j < len(heightmap[i]); j++ {
			val := heightmap[i][j]
			isLow := true
			for _, n := range neigbours(i, j, heightmap) {
				if n <= val {
					isLow = false
					break
				}
			}

			if isLow {
				lowPoints = append(lowPoints, val)
			}
		}
	}

	return lowPoints
}

func neigbours(i, j int, heightmap [][]int) []int {
	var neigbours []int
	if i > 0 {
		neigbours = append(neigbours, heightmap[i-1][j])
	}
	if i < len(heightmap)-1 {
		neigbours = append(neigbours, heightmap[i+1][j])
	}
	if j > 0 {
		neigbours = append(neigbours, heightmap[i][j-1])
	}
	if j < len(heightmap[i])-1 {
		neigbours = append(neigbours, heightmap[i][j+1])
	}
	return neigbours
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
