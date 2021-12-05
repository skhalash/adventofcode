package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type point struct {
	x, y int
}

type ventLine struct {
	from, to point
}

type ventLineGrid struct {
	from  point
	to    point
	cells []int
}

func newVentLineGrid(lines []ventLine) *ventLineGrid {
	from, to := gridSize(validate(lines))
	return &ventLineGrid{
		from:  from,
		to:    to,
		cells: make([]int, (to.x-from.x)*(to.y-from.y)),
	}
}

func validate(lines []ventLine) []ventLine {
	var result []ventLine
	for _, l := range lines {
		if l.from.x == l.to.x || l.from.y == l.to.y {
			result = append(result, l)
		}
	}
	return result
}

func gridSize(lines []ventLine) (point, point) {
	if len(lines) == 0 {
		return point{}, point{}
	}

	from := point{math.MaxInt, math.MaxInt}
	to := point{math.MinInt, math.MinInt}
	for _, l := range lines {
		from.x = minOf(l.from.x, l.to.x, from.x)
		from.y = minOf(l.from.y, l.to.y, from.y)

		to.x = maxOf(l.from.x, l.to.x, to.x)
		to.y = maxOf(l.from.y, l.to.y, to.y)
	}
	return from, to
}

func minOf(vars ...int) int {
	min := vars[0]
	for _, i := range vars {
		if min > i {
			min = i
		}
	}
	return min
}

func maxOf(vars ...int) int {
	max := vars[0]
	for _, i := range vars {
		if max < i {
			max = i
		}
	}
	return max
}

func (vlg *ventLineGrid) dangerRate() int {
	return 0
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
	lines, err := loadVentLines(filepath)
	if err != nil {
		return 0, err
	}

	grid := newVentLineGrid(lines)
	fmt.Printf("grid from (%d, %d) to (%d, %d)\n", grid.from.x, grid.from.y, grid.to.x, grid.to.y)
	return grid.dangerRate(), nil
}

func loadVentLines(filepath string) ([]ventLine, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var result []ventLine

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, " -> ")
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid format %s", line)
		}

		from, err := parsePoint(fields[0])
		if err != nil {
			return nil, err
		}

		to, err := parsePoint(fields[1])
		if err != nil {
			return nil, err
		}

		result = append(result, ventLine{from, to})
	}

	return result, nil
}

func parsePoint(s string) (point, error) {
	fields := strings.Split(s, ",")
	if len(fields) != 2 {
		return point{}, fmt.Errorf("invalid format %s", s)
	}

	x, err := strconv.Atoi(fields[0])
	if err != nil {
		return point{}, fmt.Errorf("failed to parse coordinate %s", s)
	}

	y, err := strconv.Atoi(fields[1])
	if err != nil {
		return point{}, fmt.Errorf("failed to parse coordinate %s", s)
	}

	return point{x, y}, nil
}
