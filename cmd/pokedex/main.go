// cmd/pokedex/main.go
package main

import (
	"bufio"
	"fmt"
	"os"

	"pokedex/internal/cli"
	"pokedex/internal/config"
)

func main() {
	cfg := config.New()
	commands := cli.GetCommands()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")

	for {
		if scanner.Scan() {
			input := cli.CleanInput(scanner.Text())
			if len(input) == 0 {
				fmt.Println("Please enter a command")
			} else if cmd, ok := commands[input[0]]; ok {
				if err := cmd.Callback(cfg); err != nil {
					fmt.Printf("Error: %v\n", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
			fmt.Print("Pokedex > ")
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			break
		}
	}
}
