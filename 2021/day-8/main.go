package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type task struct {
	patterns []pattern
	output   []pattern
}

type pattern string

func byLen(patterns []pattern, l int) []pattern {
	var result []pattern
	for _, s := range patterns {
		if len(s) == l {
			result = append(result, s)
		}
	}
	return result
}

func single(patterns []pattern) pattern {
	if len(patterns) != 1 {
		panic(fmt.Sprintf("must contain single elem %#v", patterns))
	}
	return patterns[0]
}

func (p pattern) union(other pattern) pattern {
	var sb strings.Builder
	sb.WriteString(string(p))
	for _, char := range string(other) {
		if !strings.ContainsRune(string(p), char) {
			sb.WriteRune(char)
		}
	}
	return pattern(sb.String())
}

func (p pattern) contains(other pattern) bool {
	for _, char := range string(other) {
		if !strings.ContainsRune(string(p), char) {
			return false
		}
	}
	return true
}

func (p pattern) equals(other pattern) bool {
	return len(p) == len(other) && p.contains(other)
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
	tasks, err := loadTasks(filepath)
	if err != nil {
		return 0, err
	}

	sum := 0
	for _, task := range tasks {
		patternByDigit := deduceDigits(task.patterns)
		sum += decodeOutput(task.output, patternByDigit)
	}

	return sum, nil
}

func deduceDigits(patterns []pattern) []pattern {
	patternByDigit := make([]pattern, 10)
	patternByDigit[1] = single(byLen(patterns, 2))
	patternByDigit[4] = single(byLen(patterns, 4))
	patternByDigit[7] = single(byLen(patterns, 3))
	patternByDigit[8] = single(byLen(patterns, 7))

	almostNine := patternByDigit[4].union(patternByDigit[7])
	for _, pattern := range byLen(patterns, 6) {
		if pattern.contains(almostNine) {
			patternByDigit[9] = pattern
		} else if pattern.contains(patternByDigit[1]) {
			patternByDigit[0] = pattern
		} else {
			patternByDigit[6] = pattern
		}
	}

	for _, pattern := range byLen(patterns, 5) {
		if patternByDigit[6].contains(pattern) {
			patternByDigit[5] = pattern
		} else if pattern.contains(patternByDigit[1]) {
			patternByDigit[3] = pattern
		} else {
			patternByDigit[2] = pattern
		}
	}

	return patternByDigit
}

func decodeOutput(output []pattern, patternByDigit []pattern) int {
	result := 0
	for _, outPattern := range output {
		digit := 0
		for i, digitPattern := range patternByDigit {
			if outPattern.equals(digitPattern) {
				digit = i
			}
		}

		result = result*10 + digit
	}
	return result
}

func loadTasks(filepath string) ([]task, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var result []task

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			return nil, fmt.Errorf("must contain two parts %s: %v", line, err)
		}

		var patterns []pattern
		for _, s := range strings.Fields(parts[0]) {
			patterns = append(patterns, pattern(s))
		}

		var output []pattern
		for _, s := range strings.Fields(parts[1]) {
			output = append(output, pattern(s))
		}

		result = append(result, task{patterns, output})
	}

	return result, nil
}
