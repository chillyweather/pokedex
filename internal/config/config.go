package config

import "github.com/chillyweather/pokedexcli/internal/pokeapi"

type Config struct {
	Next          string
	Previous      string
	CurrentArgs   []string
	CaughtPokemon map[string]pokeapi.Pokemon
}

func New() *Config {
	return &Config{
		Next:          "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Previous:      "",
		CurrentArgs:   []string{},
		CaughtPokemon: make(map[string]pokeapi.Pokemon),
	}
}
