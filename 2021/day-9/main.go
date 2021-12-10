package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type point struct {
	i, j   int
	height int
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

	var basinSizes []int
	for _, pt := range lowPoints(heightmap) {
		basinSizes = append(basinSizes, basin(pt, heightmap))
	}

	sort.Sort(sort.IntSlice(basinSizes))
	last := len(basinSizes) - 1
	return basinSizes[last] * basinSizes[last-1] * basinSizes[last-2], nil
}

func lowPoints(heightmap [][]int) []point {
	var lowPoints []point
	for i := 0; i < len(heightmap); i++ {
		for j := 0; j < len(heightmap[i]); j++ {
			height := heightmap[i][j]
			isLow := true
			for _, n := range neigbours(point{i, j, height}, heightmap) {
				if n.height <= height {
					isLow = false
					break
				}
			}

			if isLow {
				lowPoints = append(lowPoints, point{i, j, height})
			}
		}
	}

	return lowPoints
}

func basin(pt point, heightmap [][]int) int {
	visited := make(map[point]bool)
	dfs(pt, heightmap, visited)
	return len(visited)
}

func dfs(pt point, heightmap [][]int, visited map[point]bool) {
	for _, neighbour := range neigbours(pt, heightmap) {
		if neighbour.height == 9 {
			continue
		}

		if _, found := visited[neighbour]; found {
			continue
		}

		visited[neighbour] = true
		dfs(neighbour, heightmap, visited)
	}
}

func neigbours(pt point, heightmap [][]int) []point {
	var neigbours []point
	i, j := pt.i, pt.j
	if i > 0 {
		neigbours = append(neigbours, point{i - 1, j, heightmap[i-1][j]})
	}
	if i < len(heightmap)-1 {
		neigbours = append(neigbours, point{i + 1, j, heightmap[i+1][j]})
	}
	if j > 0 {
		neigbours = append(neigbours, point{i, j - 1, heightmap[i][j-1]})
	}
	if j < len(heightmap[i])-1 {
		neigbours = append(neigbours, point{i, j + 1, heightmap[i][j+1]})
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
