package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type commandType string

var forward commandType = "forward"
var up commandType = "up"
var down commandType = "down"

type command struct {
	commandType commandType
	units       int
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
	commands, err := loadCommands(filepath)
	if err != nil {
		return 0, err
	}

	return course(commands), nil
}

func course(commands []command) int {
	pos, depth, aim := 0, 0, 0

	for _, c := range commands {
		switch c.commandType {
		case down:
			aim += c.units
		case up:
			aim -= c.units
		case forward:
			pos += c.units
			depth += (c.units * aim)
		}
	}

	return pos * depth
}

func loadCommands(filepath string) ([]command, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var result []command

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("must contain two parts %s: %v", line, err)
		}

		commandType := commandType(parts[0])
		switch commandType {
		case forward, up, down:
		default:
			return nil, fmt.Errorf("invalid command type: %s", commandType)
		}

		units, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse units %s: %v", parts[1], err)
		}

		result = append(result, command{commandType: commandType, units: units})
	}

	return result, nil
}
