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
		"inspect": {
			Name:        "inspect",
			Description: "inspects the chosen pokemon",
			Callback:    commandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "lists caught pokemons",
			Callback:    commandPokedex,
		},
	}
}

var cache = pokecache.NewCache(5 * time.Second)

func commandPokedex(c *config.Config) error {
	caughtPokemon := c.CaughtPokemon
	if len(caughtPokemon) == 0 {
		return fmt.Errorf("you have nothing to show")
	}
	fmt.Println("Your Pokedex:")
	for pokemon := range caughtPokemon {
		fmt.Printf(" - %v\n", pokemon)
	}
	return nil
}

func commandInspect(c *config.Config) error {
	pokemonName := c.CurrentArgs[0]
	pokemonData, ok := c.CaughtPokemon[pokemonName]
	if !ok {
		fmt.Printf("You need to catch %s first!", pokemonName)
	}
	fmt.Printf("Name: %s\n", pokemonData.Name)
	fmt.Printf("Height: %d\n", pokemonData.Height)
	fmt.Printf("Weight: %d\n", pokemonData.Weight)
	fmt.Println("Stats:")
	fmt.Printf("  -hp: %d\n", pokemonData.Stats[0].BaseStat)
	fmt.Printf("  -attack: %d\n", pokemonData.Stats[1].BaseStat)
	fmt.Printf("  -defense: %d\n", pokemonData.Stats[2].BaseStat)
	fmt.Printf("  -special-attack: %d\n", pokemonData.Stats[3].BaseStat)
	fmt.Printf("  -special-defense: %d\n", pokemonData.Stats[4].BaseStat)
	fmt.Printf("  -speed: %d\n", pokemonData.Stats[5].BaseStat)
	fmt.Println("Types:")
	for _, t := range pokemonData.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}

	return nil
}

func commandCatch(c *config.Config) error {
	pokemonName := c.CurrentArgs[0]
	fmt.Printf("Throwing a Pokeball at %s...", pokemonName)
	pokemonData, err := pokeapi.FetchBaseExperience(pokemonName)
	if err != nil {
		return fmt.Errorf("failed to fetch pokemon %s, %v", pokemonName, err)
	}

	catchProbability := float64(pokemonData.BaseExperience) / 255.0
	if rand.Float64() > catchProbability {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	c.CaughtPokemon[pokemonName] = pokemonData
	fmt.Printf("%s was caught!\n", pokemonName)
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
