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

func (vl *ventLine) horizontal() bool {
	return vl.from.y == vl.to.y
}

func (vl *ventLine) vertical() bool {
	return vl.from.x == vl.to.x
}

type ventLineGrid struct {
	origin  point
	rows    int
	columns int
	cells   []int
}

func newVentLineGrid(lines []ventLine) *ventLineGrid {
	lines = validate(lines)
	origin, end := gridSize(lines)
	columns := end.x - origin.x + 1
	rows := end.y - origin.y + 1

	vlg := &ventLineGrid{
		origin:  origin,
		rows:    rows,
		columns: columns,
		cells:   make([]int, columns*rows),
	}
	for _, l := range lines {
		vlg.add(l)
	}
	return vlg
}

func validate(lines []ventLine) []ventLine {
	var result []ventLine
	for _, l := range lines {
		if l.horizontal() || l.vertical() {
			result = append(result, l)
		}
	}
	return result
}

func gridSize(lines []ventLine) (point, point) {
	if len(lines) == 0 {
		return point{}, point{}
	}

	origin := point{math.MaxInt, math.MaxInt}
	end := point{math.MinInt, math.MinInt}
	for _, l := range lines {
		origin.x = minOf(l.from.x, l.to.x, origin.x)
		origin.y = minOf(l.from.y, l.to.y, origin.y)

		end.x = maxOf(l.from.x, l.to.x, end.x)
		end.y = maxOf(l.from.y, l.to.y, end.y)
	}
	return origin, end
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

func (vlg *ventLineGrid) add(line ventLine) {
	if line.vertical() {
		x := line.from.x
		yFrom := minOf(line.from.y, line.to.y) - vlg.origin.y
		yTo := maxOf(line.from.y, line.to.y) - vlg.origin.y

		for y := yFrom; y <= yTo; y++ {
			vlg.cells[x+y*vlg.columns]++
		}
	} else {
		y := line.from.y
		xFrom := minOf(line.from.x, line.to.x) - vlg.origin.x
		xTo := maxOf(line.from.x, line.to.x) - vlg.origin.x

		for x := xFrom; x <= xTo; x++ {
			vlg.cells[x+y*vlg.columns]++
		}
	}
}

func (vlg *ventLineGrid) print() {
	for x := 0; x < vlg.columns; x++ {
		for y := 0; y < vlg.rows; y++ {
			fmt.Print(vlg.cells[y+x*vlg.columns])
		}

		fmt.Println()
	}
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
	fmt.Printf("origin: (%d, %d)\n", grid.origin.x, grid.origin.y)
	fmt.Printf("columns: %d\n", grid.columns)
	fmt.Printf("rows: %d\n", grid.rows)
	grid.print()
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
