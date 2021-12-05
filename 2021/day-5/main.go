package main

import (
	"bufio"
	"fmt"
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
	cells []int
}

func newVentLineGrid(lines []ventLine) *ventLineGrid {
	return nil
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
