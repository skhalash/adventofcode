package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const gridSize = 10

type coordinate struct {
	x, y int
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
	octopuses, err := loadOctopusGrid(filepath)
	if err != nil {
		return 0, err
	}

	step := 0
	for {
		allFlashed := nextStep(octopuses)
		step++
		if allFlashed {
			break
		}
	}

	return step, nil
}

func nextStep(octopuses [][]int) bool {
	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			octopuses[x][y]++
		}
	}

	flashed := make(map[coordinate]bool)

	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			if energized(x, y, octopuses) {
				flash(coordinate{x, y}, octopuses, flashed)
			}
		}
	}

	return len(flashed) == gridSize*gridSize
}

func flash(current coordinate, octopuses [][]int, flashed map[coordinate]bool) {
	flashed[current] = true

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			adjX := current.x + dx
			adjY := current.y + dy
			hasAdjacent := adjX >= 0 && adjX < gridSize && adjY >= 0 && adjY < gridSize
			if !hasAdjacent {
				continue
			}

			adjacent := coordinate{adjX, adjY}
			if flashed[adjacent] {
				continue
			}

			octopuses[adjX][adjY]++
			if energized(adjX, adjY, octopuses) {
				flash(adjacent, octopuses, flashed)
			}
		}
	}

	octopuses[current.x][current.y] = 0
}

func energized(i, j int, octopuses [][]int) bool {
	return octopuses[i][j] > 9
}

func loadOctopusGrid(filepath string) ([][]int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var result [][]int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var row []int
		line := scanner.Text()
		if len(line) != gridSize {
			return nil, fmt.Errorf("must have %d columns", gridSize)
		}

		for _, char := range line {
			n, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, fmt.Errorf("failed to parse number: %v", err)
			}
			row = append(row, n)
		}

		result = append(result, row)
	}

	if len(result) != gridSize {
		return nil, fmt.Errorf("must have %d rows", gridSize)
	}

	return result, nil
}
