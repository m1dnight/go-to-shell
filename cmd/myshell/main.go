package main

import (
	"bufio"
	"fmt"

	// Uncomment this block to pass the first stage
	// "fmt"
	"os"
)

func main() {
	readCommand()
}

func readCommand() {
	// Uncomment this block to pass the first stage
	_, _ = fmt.Fprint(os.Stdout, "$ ")

	// Wait for user input
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	command := parse(input)
	result, err := evaluate(command)
	if err != nil {
		fmt.Println(printError(err))
	} else {
		fmt.Println(printSuccess(result))
	}

	readCommand()
}

func printError(errorResult *errorResult) string {
	return fmt.Sprintf("%s: %s", errorResult.cmd, errorResult.message)
}

func printSuccess(successResult *successResult) string {
	return fmt.Sprintf("%s", successResult.message)
}
