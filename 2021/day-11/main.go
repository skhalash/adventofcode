package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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

	totalFlashes := 0
	for i := 0; i < 100; i++ {
		stepFlashes := nextStep(octopuses)
		totalFlashes += stepFlashes
	}

	return totalFlashes, nil
}

func nextStep(octopuses [][]int) int {
	for x := 0; x < len(octopuses); x++ {
		for y := 0; y < len(octopuses[x]); y++ {
			octopuses[x][y]++
		}
	}

	flashed := make(map[coordinate]bool)

	for x := 0; x < len(octopuses); x++ {
		for y := 0; y < len(octopuses[x]); y++ {
			if energized(x, y, octopuses) {
				flash(coordinate{x, y}, octopuses, flashed)
			}
		}
	}

	return len(flashed)
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
			hasAdjacent := adjX >= 0 && adjX < len(octopuses) && adjY >= 0 && adjY < len(octopuses[current.x])
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
		for _, char := range line {
			n, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, fmt.Errorf("failed to parse number: %v", err)
			}
			row = append(row, n)
		}

		result = append(result, row)
	}

	return result, nil
}
