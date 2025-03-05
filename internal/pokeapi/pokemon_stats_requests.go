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
