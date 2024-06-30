package main

import (
	"os"
	"slices"
	"strconv"
)

var commands = []string{"exit"}

func evaluate(command command) (*successResult, *errorResult) {
	if !slices.Contains(commands, command.cmd) {
		result := errorResult{message: "command not found", cmd: command.cmd}
		return nil, &result
	}

	//var success successResult
	var errResult errorResult
	switch command.cmd {
	case "exit":
		exitCode := command.args[0]
		exitCodeNbr, err := strconv.Atoi(exitCode)
		if err != nil {
			errResult = errorResult{message: "invalid argument to command", cmd: command.cmd}
		} else {
			os.Exit(exitCodeNbr)
		}
	default:
		errResult = errorResult{message: "unknown error", cmd: command.cmd}

	}
	return nil, &errResult
}
