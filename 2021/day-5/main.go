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

func (vl *ventLine) diagonal() bool {
	return abs(vl.to.x-vl.from.x) == abs(vl.to.y-vl.from.y)
}

type ventLineGrid struct {
	origin  point
	rows    int
	columns int
	cells   []int
}

func newVentLineGrid(lines []ventLine) *ventLineGrid {
	lines = validate(lines)
	origin, end := gridCorners(lines)
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
		if l.horizontal() || l.vertical() || l.diagonal() {
			result = append(result, l)
		}
	}
	return result
}

func gridCorners(lines []ventLine) (point, point) {
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

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

func (vlg *ventLineGrid) add(line ventLine) {
	if line.vertical() {
		vlg.addVertical(line)
	} else if line.horizontal() {
		vlg.addHorizontal(line)
	} else if line.diagonal() {
		vlg.addDiagonal(line)
	}
}

func (vlg *ventLineGrid) addVertical(line ventLine) {
	yFrom := minOf(line.from.y, line.to.y)
	yTo := maxOf(line.from.y, line.to.y)

	for y := yFrom; y <= yTo; y++ {
		vlg.increment(line.from.x, y)
	}
}

func (vlg *ventLineGrid) addHorizontal(line ventLine) {
	xFrom := minOf(line.from.x, line.to.x)
	xTo := maxOf(line.from.x, line.to.x)

	for x := xFrom; x <= xTo; x++ {
		vlg.increment(x, line.from.y)
	}
}

func (vlg *ventLineGrid) addDiagonal(line ventLine) {
	x, y := line.from.x, line.from.y
	xStep := step(line.from.x, line.to.x)
	yStep := step(line.from.y, line.to.y)

	for {
		vlg.increment(x, y)
		if x == line.to.x && y == line.to.y {
			break
		}

		x += xStep
		y += yStep
	}
}

func step(from, to int) int {
	if from < to {
		return 1
	}
	if to < from {
		return -1
	}
	return 0
}

func (vlg *ventLineGrid) increment(x, y int) {
	x, y = vlg.normalize(x, y)
	vlg.cells[x+y*vlg.columns]++
}

func (vlg *ventLineGrid) normalize(x, y int) (int, int) {
	return x - vlg.origin.x, y - vlg.origin.y
}

func (vlg *ventLineGrid) dangerRate() int {
	rate := 0
	for _, c := range vlg.cells {
		if c >= 2 {
			rate++
		}
	}
	return rate
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
