package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchBaseExperience(name string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name

	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	if res.StatusCode > 299 {
		return Pokemon{}, fmt.Errorf("failed with status %d: %s", res.StatusCode, body)
	}

	var data Pokemon

	if err := json.Unmarshal(body, &data); err != nil {
		return Pokemon{}, err
	}

	return data, nil
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
	result := []PokemonBasicData{}
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
