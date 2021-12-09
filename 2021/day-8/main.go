package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var originalPatterns = []string{
	"abcefg",
	"cf",
	"acdeg",
	"acdfg",
	"bcdf",
	"abdfg",
	"abdefg",
	"acf",
	"abcdefg",
	"abcdfg",
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
	configs, err := loadConfigs(filepath)
	if err != nil {
		return 0, err
	}

	config := configs[0]
	patternsOfOne := byLen(config, 2)
	patternsOfSeven := byLen(config, 3)
	patternsOfFour := byLen(config, 4)
	patternsOfEight := byLen(config, 7)

	optionsByChar := map[string][]string{
		"a": split(diff(patternsOfSeven, patternsOfOne)),
		"b": split(diff(patternsOfFour, patternsOfOne)),
		"c": split(patternsOfOne),
		"d": split(diff(patternsOfFour, patternsOfOne)),
		"e": split(diff(patternsOfEight, union(patternsOfFour, patternsOfSeven))),
		"f": split(patternsOfOne),
		"g": split(diff(patternsOfEight, union(patternsOfFour, patternsOfSeven))),
	}

	for k, v := range optionsByChar {
		fmt.Printf("%s: %#v\n", k, v)
	}

	fmt.Println()

	mapping := analyze(optionsByChar)

	for k, v := range mapping {
		fmt.Printf("%s: %#v\n", k, v)
	}

	return 0, nil
}

func analyze(optionsByChar map[string][]string) map[string]string {
	mapping := make(map[string]string)
	for len(mapping) != 10 {
		for char := range optionsByChar {
			options := optionsByChar[char]
			if len(options) == 0 {
				continue
			}

			//only option
			if len(options) == 1 {
				mapping[options[0]] = char
				continue
			}

			//found itself
			idx := index(options, char)
			if idx >= 0 {
				optionsByChar[char] = append(options[:idx], options[idx+1:]...)
				continue
			}

			//found already mapped
			for i, opt := range options {
				if _, found := mapping[opt]; found {
					optionsByChar[char] = append(options[:i], options[i+1:]...)
				}
			}

			fmt.Printf("%#v\n", mapping)
		}
	}
	return mapping
}

func index(a []string, s string) int {
	for i, elem := range a {
		if elem == s {
			return i
		}
	}

	return -1
}

func byLen(patterns []string, l int) string {
	for _, p := range patterns {
		if len(p) == l {
			return p
		}
	}

	return ""
}

func diff(a, b string) string {
	var sb strings.Builder
	for _, r := range a {
		if !strings.ContainsRune(b, r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func union(a, b string) string {
	var sb strings.Builder
	sb.WriteString(a)
	for _, r := range b {
		if !strings.ContainsRune(a, r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func split(s string) []string {
	var result []string
	for _, r := range s {
		result = append(result, string(r))
	}
	return result
}

func loadConfigs(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var result [][]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			return nil, fmt.Errorf("must contain two parts %s: %v", line, err)
		}

		result = append(result, strings.Fields(parts[0]))
	}

	return result, nil
}
