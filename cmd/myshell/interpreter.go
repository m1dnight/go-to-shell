package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var commands = []string{"exit", "echo", "type"}

var mappings = map[string]string{"cat": "/bin/cat"}

func evaluate(command command) (*successResult, *errorResult) {
	if !slices.Contains(commands, command.cmd) {
		result := errorResult{message: "command not found", cmd: command.cmd}
		return nil, &result
	}

	var sucResult *successResult = nil
	var errResult *errorResult = nil
	switch command.cmd {
	case "exit":
		sucResult, errResult = evalExit(command)
		return sucResult, errResult

	case "echo":
		sucResult, errResult = evalEcho(command)
		return sucResult, errResult
	case "type":
		sucResult, errResult = evalType(command)
		return sucResult, errResult

	default:
		errResult = &errorResult{message: "unknown error", cmd: command.cmd}
		return sucResult, errResult
	}
}

func evalExit(command command) (*successResult, *errorResult) {
	var sucResult *successResult = nil
	var errResult *errorResult = nil

	// expectes 1 argument
	if len(command.args) != 1 {
		errResult = &errorResult{message: "too few arguments", cmd: command.cmd}
		return sucResult, errResult
	}

	exitCode := command.args[0]
	exitCodeNbr, err := strconv.Atoi(exitCode)
	if err != nil {
		errResult = &errorResult{message: "invalid argument to command", cmd: command.cmd}
		return sucResult, errResult
	}

	os.Exit(exitCodeNbr)
	return sucResult, errResult
}

func evalEcho(command command) (*successResult, *errorResult) {
	var sucResult *successResult = nil
	var errResult *errorResult = nil

	outputMessage := strings.Join(command.args, " ")
	sucResult = &successResult{message: outputMessage}
	return sucResult, errResult
}

func evalType(command command) (*successResult, *errorResult) {
	var sucResult *successResult = nil
	var errResult *errorResult = nil

	// expect exactly one argument
	if len(command.args) != 1 {
		errResult = &errorResult{message: "too few arguments", cmd: command.cmd}
		return sucResult, errResult
	}

	// if the commands exists, it's a builtin
	if slices.Contains(commands, command.args[0]) {
		outputMessage := fmt.Sprintf("%s is a shell builtin", command.args[0])
		sucResult = &successResult{message: outputMessage}
		return sucResult, errResult
	}

	// if the command is in the PATH, return the binding
	executablePath := findInPath(os.Getenv("PATH"), command.args[0])
	if executablePath != nil {
		outputMessage := fmt.Sprintf("%s is %s", command.args[0], *executablePath)
		sucResult = &successResult{message: outputMessage}
		return sucResult, errResult
	}

	outputMessage := fmt.Sprintf("%s: not found", command.args[0])
	sucResult = &successResult{message: outputMessage}
	return sucResult, errResult

}

func findInPath(path string, executable string) *string {
	dirs := strings.Split(path, ":")

	for _, path := range dirs {
		file, err := os.Open(path)
		if err != nil {
			continue
		}

		names, _ := file.Readdirnames(0)
		for _, name := range names {
			if name == executable {
				executablePath := fmt.Sprintf("%s/%s", path, name)
				return &executablePath
			}
		}

	}
	return nil
}
