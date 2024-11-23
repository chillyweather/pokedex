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

		command := normalizeInput(text)[0]

		switch command {
		case "exit":
			os.Exit(0)
		case "help":
			fmt.Println("")
			fmt.Println("Welcome to the Pokedex!")
			fmt.Println("")
			fmt.Println("Usage:")
			fmt.Println("")
			fmt.Println("help: Displays a help message")
			fmt.Println("exit: Exit the Pokedex")
			fmt.Println("")
		default:
			fmt.Println("Invalid command")
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
