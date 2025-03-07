package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetPokemonData(pokemon string) (PokemonDataResp, error) {
	endpoint := "/pokemon/"
	fullURL := baseURL + endpoint + pokemon

	// Check cache
	data, ok := c.cache.Get(fullURL)
	if ok {
		// fmt.Println("**cache hit**") // hide later
		pokemonDataResp := PokemonDataResp{}
		err := json.Unmarshal(data, &pokemonDataResp)
		if err != nil {
			return PokemonDataResp{}, err
		}
	} else {
		// fmt.Println("**cache miss**") // hide later
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return PokemonDataResp{}, err
	}

	// Execute the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonDataResp{}, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode > 399 {
		return PokemonDataResp{}, fmt.Errorf("bad status received: %v", resp.StatusCode)
	}

	// Read the response body
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return PokemonDataResp{}, err
	}

	// Parse the JSON response
	pokemonDataResp := PokemonDataResp{}
	err = json.Unmarshal(data, &pokemonDataResp)
	if err != nil {
		return PokemonDataResp{}, err
	}

	c.cache.Add(fullURL, data)
	return pokemonDataResp, nil
}

func (c *Client) GetPokemonCaptureRate(pokemon string) (PokemonCaptureRateResp, error) {
	endpoint := "/pokemon-species/"
	fullURL := baseURL + endpoint + pokemon

	// Check cache
	data, ok := c.cache.Get(fullURL)
	if ok {
		// fmt.Println("**cache hit**") // hide later
		pokemonDataResp := PokemonDataResp{}
		err := json.Unmarshal(data, &pokemonDataResp)
		if err != nil {
			return PokemonCaptureRateResp{}, err
		}
	} else {
		// fmt.Println("**cache miss**") // hide later
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return PokemonCaptureRateResp{}, err
	}

	// Execute the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonCaptureRateResp{}, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode > 399 {
		return PokemonCaptureRateResp{}, fmt.Errorf("bad status received: %v", resp.StatusCode)
	}

	// Read the response body
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return PokemonCaptureRateResp{}, err
	}

	// Parse the JSON response
	pokemonCaptureRateResp := PokemonCaptureRateResp{}
	err = json.Unmarshal(data, &pokemonCaptureRateResp)
	if err != nil {
		return PokemonCaptureRateResp{}, err
	}

	c.cache.Add(fullURL, data)
	return pokemonCaptureRateResp, nil
}

// GetPokemonSpecies fetches detailed species information including flavor text and descriptions
func (c *Client) GetPokemonSpecies(pokemon string) (PokemonSpeciesResp, error) {
	endpoint := "/pokemon-species/"
	fullURL := baseURL + endpoint + pokemon

	// Check cache
	data, ok := c.cache.Get(fullURL)
	if ok {
		// fmt.Println("**cache hit**") // hide later
		pokemonSpeciesResp := PokemonSpeciesResp{}
		err := json.Unmarshal(data, &pokemonSpeciesResp)
		if err != nil {
			return PokemonSpeciesResp{}, err
		}
		return pokemonSpeciesResp, nil
	} else {
		// fmt.Println("**cache miss**") // hide later
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return PokemonSpeciesResp{}, err
	}

	// Execute the HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonSpeciesResp{}, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode > 399 {
		return PokemonSpeciesResp{}, fmt.Errorf("bad status received: %v", resp.StatusCode)
	}

	// Read the response body
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return PokemonSpeciesResp{}, err
	}

	// Parse the JSON response
	pokemonSpeciesResp := PokemonSpeciesResp{}
	err = json.Unmarshal(data, &pokemonSpeciesResp)
	if err != nil {
		return PokemonSpeciesResp{}, err
	}

	c.cache.Add(fullURL, data)
	return pokemonSpeciesResp, nil
}
