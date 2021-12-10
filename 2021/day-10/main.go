package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type stack []rune

func (s *stack) push(item rune) {
	*s = append(*s, item)
}

func (s *stack) pop() (rune, bool) {
	if s.empty() {
		return 0, false
	}

	lastIndex := len(*s) - 1
	last := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return last, true
}

func (s *stack) empty() bool {
	return len(*s) == 0
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
	brackets, err := loadBrackets(filepath)
	if err != nil {
		return 0, err
	}

	var scores []int
	for _, br := range brackets {
		if score, err := autocompletionScore(br); err == nil {
			scores = append(scores, score)
		}
	}

	sort.Sort(sort.IntSlice(scores))

	return scores[len(scores)/2], nil
}

func autocompletionScore(brackets []rune) (int, error) {
	unmatched, err := balance(brackets)
	if err != nil {
		return 0, fmt.Errorf("unbalanced brackets: %v", err)
	}

	totalScore := 0
	for !unmatched.empty() {
		totalScore *= 5
		top, _ := unmatched.pop()
		totalScore += bracketScore(top)
	}

	return totalScore, nil
}

func bracketScore(b rune) int {
	switch b {
	case '(':
		return 1
	case '[':
		return 2
	case '{':
		return 3
	case '<':
		return 4
	}
	return 0
}

func balance(brackets []rune) (stack, error) {
	var unmatched stack
	for _, bracket := range brackets {
		if opening(bracket) {
			unmatched.push(bracket)
			continue
		}

		if closing(bracket) {
			if unmatched.empty() {
				return nil, fmt.Errorf("no opening bracket for %c", bracket)
			}

			top, _ := unmatched.pop()
			if !matching(top, bracket) {
				return nil, fmt.Errorf("non matching brackets %c,%c", top, bracket)
			}
		}
	}

	return unmatched, nil
}

func opening(r rune) bool {
	switch r {
	case '(', '[', '{', '<':
		return true
	}
	return false
}

func closing(r rune) bool {
	switch r {
	case ')', ']', '}', '>':
		return true
	}
	return false
}

func matching(opening, closing rune) bool {
	return opening == '(' && closing == ')' ||
		opening == '[' && closing == ']' ||
		opening == '{' && closing == '}' ||
		opening == '<' && closing == '>'
}

func loadBrackets(filepath string) ([][]rune, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var result [][]rune

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var row []rune
		for _, r := range scanner.Text() {
			row = append(row, r)
		}

		result = append(result, row)
	}

	return result, nil
}
