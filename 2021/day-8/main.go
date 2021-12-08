package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	stuff, err := loadStuff(filepath)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, s := range stuff {
		switch len(s) {
		case 2, 3, 4, 7:
			count++
		}
	}

	return count, nil
}

func loadStuff(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var result []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			return nil, fmt.Errorf("must contain two parts %s: %v", line, err)
		}

		for _, s := range strings.Fields(parts[1]) {
			result = append(result, s)
		}
	}

	return result, nil
}
