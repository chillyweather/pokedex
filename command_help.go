package main

import "fmt"

func handleHelp() error {
	fmt.Println("")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")

	availableCommands := getCommands()
	for _, command := range availableCommands {
		fmt.Printf(" - %v: %v\n", command.name, command.description)
	}
	return nil
}
