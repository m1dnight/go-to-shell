package main

import (
	"regexp"
	"strings"
)

type command struct {
	cmd  string
	args []string
}

type successResult struct {
	message string
}

type errorResult struct {
	message string
	cmd     string
}

func tokens(input string) []string {
	// remove consecutive whitespaces
	pattern := regexp.MustCompile(`\s+`)
	input = pattern.ReplaceAllString(input, " ")

	// trim leading and trailing spaces
	input = strings.TrimSpace(input)

	// split on space
	tokens := strings.Split(input, " ")
	return tokens
}

func parse(input string) command {
	inputs := tokens(input)
	return command{cmd: inputs[0], args: inputs[1:]}
}
