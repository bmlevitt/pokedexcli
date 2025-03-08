package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetPokemonData retrieves detailed information about a specific Pokemon
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
		return pokemonDataResp, nil
	} else {
		// fmt.Println("**cache miss**") // hide later
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return PokemonDataResp{}, err
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonDataResp{}, err
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return PokemonDataResp{}, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonDataResp{}, err
	}

	// Store in cache
	c.cache.Add(fullURL, body)

	// Unmarshal the response into the appropriate struct
	pokemonDataResp := PokemonDataResp{}
	err = json.Unmarshal(body, &pokemonDataResp)
	if err != nil {
		return PokemonDataResp{}, err
	}

	return pokemonDataResp, nil
}

// GetPokemonCaptureRate retrieves the capture rate for a specific Pokemon
func (c *Client) GetPokemonCaptureRate(pokemon string) (PokemonCaptureRateResp, error) {
	// First, we need to fetch the species URL from the pokemon data
	pokemonData, err := c.GetPokemonData(pokemon)
	if err != nil {
		return PokemonCaptureRateResp{}, err
	}

	// Extract the species URL
	speciesURL := pokemonData.Species.URL

	// Check cache
	data, ok := c.cache.Get(speciesURL)
	if ok {
		// fmt.Println("**cache hit**") // hide later
		var speciesResp map[string]interface{}
		err := json.Unmarshal(data, &speciesResp)
		if err != nil {
			return PokemonCaptureRateResp{}, err
		}

		// Extract the capture rate from the response
		captureRate, ok := speciesResp["capture_rate"].(float64)
		if !ok {
			return PokemonCaptureRateResp{}, fmt.Errorf("failed to extract capture rate")
		}

		return PokemonCaptureRateResp{CaptureRate: int(captureRate)}, nil
	} else {
		// fmt.Println("**cache miss**") // hide later
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", speciesURL, nil)
	if err != nil {
		return PokemonCaptureRateResp{}, err
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonCaptureRateResp{}, err
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return PokemonCaptureRateResp{}, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonCaptureRateResp{}, err
	}

	// Store in cache
	c.cache.Add(speciesURL, body)

	// Unmarshal the response into a map to extract the capture rate
	var speciesResp map[string]interface{}
	err = json.Unmarshal(body, &speciesResp)
	if err != nil {
		return PokemonCaptureRateResp{}, err
	}

	// Extract the capture rate from the response
	captureRate, ok := speciesResp["capture_rate"].(float64)
	if !ok {
		return PokemonCaptureRateResp{}, fmt.Errorf("failed to extract capture rate")
	}

	return PokemonCaptureRateResp{CaptureRate: int(captureRate)}, nil
}

// GetPokemonSpecies retrieves detailed species information about a specific Pokemon
func (c *Client) GetPokemonSpecies(pokemon string) (PokemonSpeciesResp, error) {
	endpoint := "/pokemon-species/"
	fullURL := baseURL + endpoint + pokemon

	// Check cache
	data, ok := c.cache.Get(fullURL)
	if ok {
		// fmt.Println("**cache hit**") // hide later
		speciesResp := PokemonSpeciesResp{}
		err := json.Unmarshal(data, &speciesResp)
		if err != nil {
			return PokemonSpeciesResp{}, err
		}
		return speciesResp, nil
	} else {
		// fmt.Println("**cache miss**") // hide later
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return PokemonSpeciesResp{}, err
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonSpeciesResp{}, err
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return PokemonSpeciesResp{}, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonSpeciesResp{}, err
	}

	// Store in cache
	c.cache.Add(fullURL, body)

	// Unmarshal the response into the appropriate struct
	speciesResp := PokemonSpeciesResp{}
	err = json.Unmarshal(body, &speciesResp)
	if err != nil {
		return PokemonSpeciesResp{}, err
	}

	return speciesResp, nil
}
