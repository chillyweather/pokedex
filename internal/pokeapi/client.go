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

func FetchPokemons(url string) ([]Pokemon, error) {
	result := []Pokemon{}
	res, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	if res.StatusCode > 299 {
		return result, fmt.Errorf("failed with status %d: %s", res.StatusCode, body)
	}

	var data LocationArea
	fmt.Println(body)
	if err := json.Unmarshal(body, &data); err != nil {
		return result, err
	}

	for _, encounter := range data.PokemonEncounters {
		result = append(result, encounter.Pokemon)
	}

	if len(result) > 0 {
		for key := range result {
			fmt.Println(key)
		}
	}

	return result, nil
}
