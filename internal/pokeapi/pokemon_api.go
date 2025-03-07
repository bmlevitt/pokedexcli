package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bmlevitt/pokedexcli/internal/errorhandling"
)

// GetPokemonData retrieves detailed information about a specific Pokémon from the PokeAPI.
// This function fetches comprehensive data including stats, types, moves, and more.
// Results are cached to improve performance and reduce API calls.
//
// The function first checks the cache for the requested data. If not found, it makes
// an HTTP request to the PokeAPI, processes the response, and stores it in the cache
// for future use. This caching strategy helps reduce API calls and improves performance.
//
// Parameters:
//   - pokemon: The name or ID of the Pokémon to retrieve (in lowercase with hyphens)
//
// Returns:
//   - A PokemonDataResp struct containing the Pokémon's data
//   - An error with specific error types:
//   - NetworkError: If there's an issue with creating or executing the HTTP request
//   - NotFoundError: If the requested Pokémon doesn't exist
//   - InternalError: If there's an issue parsing the API response
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

// GetPokemonCaptureRate retrieves the capture rate for a specific Pokémon.
// The capture rate is used in the game mechanics for determining how easy it is
// to catch a Pokémon, with higher values being easier to catch.
// This function fetches species data which includes the capture rate.
//
// Parameters:
//   - pokemon: The name or ID of the Pokémon (in lowercase with hyphens)
//
// Returns:
//   - A PokemonCaptureRateResp containing the capture rate value
//   - An error if the API request fails or the Pokémon doesn't exist
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

// GetPokemonSpecies retrieves detailed species information about a Pokémon.
// This includes Pokédex entries (flavor text), genus information, and evolution chain references.
// This data is used for the "describe" command and for evolution mechanics.
//
// Parameters:
//   - pokemon: The name or ID of the Pokémon species (in lowercase with hyphens)
//
// Returns:
//   - A PokemonSpeciesResp containing the species data
//   - An error if the API request fails or the species doesn't exist
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
