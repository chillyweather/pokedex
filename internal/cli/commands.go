// internal/cli/commands.go
package cli

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/chillyweather/pokedexcli/internal/config"
	"github.com/chillyweather/pokedexcli/internal/pokeapi"
	"github.com/chillyweather/pokedexcli/internal/pokecache"
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
		"explore": {
			Name:        "explore",
			Description: "displays pokemons in the area",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "making attempt to catch the choosen pokemon",
			Callback:    commandCatch,
		},
	}
}

var cache = pokecache.NewCache(5 * time.Second)

func commandCatch(c *config.Config) error {
	name := c.CurrentArgs[0]
	exp, err := pokeapi.FetchBaseExperience(name)
	if err != nil {
		return err
	}
	chance := rand.Float64()
	fmt.Printf("The %s have base experience of %d and your chance to catch it is %f \n", name, exp, chance)
	return nil
}

func commandExplore(c *config.Config) error {
	location := c.CurrentArgs[0]
	err := pokeapi.FetchPokemons(location)
	if err != nil {
		return err
	}
	return nil
}

func commandMap(c *config.Config) error {
	var data pokeapi.LocationAreaResponse

	cachedData, ok := cache.Get(c.Next)
	if !ok {
		fetchedData, err := pokeapi.FetchLocations(c.Next)
		if err != nil {
			return err
		}
		data = fetchedData

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cache.Add(c.Next, jsonData)
	} else {
		if err := json.Unmarshal(cachedData, &data); err != nil {
			return err
		}
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

	var data pokeapi.LocationAreaResponse

	cachedData, ok := cache.Get(c.Previous)
	if !ok {
		fetchedData, err := pokeapi.FetchLocations(c.Previous)
		if err != nil {
			return err
		}
		data = fetchedData

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cache.Add(c.Previous, jsonData)
	} else {
		if err := json.Unmarshal(cachedData, &data); err != nil {
			return err
		}
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
