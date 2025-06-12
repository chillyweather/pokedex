package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationArea struct {
	Name              string             `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonProperties struct {
	BaseExperience int `json:"base_experience"`
}

func FetchBaseExperience(name string) (int, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + name

	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	if res.StatusCode > 299 {
		return 0, fmt.Errorf("failed with status %d: %s", res.StatusCode, body)
	}

	var data PokemonProperties

	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	return data.BaseExperience, nil
}

func FetchLocations(url string) (LocationAreaResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	if res.StatusCode > 299 {
		return LocationAreaResponse{}, fmt.Errorf("failed with status %d: %s", res.StatusCode, body)
	}

	var data LocationAreaResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return LocationAreaResponse{}, err
	}

	return data, nil
}

func FetchPokemons(location string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + location
	result := []Pokemon{}
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("failed with status %d: %s", res.StatusCode, body)
	}

	var data LocationArea
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	for _, encounter := range data.PokemonEncounters {
		result = append(result, encounter.Pokemon)
	}

	if len(result) > 0 {
		fmt.Println("Exploring pastoria-city-area...")
		fmt.Println("Found Pokemon:")
		for _, val := range result {
			fmt.Printf(" - %s\n", val.Name)
		}
	}

	return nil
}
