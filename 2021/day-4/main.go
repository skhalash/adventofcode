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
	numbers []int
	boards  []*board
}

type board struct {
	cells []int
}

func (b *board) isValid() bool {
	return len(b.cells) == boardSize*boardSize
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

	return len(bingo.boards), nil
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

	var last *board = nil

	for scanner.Scan() {
		line := scanner.Text()

		if line == boardSeparator {
			if last != nil && last.isValid() {
				bingo.boards = append(bingo.boards, last)
			}
			last = &board{}
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
