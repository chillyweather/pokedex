package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		text := scanner.Text()

		if len(text) == 0 {
			continue
		}

		commandName := normalizeInput(text)[0]

		availableCommands := getCommands()

		command, ok := availableCommands[commandName]

		if !ok {
			fmt.Printf("Invalid command - %v\n", commandName)
			continue
		}

		err := command.callback()
		if err != nil {
			fmt.Println(err)
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "exit from the app",
			callback:    handleExit,
		},
		"help": {
			name:        "help",
			description: "prints the help menu",
			callback:    handleHelp,
		},
	}
}

func normalizeInput(text string) []string {
	normalizedString := strings.ToLower(text)
	stringArr := strings.Fields(normalizedString)
	return stringArr
}
