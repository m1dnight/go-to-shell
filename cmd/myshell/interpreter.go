package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

var commands = []string{"exit", "echo", "type", "pwd", "cd"}

var mappings = map[string]string{"cat": "/bin/cat"}

func evaluate(command command) (*successResult, *errorResult) {
	var sucResult *successResult = nil
	var errResult *errorResult = nil

	if isBuiltIn(command) {
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
		case "pwd":
			sucResult, errResult = evalPwd(command)
			return sucResult, errResult
		case "cd":
			sucResult, errResult = evalCd(command)
			return sucResult, errResult
		}
	} else {
		sucResult, errResult = evalExecutable(command)
		return sucResult, errResult
	}
	return sucResult, errResult
}

func evalCd(command command) (*successResult, *errorResult) {
	var sucResult *successResult = nil
	var errResult *errorResult = nil

	// expects 1 argument
	if len(command.args) != 1 {
		errResult = &errorResult{message: "too few arguments", cmd: command.cmd}
		return sucResult, errResult
	}

	// create the path to CD to if its relative
	absPath, err := filepath.Abs(command.args[0])
	if err != nil {
		errResult = &errorResult{message: fmt.Sprintf("%s: No such file or directory\n", command.args[0]), cmd: command.cmd}
		return sucResult, errResult
	}

	err = os.Chdir(absPath)
	if err != nil {
		errResult = &errorResult{message: fmt.Sprintf("%s: No such file or directory\n", command.args[0]), cmd: command.cmd}
		return sucResult, errResult
	}

	sucResult = &successResult{message: ""}
	return sucResult, errResult
}

func evalPwd(command command) (*successResult, *errorResult) {
	var sucResult *successResult = nil
	var errResult *errorResult = nil

	path, err := os.Getwd()
	if err != nil {
		errResult = &errorResult{message: "error getting pwd", cmd: command.cmd}
		return sucResult, errResult
	}

	sucResult = &successResult{message: path + "\n"}
	return sucResult, errResult
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
	sucResult = &successResult{message: outputMessage + "\n"}
	return sucResult, errResult
}

func evalType(command command) (*successResult, *errorResult) {
	var sucResult *successResult = nil
	var errResult *errorResult = nil

	// expect exactly one argument
	if len(command.args) != 1 {
		errResult = &errorResult{message: "too few arguments\n", cmd: command.cmd}
		return sucResult, errResult
	}

	// if the commands exists, it's a builtin
	if slices.Contains(commands, command.args[0]) {
		outputMessage := fmt.Sprintf("%s is a shell builtin\n", command.args[0])
		sucResult = &successResult{message: outputMessage}
		return sucResult, errResult
	}

	// if the command is in the PATH, return the binding
	executablePath := findInPath(command.args[0])
	if executablePath != nil {
		outputMessage := fmt.Sprintf("%s is %s\n", command.args[0], *executablePath)
		sucResult = &successResult{message: outputMessage}
		return sucResult, errResult
	}

	outputMessage := fmt.Sprintf("%s: not found\n", command.args[0])
	sucResult = &successResult{message: outputMessage}
	return sucResult, errResult

}

func evalExecutable(command command) (*successResult, *errorResult) {
	var sucResult *successResult = nil
	var errResult *errorResult = nil

	// if the command is not in the path, return an error saying it's not found
	executablePath := findInPath(command.cmd)
	if executablePath == nil {
		outputMessage := fmt.Sprintf("%s: not found\n", command.cmd)
		sucResult = &successResult{message: outputMessage}
		return sucResult, errResult
	}

	cmd := exec.Command(*executablePath, command.args...)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		errResult = &errorResult{message: fmt.Sprintf("command failed\n%s", out.String()), cmd: command.cmd}
		return sucResult, errResult
	}

	outputMessage := fmt.Sprintf("%s", out.String())
	sucResult = &successResult{message: outputMessage}
	return sucResult, errResult

}

func findInPath(executable string) *string {
	path := os.Getenv("PATH")
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

func isBuiltIn(command command) bool {
	return slices.Contains(commands, command.cmd)

}
