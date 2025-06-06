// internal/cli/commands.go
package cli

import (
	"fmt"
	"os"

	"github.com/chillyweather/pokedexcli/internal/config"
	"github.com/chillyweather/pokedexcli/internal/pokeapi"
)

type Command struct {
	Name        string
	Description string
	Callback    func(*config.Config) error
}

func GetCommands() map[string]Command {
	return map[string]Command{
		"exit": {
			Name:        "exit",
			Description: "exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			Name:        "help",
			Description: "lists available commands",
			Callback:    commandHelp,
		},
		"map": {
			Name:        "map",
			Description: "displays the next 20 locations",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "displays the previous 20 locations",
			Callback:    commandMapb,
		},
	}
}

func commandMap(c *config.Config) error {
	data, err := pokeapi.Fetch(c.Next)
	if err != nil {
		return err
	}
	c.Next = data.Next
	c.Previous = data.Previous
	for _, r := range data.Results {
		fmt.Println(r.Name)
	}
	return nil
}

func commandMapb(c *config.Config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	data, err := pokeapi.Fetch(c.Previous)
	if err != nil {
		return err
	}
	c.Next = data.Next
	c.Previous = data.Previous
	for _, r := range data.Results {
		fmt.Println(r.Name)
	}
	return nil
}

func commandHelp(c *config.Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range GetCommands() {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}

func commandExit(c *config.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
