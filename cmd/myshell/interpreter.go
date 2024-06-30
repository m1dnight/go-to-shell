package main

import (
	"slices"
)

var commands = []string{"foo"}

func evaluate(command command) (*successResult, *errorResult) {
	if !slices.Contains(commands, command.cmd) {
		result := errorResult{message: "command not found", cmd: command.cmd}
		return nil, &result
	}

	result := successResult{message: "success"}
	return &result, nil
}
