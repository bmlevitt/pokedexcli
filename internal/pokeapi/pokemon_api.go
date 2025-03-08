package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
)

// GetPokemonData retrieves detailed information about a specific Pokemon
func (c *Client) GetPokemonData(pokemon string) (PokemonDataResp, error) {
	endpoint := "/pokemon/"
	fullURL := baseURL + endpoint + pokemon

	// Check cache
	data, ok := c.cache.Get(fullURL)
	if ok {
		pokemonDataResp := PokemonDataResp{}
		err := json.Unmarshal(data, &pokemonDataResp)
		if err != nil {
			return PokemonDataResp{}, fmt.Errorf("error unmarshaling cached pokemon data: %w", err)
		}
		return pokemonDataResp, nil
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return PokemonDataResp{}, errorhandling.NewNetworkError("Failed to create HTTP request", err)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonDataResp{}, errorhandling.NewNetworkError("Failed to connect to the Pokémon API", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.StatusCode == http.StatusNotFound {
			return PokemonDataResp{}, errorhandling.PokemonNotFoundError(pokemon, fmt.Errorf("HTTP 404"))
		}
		return PokemonDataResp{}, errorhandling.NewAPIError(resp.StatusCode, endpoint+pokemon, fmt.Errorf("HTTP error: %d", resp.StatusCode))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonDataResp{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Store in cache
	c.cache.Add(fullURL, body)

	// Unmarshal the response into the appropriate struct
	pokemonDataResp := PokemonDataResp{}
	err = json.Unmarshal(body, &pokemonDataResp)
	if err != nil {
		return PokemonDataResp{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return pokemonDataResp, nil
}

// GetPokemonCaptureRate retrieves the capture rate for a specific Pokemon
func (c *Client) GetPokemonCaptureRate(pokemon string) (PokemonCaptureRateResp, error) {
	// First, we need to fetch the species URL from the pokemon data
	pokemonData, err := c.GetPokemonData(pokemon)
	if err != nil {
		return PokemonCaptureRateResp{}, fmt.Errorf("error fetching pokemon data: %w", err)
	}

	// Extract the species URL
	speciesURL := pokemonData.Species.URL

	// Check cache
	data, ok := c.cache.Get(speciesURL)
	if ok {
		var speciesResp map[string]interface{}
		err := json.Unmarshal(data, &speciesResp)
		if err != nil {
			return PokemonCaptureRateResp{}, fmt.Errorf("error unmarshaling cached species data: %w", err)
		}

		// Extract the capture rate from the response
		captureRate, ok := speciesResp["capture_rate"].(float64)
		if !ok {
			return PokemonCaptureRateResp{}, fmt.Errorf("missing capture rate in species data")
		}

		return PokemonCaptureRateResp{CaptureRate: int(captureRate)}, nil
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", speciesURL, nil)
	if err != nil {
		return PokemonCaptureRateResp{}, errorhandling.NewNetworkError("Failed to create HTTP request", err)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonCaptureRateResp{}, errorhandling.NewNetworkError("Failed to connect to the Pokémon API", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		// Extract resource type from URL for better error messages
		resourceType := errorhandling.ResourcePokemonSpecies
		resourceName := pokemon

		if resp.StatusCode == http.StatusNotFound {
			return PokemonCaptureRateResp{}, errorhandling.FormatResourceNotFoundError(resourceType, resourceName, fmt.Errorf("HTTP 404"))
		}

		endpoint := fmt.Sprintf("species URL for %s", pokemon)
		return PokemonCaptureRateResp{}, errorhandling.NewAPIError(resp.StatusCode, endpoint, fmt.Errorf("HTTP error: %d", resp.StatusCode))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonCaptureRateResp{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Store in cache
	c.cache.Add(speciesURL, body)

	// Unmarshal the response into a map to extract the capture rate
	var speciesResp map[string]interface{}
	err = json.Unmarshal(body, &speciesResp)
	if err != nil {
		return PokemonCaptureRateResp{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Extract the capture rate from the response
	captureRate, ok := speciesResp["capture_rate"].(float64)
	if !ok {
		return PokemonCaptureRateResp{}, fmt.Errorf("missing capture rate in species data")
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
		speciesResp := PokemonSpeciesResp{}
		err := json.Unmarshal(data, &speciesResp)
		if err != nil {
			return PokemonSpeciesResp{}, fmt.Errorf("error unmarshaling cached species data: %w", err)
		}
		return speciesResp, nil
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return PokemonSpeciesResp{}, errorhandling.NewNetworkError("Failed to create HTTP request", err)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonSpeciesResp{}, errorhandling.NewNetworkError("Failed to connect to the Pokémon API", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.StatusCode == http.StatusNotFound {
			return PokemonSpeciesResp{}, errorhandling.FormatResourceNotFoundError(errorhandling.ResourcePokemonSpecies, pokemon, fmt.Errorf("HTTP 404"))
		}
		return PokemonSpeciesResp{}, errorhandling.NewAPIError(resp.StatusCode, endpoint+pokemon, fmt.Errorf("HTTP error: %d", resp.StatusCode))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonSpeciesResp{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Store in cache
	c.cache.Add(fullURL, body)

	// Unmarshal the response into the appropriate struct
	speciesResp := PokemonSpeciesResp{}
	err = json.Unmarshal(body, &speciesResp)
	if err != nil {
		return PokemonSpeciesResp{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return speciesResp, nil
}
