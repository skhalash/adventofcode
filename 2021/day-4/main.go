package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const boardSize = 5

type bingo struct {
	numbers          []int
	boards           []*board
	winnerBoardCount int
}

func (b *bingo) play() int {
	for _, n := range b.numbers {
		for _, board := range b.boards {
			if board.isWinner() {
				continue
			}

			board.mark(n)
			if board.isWinner() {
				b.winnerBoardCount++
				if b.winnerBoardCount == len(b.boards) {
					return board.score()
				}
			}
		}
	}

	return 0
}

type board struct {
	cells                []int
	marked               []bool
	lastMarked           int
	markedCountPerRow    []int
	markedCountPerColumn []int
}

func newBoard() *board {
	size := boardSize
	return &board{
		cells:                make([]int, 0, size*size),
		marked:               make([]bool, size*size, size*size),
		markedCountPerRow:    make([]int, size, size),
		markedCountPerColumn: make([]int, size, size),
	}
}

func (b *board) isValid() bool {
	return len(b.cells) == boardSize*boardSize
}

func (b *board) mark(n int) {
	for i := range b.cells {
		if b.cells[i] != n || b.marked[i] {
			continue
		}

		b.marked[i] = true
		row := i / boardSize
		column := i % boardSize
		b.markedCountPerRow[row]++
		b.markedCountPerColumn[column]++
		b.lastMarked = n
	}
}

func (b *board) isWinner() bool {
	for _, count := range b.markedCountPerRow {
		if count == boardSize {
			return true
		}
	}

	for _, count := range b.markedCountPerColumn {
		if count == boardSize {
			return true
		}
	}

	return false
}

func (b *board) score() int {
	unmarkedSum := 0
	for i, marked := range b.marked {
		if !marked {
			unmarkedSum += b.cells[i]
		}
	}

	return unmarkedSum * b.lastMarked
}

func (b *board) print() {
	for x := 0; x < boardSize; x++ {
		for y := 0; y < boardSize; y++ {
			i := y + x*boardSize
			if b.marked[i] {
				fmt.Printf("X ")
			} else {
				fmt.Printf("0 ")
			}
		}

		fmt.Println()
	}
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
	bingo, err := loadBingo(filepath)
	if err != nil {
		return 0, err
	}

	return bingo.play(), nil
}

const numberSeparator = ","
const boardSeparator = ""

func loadBingo(filepath string) (*bingo, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var bingo bingo

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	header := scanner.Text()
	for _, s := range strings.Split(header, numberSeparator) {
		number, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number %s: %v", s, err)
		}
		bingo.numbers = append(bingo.numbers, number)
	}

	var last *board
	for scanner.Scan() {
		line := scanner.Text()

		if line == boardSeparator {
			if last != nil && last.isValid() {
				bingo.boards = append(bingo.boards, last)
			}
			last = newBoard()
			continue
		}

		for _, s := range strings.Fields(line) {
			cell, err := strconv.Atoi(s)
			if err != nil {

				return nil, fmt.Errorf("failed to parse board cell %s: %v", s, err)
			}

			last.cells = append(last.cells, cell)
		}
	}

	if last != nil && last.isValid() {
		bingo.boards = append(bingo.boards, last)
	}

	return &bingo, nil
}
