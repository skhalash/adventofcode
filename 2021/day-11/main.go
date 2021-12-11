package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	octopuses, err := loadOctopusGrid(filepath)
	if err != nil {
		return 0, err
	}

	totalFlashes := 0
	for i := 0; i < 3; i++ {
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

	for x := 0; x < len(octopuses); x++ {
		for y := 0; y < len(octopuses[x]); y++ {
			if energized(x, y, octopuses) {
				flash(x, y, octopuses)
			}
		}
	}

	flashes := 0
	for x := 0; x < len(octopuses); x++ {
		for y := 0; y < len(octopuses[x]); y++ {
			if octopuses[x][y] == 0 {
				flashes++
			}
			fmt.Printf("%d ", octopuses[x][y])
		}
		fmt.Println()
	}
	fmt.Println()
	return flashes
}

func flash(x, y int, octopuses [][]int) {
	octopuses[x][y] = 0

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			adjX := x + dx
			adjY := y + dy
			hasAdjacent := adjX >= 0 && adjX < len(octopuses) && adjY >= 0 && adjY < len(octopuses[x])
			if !hasAdjacent {
				continue
			}

			octopuses[adjX][adjY]++
			if energized(adjX, adjY, octopuses) {
				flash(adjX, adjY, octopuses)
			}
		}
	}
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
