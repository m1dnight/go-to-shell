package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var commands = []string{"exit", "echo", "type"}

func evaluate(command command) (*successResult, *errorResult) {
	if !slices.Contains(commands, command.cmd) {
		result := errorResult{message: "command not found", cmd: command.cmd}
		return nil, &result
	}

	var sucResult *successResult = nil
	var errResult *errorResult = nil
	switch command.cmd {
	case "exit":
		if len(command.args) != 1 {
			errResult = &errorResult{message: "too few arguments", cmd: command.cmd}
			return sucResult, errResult
		}
		exitCode := command.args[0]
		exitCodeNbr, err := strconv.Atoi(exitCode)
		if err != nil {
			errResult = &errorResult{message: "invalid argument to command", cmd: command.cmd}
			return sucResult, errResult
		} else {
			os.Exit(exitCodeNbr)
		}
	case "echo":
		outputMessage := strings.Join(command.args, " ")
		sucResult = &successResult{message: outputMessage}
		return sucResult, errResult
	case "type":
		if len(command.args) != 1 {
			errResult = &errorResult{message: "too few arguments", cmd: command.cmd}
			return sucResult, errResult
		}
		if slices.Contains(commands, command.args[0]) {
			outputMessage := fmt.Sprintf("%s is a shell builtin", command.args[0])
			sucResult = &successResult{message: outputMessage}
			return sucResult, errResult
		}
		outputMessage := fmt.Sprintf("%s: not found", command.args[0])
		sucResult = &successResult{message: outputMessage}
		return sucResult, errResult

	default:
		errResult = &errorResult{message: "unknown error", cmd: command.cmd}
		return sucResult, errResult

	}
	return sucResult, errResult
}
