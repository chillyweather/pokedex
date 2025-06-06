package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
}

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "lists available commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays the previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func commandMap(c *Config) error {
	res, err := http.Get(c.Next)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	locationAreaResponse := LocationAreaResponse{}
	json.Unmarshal(body, &locationAreaResponse)

	c.Next = locationAreaResponse.Next
	c.Previous = locationAreaResponse.Previous

	results := locationAreaResponse.Results
	for _, r := range results {
		fmt.Println(r.Name)
	}

	return nil
}

func commandMapb(c *Config) error {
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	res, err := http.Get(c.Previous)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	locationAreaResponse := LocationAreaResponse{}
	json.Unmarshal(body, &locationAreaResponse)

	c.Next = locationAreaResponse.Next
	c.Previous = locationAreaResponse.Previous

	results := locationAreaResponse.Results
	for _, r := range results {
		fmt.Println(r.Name)
	}

	return nil
}

func commandHelp(c *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for i := range commands {
		fmt.Printf("%v: %v\n", i, commands[i].description)
	}
	return nil
}

func commandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

var config = Config{
	Next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
	Previous: "",
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")

	for {
		if scanner.Scan() {
			text := scanner.Text()
			cleanedTextSlice := cleanInput(text)

			if len(cleanedTextSlice) == 0 {
				fmt.Println("Please enter a command")
				fmt.Print("Pokedex > ")
				continue
			}

			cmd, ok := commands[cleanedTextSlice[0]]
			if !ok {
				fmt.Println("Unknown command")
			} else {
				err := cmd.callback(&config)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
			}

			fmt.Print("Pokedex > ")
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			break
		}
	}
}

func cleanInput(text string) []string {
	lowerCaseString := strings.ToLower(text)
	result := strings.Fields(lowerCaseString)
	return result
}
